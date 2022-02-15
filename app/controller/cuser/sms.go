package cuser

import (
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/muser"
	"sports_service/server/models/sms"
	"sports_service/server/tools/tencentCloud"
	"sports_service/server/util"
	"strings"
	"time"
	"fmt"
)

var (
	randSource = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// 发送短信验证码
func (svc *UserModule) SendSmsCode(params *sms.SendSmsCodeParams) int {
	// 校验手机合法性
	if b := svc.user.CheckCellPhoneNumber(params.MobileNum); !b {
		log.Log.Errorf("sms_trace: invalid mobile num %v", params.MobileNum)
		return errdef.USER_INVALID_MOBILE_NUM
	}

	// 获取短信模版
	//mod := svc.sms.GetSendMod(params.SendType)
	//if mod == "" {
	//	log.Log.Error("sms_trace: invalid send type" )
	//	return errdef.SMS_CODE_INVALID_SEND_TYPE
	//}

	// 获取当前手机号发送的短信验证码次数
	limitNum, err := svc.sms.GetSendSmsLimitNum(params.MobileNum)
	if err != nil && err != redis.ErrNil {
		log.Log.Errorf("sms_trace: get send sms limit num err:%s", err)
		return errdef.ERROR
	}

	var newDay bool
	// 限制短信次数的key不存在 则表示为新的一天
	if err != nil && err == redis.ErrNil {
		newDay = true
	}

	// 最多一天内十条 todo: 需改回限制10条
	if limitNum >= consts.SMS_INTERVAL_NUM {
		log.Log.Errorf("sms_trace: send sms code interval, mobileNum:%s", params.MobileNum)
		return errdef.SMS_CODE_INTERVAL_ERROR
	}

	// 是否已过重发验证码的间隔时间
	ok, err := svc.sms.HasSmsIntervalPass(params.SendType, params.MobileNum)
	if err != nil && err != redis.ErrNil {
		log.Log.Errorf("sms_trace: has sms interval pass err:%s", err)
		return errdef.ERROR
	}

	if ok {
		return errdef.SMS_CODE_INTERVAL_SHORT
	}

	// 获取reids存储的验证码
	code, err := svc.sms.GetSmsCodeByRds(params.SendType, params.MobileNum)
	if err != nil && err != redis.ErrNil {
		log.Log.Errorf("sms_trace: get sms code by redis err:%s", err)
		return errdef.ERROR
	}

	// 重新获取验证码 并存储到redis
	code = svc.sms.GetSmsCode()
	if err := svc.sms.SaveSmsCodeByRds(params.SendType, params.MobileNum, code); err != nil {
		log.Log.Errorf("sms_trace: save sms code by redis err:%s", err)
		return errdef.ERROR
	}

  // 发送短信验证码
  if err := svc.sms.Send(params.MobileNum, code); err != nil {
    log.Log.Errorf("sms_trace: send sms code err:%s", err)
    return errdef.SMS_CODE_SEND_FAIL
  }

  // 重发验证码的限制时间
  if err := svc.sms.SetSmsIntervalTm(params.SendType, params.MobileNum); err != nil {
    log.Log.Errorf("sms_trace: set sms interval tm err:%s", err)
    return errdef.ERROR
  }

	// 记录发送验证码的次数
	if err := svc.sms.IncrSendSmsNum(params.MobileNum); err != nil {
		log.Log.Errorf("sms_trace: incr send sms code num err:%s", err)
		return errdef.ERROR
	}

	// 新的一天 给记录每日短信数量的key过期时间
	if newDay {
		if _, err := svc.sms.SetSmsIntervalExpire(params.MobileNum); err != nil {
			log.Log.Errorf("sms_trace: set sms interval expire err:%s", err)
		}
	}

	return errdef.SUCCESS
}

// 短信验证码登陆
func (svc *UserModule) SmsCodeLogin(params *sms.SmsCodeLoginParams) (int, string, *models.User) {
	if len(params.MobileNum) == 11 && params.MobileNum[0:8] == "18888888" {
		return svc.RegisterOfficialAccount(params)
	}

	// 校验手机号合法性
	if b := svc.user.CheckCellPhoneNumber(params.MobileNum); !b {
		log.Log.Errorf("sms_trace: invalid mobile num %v", params.MobileNum)
		return errdef.USER_INVALID_MOBILE_NUM, "", nil
	}

	// 校验手机验证码
	if syscode := svc.CheckSmsCode(params.Code, params.MobileNum); syscode != errdef.SUCCESS {
		return syscode, "", nil
	}

	// 根据手机号查询用户 不存在 注册用户 用户存在 为登陆
	user := svc.user.FindUserByPhone(params.MobileNum)
	if user == nil {
		// 开启事务
		if err := svc.engine.Begin(); err != nil {
			log.Log.Errorf("user_trace: session begin err:%s", err)
			return errdef.ERROR, "", nil
		}

		// 注册
		reg := muser.NewMobileRegister()
		if err := reg.Register(svc.user, params.Platform, params.MobileNum, svc.context.ClientIP()); err != nil {
			log.Log.Errorf("user_trace: register err:%s", err)
			svc.engine.Rollback()
			return errdef.USER_REGISTER_FAIL, "", nil
		}

		log.Log.Debugf("userInfo:%s", user)

		// 添加用户
		if err := svc.user.AddUser(); err != nil {
			log.Log.Errorf("reg_trace: add user info err:%s", err)
			svc.engine.Rollback()
			return errdef.USER_ADD_INFO_FAIL, "", nil
		}

		// 添加用户初始化通知设置
		if err := svc.notify.AddUserNotifySetting(svc.user.User.UserId, int(time.Now().Unix())); err != nil {
			log.Log.Errorf("user_trace: add user notify setting err:%s", err)
			svc.engine.Rollback()
			return errdef.USER_ADD_NOTIFY_SET_FAIL, "", nil
		}

		svc.engine.Commit()

		// 生成token
		token := svc.user.GenUserToken(svc.user.User.UserId, util.Md5String(svc.user.User.Password))
		// 保存到redis
		if err := svc.user.SaveUserToken(svc.user.User.UserId, token); err != nil {
			log.Log.Errorf("user_trace: save user token err:%s", err)
		}
		
		avatar := tencentCloud.BucketURI(svc.user.User.Avatar)
		svc.user.User.Avatar = string(avatar)
		return errdef.SUCCESS, token, svc.user.User

	}

	
	// 登陆的时候 检查用户状态
	if !svc.CheckUserStatus(svc.user.User.Status) {
		log.Log.Errorf("user_trace: forbid status, userId:%s", svc.user.User.UserId)
		return errdef.USER_FORBID_STATUS, "", nil
	}
	
	avatar := tencentCloud.BucketURI(svc.user.User.Avatar)
	svc.user.User.Avatar = string(avatar)
	// 用户已注册过, 则直接从redis中获取token并返回
	token, err := svc.user.GetUserToken(svc.user.User.UserId)
	if err != nil && err == redis.ErrNil {
		// redis 没有，重新生成token
		token = svc.user.GenUserToken(svc.user.User.UserId, util.Md5String(svc.user.User.Password))
		// 重新保存到redis
		if err := svc.user.SaveUserToken(svc.user.User.UserId, token); err != nil {
			log.Log.Errorf("user_trace: save user token err:%s", err)
		}
	}

	return errdef.SUCCESS, token, user
}

// todo: 临时注册官方账号 用于app内容填充 后续剔除
func (svc *UserModule) RegisterOfficialAccount(params *sms.SmsCodeLoginParams) (int, string, *models.User) {
	if params.Code != fmt.Sprintf("0110%s", params.MobileNum[9:11]) {
		return errdef.SMS_CODE_NOT_MATCH, "", nil
	}

	// 根据手机号查询用户 不存在 注册用户 用户存在 为登陆
	user := svc.user.FindUserByPhone(params.MobileNum)
	if user != nil {
		// 用户已注册过, 则直接从redis中获取token并返回
		token, err := svc.user.GetUserToken(svc.user.User.UserId)
		if err != nil && err == redis.ErrNil {
			// redis 没有，重新生成token
			token = svc.user.GenUserToken(svc.user.User.UserId, util.Md5String(svc.user.User.Password))
			// 重新保存到redis
			if err := svc.user.SaveUserToken(svc.user.User.UserId, token); err != nil {
				log.Log.Errorf("user_trace: save user token err:%s", err)
			}
		}

		return errdef.SUCCESS, token, user
	}

	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("user_trace: session begin err:%s", err)
		return errdef.ERROR, "", nil
	}
	// 注册
	reg := muser.NewMobileRegister()
	svc.user.User.AccountType = 1
	if err := reg.Register(svc.user, params.Platform, params.MobileNum, svc.context.ClientIP()); err != nil {
		log.Log.Errorf("user_trace: register err:%s", err)
		svc.engine.Rollback()
		return errdef.USER_REGISTER_FAIL, "", nil
	}

	// 添加用户
	if err := svc.user.AddUser(); err != nil {
		log.Log.Errorf("reg_trace: add user info err:%s", err)
		svc.engine.Rollback()
		return errdef.USER_ADD_INFO_FAIL, "", nil
	}

	// 添加用户初始化通知设置
	if err := svc.notify.AddUserNotifySetting(svc.user.User.UserId, int(time.Now().Unix())); err != nil {
		log.Log.Errorf("user_trace: add user notify setting err:%s", err)
		svc.engine.Rollback()
		return errdef.USER_ADD_NOTIFY_SET_FAIL, "", nil
	}

	svc.engine.Commit()

	// 生成token
	token := svc.user.GenUserToken(svc.user.User.UserId, util.Md5String(svc.user.User.Password))
	// 保存到redis
	if err := svc.user.SaveUserToken(svc.user.User.UserId, token); err != nil {
		log.Log.Errorf("user_trace: save user token err:%s", err)
	}

	return errdef.SUCCESS, token, svc.user.User
}

// 检查用户状态
func (svc *UserModule) CheckUserStatus(status int) bool {
  // 如果用户状态不正常
  if status != consts.USER_NORMAL {
    return false
  }

  return true
}

// 校验短信验证码
func (svc *UserModule) CheckSmsCode(code, mobileNum string) int {
	if len(code) != 6 || code == "" {
		log.Log.Errorf("sms_trace: invalid sms code, code:%s, mobileNum:%s", code, mobileNum)
		return errdef.SMS_INVALID_CODE
	}

	// 从redis获取发送的验证码
	rCode, err := svc.sms.GetSmsCodeByRds(consts.ACCOUNT_OPT_TYPE, mobileNum)
	if err != nil && err == redis.ErrNil {
		log.Log.Errorf("sms_trace: get sms code by redis err:%s", err)
		return errdef.SMS_CODE_NOT_SEND
	}

	if strings.Compare(rCode, code) != 0 {
		log.Log.Errorf("sms_trace: sms code not match, redis code:%s, code:%s, mobileNum:%s", rCode, code, mobileNum)
		return errdef.SMS_CODE_NOT_MATCH
	}

	// 删除存储验证码的key
	if err := svc.sms.DelSmsCodeKey(consts.ACCOUNT_OPT_TYPE, mobileNum); err != nil {
		log.Log.Errorf("sms_trace: del sms code key err:%s", err)
		return errdef.ERROR
	}
	// 删除限制重发验证码间隔时间的key
	if err := svc.sms.DelSmsIntervalTmKey(consts.ACCOUNT_OPT_TYPE, mobileNum); err != nil {
		log.Log.Errorf("sms_trace: del sms interval time key err:%s", err)
		return errdef.ERROR
	}

	return errdef.SUCCESS

}

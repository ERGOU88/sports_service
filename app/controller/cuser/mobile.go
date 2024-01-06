package cuser

import (
	"github.com/garyburd/redigo/redis"
	"sports_service/global/app/errdef"
	"sports_service/global/app/log"
	"sports_service/global/consts"
	"sports_service/models/muser"
	"sports_service/tools/tencentCloud"
	"sports_service/util"
	"time"
)

// 手机一键登陆/注册
func (svc *UserModule) MobileLoginOrReg(param *muser.LoginParams) (int, string, *muser.User) {
	//mob := mobTech.NewMobTech()
	// 校验客户端token 并从mob获取手机号码
	//mobileNum, err := mob.FreeLogin(param.Token, param.OpToken, param.Operator)
	//if err != nil {
	//	log.Log.Errorf("user_trace: mob free login err:%s", err)
	//	return errdef.USER_FREE_LOGIN_FAIL, "", nil
	//}
	// 替换为腾讯一键登录
	client := tencentCloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, "")
	mobileNum, err := client.FreeLogin(param.OpToken, param.Operator, "86")
	if err != nil {
		log.Log.Errorf("user_trace: tencent free login err:%s", err)
		return errdef.USER_FREE_LOGIN_FAIL, "", nil
	}

	log.Log.Infof("user_trace: mobileNum:%s", mobileNum)

	//mobileNum := "13177666666"

	// 校验手机号合法性
	if b := svc.user.CheckCellPhoneNumber(mobileNum); !b {
		log.Log.Errorf("user_trace: invalid mobile num %v", mobileNum)
		return errdef.USER_INVALID_MOBILE_NUM, "", nil
	}

	// 根据手机号查询用户 不存在 注册用户 用户存在 为登陆
	user := svc.user.FindUserByPhone(mobileNum)
	if user == nil {
		// 开启事务
		if err := svc.engine.Begin(); err != nil {
			log.Log.Errorf("user_trace: session begin err:%s", err)
			return errdef.ERROR, "", nil
		}
		// 注册
		reg := muser.NewMobileRegister()
		if err := reg.Register(svc.user, param.Platform, mobileNum, svc.context.ClientIP()); err != nil {
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

		return errdef.SUCCESS, token, svc.UserInfoResp()
	}

	// 登陆的时候 检查用户状态
	if !svc.CheckUserStatus(user.Status) {
		log.Log.Errorf("user_trace: forbid status, userId:%s", user.Status)
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

	return errdef.SUCCESS, token, svc.UserInfoResp()
}

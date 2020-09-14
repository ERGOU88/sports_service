package cuser

import (
	"github.com/garyburd/redigo/redis"
	"sports_service/server/global/app/errdef"
	"sports_service/server/models"
	"sports_service/server/models/muser"
	"sports_service/server/util"
	"sports_service/server/global/app/log"
	"sports_service/server/tools/mobtech"
)

// 手机一键登陆/注册
func (svc *UserModule) MobileLoginOrReg(param *muser.LoginParams) (int, string, *models.User) {
	mob := new(mobtech.MobTech)
	// 从mob获取手机号码
	mobileNum, err := mob.FreeLogin(param.Token, param.OpToken, param.Operator)
	if err != nil {
		return errdef.FREE_LOGIN_FAIL, "", nil
	}

	// 校验手机号合法性
	if b := svc.user.CheckCellPhoneNumber(mobileNum); !b {
		log.Log.Errorf("user_trace: invalid mobile num %v", mobileNum)
		return errdef.INVALID_MOBILE_NUM, "", nil
	}

	// 根据手机号查询用户 不存在 注册用户 用户存在 为登陆
	if user := svc.user.FindUserByPhone(mobileNum); user == nil {
		// 注册
		reg := muser.NewMobileRegister()
		if err := reg.Register(svc.user, param, mobileNum); err != nil {
			log.Log.Errorf("user_trace: register err:%s", err)
			return errdef.USER_REGISTER_FAIL, "", nil
		}

		// 添加用户
		if err := svc.user.AddUser(); err != nil {
			log.Log.Errorf("reg_trace: add user info err:%s", err)
			return errdef.USER_ADD_INFO_FAIL, "", nil
		}

		// 生成token
		token := svc.user.GenUserToken(svc.user.User.UserId, util.Md5String(svc.user.User.Password))
		// 保存到redis
		if err := svc.user.SaveUserToken(svc.user.User.UserId, token); err != nil {
			log.Log.Errorf("user_trace: save user token err:%s", err)
		}

		return errdef.SUCCESS, token, svc.user.User

	}

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

	return errdef.SUCCESS, token, svc.user.User
}

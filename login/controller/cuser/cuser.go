package cuser

import (
	"github.com/garyburd/redigo/redis"
	"sports_service/server/global/login/log"
	"sports_service/server/models/muser"
	"sports_service/server/global/login/errdef"
	"sports_service/server/util"
	"sports_service/server/models"
)

// 手机一键登陆/注册
func MobileLoginOrReg(param *muser.LoginParams) (int, string, *models.User) {
	model := muser.NewUserModel()
	if b := model.CheckCellPhoneNumber(param.MobileNum); !b {
		log.Log.Errorf("user_trace: invalid mobile num %v", param.MobileNum)
		return errdef.INVALID_MOBILE_NUM, "", nil
	}

	// 根据手机号查询用户 不存在 注册用户 用户存在 为登陆
	if user := model.FindUserByPhone(param.MobileNum); user == nil {
		// 登陆
		reg := muser.NewPhoneRegister()
		if err := reg.Register(model, param); err != nil {
			log.Log.Errorf("user_trace: register err:%s", err)
			return errdef.USER_REGISTER_FAIL, "", nil
		}

		// 生成token
		token := model.GenUserToken(model.User.UserId, util.Md5String(model.User.Password))
		// 保存到redis
		if err := model.SaveUserToken(model.User.UserId, token); err != nil {
			log.Log.Errorf("user_trace: save user token err:%s", err)
		}

		return errdef.SUCCESS, token, model.User

	}

	// 用户已注册过, 则直接从redis中获取token并返回
	token, err := model.GetUserToken(model.User.UserId)
	if err != nil && err == redis.ErrNil {
		// redis 没有，重新生成token
		token = model.GenUserToken(model.User.UserId, util.Md5String(model.User.Password))
		// 重新保存到redis
		model.SaveUserToken(model.User.UserId, token)
	}

	return errdef.SUCCESS, token, model.User



}

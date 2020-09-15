package cuser

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sports_service/server/global/consts"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/models"
	"sports_service/server/models/muser"
	"sports_service/server/util"
	third "sports_service/server/tools/thirdLogin"
)

// 微信登陆/注册
func (svc *UserModule) WeiboLoginOrReg(params *muser.WeiboLoginParams) (int, string, *models.User) {
	info := svc.social.GetSocialAccountByType(consts.TYPE_WEIBO, fmt.Sprint(params.Uid))
	// 微博信息不存在 则注册
	if info == nil {
		weibo := third.NewWeibo()
		weiboInfo := weibo.GetWeiboUserInfo(params.Uid, params.AccessToken)
		if weiboInfo == nil {
			log.Log.Errorf("weibo_trace: get weibo user info error")
			return errdef.WEIBO_USER_INFO_FAIL, "", nil
		}

		r := muser.NewWeiboRegister()
		// 替换成app性别标识
		gender := svc.WeiboInterChangeGender(weiboInfo.Gender)
		// 注册
		if err := r.Register(svc.user, svc.social, svc.context, fmt.Sprint(params.Uid), weiboInfo.AvatarLarge,
			weiboInfo.Name, consts.TYPE_WEIBO, gender); err != nil {
			log.Log.Errorf("weibo_trace: register err:%s", err)
			return errdef.WEIBO_REGISTER_FAIL, "", nil
		}

		// 开启事务
		svc.engine.Begin()
		// 添加用户
		if err := svc.user.AddUser(); err != nil {
			log.Log.Errorf("weibo_trace: add user err:%s", err)
			svc.engine.Rollback()
			return errdef.USER_ADD_INFO_FAIL, "", nil
		}

		// 添加社交帐号（微博）
		if err := svc.social.AddSocialAccountInfo(); err != nil {
			log.Log.Errorf("weibo_trace: add wx account err:%s", err)
			svc.engine.Rollback()
			return errdef.WEIBO_ADD_ACCOUNT_FAIL, "", nil
		}

		// 提交事务
		if err := svc.engine.Commit(); err != nil {
			log.Log.Errorf("weibo_trace: commit transaction err:%s", err)
			return errdef.ERROR, "", nil
		}

		// 生成token
		token := svc.user.GenUserToken(svc.user.User.UserId, util.Md5String(svc.user.User.Password))
		// 保存到redis
		if err := svc.user.SaveUserToken(svc.user.User.UserId, token); err != nil {
			log.Log.Errorf("weibo_trace: save user token err:%s", err)
		}

		return errdef.SUCCESS, token, svc.user.User
	}

	// 通过uid查询用户信息
	if user := svc.user.FindUserByUserid(info.UserId); user == nil {
		return errdef.USER_GET_INFO_FAIL, "", nil
	}
	// 用户已注册过, 则直接从redis中获取token并返回
	token, err := svc.user.GetUserToken(svc.user.User.UserId)
	if err != nil && err == redis.ErrNil {
		// redis 没有，重新生成token
		token = svc.user.GenUserToken(svc.user.User.UserId, util.Md5String(svc.user.User.Password))
		// 重新保存到redis
		if err := svc.user.SaveUserToken(svc.user.User.UserId, token); err != nil {
			log.Log.Errorf("weibo_trace: save user token err:%s", err)
		}
	}

	return errdef.SUCCESS, token, svc.user.User
}

// 微博性别替换成app的性别标识
func (svc *UserModule) WeiboInterChangeGender(gender string) int {
	switch gender {
	case "m":
		return consts.BOY
	case "f":
		return consts.GIRL
	default:
		return consts.BOY_OR_GIRL
	}
}

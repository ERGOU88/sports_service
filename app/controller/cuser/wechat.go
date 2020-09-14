package cuser

import (
	"github.com/garyburd/redigo/redis"
	"sports_service/server/global/consts"
	"sports_service/server/global/app/errdef"
	"sports_service/server/models"
	"sports_service/server/models/muser"
	"sports_service/server/global/app/log"
	"sports_service/server/util"
	third "sports_service/server/tools/third-login"
)

// 微信登陆/注册
func (svc *UserModule) WechatLoginOrReg(code string) (int, string, *models.User) {
	wechat := third.NewWechat()
	// 获取微信 access token
	accessToken := wechat.GetWechatAccessToken(code)
	if accessToken == nil {
		log.Log.Error("wx_trace: get wx access info error")
		return errdef.WX_ACCESS_TOKEN_FAIL, "", nil
	}

	// 数据库获取社交平台帐号信息
	info := svc.social.GetSocialAccountByType(consts.TYPE_WECHAT, accessToken.Unionid)
	// 微信信息不存在 则注册
	if info == nil {
		// 获取微信用户信息
		wxinfo := wechat.GetWechatUserInfo(accessToken)
		if wxinfo == nil {
			log.Log.Error("wx_trace: get wx user info error")
			return errdef.WX_USER_INFO_FAIL, "", nil
		}

		r := muser.NewWechatRegister()
		// 注册
		if err := r.Register(svc.user, svc.social, svc.context, accessToken.Unionid, wxinfo.Headimgurl, wxinfo.Nickname,
			consts.TYPE_WECHAT, wxinfo.Sex); err != nil {
			log.Log.Errorf("wx_trace: register err:%s", err)
			return errdef.WX_REGISTER_FAIL, "", nil
		}

		// 开启事务
		svc.engine.Begin()
		// 添加用户
		if err := svc.user.AddUser(); err != nil {
			log.Log.Errorf("wx_trace: add user err:%s", err)
			svc.engine.Rollback()
			return errdef.USER_ADD_INFO_FAIL, "", nil
		}

		// 添加社交帐号（微信）
		if err := svc.social.AddSocialAccountInfo(); err != nil {
			log.Log.Errorf("wx_trace: add wx account err:%s", err)
			svc.engine.Rollback()
			return errdef.WX_ADD_ACCOUNT_FAIL, "", nil
		}

		// 提交事务
		if err := svc.engine.Commit(); err != nil {
			log.Log.Errorf("wx_trace: commit transaction err:%s", err)
			return errdef.ERROR, "", nil
		}

		// 生成token
		token := svc.user.GenUserToken(svc.user.User.UserId, util.Md5String(svc.user.User.Password))
		// 保存到redis
		if err := svc.user.SaveUserToken(svc.user.User.UserId, token); err != nil {
			log.Log.Errorf("wx_trace: save user token err:%s", err)
		}

		return errdef.SUCCESS, token, svc.user.User
	}

	// 通过uid查询用户信息
	if user := svc.user.FindUserByUserid(info.UserId); user == nil {
		log.Log.Error("wx_trace: find user by uid error")
		return errdef.USER_GET_INFO_FAIL, "", nil
	}
	// 用户已注册过, 则直接从redis中获取token并返回
	token, err := svc.user.GetUserToken(svc.user.User.UserId)
	if err != nil && err == redis.ErrNil {
		// redis 没有，重新生成token
		token = svc.user.GenUserToken(svc.user.User.UserId, util.Md5String(svc.user.User.Password))
		// 重新保存到redis
		if err := svc.user.SaveUserToken(svc.user.User.UserId, token); err != nil {
			log.Log.Errorf("wx_trace: save user token err:%s", err)
		}
	}

	return errdef.SUCCESS, token, svc.user.User
}

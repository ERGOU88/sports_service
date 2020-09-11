package cuser

import (
	"github.com/garyburd/redigo/redis"
	"github.com/parnurzeal/gorequest"
	"net/url"
	"sports_service/server/global/consts"
	"sports_service/server/global/login/errdef"
	"sports_service/server/models"
	"sports_service/server/models/muser"
	"sports_service/server/login/config"
	"sports_service/server/global/login/log"
	"sports_service/server/util"
)

// 微信登陆/注册
func (svc *UserModule) WechatLoginOrReg(code string) (int, string, *models.User) {
	// 获取微信 access token
	accessToken := svc.WechatAccessToken(code)
	if accessToken == nil {
		return errdef.WX_ACCESS_TOKEN_FAIL, "", nil
	}

	// 数据库获取社交平台帐号信息
	info := svc.social.GetSocialAccountByType(consts.TYPE_WECHAT, accessToken.Unionid)
	// 微信信息不存在 则注册
	if info == nil {
		// 获取微信用户信息
		wxinfo := svc.WechatInfo(accessToken)
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

// WxAccessToken 获取微信accessToken
func (svc *UserModule) WechatAccessToken(code string) *muser.AccessToken {
	v := url.Values{}
	v.Set("code", code)
	// 开放平台appid
	v.Set("appid", config.Global.WechatAppid)
	// 开放平台secret
	v.Set("secret", config.Global.WechatSecret)
	v.Set("grant_type", "authorization_code")
	// 返回值
	accessToken := muser.AccessToken{}
	resp, body, errs := gorequest.New().Get(consts.WECHAT_ACCESS_TOKEN_URL + v.Encode()).EndStruct(&accessToken)
	if errs != nil {
		log.Log.Errorf("%+v", errs)
		return nil
	}

	if accessToken.Unionid == "" {
		log.Log.Errorf("err body: %s, resp: %+v", string(body), resp)
		return nil
	}

	return &accessToken
}

// 微信用户信息
func (svc *UserModule) WechatInfo(accessToken *muser.AccessToken) *muser.WechatUserInfo {
	v := url.Values{}
	v.Set("access_token", accessToken.AccessToken)
	v.Set("openid", accessToken.Openid)
	wxinfo := muser.WechatUserInfo{}
	resp, body, errs := gorequest.New().Get(consts.WECHAT_USER_INFO_URL + v.Encode()).EndStruct(&wxinfo)
	if errs != nil {
		log.Log.Errorf("wx_trace: get wxinfo err %+v", errs)
		return nil
	}

	log.Log.Debugf("wxUserInfo: %+v", wxinfo)
	log.Log.Debugf("resp : %+v", resp)
	log.Log.Debugf("body : %+v", string(body))

	if wxinfo.Errcode != 0 || resp.StatusCode != 200 {
		log.Log.Errorf("wx_trace: request failed, errCode:%d, statusCode:%d", wxinfo.Errcode, resp.StatusCode)
		return nil
	}

	return &wxinfo
}

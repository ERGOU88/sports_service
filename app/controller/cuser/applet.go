package cuser

import (
	"github.com/garyburd/redigo/redis"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/muser"
	third "sports_service/server/tools/thirdLogin"
	"sports_service/server/util"
	"strings"
	"time"
)

// 微信小程序 登陆/注册
func (svc *UserModule) AppletLoginOrReg(param *muser.AppletLoginParam) (int, string, *models.User) {
	wechat := third.NewWechat()
	//wxToken, err := svc.user.GetAppletAccessToken()
	//if err != nil || wxToken == "" {
	//	accessToken := wechat.GetAppletAccessToken()
	//	if accessToken == nil {
	//		log.Log.Error("applet_trace: get applet access token failed")
	//		return errdef.ERROR, "", nil
	//	}
	//
	//	wxToken = accessToken.AccessToken
	//}
	
	resp, err := wechat.AppletCode2Session(param.Code)
	if err != nil {
		log.Log.Errorf("applet_trace: applet code2 session failed, err:%s", err)
		return errdef.ERROR, "", nil
	}
	
	mobileNum, err := wechat.DecryptPhoneData(param.PhoneData, resp.SessionKey, param.Iv)
	if err != nil {
		log.Log.Errorf("applet_trace: decrypt phone data, err:%s", err)
		return errdef.ERROR, "", nil
	}
	
	//mobileNum, err := wechat.GetUserPhoneNumber(param.CodeByPhone, wxToken)
	//if err != nil {
	//	log.Log.Errorf("applet_trace: get user phone number failed, err:%s", err)
	//	return errdef.ERROR, "", nil
	//}
	
	// 校验手机号合法性
	if b := svc.user.CheckCellPhoneNumber(mobileNum); !b {
		log.Log.Errorf("applet_trace: invalid mobile num %v", mobileNum)
		return errdef.USER_INVALID_MOBILE_NUM, "", nil
	}
	
	// 根据手机号查询用户 不存在 注册用户 用户存在 为登陆
	user := svc.user.FindUserByPhone(mobileNum)
	if user == nil {
		// 开启事务
		if err := svc.engine.Begin(); err != nil {
			log.Log.Errorf("applet_trace: session begin err:%s", err)
			return errdef.ERROR, "", nil
		}
		// 注册
		reg := muser.NewMobileRegister()
		if err := reg.Register(svc.user, int(consts.APPLET_PLATFORM), mobileNum, svc.context.ClientIP()); err != nil {
			log.Log.Errorf("applet_trace: register err:%s", err)
			svc.engine.Rollback()
			return errdef.USER_REGISTER_FAIL, "", nil
		}
		
		// 添加用户
		if err := svc.user.AddUser(); err != nil {
			log.Log.Errorf("applet_trace: add user info err:%s", err)
			svc.engine.Rollback()
			return errdef.USER_ADD_INFO_FAIL, "", nil
		}
		
		// 添加用户初始化通知设置
		if err := svc.notify.AddUserNotifySetting(svc.user.User.UserId, int(time.Now().Unix())); err != nil {
			log.Log.Errorf("applet_trace: add user notify setting err:%s", err)
			svc.engine.Rollback()
			return errdef.USER_ADD_NOTIFY_SET_FAIL, "", nil
		}
		
		svc.user.NewSocialAccount(svc.social, consts.TYPE_APPLET, svc.user.User.UserId, resp.Unionid, resp.Openid)
		// 添加社交帐号（微信小程序）
		if err := svc.social.AddSocialAccountInfo(); err != nil {
			log.Log.Errorf("applet_trace: add wx account err:%s", err)
			svc.engine.Rollback()
			return errdef.WX_ADD_ACCOUNT_FAIL, "", nil
		}
		
		svc.engine.Commit()
		
		// 生成token
		token := svc.user.GenUserToken(svc.user.User.UserId, util.Md5String(svc.user.User.Password))
		// 保存到redis
		if err := svc.user.SaveUserToken(svc.user.User.UserId, token); err != nil {
			log.Log.Errorf("applet_trace: save user token err:%s", err)
		}
		
		return errdef.SUCCESS, token, svc.user.User
	}
	
	has, err := svc.social.HasExistsSocialAccount(consts.TYPE_APPLET, svc.user.User.UserId)
	if !has && err == nil {
		svc.user.NewSocialAccount(svc.social, consts.TYPE_APPLET, svc.user.User.UserId, resp.Unionid, resp.Openid)
		// 添加社交帐号（微信小程序）
		if err := svc.social.AddSocialAccountInfo(); err != nil {
			log.Log.Errorf("applet_trace: add wx account err:%s", err)
			return errdef.WX_ADD_ACCOUNT_FAIL, "", nil
		}
	}
	
	// 登陆的时候 检查用户状态
	if !svc.CheckUserStatus(user.Status) {
		log.Log.Errorf("applet_trace: forbid status, userId:%s", user.Status)
		return errdef.USER_FORBID_STATUS, "", nil
	}
	
	// 用户已注册过, 则直接从redis中获取token并返回
	token, err := svc.user.GetUserToken(svc.user.User.UserId)
	if err != nil && err == redis.ErrNil {
		// redis 没有，重新生成token
		token = svc.user.GenUserToken(svc.user.User.UserId, util.Md5String(svc.user.User.Password))
		// 重新保存到redis
		if err := svc.user.SaveUserToken(svc.user.User.UserId, token); err != nil {
			log.Log.Errorf("applet_trace: save user token err:%s", err)
		}
	}
	
	return errdef.SUCCESS, token, user
}

// 绑定微信
func (svc *UserModule) BindWechat(param *muser.BindWechatParam) int {
	user := svc.user.FindUserByUserid(param.UserId)
	if user == nil {
		return errdef.USER_NOT_EXISTS
	}
	
	wechat := third.NewWechat()
	resp, err := wechat.AppletCode2Session(param.Code)
	if err != nil {
		log.Log.Errorf("applet_trace: applet code2 session failed, err:%s", err)
		return errdef.ERROR
	}
	
	has, err := svc.social.HasExistsSocialAccount(consts.TYPE_APPLET, user.UserId)
	if !has && err == nil {
		svc.user.NewSocialAccount(svc.social, consts.TYPE_APPLET, user.UserId, resp.Unionid, resp.Openid)
		// 添加社交帐号（微信小程序）
		if err := svc.social.AddSocialAccountInfo(); err != nil {
			log.Log.Errorf("applet_trace: add wx account err:%s", err)
			return errdef.WX_ADD_ACCOUNT_FAIL
		}
	}
	
	return errdef.SUCCESS
}

func (svc *UserModule) VerifyToken(param *muser.VerifyTokenParam) int {
	var userid, hashcode string
	ks := strings.Split(param.Token, "_")
	if len(ks) != 2 {
		log.Log.Errorf("len(ks) != 2")
		ks = strings.Split(param.Token, "%09")
	}
	
	if len(ks) == 2 {
		userid = ks[0]
		hashcode = ks[1]
	}
	
	if len(hashcode) <= 0 {
		log.Log.Errorf("len(hashcode) <= 0")
		return errdef.UNAUTHORIZED
	}
	
	
	token, err := svc.user.GetUserToken(userid)
	if err != nil && err == redis.ErrNil {
		log.Log.Errorf("token_trace: get user token by redis err:%s", err)
		return errdef.UNAUTHORIZED
	}
	
	// token是否和redis存储的一致
	if res := strings.Compare(param.Token, token); res != 0 {
		log.Log.Errorf("token_trace: token not match, server token:%s, client token:%s", token, param.Token)
		return errdef.INVALID_TOKEN
	}
	
	log.Log.Debugf("client token:%s, server token:%s", param.Token, token)
	
	if userid != "" {
		// 给token续约
		if err := svc.user.SaveUserToken(userid, param.Token); err != nil {
			log.Log.Errorf("token_trace: save user token err:%s", err)
		}
	}
	
	return errdef.SUCCESS
}

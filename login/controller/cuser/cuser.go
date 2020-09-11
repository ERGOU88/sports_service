package cuser

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"net/url"
	"sports_service/server/dao"
	"sports_service/server/global/login/log"
	"sports_service/server/login/config"
	"sports_service/server/models/muser"
	"sports_service/server/global/login/errdef"
	"sports_service/server/util"
	"sports_service/server/models"
	"github.com/parnurzeal/gorequest"
	"sports_service/server/global/consts"
)

type UserModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	social      *muser.SocialModel
}

func New(c *gin.Context) UserModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return UserModule{
		context: c,
		user: muser.NewUserModel(socket),
		social: muser.NewSocialPlatform(socket),
		engine:  socket,
	}
}

// 手机一键登陆/注册
func (svc *UserModule) MobileLoginOrReg(param *muser.LoginParams) (int, string, *models.User) {
	if b := svc.user.CheckCellPhoneNumber(param.MobileNum); !b {
		log.Log.Errorf("user_trace: invalid mobile num %v", param.MobileNum)
		return errdef.INVALID_MOBILE_NUM, "", nil
	}

	// 根据手机号查询用户 不存在 注册用户 用户存在 为登陆
	if user := svc.user.FindUserByPhone(param.MobileNum); user == nil {
		// 登陆
		reg := muser.NewMobileRegister()
		if err := reg.Register(svc.user, param); err != nil {
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

// 微信登陆/注册
func (svc *UserModule) WechatLoginOrReg() {

}

// WxAccessToken 获取微信accessToken
func (svc *UserModule) WxAccessToken(code string) *muser.AccessToken {
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

func WechatInfo(accessToken *muser.AccessToken) *muser.WechatUserInfo {
	v := url.Values{}
	v.Set("access_token", accessToken.AccessToken)
	v.Set("openid", accessToken.Openid)
	wxinfo := muser.WechatUserInfo{}
	resp, body, errs := gorequest.New().Get(consts.WECHAT_USER_INFO_URL + v.Encode()).EndStruct(&wxinfo)
	if errs != nil {
		log.Log.Errorf("get wxinfo err %+v", errs)
		return nil
	}

	log.Log.Debugf("wxUserInfo: %+v", wxinfo)
	log.Log.Debugf("gorequest resp : %+v", resp)
	log.Log.Debugf("gorequest body : %+v", string(body))

	if wxinfo.Errcode != 0 || resp.StatusCode != 200 {
		return nil
	}

	return &wxinfo
}


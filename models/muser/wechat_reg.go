package muser

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/dao"
	"sports_service/server/global/consts"
	"sports_service/server/global/rdskey"
	"sports_service/server/util"
	"time"
	"sports_service/server/global/login/log"
	"errors"
)

type wechatRegister struct {
	*base
}

// 微信登陆 请求参数
type WxLoginParam struct {
	Code    string     `binding:"required" json:"code" example:"code码"`
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`
}

type WechatUserInfo struct {
	Errcode    int      `json:"errcode"`    //
	Errmsg     string   `json:"errmsg"`     //
	Openid     string   `json:"openid"`     // 普通用户的标识，对当前开发者帐号唯一
	Nickname   string   `json:"nickname"`   // 普通用户昵称
	Sex        int      `json:"sex"`        // 普通用户性别，1为男性，2为女性
	Province   string   `json:"province"`   // 普通用户个人资料填写的省份
	City       string   `json:"city"`       // 普通用户个人资料填写的城市
	Country    string   `json:"country"`    // 国家，如中国为CN
	Headimgurl string   `json:"headimgurl"` // 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空
	Privilege  []string `json:"privilege"`  // 用户特权信息，json数组
	Unionid    string   `json:"unionid"`    // 用户统一标识。针对一个微信开放平台帐号下的应用，同一用户的unionid是唯一的。
}

// 实栗
func NewWechatRegister() *wechatRegister {
	return &wechatRegister{
		&base{},
	}
}

// wechat注册
func (r *wechatRegister) Register(u *UserModel, s *SocialModel, c *gin.Context, wxUnionId, avatar, nickName string, gender int) error {
	key := rdskey.MakeKey(rdskey.LOGIN_REPEAT, consts.TYPE_WECHAT, wxUnionId)
	ok, err:= IsReapeat(key)
	if err != nil {
		log.Log.Errorf("wx_trace: redis err %s",err)
		return err
	}

	if !ok {
		log.Log.Errorf("wx_trace: 微信用户重复注册 unionID:%s", wxUnionId)
		return errors.New("wx_trace: 微信用户重复注册")
	}

	rds:=dao.NewRedisDao()
	rds.EXPIRE64(key, rdskey.KEY_EXPIRE_MIN)
	r.newUser(u, c, avatar, nickName, gender)
	r.newSocialAccount(s, consts.TYPE_WECHAT, u.User.UserId, wxUnionId)
	return nil
}

// 设置用户社交帐号信息
func (r *wechatRegister) newSocialAccount(s *SocialModel, socialType int, userid, unionid string) {
	s.SetCreateAt(time.Now().Unix())
	s.SetSocialType(socialType)
	s.SetUnionId(unionid)
	s.SetUserId(userid)
	return
}

// 设置用户信息
func (r *wechatRegister) newUser(u *UserModel, c *gin.Context, avatar, nickName string, gender int) {
	now := time.Now().Unix()
	r.setDefaultInfo(u, avatar, gender)
	u.SetUserType(consts.TYPE_WECHAT)
	u.SetDeviceType(r.getDeviceType(c))
	u.SetNickName(r.getNickName(nickName))
	u.SetLastLoginTime(now)
	// todo 暂时先使用时间 + 4位随机数 生成uid
	u.SetUid(util.NewUserId())
	u.SetCreateAt(now)
	u.SetUpdateAt(now)
	u.SetPassword("")
	return
}

// 设置默认信息
func (r *wechatRegister) setDefaultInfo(u *UserModel, avatar string, gender int) {
	u.SetAvatar(consts.DEFAULT_AVATAR)
	if avatar != "" {
		u.SetAvatar(avatar)
	}

	if gender != 0 {
		u.SetGender(gender)
	}
}





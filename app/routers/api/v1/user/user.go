package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cuser"
	_ "sports_service/server/app/routers/api/v1/swag"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/muser"
	"sports_service/server/models/sms"
)

// @Summary 获取短信验证码 (ok)
// @Tags 账号体系
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   SendSmsCodeParams  body sms.SendSmsCodeParams true "获取短信验证码请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/user/smscode [post]
// 获取短信验证码
func SmsCode(c *gin.Context) {
	reply := errdef.New(c)
	params := new(sms.SendSmsCodeParams)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("sms_trace: 发送短信 参数错误, params:%+v", params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cuser.New(c)
	syscode := svc.SendSmsCode(params)

	reply.Response(http.StatusOK, syscode)
}

// @Summary 短信验证码注册/登陆 (ok)
// @Tags 账号体系
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   SmsCodeLoginParams  body sms.SmsCodeLoginParams true "短信验证码登陆/注册 请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/user/smscode/login [post]
// 短信验证码登陆
func SmsCodeLogin(c *gin.Context) {
	reply := errdef.New(c)
	params := new(sms.SmsCodeLoginParams)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("sms_trace: 短信验证码登陆 参数错误, params:%+v", params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cuser.New(c)
	// 短信验证码登陆/注册
	syscode, token, user := svc.SmsCodeLogin(params)
	if syscode != errdef.SUCCESS {
		log.Log.Errorf("sms_trace: 用户短信验证码登陆/注册失败，params:%+v", params)
		reply.Response(http.StatusOK, syscode)
		return
	}

	reply.Data["token"] = token
	reply.Data["user_info"] = user
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 手机一键注册/登陆 (ok)
// @Tags 账号体系
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   LoginParams  body muser.LoginParams true "手机号登陆/注册 请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/user/mobile/login [post]
// 手机一键登陆
func MobilePhoneLogin(c *gin.Context) {
	reply := errdef.New(c)
	params := new(muser.LoginParams)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("user_trace: 参数错误, params:%+v", params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cuser.New(c)
	// 手机一键登陆/注册
	syscode, token, user := svc.MobileLoginOrReg(params)
	if syscode != errdef.SUCCESS {
		log.Log.Errorf("user_trace: 用户登陆/注册失败，params:%+v", params)
		reply.Response(http.StatusOK, syscode)
		return
	}

	log.Log.Debugf("freeLogin: user info: %+v", user)

	reply.Data["token"] = token
	reply.Data["user_info"] = user
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 微信注册/登陆 (ok)
// @Tags 账号体系
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   WxLoginParam  body muser.WxLoginParam true "微信登陆/注册 请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/user/wechat/login [post]
// 用户微信登陆
func UserWechatLogin(c *gin.Context) {
	reply := errdef.New(c)
	param := new(muser.WxLoginParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("wx_trace: 参数错误, params:%+v", param)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cuser.New(c)
	// 微信登陆 / 注册
	syscode, token, user := svc.WechatLoginOrReg(param.Code)
	if syscode != errdef.SUCCESS {
		log.Log.Errorf("wx_trace: 用户登陆/注册失败，code:%s", param.Code)
		reply.Response(http.StatusOK, syscode)
		return
	}

	reply.Data["token"] = token
	reply.Data["userInfo"] = user
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 微博注册/登陆 (ok)
// @Tags 账号体系
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   WeiboLoginParams  body muser.WeiboLoginParams true "微博登陆/注册 请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/user/weibo/login [post]
// 用户微博登陆
func UserWeiboLogin(c *gin.Context) {
	reply := errdef.New(c)
	params := new(muser.WeiboLoginParams)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("weibo_trace: 参数错误, params:%+v", params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cuser.New(c)
	// 微博登陆 / 注册
	syscode, token, user := svc.WeiboLoginOrReg(params)
	if syscode != errdef.SUCCESS {
		log.Log.Errorf("weibo_trace: 用户登陆/注册失败, params:%+v", params)
		reply.Response(http.StatusOK, syscode)
		return
	}

	reply.Data["token"] = token
	reply.Data["userInfo"] = user
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary QQ注册/登陆 (ok)
// @Tags 账号体系
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   WeiboLoginParams  body muser.QQLoginParams true "QQ登陆/注册 请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/user/qq/login [post]
// 用户微博登陆
func UserQQLogin(c *gin.Context) {
	reply := errdef.New(c)
	params := new(muser.QQLoginParams)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("qq_trace: 参数错误, params:%+v", params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cuser.New(c)
	// QQ登陆 / 注册
	syscode, token, user := svc.QQLoginOrReg(params)
	if syscode != errdef.SUCCESS {
		log.Log.Errorf("qq_trace: 用户登陆/注册失败, params:%+v", params)
		reply.Response(http.StatusOK, syscode)
		return
	}

	reply.Data["token"] = token
	reply.Data["userInfo"] = user
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 用户信息 (ok)
// @Tags 账号体系
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/user/info [get]
// 用户信息
func UserInfo(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("user_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	svc := cuser.New(c)
	// 获取用户信息
	syscode, userInfo := svc.GetUserInfoByUserid(userId.(string))
	if syscode != errdef.SUCCESS {
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	reply.Data["user_info"] = userInfo
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 编辑用户信息 (ok)
// @Tags 账号体系
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   EditUserInfoParams  body muser.EditUserInfoParams true "编辑用户信息请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/user/edit/info [post]
// 修改用户信息
func EditUserInfo(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("user_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	params := new(muser.EditUserInfoParams)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("user_trace: param err, params:%+v", params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	if params.CountryId == 0 && params.Born == "" && params.Signature == "" && params.Gender == 0 && params.NickName == "" && params.Avatar == "" {
		log.Log.Errorf("user_trace: invalid param, params:%+v", params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cuser.New(c)
	// 修改用户信息
	syscode := svc.EditUserInfo(userId.(string), params)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 用户反馈问题 (ok)
// @Tags 账号体系
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   FeedbackParam  body muser.FeedbackParam true "用户反馈问题 请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/user/feedback [post]
// 用户反馈
func UserFeedback(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("user_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	params := new(muser.FeedbackParam)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("user_trace: feedback param err:%s, params:%+v", err, params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cuser.New(c)
	// 记录用户反馈的问题
	syscode := svc.RecordUserFeedback(userId.(string), params)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 个人空间用户信息 (ok)
// @Tags 账号体系
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	    header    string 	true  "版本" default(1.0.0)
// @Param   user_id 	    query     string 	true  "用户id"
// @Param   to_user_id 	  query     string 	true  "被查看人的用户id"
// @Param   UserZoneInfoParam  body muser.UserZoneInfoParam true "个人空间用户信息 请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/user/zone/info [get]
// 个人空间用户信息
func UserZoneInfo(c *gin.Context) {
	reply := errdef.New(c)
	userId := c.Query("user_id")

	toUserId := c.Query("to_user_id")
	if toUserId == "" {
		log.Log.Errorf("user_trace: request toUserId is empty, toUserId:%s", userId)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cuser.New(c)
	// 获取用户个人空间信息
	syscode, userInfo, zoneInfo := svc.GetUserZoneInfo(userId, toUserId)
	reply.Data["user_info"] = userInfo
	reply.Data["zone_info"] = zoneInfo
	reply.Data["code"] = syscode
	reply.Response(http.StatusOK, syscode)
}

// 绑定设备token
func BindDeviceToken(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	if userId == "" {
		log.Log.Errorf("user_trace: user not found, uid:%s", userId)
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	param := new(muser.BindDeviceTokenParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("user_trace: invalid param, param:%+v", param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cuser.New(c)
	syscode := svc.BindDeviceToken(userId.(string), param)
	reply.Response(http.StatusOK, syscode)
}


// 版本更新(load数据库)
func VersionUp(c *gin.Context) {
	reply := errdef.New(c)
	var versions []string
	versions = c.Request.Header["Version"]
	var version string
	if len(versions) > 0 {
		version = versions[0]
	}

	log.Log.Infof("configure_trace: cur client version:%s", version)

	// debug模式不校验强更版本
	//if conf.Global.Debug {
	//	return
	//}

	svc := cuser.New(c)
	syscode, upgrade := svc.VersionUp(version)
	if syscode != errdef.SUCCESS {
		reply.Response(http.StatusOK, syscode)
		return
	}

	reply.Data["upgrade"] = upgrade
	reply.Response(http.StatusOK, errdef.SUCCESS)
	return
}

// 用户卡包
func UserKabaw(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	svc := cuser.New(c)
	code, kabaw := svc.GetKabawInfo(userId.(string))
	if code == errdef.SUCCESS {
		reply.Data["detail"] = kabaw
	}

	reply.Response(http.StatusOK, code)
}

// 更新腾讯im签名
func UpdateTencentImSign(c *gin.Context) {
	reply := errdef.New(c)
	param := &muser.UpdateTencentImSign{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("user_trace: invalid param, param:%+v", param)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cuser.New(c)
	code, sign := svc.UpdateTencentImSign(param.UserId)
	reply.Data["sign"] = sign
	reply.Response(http.StatusOK, code)
}

// 添加游客[腾讯im]
func TencentImAddGuest(c *gin.Context) {
	reply := errdef.New(c)
	svc := cuser.New(c)
	code, info := svc.AddGuestByTencentIm()
	reply.Data["detail"] = info
	reply.Response(http.StatusOK, code)
}

// 腾讯im 添加用户
func TencentImAddUser(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	svc := cuser.New(c)
	code, sign := svc.GetTencentImSignByUser(userId.(string))
	reply.Data["sign"] = sign
	reply.Response(http.StatusOK, code)
}

// 获取腾讯im签名
func GetTencentImSign(c *gin.Context) {
	reply := errdef.New(c)
	userId := c.Query("user_id")
	svc := cuser.New(c)
	code, info := svc.GetTencentImSign(userId)
	reply.Data["detail"] = info
	reply.Response(http.StatusOK, code)
}

func VerifyWxCode(c *gin.Context) {

}

func AppletLogin(c *gin.Context) {
	reply := errdef.New(c)
	param := &muser.AppletLoginParam{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}
	
	svc := cuser.New(c)
	code, token, user := svc.AppletLoginOrReg(param)
	reply.Data["token"] = token
	reply.Data["user"] = user
	reply.Response(http.StatusOK, code)
}

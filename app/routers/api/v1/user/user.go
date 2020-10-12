package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cuser"
	_ "sports_service/server/app/routers/api/v1/swag"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/sms"
	"sports_service/server/models/muser"
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
// @Success 200 {object} swag.LoginSwag
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
// @Success 200 {object} swag.LoginSwag
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
// @Success 200 {object} swag.LoginSwag
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
// @Success 200 {object} swag.LoginSwag
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
// @Success 200 {object} swag.LoginSwag
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
// @Success 200 {object}  muser.UserInfoResp
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

	reply.Data["userInfo"] = userInfo
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

	if params.CountryId == 0 && params.Born == "" && params.Signature == "" && params.Gender == 0 && params.NickName == "" && params.AvatarId == 0 {
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
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   UserZoneInfoParam  body muser.UserZoneInfoParam true "个人空间用户信息 请求参数"
// @Success 200 {object}  swag.ZoneInfoSwag
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/user/zone/info [get]
// 个人空间用户信息
func UserZoneInfo(c *gin.Context) {
	reply := errdef.New(c)
	userId := c.Query("user_id")
	if userId == "" {
		log.Log.Errorf("user_trace: request userId is empty, userId:%s", userId)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cuser.New(c)
	// 获取用户个人空间信息
	syscode, userInfo, zoneInfo := svc.GetUserZoneInfo(userId)
	reply.Data["user_info"] = userInfo
	reply.Data["zone_info"] = zoneInfo
	reply.Response(http.StatusOK, syscode)
}



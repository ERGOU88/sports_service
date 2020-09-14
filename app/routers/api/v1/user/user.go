package user

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/global/app/errdef"
	"sports_service/server/app/controller/cuser"
	"sports_service/server/models/muser"
	"net/http"
	"sports_service/server/global/app/log"
	_ "sports_service/server/app/routers/api/v1/swag"
)

// @Summary 手机一键注册/登陆 (ok)
// @Tags 账号体系
// @Version 1.0
// @Description
// @tags 001 手机一键注册/登陆 2020-09-10
// @Accept json
// @Produce  json
// @Param   User-Agent    header    string 	true  "android" default(android)
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

	reply.Data["token"] = token
	reply.Data["userInfo"] = user
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 微信注册/登陆 (ok)
// @Tags 账号体系
// @Version 1.0
// @Description
// @tags 002 微信注册/登陆 2020-09-11
// @Accept json
// @Produce  json
// @Param   User-Agent    header    string 	true  "android" default(android)
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
// @tags 003 微博注册/登陆 2020-09-11
// @Accept json
// @Produce  json
// @Param   User-Agent    header    string 	true  "android" default(android)
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
// @tags 004 QQ注册/登陆 2020-09-14
// @Accept json
// @Produce  json
// @Param   User-Agent    header    string 	true  "android" default(android)
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



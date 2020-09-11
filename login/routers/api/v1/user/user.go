package user

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/global/login/errdef"
	"sports_service/server/login/controller/cuser"
	"sports_service/server/models/muser"
	"net/http"
	"sports_service/server/global/login/log"
	_ "sports_service/server/login/routers/api/v1/swag"
)

// @Summary 手机一键注册/登陆 (ok)
// @Version 1.0
// @Description
// @tags 001 手机一键注册/登陆 2020-09-10
// @Accept json
// @Produce  json
// @Param   User-Agent    header    string 	true  "android" default(android)
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   LoginParams  body swag.LoginParamsSwag true "登陆/注册 请求参数"
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
		log.Log.Errorf("user_trace: 用户登陆/注册失败，mobileNum:%s", params.MobileNum)
		reply.Response(http.StatusOK, syscode)
		return
	}

	reply.Data["token"] = token
	reply.Data["userInfo"] = user
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// 用户微信登陆
func UserWechatLogin(c *gin.Context) {
	//reply := errdef.New(c)

}



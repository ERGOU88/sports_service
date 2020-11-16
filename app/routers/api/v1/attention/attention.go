package attention

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cattention"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mattention"
	"sports_service/server/util"
	_ "sports_service/server/models/muser"
)

// @Summary 关注用户 (ok)
// @Tags 关注模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   AddAttentionParam  body mattention.AddAttentionParam true "关注用户请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/attention/user [post]
// 关注用户
func AttentionUser(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("attention_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	param := new(mattention.AddAttentionParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("attention_trace: invalid param, param:%+v", param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cattention.New(c)
	// 添加关注
	syscode := svc.AddAttention(userId.(string), param.UserId)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 取消关注 (ok)
// @Tags 关注模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   CancelAttentionParam  body mattention.CancelAttentionParam true "取消关注请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/attention/cancel [post]
// 取消关注
func CancelAttention(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("attention_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	param := new(mattention.CancelAttentionParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("attention_trace: invalid param, param:%+v", param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cattention.New(c)
	// 取消关注
	syscode := svc.CancelAttention(userId.(string), param.UserId)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 关注的用户列表[分页获取] (ok)
// @Tags 关注模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   page	  	  query  	string 	true  "页码 从1开始"
// @Param   size	  	  query  	string 	true  "每页展示多少 最多50条"
// @Param   user_id		  query  	string 	true  "用户id"
// @Success 200 {array}  muser.UserInfoResp
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/attention/list [get]
// 关注的用户列表
func AttentionList(c *gin.Context) {
	reply := errdef.New(c)
	//userId, ok := c.Get(consts.USER_ID)
	//if !ok {
	//	log.Log.Errorf("attention_trace: user not found, uid:%s", userId.(string))
	//	reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
	//	return
	//}
	userId := c.Query("user_id")

	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := cattention.New(c)
	// 获取关注的用户列表
	list := svc.GetAttentionUserList(userId, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 用户的粉丝列表[分页获取] (ok)
// @Tags 关注模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   page	  	  query  	string 	true  "页码 从1开始"
// @Param   size	  	  query  	string 	true  "每页展示多少 最多50条"
// @Param   user_id		  query  	string 	true  "用户id"
// @Success 200 {array}  muser.UserInfoResp
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/fans/list [get]
// 用户的粉丝列表
func FansList(c *gin.Context) {
	reply := errdef.New(c)
	//userId, ok := c.Get(consts.USER_ID)
	//if !ok {
	//	log.Log.Errorf("attention_trace: user not found, uid:%s", userId.(string))
	//	reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
	//	return
	//}
	//
	userId := c.Query("user_id")

	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := cattention.New(c)
	// 获取用户的粉丝列表
	list := svc.GetFansList(userId, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
	return
}

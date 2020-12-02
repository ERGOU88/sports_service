package notify

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cnotify"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mnotify"
	"sports_service/server/util"
	_"sports_service/server/models"
)

// @Summary 系统通知设置 (ok)
// @Tags 通知模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   NotifySettingParams  body mnotify.NotifySettingParams true "系统通知设置 请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/notify/setting [post]
// 通知设置
func NotifySetting(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("user_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	params := new(mnotify.NotifySettingParams)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("notify_trace: notify setting param err:%s, param:%+v", err, params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	log.Log.Errorf("params:%+v", params)

	svc := cnotify.New(c)
	// 保存用户通知设置
	syscode := svc.SaveUserNotifySetting(userId.(string), params)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 被点赞通知[分页获取] (ok)
// @Tags 通知模块
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
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/notify/beliked [get]
// 被点赞通知
func BeLikedNotify(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	svc := cnotify.New(c)
	list := svc.GetBeLikedList(userId.(string), page, size)
	reply.Data["list"] = list
	//reply.Data["read_index"] = readIndex
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 被@通知[分页获取] (ok)
// @Tags 通知模块
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
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/notify/receive/at [get]
// 被@通知
func ReceiveAtNotify(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	svc := cnotify.New(c)
	list, readIndex := svc.GetReceiveAtNotify(userId.(string), page, size)
	reply.Data["list"] = list
  reply.Data["read_index"] = readIndex
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 获取未读通知数量 (ok)
// @Tags 通知模块
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
// @Router /api/v1/notify/unread/quantity [get]
// 未读数量
func UnreadNum(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	svc := cnotify.New(c)
	info := svc.GetUnreadNum(userId.(string))
	reply.Data["info"] = info
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 通知设置信息 (ok)
// @Tags 通知模块
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
// @Router /api/v1/notify/setting/info [get]
// 通知设置信息
func NotifySettingInfo(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)

	svc := cnotify.New(c)
	info := svc.GetUserNotifySetting(userId.(string))
	reply.Data["notify_setting"] = info
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 系统通知消息列表 (ok)
// @Tags 通知模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	    header    string 	true  "版本" default(1.0.0)
// @Param   page	  	    query  	string 	true  "页码 从1开始"
// @Param   size	  	    query  	string 	true  "每页展示多少 最多50条"
// @Param   user_id       query   string  false "用户id 非必传"
// @Success 200 {string}  json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string}  json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/notify/system [get]
// 系统通知
func SystemNotify(c *gin.Context) {
  reply := errdef.New(c)
  userId := c.Query("user_id")
  page, size := util.PageInfo(c.Query("page"), c.Query("size"))
  svc := cnotify.New(c)
  list := svc.GetSystemNotify(userId, page, size)
  reply.Data["list"] = list
  reply.Response(http.StatusOK, errdef.SUCCESS)
}

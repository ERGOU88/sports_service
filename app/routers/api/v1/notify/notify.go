package notify

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cnotify"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mnotify"
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

	svc := cnotify.New(c)
	// 保存用户通知设置
	syscode := svc.SaveUserNotifySetting(userId.(string), params)
	reply.Response(http.StatusOK, syscode)
}

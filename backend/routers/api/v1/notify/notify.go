package notify

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/backend/controller/cnotify"
	"sports_service/global/backend/errdef"
	"sports_service/global/backend/log"
	"sports_service/models/umeng"
	"sports_service/util"
)

// 后台系统推送
func PushSystemNotify(c *gin.Context) {
	reply := errdef.New(c)

	param := new(umeng.SystemNotifyParams)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("notify_trace: system notify err:%s, param:%+v", err, param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cnotify.New(c)
	syscode := svc.PushSystemNotify(param)
	reply.Response(http.StatusOK, syscode)
}

// 后台系统推送列表
func SystemNotifyList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	// 发送状态 -1 全部 0 已发送 1 未发送
	sendStatus := c.Query("send_status")
	// 通知类型 0 指定玩家 1 全部玩家
	sendDefault := c.Query("send_default")
	svc := cnotify.New(c)
	list := svc.GetSystemNotifyList(page, size, sendStatus, sendDefault)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// 撤回系统定时推送
func CancelSystemNotify(c *gin.Context) {
	reply := errdef.New(c)
	param := &umeng.CancelSystemNotifyParam{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("notify_trace: cancel system notify param err:%s, param:%+v", err, param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cnotify.New(c)
	// 撤回系统推送
	syscode := svc.CancelSystemNotify(param.SystemId)
	reply.Response(http.StatusOK, syscode)
}

// 删除系统通知
func DelSystemNotify(c *gin.Context) {
	reply := errdef.New(c)
	param := &umeng.DelSystemNotifyParam{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("notify_trace: del system notify param err:%s, param:%+v", err, param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cnotify.New(c)
	// 删除系统推送
	syscode := svc.DelSystemNotify(param.SystemId)
	reply.Response(http.StatusOK, syscode)
}

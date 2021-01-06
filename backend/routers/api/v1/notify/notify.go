package notify

import (
  "github.com/gin-gonic/gin"
  "net/http"
  "sports_service/server/backend/controller/cnotify"
  "sports_service/server/util"
  "sports_service/server/global/backend/errdef"
  "sports_service/server/global/backend/log"
  "sports_service/server/models/umeng"
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

  svc := cnotify.New(c)
  list := svc.GetSystemNotifyList(page, size)
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
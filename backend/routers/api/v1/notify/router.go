package notify

import (
  "github.com/gin-gonic/gin"
  "sports_service/server/middleware/jwt"

  //"sports_service/server/middleware/jwt"
)

// 站内信模块
func Router(engine *gin.Engine) {
  backend := engine.Group("/backend/v1")
  notify := backend.Group("/notify")
  notify.Use(jwt.JwtAuth())
  {
    // 后台系统推送（全部用户 or 指定用户）
    notify.POST("/system", PushSystemNotify)
    // 后台系统推送列表
    notify.GET("/list", SystemNotifyList)
    // 撤回定时推送 todo: 添加发送状态 已发送 未发送 已撤回
    notify.POST("/cancel", CancelSystemNotify)
    // 删除系统消息
    notify.POST("/del", DelSystemNotify)
  }
}

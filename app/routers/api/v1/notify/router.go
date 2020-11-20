package notify

import (
  "github.com/gin-gonic/gin"
  "sports_service/server/middleware/token"
)

// 通知模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	notify := api.Group("/notify", token.TokenAuth())
	{
		// 通知设置
		notify.POST("/setting", NotifySetting)
		// 被赞通知
		notify.GET("/beliked", BeLikedNotify)
		// 被@通知
		notify.GET("/receive/at", ReceiveAtNotify)
		// 未读数量
		notify.GET("/unread/quantity", UnreadNum)
		// 通知设置信息
		notify.GET("/setting/info", NotifySettingInfo)
	}
}

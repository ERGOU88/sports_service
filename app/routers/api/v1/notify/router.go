package notify

import (
  "github.com/gin-gonic/gin"
  "sports_service/server/middleware/token"
)

// 通知模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	notify := api.Group("/notify")
	{
		// 通知设置
		notify.POST("/setting", token.TokenAuth(), NotifySetting)
		// 被赞通知
		notify.GET("/beliked", token.TokenAuth(), BeLikedNotify)
		// 被@通知
		notify.GET("/receive/at", token.TokenAuth(), ReceiveAtNotify)
		// 未读数量
		notify.GET("/unread/quantity", token.TokenAuth(), UnreadNum)
		// 通知设置信息
		notify.GET("/setting/info", token.TokenAuth(), NotifySettingInfo)
		// 系统通知列表
		notify.GET("/system", SystemNotify)
    // 系统消息详情
    notify.GET("/system/message/detail", SystemMessageDetail)
		// 首页通知信息
		notify.GET("/homepage/info", HomePageNotify)
	}
}

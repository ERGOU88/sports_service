package notify

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
	"sports_service/server/middleware/token"
)

// 通知模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	notify := api.Group("/notify", sign.CheckSign(), token.TokenAuth())
	{
		// 通知设置
		notify.GET("/setting", NotifySetting)
	}
}

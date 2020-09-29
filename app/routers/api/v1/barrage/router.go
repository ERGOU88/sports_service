package barrage

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/token"
)

// 弹幕模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	barrage := api.Group("/barrage", token.TokenAuth())
	{
		// 发送弹幕
		barrage.POST("/send", SendBarrage)
		// 获取视频弹幕内容（按时长区间获取）
		barrage.GET("/video/list", VideoBarrage)
	}
}


package video

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
	"sports_service/server/middleware/token"
)

// 视频点播模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	video := api.Group("/video", sign.CheckSign(), token.TokenAuth())
	{
		// 用户发布的视频列表
		video.POST("/publish", VideoPublish)
		// 用户视频浏览记录
		video.GET("/browse/history", BrowseHistory)

	}
}

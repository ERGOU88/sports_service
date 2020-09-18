package video

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/token"
)

// 视频点播模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	video := api.Group("/video", token.TokenAuth())
	{
		// 用户发布视频
		video.POST("/publish", VideoPublish)
		// 用户视频浏览记录
		video.GET("/browse/history", BrowseHistory)
		// 用户发布的视频列表
		video.POST("/publish/list", VideoPublishList)
		// 删除浏览记录
		video.POST("/delete/history", DeleteHistory)
		// 删除发布的记录
		video.POST("/delete/publish", DeletePublish)

	}
}

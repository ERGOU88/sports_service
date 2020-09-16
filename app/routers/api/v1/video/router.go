package video

import (
	"github.com/gin-gonic/gin"
)

func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	video := api.Group("/video")
	{
		// 用户发布的视频列表
		video.GET("/user/publish")

	}
}

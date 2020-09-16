package collect

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
)

// 路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	collect := api.Group("/collect")
	collect.Use(sign.CheckSign())
	{
		// 收藏视频
		collect.POST("/video", CollectVideo)
		// 取消收藏
		collect.POST("/video/cancel", CancelCollect)
		// 用户收藏的视频列表
		collect.GET("/video/list", CollectVideoList)
	}
}

package collect

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
	"sports_service/server/middleware/token"
)

// 收藏模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	collect := api.Group("/collect", sign.CheckSign(), token.TokenAuth())
	{
		// 收藏视频
		collect.POST("/video", CollectVideo)
		// 取消收藏
		collect.POST("/video/cancel", CancelCollect)
		// 用户收藏的视频列表
		collect.GET("/video/list", CollectVideoList)
		// 删除收藏
		collect.POST("/delete", DeleteCollect)
	}
}

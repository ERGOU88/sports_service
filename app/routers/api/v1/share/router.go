package share

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
	"sports_service/server/middleware/token"
)

// 分享/转发模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	share := api.Group("/share")
	share.Use(sign.CheckSign())
	{
		// 分享/转发数据 [帖子、视频]
		share.POST("/data", token.TokenAuth(), ShareData)

	}
}

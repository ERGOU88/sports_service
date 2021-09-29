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
		// 分享/转发到社交平台
		share.POST("/social", ShareWithSocialPlatform)
        // 分享/转发到app社区
		share.POST("/community", token.TokenAuth(), ShareWithCommunity)
		// 获取分享链接
		share.GET("/url", GetShareUrl)
	}
}

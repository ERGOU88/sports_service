package posting

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
)

// 贴子模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	posting := api.Group("/post")
	posting.Use(sign.CheckSign())
	{
		// 发布贴子
		posting.POST("/publish", PublishPosting)
		// 帖子详情
		posting.GET("/detail", PostDetail)
	}
}

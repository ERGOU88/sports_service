package posting

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
	"sports_service/server/middleware/token"
)

// 贴子模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	posting := api.Group("/post")
	posting.Use(sign.CheckSign())
	{
		// 发布贴子
		posting.POST("/publish", token.TokenAuth(), PublishPosting)
		// 帖子详情
		posting.GET("/detail", PostDetail)
		// 用户发布的帖子列表
		posting.GET("/publish/list", token.TokenAuth(), PostPublishList)
		// 用户删除发布的帖子
		posting.POST("/delete/publish", token.TokenAuth(), DeletePublishPost)
	}
}

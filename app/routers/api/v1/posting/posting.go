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
		posting.GET("/publish/list", PostPublishList)
		// 查看其他用户发布的帖子
		posting.GET("/other/publish/list", OtherPublishPost)
		// 用户删除发布的帖子
		posting.POST("/delete/publish", token.TokenAuth(), DeletePublishPost)
        // 帖子申请精华
		posting.POST("/apply/cream", token.TokenAuth(), ApplyPostCream)
		// 举报帖子
		posting.POST("/report", PostReport)
	}
}

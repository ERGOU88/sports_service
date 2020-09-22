package comment

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/token"
)

// 评论模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	comment := api.Group("/comment", token.TokenAuth())
	{
		// 发布评论
		comment.POST("/publish", PublishComment)
		// 回复评论
		comment.POST("/reply", PublishReply)
		// 评论列表
		comment.GET("/list", CommentList)
		// 回复列表
		comment.GET("/reply/list", ReplyList)
	}
}

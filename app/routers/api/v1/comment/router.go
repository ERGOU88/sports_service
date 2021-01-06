package comment

import (
  "github.com/gin-gonic/gin"
  "sports_service/server/middleware/sign"
  "sports_service/server/middleware/token"
)

// 评论模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	comment := api.Group("/comment")
  comment.Use(sign.CheckSign())
	{
		// 发布评论
		comment.POST("/publish", token.TokenAuth(), PublishComment)
		// 回复评论
		comment.POST("/reply", token.TokenAuth(), PublishReply)
		// 评论列表
		comment.GET("/list", CommentList)
		// 回复列表
		comment.GET("/reply/list", ReplyList)
		// 举报评论
		comment.POST("/report", CommentReport)
	}
}

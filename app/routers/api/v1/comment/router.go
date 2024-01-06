package comment

import (
	"github.com/gin-gonic/gin"
	"sports_service/middleware/sign"
	"sports_service/middleware/token"
)

// 评论模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	comment := api.Group("/comment")
	{
		// 发布评论
		comment.POST("/publish", sign.CheckSign(), token.TokenAuth(), PublishComment)
		// 发布评论v2
		comment.POST("/publish/v2", sign.CheckSign(), token.TokenAuth(), V2PublishComment)
		// 回复评论
		comment.POST("/reply", sign.CheckSign(), token.TokenAuth(), PublishReply)
		// 评论列表
		comment.GET("/list", sign.CheckSign(), CommentList)
		// 回复列表
		comment.GET("/reply/list", sign.CheckSign(), ReplyList)
		// 举报评论
		comment.POST("/report", CommentReport)
		// 删除评论
		comment.POST("/del", sign.CheckSign(), token.TokenAuth(), DelComment)
	}
}

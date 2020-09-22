package like

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/token"
)

// 点赞模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	like := api.Group("/like", token.TokenAuth())
	{
		// 视频点赞
		like.POST("/video", GiveLikeForVideo)
		// 视频取消点赞
		like.POST("/video/cancel", CancelLikeForVideo)
		// 用户点赞的视频列表
		like.GET("/video/list", LikeVideoList)
		// 评论点赞
		like.POST("/comment", GiveLikeForComment)
		// 评论取消点赞
		like.POST("/comment/cancel", CancelLikeForComment)
	}
}

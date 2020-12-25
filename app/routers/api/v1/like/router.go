package like

import (
  "github.com/gin-gonic/gin"
  "sports_service/server/middleware/token"
)

// 点赞模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	like := api.Group("/like")
	{
		// 视频点赞
		like.POST("/video", token.TokenAuth(), GiveLikeForVideo)
		// 视频取消点赞
		like.POST("/video/cancel", token.TokenAuth(), CancelLikeForVideo)
		// 用户点赞的视频列表
		like.GET("/video/list", LikeVideoList)
		// 查看其他用户点赞的视频列表
		like.GET("/other/list", OtherUserLikeVideoList)
		// 评论点赞
		like.POST("/comment", token.TokenAuth(), GiveLikeForComment)
		// 评论取消点赞
		like.POST("/comment/cancel", token.TokenAuth(), CancelLikeForComment)
	}
}

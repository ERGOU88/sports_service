package like

import (
	"github.com/gin-gonic/gin"
	"sports_service/middleware/sign"
	"sports_service/middleware/token"
)

// 点赞模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	like := api.Group("/like")
	like.Use(sign.CheckSign())
	{
		// 视频点赞
		like.POST("/video", token.TokenAuth(), GiveLikeForVideo)
		// 视频取消点赞
		like.POST("/video/cancel", token.TokenAuth(), CancelLikeForVideo)
		// 用户点赞的视频列表
		like.GET("/video/list", LikeVideoList)
		// 查看其他用户点赞的视频列表
		like.GET("/other/list", OtherUserLikeVideoList)
		// 视频/帖子/资讯 评论点赞
		like.POST("/comment", token.TokenAuth(), GiveLikeForVideoComment)
		// 视频/帖子/资讯 评论取消点赞
		like.POST("/comment/cancel", token.TokenAuth(), CancelLikeForVideoComment)
		// 帖子点赞
		like.POST("/post", token.TokenAuth(), GiveLikeForPost)
		// 帖子取消点赞
		like.POST("/post/cancel", token.TokenAuth(), CancelLikeForPost)
		// 资讯点赞
		like.POST("/information", token.TokenAuth(), GiveLikeForInformation)
		// 资讯取消点赞
		like.POST("/information/cancel", token.TokenAuth(), CancelLikeForInformation)
	}
}

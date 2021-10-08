package comment

import "github.com/gin-gonic/gin"

// 评论模块后台路由
func Router(engine *gin.Engine) {
	api := engine.Group("/backend/v1")
	comment := api.Group("/comment")
	{
		// 获取视频评论列表(后台)
		comment.GET("/list", VideoCommentList)
		// 删除评论（物理删除）todo: 改为逻辑删除 且 只删除单条
		comment.POST("/delete", DelVideoComments)
		// 视频弹幕列表
		comment.GET("/barrage", VideoBarrageList)
		// 删除视频弹幕
		comment.POST("/barrage/delete", DelVideoBarrage)
		// 帖子评论列表
		comment.GET("/post/list", PostCommentList)
		// 资讯评论列表
		comment.GET("/information/list", InformationCommentList)
	}
}

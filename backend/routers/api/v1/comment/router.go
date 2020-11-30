package comment

import "github.com/gin-gonic/gin"

// 评论模块后台路由
func Router(engine *gin.Engine) {
	api := engine.Group("/backend/v1")
	comment := api.Group("/comment")
	{
		// 获取用户的评论列表(后台)
		comment.GET("/list", VideoCommentList)
		// 删除评论（物理删除）
		comment.POST("/delete", DelVideoComments)
		// 视频弹幕列表
		comment.GET("/barrage/list", VideoBarrageList)
		// 删除视频弹幕
		comment.POST("/barrage/delete", DelVideoBarrage)
	}
}

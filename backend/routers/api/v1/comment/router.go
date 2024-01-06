package comment

import (
	"github.com/gin-gonic/gin"
	"sports_service/middleware/jwt"
)

// 评论模块后台路由
func Router(engine *gin.Engine) {
	api := engine.Group("/backend/v1")
	comment := api.Group("/comment")
	comment.Use(jwt.JwtAuth())
	{
		// 获取视频评论列表(后台)
		comment.GET("/list", VideoCommentList)
		// 删除评论（软删）todo: 改为逻辑删除 且 只删除单条
		comment.POST("/delete", DelComments)
		// 弹幕列表 视频/直播/回放
		comment.GET("/barrage", BarrageList)
		// 删除视频/直播/回放 弹幕
		comment.POST("/barrage/delete", DelVideoBarrage)
		// 帖子评论列表
		comment.GET("/post/list", PostCommentList)
		// 资讯评论列表
		comment.GET("/information/list", InformationCommentList)
	}
}

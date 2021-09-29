package post

import "github.com/gin-gonic/gin"

// 帖子模块
func Router(engine *gin.Engine) {
	backend := engine.Group("/backend/v1")
	post := backend.Group("/post")
	{
		// 后台审核帖子
		post.POST("/audit", AuditPost)
		// 后台帖子列表
		post.GET("/list", PostList)
		// 添加板块
		post.POST("/add/section", AddSection)
	}
}

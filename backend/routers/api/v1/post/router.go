package post

import (
	"github.com/gin-gonic/gin"
)

// 帖子模块
func Router(engine *gin.Engine) {
	backend := engine.Group("/backend/v1")
	post := backend.Group("/post")
	//post.Use(jwt.JwtAuth())
	{
		// 后台审核帖子
		post.POST("/audit", AuditPost)
		// 后台帖子列表
		post.GET("/list", PostList)
		// 帖子板块列表
		post.GET("/section/list", SectionList)
		// 添加板块
		post.POST("/add/section", AddSection)
		// 删除板块
		post.POST("/del/section", DelSection)
		// 帖子话题列表
		post.GET("/topic/list", TopicList)
		// 添加话题
		post.POST("/add/topic", AddTopic)
		// 删除话题
		post.POST("/del/topic", DelTopic)
		// 帖子设置 置顶/精华
		post.POST("/setting", PostSetting)
		// 申精列表
		post.GET("/apply/cream", ApplyCreamList)
		// 批量修改帖子信息
		post.POST("/batch/edit", BatchEditPostInfo)
	}
}

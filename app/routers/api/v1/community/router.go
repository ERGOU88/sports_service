package community

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
)

// 社区模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	community := api.Group("/community")
	community.Use(sign.CheckSign())
	{
		// 社区板块
		community.GET("/section/list", CommunitySections)
		// 社区话题
		community.GET("/topic/list", CommunityTopics)
		// 通过id获取社区话题信息
		community.GET("/topic", CommunityTopicById)
		// 社区板块下的帖子列表
		community.GET("/section/post", SectionPostList)
		// 社区话题下的帖子列表
		community.GET("/topic/post", TopicPostList)
		// 关注的人发布的帖子列表
		community.GET("/post/attention", PostListByAttention)
	}
}

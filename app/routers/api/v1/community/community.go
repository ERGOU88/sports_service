package community

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/community"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/util"
)

// /api/v1/community/section/list
// 社区板块
func CommunitySections(c *gin.Context) {
	reply := errdef.New(c)
	svc := community.New(c)
	code, list := svc.GetCommunitySections()
	if code == errdef.SUCCESS {
		reply.Data["list"] = list
	}

	reply.Response(http.StatusOK, code)
}

// /api/v1/community/topic/list
// 社区话题
func CommunityTopics(c *gin.Context) {
	reply := errdef.New(c)
	isHot := c.Query("is_hot")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := community.New(c)
	code, list := svc.GetCommunityTopics(isHot, page, size)
	if code == errdef.SUCCESS {
		reply.Data["list"] = list
	}

	reply.Response(http.StatusOK, code)
}

// /api/v1/community/topic
// 通过id获取社区话题
func CommunityTopicById(c *gin.Context) {
	reply := errdef.New(c)

	id := c.Query("id")
	svc := community.New(c)
	code, info := svc.GetCommunityTopicById(id)
	if code == errdef.SUCCESS {
		reply.Data["info"] = info
	}

	reply.Response(http.StatusOK, code)
}

// /api/v1/community/section/post
func SectionPostList(c *gin.Context) {
	reply := errdef.New(c)
	// 板块id 默认综合
	sectionId := c.DefaultQuery("section_id", "1")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	userId := c.Query("user_id")

	svc := community.New(c)
	code, list := svc.GetPostListBySection(page, size, userId, sectionId)
	reply.Data["list"] = list


	_, topicList := svc.GetCommunityTopics(consts.POST_SORT_HOT, 1, 4)
	reply.Data["topic_list"] = topicList
	// 默认取两个置顶帖
	_, topList := svc.GetTopPostBySectionId(1, 2, sectionId)
	reply.Data["top_list"] = topList

	// 该板块下的帖子数
	reply.Data["post_num"] = svc.GetPostNumBySection(sectionId)

	reply.Response(http.StatusOK, code)
}

// /api/v1/community/topic/post
func TopicPostList(c *gin.Context) {
	reply := errdef.New(c)
	topicId := c.Query("topic_id")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	userId := c.Query("user_id")
	sortHot := c.Query("sort_hot")

	svc := community.New(c)
	code, list := svc.GetPostListByTopic(page, size, userId, topicId, sortHot)
	if code == errdef.SUCCESS {
		reply.Data["list"] = list
	}

	reply.Response(http.StatusOK, code)
}

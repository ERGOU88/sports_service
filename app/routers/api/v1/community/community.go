package community

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/community"
	"sports_service/server/global/app/errdef"
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

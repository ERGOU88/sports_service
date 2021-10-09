package post

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/backend/controller/cpost"
	"sports_service/server/global/app/log"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/models/mcommunity"
	"sports_service/server/models/mposting"
	"sports_service/server/util"
)

// 帖子审核
func AuditPost(c *gin.Context) {
	reply := errdef.New(c)
	param := &mposting.AudiPostParam{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("post_trace: invalid param, err:%s", err)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cpost.New(c)
	syscode := svc.AudiPost(param)
	reply.Response(http.StatusOK, syscode)
}

// 帖子列表 todo：展示数据待确认
func PostList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := cpost.New(c)
	code, list := svc.GetPostList(page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

func AddSection(c *gin.Context) {
	reply := errdef.New(c)
	param := &mcommunity.AddSection{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cpost.New(c)
	reply.Response(http.StatusOK, svc.AddSection(param))
}

func DelSection(c *gin.Context) {
	reply := errdef.New(c)
	param := &mcommunity.DelSection{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cpost.New(c)
	reply.Response(http.StatusOK, svc.DelSection(param))
}

func AddTopic(c *gin.Context) {
	reply := errdef.New(c)
	param := &mcommunity.AddTopic{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cpost.New(c)
	reply.Response(http.StatusOK, svc.AddTopic(param))
}

func DelTopic(c *gin.Context) {
	reply := errdef.New(c)
	param := &mcommunity.DelTopic{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cpost.New(c)
	reply.Response(http.StatusOK, svc.DelTopic(param))
}

func PostSetting(c *gin.Context) {
	reply := errdef.New(c)
	param := &mposting.SettingParam{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cpost.New(c)
	reply.Response(http.StatusOK, svc.PostSetting(param))
}

func ApplyCreamList(c *gin.Context) {
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	reply := errdef.New(c)
	svc := cpost.New(c)
	code, list := svc.GetApplyCreamList(page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

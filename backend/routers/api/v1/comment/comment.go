package comment

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/backend/controller/comment"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/models/mcomment"
	"sports_service/server/util"
)

// 获取视频评论列表（后台）
func VideoCommentList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	sortType := c.Query("sort_type")

	svc := comment.New(c)
	list := svc.GetVideoComments(sortType, page, size)
	reply.Data["list"] = list
	reply.Data["total"] = svc.GetCommentTotal()
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// 删除视频评论
func DelVideoComments(c *gin.Context) {
	reply := errdef.New(c)
	param := new(mcomment.DelCommentParam)
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := comment.New(c)
	syscode := svc.DelVideoComments(param)
	reply.Response(http.StatusOK, syscode)
}

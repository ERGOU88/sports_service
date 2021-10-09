package comment

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/backend/controller/comment"
	"sports_service/server/global/backend/errdef"
  "sports_service/server/models/mbarrage"
  "sports_service/server/models/mcomment"
	"sports_service/server/util"
)

// 获取视频评论列表（后台）
func VideoCommentList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	condition := c.Query("condition")
	sortType := c.Query("sort_type")
	queryId := c.Query("query_id")

	svc := comment.New(c)
	list, total := svc.GetVideoComments(queryId, sortType, condition, page, size)
	reply.Data["list"] = list
	reply.Data["total"] = total
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// 删除评论
func DelComments(c *gin.Context) {
	reply := errdef.New(c)
	param := new(mcomment.DelCommentParam)
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := comment.New(c)
	syscode := svc.DelComment(param)
	reply.Response(http.StatusOK, syscode)
}

// 弹幕列表
func BarrageList(c *gin.Context) {
  reply := errdef.New(c)
  page, size := util.PageInfo(c.Query("page"), c.Query("size"))
  barrageType := c.DefaultQuery("barrage_type", "0")

  svc := comment.New(c)
  list := svc.GetBarrageList(barrageType, page, size)
  reply.Data["list"] = list
  reply.Data["total"] = svc.GetVideoBarrageTotal(barrageType)
  reply.Response(http.StatusOK, errdef.SUCCESS)
}

// 删除视频/直播 弹幕
func DelVideoBarrage(c *gin.Context) {
  reply := errdef.New(c)
  param := new(mbarrage.DelBarrageParam)
  if err := c.BindJSON(param); err != nil {
    reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
    return
  }

  svc := comment.New(c)
  if err := svc.DelVideoBarrage(param); err != nil {
    reply.Response(http.StatusOK, errdef.VIDEO_BARRAGE_DELETE_FAIL)
    return
  }

  reply.Response(http.StatusOK, errdef.SUCCESS)
}

// 获取帖子评论列表（后台）
func PostCommentList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	condition := c.Query("condition")
	sortType := c.Query("sort_type")
	queryId := c.Query("query_id")

	svc := comment.New(c)
	list, total := svc.GetPostComments(queryId, sortType, condition, page, size)
	reply.Data["list"] = list
	reply.Data["total"] = total
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

func InformationCommentList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	condition := c.Query("condition")
	sortType := c.Query("sort_type")
	queryId := c.Query("query_id")

	svc := comment.New(c)
	list, total := svc.GetPostComments(queryId, sortType, condition, page, size)
	reply.Data["list"] = list
	reply.Data["total"] = total
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

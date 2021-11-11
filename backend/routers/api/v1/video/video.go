package video

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/backend/controller/cvideo"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/middleware/jwt"
	"sports_service/server/models/mlabel"
	"sports_service/server/models/mvideo"
	"sports_service/server/util"
)

// 视频审核 修改视频状态
func EditVideoStatus(c *gin.Context) {
	reply := errdef.New(c)
	param := new(mvideo.EditVideoStatusParam)
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cvideo.New(c)
	syscode := svc.EditVideoStatus(param)
	reply.Response(http.StatusOK, syscode)
}

// 分页获取视频列表（已审核通过的）
func VideoList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	svc := cvideo.New(c)
	list := svc.GetVideoList(page, size)
	total := svc.GetVideoTotalCount()
	reply.Data["list"] = list
	reply.Data["total"] = total
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// 编辑视频是否置顶 0 不置顶 1 置顶
func EditVideoTop(c *gin.Context) {
	reply := errdef.New(c)
	param := new(mvideo.EditTopStatusParam)
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cvideo.New(c)
	syscode := svc.EditVideoTopStatus(param)
	reply.Response(http.StatusOK, syscode)
}

// 编辑视频是否推荐 0 不推荐 1 推荐
func EditVideoRecommend(c *gin.Context) {
	reply := errdef.New(c)
	param := new(mvideo.EditRecommendStatusParam)
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cvideo.New(c)
	syscode := svc.EditVideoRecommendStatus(param)
	reply.Response(http.StatusOK, syscode)
}

// 审核中/审核失败的视频列表
func VideoReviewList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	svc := cvideo.New(c)
	list := svc.GetVideoReviewList(page, size)
	total := svc.GetVideoReviewTotalCount()
	reply.Data["list"] = list
	reply.Data["total"] = total
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// 获取视频标签列表
func VideoLabelList(c *gin.Context) {
	reply := errdef.New(c)
	svc := cvideo.New(c)
	list := svc.GetVideoLabelList()
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// 添加视频标签
func AddVideoLabel(c *gin.Context) {
	reply := errdef.New(c)
	param := new(mlabel.AddVideoLabelParam)
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cvideo.New(c)
	syscode := svc.AddVideoLabel(param)
	reply.Response(http.StatusOK, syscode)
}

func EditVideoLabel(c *gin.Context) {
	reply := errdef.New(c)
	param := &mlabel.AddVideoLabelParam{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cvideo.New(c)
	syscode := svc.EditVideoLabel(param)
	reply.Response(http.StatusOK, syscode)
}

// 删除视频标签
func DelVideoLabel(c *gin.Context) {
	reply := errdef.New(c)
	param := new(mlabel.DelVideoLabelParam)
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cvideo.New(c)
	syscode := svc.DelVideoLabel(param.LabelId)
	reply.Response(http.StatusOK, syscode)
}

// 添加视频分区配置
func AddVideoSubareaConf(c *gin.Context) {
	reply := errdef.New(c)
	param := new(mvideo.AddSubarea)
	if err := c.Bind(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	param.SysId, _ = util.StringToInt(jwt.GetUserInfo(c, consts.IDENTIFY))
	param.SysUser = jwt.GetUserInfo(c, consts.USER_NAME)
	svc := cvideo.New(c)
	syscode := svc.AddVideoSubareaConf(param)
	reply.Response(http.StatusOK, syscode)
}

func EditVideoSubareaConf(c *gin.Context) {
	reply := errdef.New(c)
	param := new(mvideo.AddSubarea)
	if err := c.Bind(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	param.SysId, _ = util.StringToInt(jwt.GetUserInfo(c, consts.IDENTIFY))
	param.SysUser = jwt.GetUserInfo(c, consts.USER_NAME)
	svc := cvideo.New(c)
	syscode := svc.EditVideoSubareaConf(param)
	reply.Response(http.StatusOK, syscode)
}

func DelVideoSubareaConf(c *gin.Context) {
	reply := errdef.New(c)
	param := &mvideo.DelSubarea{}
	if err := c.Bind(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cvideo.New(c)
	syscode := svc.DelVideoSubareaConf(param.Id)
	reply.Response(http.StatusOK, syscode)
}

func VideoSubareaList(c *gin.Context) {
	reply := errdef.New(c)
	svc := cvideo.New(c)
	code, list := svc.GetVideoSubareaList()
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

func BatchEditVideoInfo(c *gin.Context) {
	reply := errdef.New(c)
	param := &mvideo.BatchEditVideos{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cvideo.New(c)
	reply.Response(http.StatusOK, svc.BatchEditVideos(param))
}

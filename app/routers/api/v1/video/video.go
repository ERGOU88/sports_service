package video

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cvideo"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mvideo"
	"sports_service/server/util"
	_ "sports_service/server/models"
	"sports_service/server/tools/tencentCloud/vod"
)

// @Summary 视频发布 (ok)
// @Tags 视频模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   VideoPublishParams  body mvideo.VideoPublishParams true "发布视频请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/publish [post]
// 视频发布
func VideoPublish(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("video_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	params := new(mvideo.VideoPublishParams)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("video_trace: video publish params err:%s, params:%+v", err, params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cvideo.New(c)
	// 用户发布视频
	//if err := svc.UserPublishVideo(userId.(string), params); err != nil {
	//	log.Log.Errorf("video_trace: video publish failed, err:%s", err)
	//	reply.Response(http.StatusOK, errdef.VIDEO_PUBLISH_FAIL)
	//	return
	//}

	// 修改逻辑 用户发布视频 先记录到缓存
	syscode := svc.RecordPubVideoInfo(userId.(string), params)
	reply.Response(http.StatusOK, syscode)

}

// @Summary 用户浏览过的视频记录[分页获取] (ok)
// @Tags 视频模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   page	  	  query  	string 	true  "页码 从1开始"
// @Param   size	  	  query  	string 	true  "每页展示多少 最多50条"
// @Success 200 {array}  mvideo.VideosInfoResp
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/browse/history [get]
// 视频浏览记录
func BrowseHistory(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("video_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := cvideo.New(c)
	// 获取用户浏览过的视频列表
	list := svc.UserBrowseVideosRecord(userId.(string), page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 用户发布的视频记录[分页获取] (ok)
// @Tags 视频模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   page	  	  query  	string 	true  "页码 从1开始"
// @Param   size	  	  query  	string 	true  "每页展示多少 最多50条"
// @Param   status	  	  query  	string 	true  "status 状态 -1 查询全部 0 审核中 1 已发布 2 不通过"
// @Param   condition	  query  	string 	true  "条件 -1 默认时间排序 0 播放数 1 弹幕数 2 点赞数 3 评论数 4 分享数"
// @Success 200 {array}  mvideo.VideosInfo
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/publish/list [get]
// 用户发布的视频列表
func VideoPublishList(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("video_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	status := c.DefaultQuery("status", "-1")
	condition := c.DefaultQuery("condition", "-1")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	svc := cvideo.New(c)
	// 获取用户发布的内容列表
	list := svc.GetUserPublishList(userId.(string), status, condition, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 删除浏览的历史记录 (ok)
// @Tags 视频模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   DeleteHistoryParam  body mvideo.DeleteHistoryParam true "删除浏览历史记录 请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/delete/history [post]
// 删除历史记录
func DeleteHistory(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("video_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	param := new(mvideo.DeleteHistoryParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("video_trace: delete history params err:%s, params:%+v", err, param)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cvideo.New(c)
	// 删除历史记录
	syscode := svc.DeleteHistoryByIds(userId.(string), param)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 删除发布的视频 (ok)
// @Tags 视频模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   DeletePublishParam  body mvideo.DeletePublishParam true "删除发布的视频 请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/delete/publish [post]
// 删除发布的视频
func DeletePublish(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("video_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	param := new(mvideo.DeletePublishParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("video_trace: delete publish params err:%s, params:%+v", err, param)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cvideo.New(c)
	// 删除发布的视频
	syscode := svc.DeletePublishVideo(userId.(string), param.ComposeIds)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 首页推荐的视频列表[分页获取] (ok)
// @Tags 视频模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   page	  	  query  	string 	true  "页码 从1开始"
// @Param   size	  	  query  	string 	true  "每页展示多少 最多50条"
// @Success 200 {array}  mvideo.VideoDetailInfo
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/recommend [get]
// 首页推荐列表
func RecommendVideos(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	userId, _ := c.Get(consts.USER_ID)

	svc := cvideo.New(c)
	list := svc.GetRecommendVideos(userId.(string), page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 视频首页推荐的banner (ok)
// @Tags 视频模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Success 200 {array}  models.Banner
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/homepage/banner [get]
// 首页推荐的banner
func RecommendBanners(c *gin.Context) {
	reply := errdef.New(c)
	svc := cvideo.New(c)
	banners := svc.GetRecommendBanners()
	reply.Data["banners"] = banners
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 首页关注的用户发布的视频列表[分页获取] (ok)
// @Tags 视频模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   page	  	  query  	string 	true  "页码 从1开始"
// @Param   size	  	  query  	string 	true  "每页展示多少 最多50条"
// @Success 200 {array}  mvideo.VideoDetailInfo
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/attention [get]
// 关注的人发布的视频
func AttentionVideos(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	userId, _ := c.Get(consts.USER_ID)

	svc := cvideo.New(c)
	list := svc.GetAttentionVideos(userId.(string), page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 视频详情 (ok)
// @Tags 视频模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   video_id	  query  	string 	true  "视频id"
// @Success 200 {object}  mvideo.VideoDetailInfo
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/detail [get]
// 视频详情信息
func VideoDetail(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	videoId := c.Query("video_id")

	svc := cvideo.New(c)
	detail := svc.GetVideoDetail(userId.(string), videoId)
	reply.Data["detail"] = detail
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 视频详情 (ok)
// @Tags 视频模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   video_id	  query  	string 	true  "视频id"
// @Param   page	  	  query  	string 	true  "页码 从1开始"
// @Param   size	  	  query  	string 	true  "每页展示多少 最多50条"
// @Success 200 {array}  mvideo.VideoDetailInfo
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/detail/recommend [get]
// 详情页推荐视频（同标签推荐）
func DetailRecommend(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	videoId := c.Query("video_id")

	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := cvideo.New(c)
	// 获取详情页推荐视频（同标签推荐）
	list := svc.GetDetailRecommend(userId.(string), videoId, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// 获取上传签名（腾讯云）
func UploadSign(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)

	svc := cvideo.New(c)
	syscode, sign, taskId := svc.GetUploadSign(userId.(string))
	reply.Data["sign"] = sign
	reply.Data["task_id"] = taskId
	reply.Response(http.StatusOK, syscode)
}

// 事件回调（腾讯云）
func EventCallback(c *gin.Context) {
	reply := errdef.New(c)
	callback := new(vod.EventNotify)
	if err := c.BindJSON(callback); err != nil {
		log.Log.Errorf("video_trace: callback params err:%s, params:%+v", err, callback)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}


}

// 用户上传自定义标签
func CheckCustomLabels(c *gin.Context) {
  reply := errdef.New(c)
  userId, _ := c.Get(consts.USER_ID)
  params := new(mvideo.CustomLabelParams)
  if err := c.BindJSON(params); err != nil {
    log.Log.Errorf("video_trace: custom labels params err:%s, params:%+v", err, params)
    reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
    return
  }

  svc := cvideo.New(c)
  // 检测自定义标签文本
  syscode := svc.CheckCustomLabel(userId.(string), params)
  reply.Response(http.StatusOK, syscode)
}

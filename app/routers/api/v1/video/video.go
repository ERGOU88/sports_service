package video

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cvideo"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mvideo"
	_ "sports_service/server/models/mlabel"
	cloud "sports_service/server/tools/tencentCloud"
	"sports_service/server/util"
	_ "sports_service/server/models"
	"sports_service/server/tools/tencentCloud/vod"
	"strconv"
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

	log.Log.Errorf("##### publish video params:%+v", params)

	svc := cvideo.New(c)
	// 修改逻辑 用户发布视频 先记录到缓存
	syscode := svc.RecordPubVideoInfo(userId.(string), params)
	log.Log.Errorf("##### publish video syscode:%d", syscode)
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
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/browse/history [get]
// 视频浏览记录 todo: 分页数据重复问题 客户端传递最后一条记录创建时间
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
// @Param   user_id	    query  	string 	true  "查看的用户id"
// @Param   page	  	  query  	string 	true  "页码 从1开始"
// @Param   size	  	  query  	string 	true  "每页展示多少 最多50条"
// @Param   status	  	  query  	string 	true  "status 状态 -1 查询全部 0 审核中 1 已发布 2 不通过"
// @Param   condition	  query  	string 	true  "条件 -1 默认时间排序 0 播放数 1 弹幕数 2 点赞数 3 评论数 4 分享数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/publish/list [get]
// 用户发布的视频列表
func VideoPublishList(c *gin.Context) {
	reply := errdef.New(c)
	//userId, ok := c.Get(consts.USER_ID)
	//if !ok {
	//	log.Log.Errorf("video_trace: user not found, uid:%s", userId.(string))
	//	reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
	//	return
	//}
	uid := c.Query("user_id")
	status := c.DefaultQuery("status", "-1")
	condition := c.DefaultQuery("condition", "-1")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	svc := cvideo.New(c)
	// 获取用户发布的内容列表
	list := svc.GetUserPublishList(uid, status, condition, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 查看其他用户发布的视频记录[分页获取] (ok)
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
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/other/publish [get]
// 获取其他用户发布的视频列表
func OtherUserPublishList(c *gin.Context) {
	reply := errdef.New(c)
	userId := c.Query("user_id")

	status := c.DefaultQuery("status", "-1")
	condition := c.DefaultQuery("condition", "-1")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	svc := cvideo.New(c)
	// 获取其他用户发布的内容列表
	list := svc.GetUserPublishList(userId, status, condition, page, size)
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
// @Param   DeleteHistoryParam  body mvideo.DeleteHistoryParam true "删除浏览历史记录 请求参数 -1删除所有"
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
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/recommend [get]
// 首页推荐列表
func RecommendVideos(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	//userId, _ := c.Get(consts.USER_ID)
	// 视频id
	index := c.Query("id")
	if index == "-1" {
		index = fmt.Sprint(1e6)
	}

	userId := c.Query("user_id")
	svc := cvideo.New(c)
	minId, list := svc.GetRecommendVideos(userId, index, page, size)
	reply.Data["id"] = minId
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
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
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
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
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
// @Param   user_id	    query  	string 	true  "用户id"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/detail [get]
// 视频详情信息
func VideoDetail(c *gin.Context) {
	reply := errdef.New(c)
	//userId, _ := c.Get(consts.USER_ID)
	videoId := c.Query("video_id")
	userId := c.Query("user_id")

	svc := cvideo.New(c)
	detail, syscode := svc.GetVideoDetail(userId, videoId)
	if syscode != errdef.SUCCESS {
		reply.Response(http.StatusOK, syscode)
		return
	}

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
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/detail/recommend [get]
// 详情页推荐视频（同标签推荐）
func DetailRecommend(c *gin.Context) {
	reply := errdef.New(c)
	//userId, _ := c.Get(consts.USER_ID)
	videoId := c.Query("video_id")
	userId := c.Query("user_id")

	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := cvideo.New(c)
	// 获取详情页推荐视频（同标签推荐）
	list := svc.GetDetailRecommend(userId, videoId, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 获取上传签名（腾讯云） (ok)
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
// @Param   bite_rate	  query  	string 	true  "视频码率"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/upload/sign [get]
// 获取上传签名（腾讯云）
func UploadSign(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	biteRate, err := strconv.Atoi(c.Query("bite_rate"))
	if err != nil {
		biteRate = 0
	}

	svc := cvideo.New(c)
	syscode, sign, taskId := svc.GetUploadSign(userId.(string), int64(biteRate))
	reply.Data["sign"] = sign
	reply.Data["task_id"] = taskId

	log.Log.Errorf("#####taskId:%s", taskId)
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

// @Summary 视频标签 (ok)
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
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/label/list [get]
// 获取视频标签列表
func VideoLabelList(c *gin.Context) {
	reply := errdef.New(c)
	svc := cvideo.New(c)
	list := svc.GetVideoLabelList()
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 举报视频 (ok)
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
// @Param   VideoReportParam  body mvideo.VideoReportParam true "举报视频请求参数 游客userid传空字符串"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/report [post]
// 举报视频
func VideoReport(c *gin.Context) {
	reply := errdef.New(c)
	param := new(mvideo.VideoReportParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("video_trace: video report params err:%s, params:%+v", err, param)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cvideo.New(c)
	syscode := svc.AddVideoReport(param)
	reply.Response(http.StatusOK, syscode)
}

func TestUpload(c *gin.Context) {
	reply := errdef.New(c)
	svc := cvideo.New(c)
	syscode, _, taskId := svc.GetUploadSign("202009101933004667", 2500)
	if syscode != errdef.SUCCESS {
		reply.Response(http.StatusOK, syscode)
		return
	}

	client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
	resp, err := client.Upload(taskId,"202009101933004667", "", "/Users/jelly/go/src/sports_service/server/tools/tencentCloud/test.mp4",
		"ap-shanghai", consts.VOD_PROCEDURE_NAME)
	if err != nil {
		fmt.Printf("upload err:%s", err)
		reply.Response(http.StatusOK, errdef.ERROR)
		return
	}

	params := new(mvideo.VideoPublishParams)
	params.Title = "test"
	params.Describe = "test"
	params.FileId = *resp.Response.FileId
	params.VideoAddr = *resp.Response.MediaUrl
	params.Cover = *resp.Response.CoverUrl
	params.VideoLabels = "1,2"
	params.TaskId = taskId
	if syscode := svc.RecordPubVideoInfo("202009101933004667", params); syscode != errdef.SUCCESS {
		reply.Response(http.StatusOK, errdef.VIDEO_PUBLISH_FAIL)
		return
	}


	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 记录用户视频播放时长 (ok)
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
// @Param   PlayDurationParams  body mvideo.PlayDurationParams true "用户播放视频的时长"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/video/record/play/duration [post]
// 记录视频播放时长
func RecordPlayDuration(c *gin.Context) {
	reply := errdef.New(c)
	param := new(mvideo.PlayDurationParams)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("video_trace: invalid param, param:%+v, err:%s", param, err)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cvideo.New(c)
	syscode := svc.RecordPlayDuration(param)
	reply.Response(http.StatusOK, syscode)
}

// /api/v1/video/subarea
// 视频分区
func VideoSubarea(c *gin.Context) {
	reply := errdef.New(c)
	svc := cvideo.New(c)

	code, list := svc.GetVideoSubarea()
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

// /api/v1/video/create/album
// 添加视频专辑
func CreateVideoAlbum(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	param := new(mvideo.CreateAlbumParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("video_trace: create album param fail, err:%s, param:%+v", err, param)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cvideo.New(c)
	syscode, album := svc.CreateVideoAlbum(userId.(string), param)
	if syscode == errdef.SUCCESS {
		reply.Data["album"] = album
	}

	reply.Response(http.StatusOK, syscode)
}

// /api/v1/video/add/album
// 将视频添加到专辑内
func AddVideoToAlbum(c *gin.Context) {
	reply := errdef.New(c)
	svc := cvideo.New(c)
	param := new(mvideo.AddVideoToAlbumParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("video_trace: add video to album param fail, err:%s, param:%+v", err, param)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	userId, _ := c.Get(consts.USER_ID)
	syscode := svc.AddVideoToAlbum(userId.(string), param)
	reply.Response(http.StatusOK, syscode)
}

// 获取分区下的视频列表
func VideoListBySubarea(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	subareaId := c.Query("subarea_id")
	userId := c.Query("user_id")

	svc := cvideo.New(c)
	syscode, list := svc.GetVideoListBySubarea(subareaId, userId, page, size)
	reply.Data["list"] = list

	reply.Response(http.StatusOK, syscode)
}

// 用户发布的视频专辑列表
func VideoAlbumList(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := cvideo.New(c)
	code, list := svc.GetVideoAlbumByUserId(userId.(string), page, size)
	reply.Data["list"] = list

	reply.Response(http.StatusOK, code)

}

func HomePageSectionInfo(c *gin.Context) {
	reply := errdef.New(c)
	svc := cvideo.New(c)
	sectionType := c.DefaultQuery("section_type", "0")
	code, list := svc.GetHomepageSectionInfo(sectionType)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

func SectionRecommendInfo(c *gin.Context) {
	reply := errdef.New(c)
	svc := cvideo.New(c)
	userId := c.Query("user_id")
	sectionId := c.Query("section_id")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	code, list := svc.GetRecommendInfoBySection(userId, sectionId, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}



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
	if err := svc.UserPublishVideo(userId.(string), params); err != nil {
		log.Log.Errorf("video_trace: video publish failed, err:%s", err)
		reply.Response(http.StatusOK, errdef.VIDEO_PUBLISH_FAIL)
		return
	}

	reply.Response(http.StatusOK, errdef.SUCCESS)

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

// 发布的视频列表
// status 状态 -1 查询全部 0 审核中 1 已发布 2 不通过
// condition 条件 -1 默认时间排序 1 播放数 2 弹幕数 3 评论数 4 点赞数 5 分享数
func VideoPublishList(c *gin.Context) {
	//status := c.DefaultQuery("status", "-1")
	//condition := c.DefaultQuery("condition", "-1")


}

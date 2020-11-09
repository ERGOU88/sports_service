package barrage

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cbarrage"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mbarrage"
	_ "sports_service/server/models"
)

// @Summary 发送视频弹幕 (ok)
// @Tags 弹幕模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   SendBarrageParams  body mbarrage.SendBarrageParams true "发送弹幕请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 400 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/barrage/send [post]
// 发送弹幕
func SendBarrage(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok || userId == "" {
		log.Log.Errorf("barrage_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	params := new(mbarrage.SendBarrageParams)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("barrage_trace: invalid params, params:%+v", params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cbarrage.New(c)
	// 发送弹幕
	syscode := svc.SendVideoBarrage(userId.(string), params)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 视频弹幕列表[视频时长区间获取] (ok)
// @Tags 弹幕模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   min_duration  query     string  true  "时长区间 最小时长"
// @Param   max_duration  query     string  true  "时长区间 最大时长"
// @Param   video_id      query     string  true  "视频id"
// @Success 200 {object}  models.VideoBarrage
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/barrage/video/list [get]
// 视频弹幕列表（通过时长区间查询）
func VideoBarrage(c *gin.Context) {
	reply := errdef.New(c)
	minDuration := c.Query("min_duration")
	maxDuration := c.Query("max_duration")
	videoId := c.Query("video_id")

	svc := cbarrage.New(c)
	// 获取视频弹幕列表
	syscode, list := svc.GetVideoBarrageList(videoId, minDuration, maxDuration)
	if syscode != errdef.SUCCESS {
	  reply.Response(http.StatusOK, errdef.BARRAGE_VIDEO_LIST_FAIL)
	  return
  }

	reply.Data["list"] = list
	reply.Response(http.StatusOK, syscode)
}

package collect

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/collect"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mcollect"
	"sports_service/server/util"
	_ "sports_service/server/models/mvideo"
)

// @Summary 收藏视频 (ok)
// @Tags 收藏模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   AddCollectParam  body mcollect.AddCollectParam true "收藏视频请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/collect/video [post]
// 收藏视频
func CollectVideo(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("collect_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	param := new(mcollect.AddCollectParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("collect_trace: invalid param, param:%+v", param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := collect.New(c)
	// 添加收藏
	syscode := svc.AddCollect(userId.(string), param.VideoId)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 取消收藏 (ok)
// @Tags 收藏模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   CancelCollectParam  body mcollect.CancelCollectParam true "取消视频收藏请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/collect/video/cancel [post]
// 取消收藏
func CancelCollect(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("collect_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	param := new(mcollect.CancelCollectParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("collect_trace: invalid param, param:%+v", param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := collect.New(c)
	// 取消收藏
	syscode := svc.CancelCollect(userId.(string), param.VideoId)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 用户收藏的视频列表[分页获取] (ok)
// @Tags 收藏模块
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
// @Success 200 {array}   mvideo.VideosInfoResp
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/collect/video/list [get]
// 用户收藏的视频列表
func CollectVideoList(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("collect_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := collect.New(c)
	// 获取用户收藏的视频列表
	list := svc.GetUserCollectVideos(userId.(string), page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

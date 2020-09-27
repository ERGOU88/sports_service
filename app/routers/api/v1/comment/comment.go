package comment

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/comment"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mcomment"
	"sports_service/server/util"
)

// @Summary 发布评论 (ok)
// @Tags 评论模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   PublishCommentParams  body mcomment.PublishCommentParams true "发布评论请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/comment/publish [post]
// 发布评论
func PublishComment(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok || userId == "" {
		log.Log.Error("comment_trace: need login")
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	params := new(mcomment.PublishCommentParams)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("comment_trace: publish comment params err:%s, params:%+v", err, params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := comment.New(c)
	syscode := svc.PublishComment(userId.(string), params)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 回复评论 (ok)
// @Tags 评论模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   ReplyCommentParams  body mcomment.ReplyCommentParams true "回复评论请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/comment/reply [post]
// 回复评论
func PublishReply(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok || userId == "" {
		log.Log.Error("comment_trace: need login")
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	params := new(mcomment.ReplyCommentParams)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("comment_trace: publish reply params err:%s, params:%+v", err, params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := comment.New(c)
	syscode := svc.PublishReply(userId.(string), params)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 评论列表[分页获取] (ok)
// @Tags 评论模块
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
// @Param   video_id	  query  	string 	true  "视频id"
// @Param   sort_type     query     string  true  "排序规则 0 时间 1 热度"
// @Success 200 {array}  mcomment.VideoComments
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/comment/list [get]
// 评论列表
func CommentList(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	videoId := c.Query("video_id")
	sortType := c.DefaultQuery("sort_type", "0")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	svc := comment.New(c)
	list := svc.GetVideoComments(userId.(string), videoId, sortType, page, size)

	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 回复列表[分页获取] (ok)
// @Tags 评论模块
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
// @Param   video_id	  query  	string 	true  "视频id"
// @Param   comment_id	  query  	string 	true  "评论id"
// @Success 200 {array}  mcomment.ReplyComment
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/comment/reply/list [get]
// 回复列表
func ReplyList(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	commentId := c.Query("comment_id")
	videoId := c.Query("video_id")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	svc := comment.New(c)
	// 获取评论回复列表
	syscode, list := svc.GetCommentReplyList(userId.(string), videoId, commentId, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, syscode)
}
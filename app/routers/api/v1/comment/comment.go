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
	"strconv"
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
	syscode, commentId := svc.PublishComment(userId.(string), params)
	reply.Data["comment_id"] = commentId
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
	syscode, commentId := svc.PublishReply(userId.(string), params)
	reply.Data["comment_id"] = commentId
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
// @Param   user_id	    query  	string 	true  "用户id"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/comment/list [get]
// 评论列表
func CommentList(c *gin.Context) {
	reply := errdef.New(c)
	userId := c.Query("user_id")
	composeId := c.Query("compose_id")
	sortType := c.DefaultQuery("sort_type", "0")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	commentId := c.Query("comment_id")
	// 1 视频评论 2 帖子评论
	commentType, _ := strconv.Atoi(c.DefaultQuery("comment_type", "0"))

	svc := comment.New(c)
	syscode, list := svc.GetComments(userId, composeId, sortType, commentType, page, size)
	if syscode != errdef.SUCCESS {
		reply.Response(http.StatusOK, syscode)
		return
	}

	first := svc.GetFirstComment(userId, commentId)

	reply.Data["first"] = first
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
// @Param   user_id	    query  	string 	true  "用户id"
// @Param   comment_id	  query  	string 	true  "评论id"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/comment/reply/list [get]
// 回复列表
func ReplyList(c *gin.Context) {
	reply := errdef.New(c)
	userId := c.Query("user_id")
	commentId := c.Query("comment_id")
	videoId := c.Query("video_id")
	// todo: 替换videoId
	//composeId := c.Query("compose_id")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	// 1 视频回复 2 帖子回复
	commentType, _ := strconv.Atoi(c.DefaultQuery("comment_type", "1"))

	svc := comment.New(c)
	// 获取评论回复列表
	syscode, list := svc.GetCommentReplyList(userId, videoId, commentId, commentType, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, syscode)
}

// @Summary 举报评论 (ok)
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
// @Param   CommentReportParam  body mcomment.CommentReportParam true "举报评论请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/comment/report [post]
// 举报评论
func CommentReport(c *gin.Context) {
	reply := errdef.New(c)
	params := new(mcomment.CommentReportParam)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("comment_trace: comment report params err:%s, params:%+v", err, params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := comment.New(c)
	syscode := svc.AddCommentReport(params)
	reply.Response(http.StatusOK, syscode)
}

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
// @Param   PublishCommentParams  body mcomment.V2PubCommentParams true "发布评论请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/comment/publish [post]
// 新版发布评论
func V2PublishComment(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok || userId == "" {
		log.Log.Error("comment_trace: need login")
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	params := new(mcomment.V2PubCommentParams)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("comment_trace: publish comment params err:%s, params:%+v", err, params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := comment.New(c)
	syscode, commentId := svc.V2PublishComment(userId.(string), params)
	reply.Data["comment_id"] = commentId
	reply.Response(http.StatusOK, syscode)
}

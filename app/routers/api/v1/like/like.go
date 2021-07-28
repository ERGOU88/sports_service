package like

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/clike"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mlike"
	"sports_service/server/util"
	_ "sports_service/server/models/mvideo"
)

// @Summary 视频点赞 (ok)
// @Tags 点赞模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   GiveLikeParam  body mlike.GiveLikeParam true "点赞视频请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/like/video [post]
// 视频点赞
func GiveLikeForVideo(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("like_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	param := new(mlike.GiveLikeParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("like_trace: invalid param, param:%+v", param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := clike.New(c)
	// 视频点赞
	syscode := svc.GiveLikeForVideo(userId.(string), param.ComposeId)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 取消视频点赞 (ok)
// @Tags 点赞模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   CancelLikeParam  body mlike.CancelLikeParam true "取消视频点赞请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/like/video/cancel [post]
// 取消点赞
func CancelLikeForVideo(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("like_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	param := new(mlike.CancelLikeParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("like_trace: invalid param, param:%+v", param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := clike.New(c)
	// 取消点赞
	syscode := svc.CancelLikeForVideo(userId.(string), param.ComposeId)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 用户点赞的视频列表[分页获取] (ok)
// @Tags 点赞模块
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
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/like/video/list [get]
// 用户点赞的视频列表
func LikeVideoList(c *gin.Context) {
	reply := errdef.New(c)
	//userId, ok := c.Get(consts.USER_ID)
	//if !ok {
	//	log.Log.Errorf("like_trace: user not found, uid:%s", userId.(string))
	//	reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
	//	return
	//}

	uid := c.Query("user_id")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := clike.New(c)
	// 获取用户点赞的视频列表
	list := svc.GetUserLikeVideos(uid, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 查看其他用户点赞的视频列表[分页获取] (ok)
// @Tags 点赞模块
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
// @Router /api/v1/like/other/list [get]
// 用户点赞的视频列表
func OtherUserLikeVideoList(c *gin.Context) {
	reply := errdef.New(c)
	userId := c.Query("user_id")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := clike.New(c)
	// 获取用户点赞的视频列表
	list := svc.GetUserLikeVideos(userId, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 视频评论点赞 (ok)
// @Tags 点赞模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   GiveLikeParam  body mlike.GiveLikeParam true "评论点赞请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/like/comment [post]
// 评论点赞
func GiveLikeForVideoComment(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("like_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	param := new(mlike.GiveLikeParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("like_trace: invalid param, param:%+v", param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := clike.New(c)
	// 视频/帖子 评论点赞
	syscode := svc.GiveLikeForComment(userId.(string), param.ComposeId, param.CommentType)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 取消评论点赞 (ok)
// @Tags 点赞模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   CancelLikeParam  body mlike.CancelLikeParam true "取消评论点赞请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/like/comment/cancel [post]
// 评论取消点赞
func CancelLikeForVideoComment(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("like_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	param := new(mlike.CancelLikeParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("like_trace: invalid param, param:%+v", param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := clike.New(c)
	// 取消点赞
	syscode := svc.CancelLikeForComment(userId.(string), param.ComposeId, param.CommentType)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 帖子点赞 (ok)
// @Tags 点赞模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   GiveLikeParam  body mlike.GiveLikeParam true "点赞视频请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/like/post [post]
// 帖子点赞
func GiveLikeForPost(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("like_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	param := new(mlike.GiveLikeParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("like_trace: invalid param, param:%+v", param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := clike.New(c)
	// 视频点赞
	syscode := svc.GiveLikeForPost(userId.(string), param.ComposeId)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 取消帖子点赞 (ok)
// @Tags 点赞模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   CancelLikeParam  body mlike.CancelLikeParam true "取消视频点赞请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/like/post/cancel [post]
// 取消帖子点赞
func CancelLikeForPost(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("like_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	param := new(mlike.CancelLikeParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("like_trace: invalid param, param:%+v", param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := clike.New(c)
	// 取消帖子点赞
	syscode := svc.CancelLikeForPost(userId.(string), param.ComposeId)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 帖子评论点赞 (ok)
// @Tags 点赞模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   GiveLikeParam  body mlike.GiveLikeParam true "评论点赞请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/like/comment/post [post]
// 帖子评论点赞
func GiveLikeForPostComment(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("like_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	param := new(mlike.GiveLikeParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("like_trace: invalid param, param:%+v", param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := clike.New(c)
	// 帖子评论点赞
	syscode := svc.GiveLikeForPostComment(userId.(string), param.ComposeId)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 取消帖子评论点赞 (ok)
// @Tags 点赞模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   CancelLikeParam  body mlike.CancelLikeParam true "取消评论点赞请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/like/comment/post/cancel [post]
// 帖子评论取消点赞
func CancelLikeForPostComment(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("like_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	param := new(mlike.CancelLikeParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("like_trace: invalid param, param:%+v", param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := clike.New(c)
	// 取消帖子评论点赞
	syscode := svc.CancelLikeForPostComment(userId.(string), param.ComposeId)
	reply.Response(http.StatusOK, syscode)
}

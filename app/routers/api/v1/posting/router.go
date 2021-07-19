package posting

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cposting"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mposting"
	"sports_service/server/util"
	"fmt"
)

// /api/v1/post/publish
// 发布帖子
func PublishPosting(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("post_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}
	//userId := "13918242"

	params := new(mposting.PostPublishParam)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("post_trace: post publish params err:%s, params:%+v", err, params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cposting.New(c)
	code := svc.PublishPosting(userId.(string), params)
	reply.Response(http.StatusOK, code)
}

// /api/v1/post/publish/detail
// 帖子详情
func PostDetail(c *gin.Context) {
	reply := errdef.New(c)
	userId := c.Query("user_id")

	postId := c.Query("post_id")
	svc := cposting.New(c)
	detail, code := svc.GetPostDetail(userId, postId)
	if code == errdef.SUCCESS {
		reply.Data["detail"] = detail
	}

	reply.Response(http.StatusOK, code)
}

// /api/v1/post/publish/list
// 用户发布的视频列表
func PostPublishList(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("post_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	//userId := "13918242"

	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := cposting.New(c)
	list := svc.GetPostPublishListByUser(userId.(string), page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)

}

// 用户删除发布的帖子
func DeletePublishPost(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		log.Log.Errorf("post_trace: user not found, uid:%s", userId.(string))
		reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
		return
	}

	param := new(mposting.DeletePostParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("post_trace: delete publish params err:%s, params:%+v", err, param)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cposting.New(c)
	// 删除发布的帖子
	syscode := svc.DeletePublishPost(userId.(string), fmt.Sprint(param.PostId))
	reply.Response(http.StatusOK, syscode)
}

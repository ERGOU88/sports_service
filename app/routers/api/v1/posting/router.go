package posting

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cposting"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/models/mposting"
)

// 发布帖子
func PublishPosting(c *gin.Context) {
	reply := errdef.New(c)
	//userId, ok := c.Get(consts.USER_ID)
	//if !ok {
	//	log.Log.Errorf("post_trace: user not found, uid:%s", userId.(string))
	//	reply.Response(http.StatusOK, errdef.USER_NOT_EXISTS)
	//	return
	//}
	userId := "13918242"

	params := new(mposting.PostPublishParam)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("post_trace: post publish params err:%s, params:%+v", err, params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cposting.New(c)
	code := svc.PublishPosting(userId, params)
	reply.Response(http.StatusOK, code)
}

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

package post

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/backend/controller/cpost"
	"sports_service/server/global/app/log"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/models/mposting"
)

// 帖子审核
func AuditPost(c *gin.Context) {
	reply := errdef.New(c)
	param := &mposting.AudiPostParam{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("post_trace: invalid param, err:%s", err)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cpost.New(c)
	syscode := svc.AudiPost(param)
	reply.Response(http.StatusOK, syscode)
}

// 帖子列表 todo：展示数据待确认
func PostList(c *gin.Context) {

}

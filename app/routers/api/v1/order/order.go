package order

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/corder"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/util"
)

//
func OrderList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	userId, _ := c.Get(consts.USER_ID)
	status := c.Query("status")
	svc := corder.New(c)
	code, list := svc.GetOrderList(userId.(string), status, page, size)
	if code == errdef.SUCCESS {
		reply.Data["list"] = list
	}

	reply.Response(http.StatusOK, code)
}

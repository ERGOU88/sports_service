package information

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cinformation"
	"sports_service/server/global/app/errdef"
	"sports_service/server/util"
)

func InformationList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := cinformation.New(c)
	code, list := svc.GetInformationList(page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

func InformationDetail(c *gin.Context) {
	reply := errdef.New(c)
	id := c.Query("id")
	userId := c.Query("user_id")
	svc := cinformation.New(c)
	code, detail := svc.GetInformationDetail(id, userId)
	reply.Data["detail"] = detail
	reply.Response(http.StatusOK, code)
}

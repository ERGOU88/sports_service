package information

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/util"
	"sports_service/server/backend/controller/cinformation"
)

func InformationList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := cinformation.New(c)
	code, list := svc.GetInformationList(page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

func DeleteInformation(c *gin.Context) {
	reply := errdef.New(c)
	id := c.Query("id")
	svc := cinformation.New(c)
	reply.Response(http.StatusOK, svc.DeleteInformation(id))
}

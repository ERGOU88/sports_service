package information

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/backend/controller/cinformation"
	"sports_service/global/backend/errdef"
	"sports_service/util"
)

func InformationList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := cinformation.New(c)
	code, list := svc.GetInformationList(page, size)
	reply.Data["list"] = list
	reply.Data["total"] = svc.GetTotalNumByInformation()
	reply.Response(http.StatusOK, code)
}

func DeleteInformation(c *gin.Context) {
	reply := errdef.New(c)
	id := c.Query("id")
	svc := cinformation.New(c)
	reply.Response(http.StatusOK, svc.DeleteInformation(id))
}

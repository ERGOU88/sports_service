package finance

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/backend/controller/cfinance"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/util"
)

func OrderList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := cfinance.New(c)
	code, list := svc.GetOrderList(page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

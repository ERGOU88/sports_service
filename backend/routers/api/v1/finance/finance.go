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
	reply.Data["total"] = svc.GetOrderTotal()
	reply.Response(http.StatusOK, code)
}

func RefundList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := cfinance.New(c)
	code, list := svc.GetRefundList(c.Query("order_id"), page, size)
	reply.Data["list"] = list
	reply.Data["total"] = svc.GetRefundRecordTotal()
	reply.Response(http.StatusOK, code)
}

func RevenueFlow(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	minDate := c.Query("min_date")
	maxDate := c.Query("max_date")
	svc := cfinance.New(c)
	code, total, list := svc.GetRevenueFlow(minDate, maxDate, c.Query("order_id"), page, size)
	reply.Data["list"] = list
	reply.Data["total_revenue"] = total
	reply.Response(http.StatusOK, code)
}

func TopStat(c *gin.Context) {
	reply := errdef.New(c)
	svc := cfinance.New(c)
	code, stat := svc.TopStat()
	reply.Data["detail"] = stat
	reply.Response(http.StatusOK, code)
}

func ChartStat(c *gin.Context) {
	reply := errdef.New(c)
	svc := cfinance.New(c)
	code, detail := svc.GetChartStat(c.Query("min_date"), c.Query("max_date"))
	reply.Data["detail"] = detail
	reply.Response(http.StatusOK, code)
}

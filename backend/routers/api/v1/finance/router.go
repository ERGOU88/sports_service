package finance

import (
	"github.com/gin-gonic/gin"
	"sports_service/middleware/jwt"
)

func Router(engine *gin.Engine) {
	api := engine.Group("/backend/v1")
	finance := api.Group("/finance")
	finance.Use(jwt.JwtAuth())
	{
		// 订单流水列表
		finance.GET("/order/list", OrderList)
		// 退款流水列表
		finance.GET("/refund/list", RefundList)
		// 收益流水
		finance.GET("/revenue/flow", RevenueFlow)
		// 财务首页 顶部栏统计
		finance.GET("/top/stat", TopStat)
		// 财务首页 图表统计
		finance.GET("/chart/stat", ChartStat)
	}
}

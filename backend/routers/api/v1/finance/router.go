package finance

import "github.com/gin-gonic/gin"

func Router(engine *gin.Engine) {
	api := engine.Group("/backend/v1")
	finance := api.Group("/finance")
	{
		// 订单流水列表
		finance.GET("/order/list", OrderList)
	}
}

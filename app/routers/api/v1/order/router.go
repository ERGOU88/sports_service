package order

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
	"sports_service/server/middleware/token"
)

// 订单模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	order := api.Group("/order")
	order.Use(sign.CheckSign(), token.TokenAuth())
	{
		// 订单列表
		order.GET("/list", OrderList)
		// 订单详情
		order.GET("/detail", OrderDetail)
	}
}


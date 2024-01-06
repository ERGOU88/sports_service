package order

import (
	"github.com/gin-gonic/gin"
	"sports_service/middleware/sign"
	"sports_service/middleware/token"
)

// 订单模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	api.GET("/order/gift/detail", sign.CheckSign(), GiftDetail)

	order := api.Group("/order")
	order.Use(sign.CheckSign(), token.TokenAuth())
	{
		// 订单列表
		order.GET("/list", OrderList)
		// 订单详情
		order.GET("/detail", OrderDetail)
		// 订单退款
		order.POST("/refund", OrderRefund)
		// 删除订单
		order.POST("/delete", OrderDelete)
		// 查看券码
		order.GET("/coupon/code", OrderCouponCode)
		// 取消订单
		order.POST("/cancel", OrderCancel)
		// 退款规则
		order.GET("/refund/rules", RefundRules)
		// 领取订单赠礼
		order.POST("/receive/gift", ReceiveGift)
	}
}

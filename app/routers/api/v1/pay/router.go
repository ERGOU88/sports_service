package pay

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
	"sports_service/server/middleware/token"
)

// 支付模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	pay := api.Group("/pay")
	{
		// app支付
		pay.POST("/app/trade", sign.CheckSign(), token.TokenAuth(), AppPay)
		// 支付宝支付/退款 回调通知
		pay.POST("/alipay/notify", AliPayNotify)
		// 微信支付回调通知
		pay.POST("/wechat/notify", WechatNotify)
		// 微信退款回调通知
		pay.POST("/wechat/refund", WechatRefundNotify)
	}
}


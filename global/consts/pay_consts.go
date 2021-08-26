package consts

const (
	TradeSuccess  = "TRADE_SUCCESS"
	TradeClosed   = "TRADE_CLOSED"
	TradeFinished = "TRADE_FINISHED"
	WaitBuyerPay  = "WAIT_BUYER_PAY"
)

// 1 支付回调
// 2 退款回调
// 3 退款申请
// 4 取消订单
const (
	PAY_NOTIFY    = 1
	REFUND_NOTIFY = 2
	APPLY_REFUND  = 3
	CANCEL_ORDER  = 4
)

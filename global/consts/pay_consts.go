package consts

const (
	TradeSuccess  = "TRADE_SUCCESS"
	TradeClosed   = "TRADE_CLOSED"
	TradeFinished = "TRADE_FINISHED"
	WaitBuyerPay  = "WAIT_BUYER_PAY"
)

// 1 支付回调
// 2 退款回调
const (
	PAY_NOTIFY    = 1
	REFUND_NOTIFY = 2
)

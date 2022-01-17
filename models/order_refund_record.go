package models

type OrderRefundRecord struct {
	Id            int    `json:"id" xorm:"not null pk autoincr comment('自增ID') INT(11)"`
	RefundAmount  int    `json:"refund_amount" xorm:"not null default 0 comment('退款金额') INT(10)"`
	RefundPayType int    `json:"refund_pay_type" xorm:"not null comment('退款类型 2001 支付宝  3001 微信') INT(8)"`
	RefundType    int    `json:"refund_type" xorm:"not null comment('退款形式：1001 APP退款 1002 小程序退款') INT(8)"`
	RefundAddress string `json:"refund_address" xorm:"not null default '' comment('退货地址') VARCHAR(255)"`
	OrderId       string `json:"order_id" xorm:"not null comment('订单id') VARCHAR(150)"`
	RefundTradeNo string `json:"refund_trade_no" xorm:"not null comment('退款交易号') VARCHAR(150)"`
	UserId        string `json:"user_id" xorm:"not null comment('用户id') VARCHAR(60)"`
	Remark        string `json:"remark" xorm:"not null default '' comment('备注') VARCHAR(255)"`
	Status        int    `json:"status" xorm:"not null default 0 comment('0 退款中 1 已退款') TINYINT(1)"`
	RefundTime    int    `json:"refund_time" xorm:"not null default 0 comment('退款时间') INT(11)"`
	CreateAt      int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
}

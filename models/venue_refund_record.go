package models

type VenueRefundRecord struct {
	Id              int    `json:"id" xorm:"not null pk autoincr comment('自增ID') INT(20)"`
	RefundAmount    int    `json:"refund_amount" xorm:"not null default 0 comment('退款金额') INT(11)"`
	RefundChannelId int    `json:"refund_channel_id" xorm:"not null comment('退款的渠道id 关联payment_channel表') INT(10)"`
	RefundType      int    `json:"refund_type" xorm:"not null comment('退款形式：1001 APP退款 1002 线下退款') INT(8)"`
	PayOrderId      string `json:"pay_order_id" xorm:"not null comment('订单号') VARCHAR(150)"`
	RefundTradeNo   string `json:"refund_trade_no" xorm:"not null comment('退款交易号') VARCHAR(150)"`
	UserId          string `json:"user_id" xorm:"not null comment('用户id') VARCHAR(60)"`
	Remark          string `json:"remark" xorm:"not null default '' comment('备注') VARCHAR(255)"`
	FeeRate         int    `json:"fee_rate" xorm:"not null default 0 comment('手续费比例 例如：1.55% 则 值为155 [乘以100存储]') INT(5)"`
	MinimumCharge   int    `json:"minimum_charge" xorm:"not null default 0 comment('最低收取手续费金额 单位[分]') INT(10)"`
	RefundFee       int    `json:"refund_fee" xorm:"not null default 0 comment('退款手续费') INT(11)"`
	Status          int    `json:"status" xorm:"not null default 0 comment('0 退款中 1 已退款') TINYINT(1)"`
	RefundTime      int    `json:"refund_time" xorm:"not null default 0 comment('退款时间') INT(11)"`
	CreateAt        int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt        int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}

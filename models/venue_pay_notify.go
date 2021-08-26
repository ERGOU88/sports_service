package models

type VenuePayNotify struct {
	Id          int64  `json:"id" xorm:"pk autoincr comment('自增主键') BIGINT(20)"`
	NotifyInfo  string `json:"notify_info" xorm:"not null default '' comment('回调通知信息') VARCHAR(1000)"`
	Transaction string `json:"transaction" xorm:"not null default '' comment('第三方订单号') VARCHAR(200)"`
	PayType     int    `json:"pay_type" xorm:"not null comment('1 支付宝 2 微信 ') TINYINT(1)"`
	PayOrderId  string `json:"pay_order_id" xorm:"not null comment('订单号') index VARCHAR(150)"`
	CreateAt    int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt    int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	NotifyType  int    `json:"notify_type" xorm:"not null default 1 comment('1 付款回调 2 退款回调') TINYINT(1)"`
}

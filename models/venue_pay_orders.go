package models

type VenuePayOrders struct {
	Id             int64  `json:"id" xorm:"pk autoincr comment('自增主键') BIGINT(20)"`
	PayOrderId     string `json:"pay_order_id" xorm:"not null comment('订单号') index VARCHAR(150)"`
	UserId         string `json:"user_id" xorm:"not null comment('用户id') index(user_id) VARCHAR(60)"`
	Amount         int    `json:"amount" xorm:"not null comment('商品总价[优惠后的金额]（分）') INT(11)"`
	Status         int    `json:"status" xorm:"not null default 0 comment(' 0 待支付 1 订单超时/未支付/已取消 2 已支付 3 已完成  4 退款中 5 已退款 6 已过期') index(user_id) TINYINT(4)"`
	Extra          string `json:"extra" xorm:"comment('记录订单相关扩展数据') TEXT"`
	Transaction    string `json:"transaction" xorm:"not null default '' comment('第三方订单号') VARCHAR(200)"`
	ProductType    int    `json:"product_type" xorm:"default 0 comment('-1组合购买 1001 场馆预约 3001 私教预约 3002 课程预约 2101 临时卡 2201 次卡 2311 月卡 2321 季卡 2331 半年卡 2341 年卡 4001 储值卡  5101 线下实体商品 5102 离场结算') INT(8)"`
	ErrorCode      string `json:"error_code" xorm:"not null default '' comment('错误码') VARCHAR(20)"`
	OrderType      int    `json:"order_type" xorm:"not null comment('下单方式：1001 APP下单，1002 前台购买，1003第三方推广渠道购买') index(user_id) INT(8)"`
	PayTime        int    `json:"pay_time" xorm:"not null default 0 comment('用户支付时间') INT(11)"`
	ChannelId      int    `json:"channel_id" xorm:"comment('购买渠道，1001 android ; 1002 ios') INT(10)"`
	IsCallback     int    `json:"is_callback" xorm:"not null default 0 comment('是否接收到第三方回调 0 未接收到回调 1 已接收回调') TINYINT(1)"`
	Subject        string `json:"subject" xorm:"not null default '' comment('商品名称') VARCHAR(150)"`
	RefundAmount   int    `json:"refund_amount" xorm:"not null default 0 comment('退款金额（分）') INT(11)"`
	IsDelete       int    `json:"is_delete" xorm:"not null default 0 comment('是否删除0正常 1删除') TINYINT(4)"`
	OriginalAmount int    `json:"original_amount" xorm:"not null comment('订单原始金额（分）') INT(11)"`
	RefundFee      int    `json:"refund_fee" xorm:"not null default 0 comment('退款手续费（分）') INT(11)"`
	PayChannelId   int    `json:"pay_channel_id" xorm:"not null default 0 comment('支付渠道id') INT(10)"`
	AdminId        int    `json:"admin_id" xorm:"default 0 comment('操作原') INT(11)"`
	VenueId        int64  `json:"venue_id" xorm:"not null default 0 comment('场馆id') BIGINT(20)"`
	CreateAt       int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt       int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	IsGift         int    `json:"is_gift" xorm:"not null default 0 comment('是否为赠品 0 不是 1 是') TINYINT(4)"`
	GiftStatus     int    `json:"gift_status" xorm:"not null default 0 comment('礼物赠送状态 0 未赠送 1 已过期 2 已赠送/已领取') TINYINT(4)"`
	ReceiveTm      int    `json:"receive_tm" xorm:"not null default 0 comment('礼物领取时间') INT(11)"`
}

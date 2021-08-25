package models

type VenuePayOrders struct {
	Id           int64  `json:"id" xorm:"pk autoincr comment('自增主键') BIGINT(20)"`
	UserId       string `json:"user_id" xorm:"not null comment('用户id') index(user_id) VARCHAR(60)"`
	Amount       int    `json:"amount" xorm:"not null comment('商品总价（分）') INT(11)"`
	Status       int    `json:"status" xorm:"not null default 0 comment(' 0 待支付 1 订单超时/未支付/已取消 2 已支付 3 已完成  4 退款中 5 已退款 6 已过期') index(user_id) TINYINT(4)"`
	Extra        string `json:"extra" xorm:"not null default '' comment('记录订单相关扩展数据') VARCHAR(1000)"`
	Transaction  string `json:"transaction" xorm:"not null default '' comment('第三方订单号') VARCHAR(200)"`
	PayType      int    `json:"pay_type" xorm:"not null comment('1 支付宝 2 微信 3 钱包 4 苹果内购') TINYINT(1)"`
	ProductType  int    `json:"product_type" xorm:"not null comment('1001 场馆预约 2001 购买月卡 2002 购买季卡 2003 购买年卡 2004 体验券 3001 私教（教练）订单 3002 课程订单 4001 充值订单') INT(8)"`
	ErrorCode    string `json:"error_code" xorm:"not null default '' comment('错误码') VARCHAR(20)"`
	PayOrderId   string `json:"pay_order_id" xorm:"not null comment('订单号') index VARCHAR(150)"`
	OrderType    int    `json:"order_type" xorm:"not null comment('下单方式：1001 APP下单，1002 前台购买，1003第三方推广渠道购买') index(user_id) INT(8)"`
	PayTime      int    `json:"pay_time" xorm:"not null default 0 comment('用户支付时间') INT(11)"`
	ChannelId    int    `json:"channel_id" xorm:"comment('购买渠道，1001 android ; 1002 ios') INT(10)"`
	IsCallback   int    `json:"is_callback" xorm:"not null default 0 comment('是否接收到第三方回调 0 未接收到回调 1 已接收回调') TINYINT(1)"`
	CreateAt     int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt     int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	Subject      string `json:"subject" xorm:"not null default '' comment('商品名称') VARCHAR(150)"`
	WriteOffCode string `json:"write_off_code" xorm:"not null default '' comment('核销码') VARCHAR(200)"`
	RefundAmount int    `json:"refund_amount" xorm:"not null default 0 comment('退款金额（分）') INT(11)"`
	IsDelete     int    `json:"is_delete" xorm:"not null default 0 comment('是否删除0正常 1删除') TINYINT(2)"`
}

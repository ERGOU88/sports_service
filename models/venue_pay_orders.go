package models

type VenuePayOrders struct {
	Id          int64  `json:"id" xorm:"pk autoincr comment('自增主键') BIGINT(20)"`
	UserId      string `json:"user_id" xorm:"not null comment('用户id') index(user_id) VARCHAR(60)"`
	Amount      int    `json:"amount" xorm:"not null comment('商品总价（分）') INT(11)"`
	Status      int    `json:"status" xorm:"not null default 0 comment('0 待支付 1 订单超时/未支付 2 已支付 4 已完成 5 已取消  6  退款中  7  已退款 8 软删除') index(user_id) TINYINT(4)"`
	Extra       string `json:"extra" xorm:"not null default '' comment('记录订单相关扩展数据') VARCHAR(1000)"`
	Transaction string `json:"transaction" xorm:"not null default '' comment('第三方订单号') VARCHAR(200)"`
	PayType     int    `json:"pay_type" xorm:"not null comment('1 支付宝 2 微信 3 钱包 4 苹果内购') TINYINT(1)"`
	ErrorCode   string `json:"error_code" xorm:"not null default '' comment('错误码') VARCHAR(20)"`
	PayOrderId  string `json:"pay_order_id" xorm:"not null comment('订单号') index VARCHAR(150)"`
	IsOnline    int    `json:"is_online" xorm:"not null default 0 comment('0线上购买 1线下购买') TINYINT(1)"`
	OrderType   int    `json:"order_type" xorm:"not null comment('1 场馆预约 2 购买月卡 3 购买季卡 4 购买年卡 5 私教（教练）订单 6 课程订单 7 充值订单 8 体验券') index(user_id) TINYINT(2)"`
	PayTime     int    `json:"pay_time" xorm:"not null default 0 comment('用户支付时间') INT(11)"`
	IsCallback  int    `json:"is_callback" xorm:"not null default 0 comment('是否接收到第三方回调 0 未接收到回调 1 已接收回调') TINYINT(1)"`
	CreateAt    int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt    int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}

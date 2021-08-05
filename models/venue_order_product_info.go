package models

type VenueOrderProductInfo struct {
	Id             int64  `json:"id" xorm:"pk autoincr comment('自增ID') BIGINT(20)"`
	PayOrderId     string `json:"pay_order_id" xorm:"not null comment('订单号') index VARCHAR(150)"`
	ProductId      string `json:"product_id" xorm:"not null comment('商品id') VARCHAR(60)"`
	OrderType      int    `json:"order_type" xorm:"not null comment('1 场馆预约 2 购买月卡 3 购买季卡 4 购买年卡 5 体验券 6 私教（教练）订单 7 课程订单 8 充值订单 ') TINYINT(2)"`
	Count          int    `json:"count" xorm:"not null comment('购买数量') INT(11)"`
	RealAmount     int    `json:"real_amount" xorm:"not null comment('真实价格（单位：分）') INT(11)"`
	CurAmount      int    `json:"cur_amount" xorm:"not null comment('当前价格 (包含真实价格、 折扣价格（单位：分）') INT(11)"`
	DiscountRate   int    `json:"discount_rate" xorm:"not null default 0 comment('折扣率') INT(11)"`
	DiscountAmount int    `json:"discount_amount" xorm:"not null default 0 comment('优惠的金额') INT(11)"`
	Amount         int    `json:"amount" xorm:"not null comment('商品总价') INT(11)"`
	ReceiveAmount  int    `json:"receive_amount" xorm:"not null default 0 comment('充值金额（钱包）') INT(11)"`
	Duration       int    `json:"duration" xorm:"not null default 0 comment('购买相关服务时长') INT(11)"`
	CreateAt       int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt       int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}

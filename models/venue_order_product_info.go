package models

type VenueOrderProductInfo struct {
	Id              int64  `json:"id" xorm:"pk autoincr comment('自增ID') BIGINT(20)"`
	PayOrderId      string `json:"pay_order_id" xorm:"not null comment('订单号') index VARCHAR(150)"`
	ProductId       int64  `json:"product_id" xorm:"not null comment('商品id') BIGINT(20)"`
	ProductType     int    `json:"product_type" xorm:"not null comment('1001 场馆预约 2101 临时卡 2201 次卡 2311 月卡 2321 季卡 2331 半年卡 2341 年卡 3001 私教（教练）订单 3002 课程订单 4001 充值订单  5101 线下实体商品 5102线下结算') INT(8)"`
	Count           int    `json:"count" xorm:"not null comment('购买数量') INT(11)"`
	RealAmount      int    `json:"real_amount" xorm:"not null comment('[单个商品]定价（单位：分）') INT(11)"`
	CurAmount       int    `json:"cur_amount" xorm:"not null comment('[单个商品]当前价格 [售价](包含真实价格、 折扣价格（单位：分）') INT(11)"`
	DiscountRate    int    `json:"discount_rate" xorm:"not null default 0 comment('折扣率') INT(11)"`
	DiscountAmount  int    `json:"discount_amount" xorm:"not null default 0 comment('[单个商品]优惠的金额') INT(11)"`
	Amount          int    `json:"amount" xorm:"not null comment('商品总价') INT(11)"`
	CreateAt        int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt        int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	Status          int    `json:"status" xorm:"not null default 0 comment('0 待支付 1 订单超时/未支付/已取消 2 已支付 3 已完成  4 退款中 5 已退款 6 已过期') TINYINT(2)"`
	ProductCategory int    `json:"product_category" xorm:"not null default 0 comment('商品类别 1000 预约类 2000 卡类 5000 实物类') INT(8)"`
	VenueId         int64  `json:"venue_id" xorm:"not null default 0 comment('场馆id') BIGINT(20)"`
	IsWriteOff      int    `json:"is_write_off" xorm:"not null default 0 comment('是否核销 0 未核销 1 已核销') TINYINT(1)"`
	SnapshotId      int64  `json:"snapshot_id" xorm:"not null default 0 comment('商品快照id') BIGINT(20)"`
}

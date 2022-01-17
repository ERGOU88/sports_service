package models

type OrderProduct struct {
	Id             int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	OrderId        string `json:"order_id" xorm:"not null default '' comment('订单id') index VARCHAR(50)"`
	UserId         string `json:"user_id" xorm:"not null default '' comment('买家id') index VARCHAR(60)"`
	ProductId      int    `json:"product_id" xorm:"not null default 0 comment('商品id') index INT(11)"`
	SkuId          int    `json:"sku_id" xorm:"not null default 0 comment('商品skuid') index INT(11)"`
	SkuName        string `json:"sku_name" xorm:"not null default '' comment('商品实体名称') VARCHAR(255)"`
	SkuImage       string `json:"sku_image" xorm:"not null default '' comment('商品实体图片') VARCHAR(255)"`
	SkuNo          string `json:"sku_no" xorm:"not null default '' comment('商品实体编码') VARCHAR(255)"`
	ProductAmount  int    `json:"product_amount" xorm:"not null default 0 comment('商品总金额 (分) ') INT(10)"`
	DeliveryAmount int    `json:"delivery_amount" xorm:"not null default 0 comment('配送费用（分）') INT(10)"`
	OrderAmount    int    `json:"order_amount" xorm:"not null default 0 comment('合计金额 (分)') INT(10)"`
	DiscountAmount int    `json:"discount_amount" xorm:"not null default 0 comment('优惠金额 (分)') INT(10)"`
	PayAmount      int    `json:"pay_amount" xorm:"not null default 0 comment('应付金额 (分)') INT(10)"`
	Count          int    `json:"count" xorm:"not null default 0 comment('购买数量') INT(10)"`
	DeliveryNo     string `json:"delivery_no" xorm:"not null default '' comment('配送单号') VARCHAR(50)"`
	ProductName    string `json:"product_name" xorm:"not null default '' comment('商品名称') VARCHAR(400)"`
	SkuSpec        string `json:"sku_spec" xorm:"not null default '' comment('sku规格格式') VARCHAR(1000)"`
}

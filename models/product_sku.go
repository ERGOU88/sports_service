package models

type ProductSku struct {
	Id            int    `json:"id" xorm:"not null pk autoincr comment('sku id') INT(11)"`
	ProductId     int64  `json:"product_id" xorm:"not null comment('商品id') index BIGINT(20)"`
	Title         string `json:"title" xorm:"not null comment('商品标题') VARCHAR(255)"`
	SkuImage      string `json:"sku_image" xorm:"not null default '' comment('sku主图') VARCHAR(255)"`
	SkuNo         string `json:"sku_no" xorm:"not null default '' comment('商品sku编码') VARCHAR(255)"`
	Images        string `json:"images" xorm:"default '' comment('sku图片[多张]') VARCHAR(1000)"`
	CurPrice      int    `json:"cur_price" xorm:"not null default 0 comment('商品价格（分）') INT(10)"`
	MarketPrice   int    `json:"market_price" xorm:"not null default 0 comment('划线价格（分）') INT(10)"`
	IsFreeShip    int    `json:"is_free_ship" xorm:"not null default 0 comment('是否免邮 0 免邮') TINYINT(2)"`
	Indexes       string `json:"indexes" xorm:"comment('特有规格属性在商品属性模板中的对应下标组合 例如0_1_0 则可能表示 颜色选项 0 远峰蓝 内存选项 1 3GB  机身存储 0  16GB') index VARCHAR(100)"`
	OwnSpec       string `json:"own_spec" xorm:"comment('商品实体的特有规格参数，json格式，反序列化时保证有序') VARCHAR(1000)"`
	Status        int    `json:"status" xorm:"not null default 0 comment('是否有效，0 有效, 1 无效') TINYINT(1)"`
	DiscountPrice int    `json:"discount_price" xorm:"not null default 0 comment('活动折扣价（默认等于单价）单位：分') INT(10)"`
	StartTime     int    `json:"start_time" xorm:"not null default 0 comment('活动开始时间') INT(11)"`
	EndTime       int    `json:"end_time" xorm:"not null default 0 comment('活动结束时间') INT(11)"`
	CreateAt      int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt      int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
	Sortorder     int    `json:"sortorder" xorm:"not null default 0 comment('排序权重') INT(11)"`
}

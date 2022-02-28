package models

type Products struct {
	Id             int    `json:"id" xorm:"not null pk autoincr comment('商品id') INT(11)"`
	ProductName    string `json:"product_name" xorm:"not null default '' comment('商品名称') VARCHAR(255)"`
	ProductImage   string `json:"product_image" xorm:"not null default '' comment('商品主图路径') VARCHAR(512)"`
	ProductDetail  string `json:"product_detail" xorm:"comment('商品详情') TEXT"`
	Status         int    `json:"status" xorm:"not null default 1 comment('商品状态（0. 正常 1. 下架）') TINYINT(2)"`
	IsFreeShip     int    `json:"is_free_ship" xorm:"not null default 0 comment('是否免邮 0 免邮') TINYINT(2)"`
	IsDelete       int    `json:"is_delete" xorm:"not null default 0 comment('是否已经删除') index TINYINT(4)"`
	Introduction   string `json:"introduction" xorm:"not null default '' comment('促销语') VARCHAR(255)"`
	Keywords       string `json:"keywords" xorm:"not null default '' comment('关键词') VARCHAR(255)"`
	Sortorder      int    `json:"sortorder" xorm:"not null default 0 comment('排序') index INT(11)"`
	VideoUrl       string `json:"video_url" xorm:"not null default '' comment('视频地址') VARCHAR(512)"`
	SaleNum        int    `json:"sale_num" xorm:"not null default 0 comment('销量') INT(11)"`
	CurPrice       int    `json:"cur_price" xorm:"not null default 0 comment('商品价格（分）') INT(10)"`
	MarketPrice    int    `json:"market_price" xorm:"not null default 0 comment('划线价格（分）') INT(10)"`
	EvaluateNum    int    `json:"evaluate_num" xorm:"not null default 0 comment('总评价数') INT(11)"`
	IsRecommend    int    `json:"is_recommend" xorm:"not null default 0 comment('推荐（0：不推荐；1：推荐）') TINYINT(1)"`
	IsTop          int    `json:"is_top" xorm:"not null default 0 comment('置顶（0：不置顶；1：置顶；）') TINYINT(1)"`
	IsCream        int    `json:"is_cream" xorm:"not null default 0 comment('是否精华（0: 不是 1: 是）') TINYINT(1)"`
	Specifications string `json:"specifications" xorm:"not null default '' comment('全部规格参数数据') VARCHAR(3000)"`
	SpecTemplate   string `json:"spec_template" xorm:"not null default '' comment('特有规格参数及可选值信息，json格式') VARCHAR(1000)"`
	AfterService   string `json:"after_service" xorm:"default '' comment('售后服务') VARCHAR(1000)"`
	TimerOn        int    `json:"timer_on" xorm:"not null default 0 comment('上架时间') INT(11)"`
	TimerOff       int    `json:"timer_off" xorm:"not null default 0 comment('下架时间') INT(11)"`
	CreateAt       int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt       int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
}

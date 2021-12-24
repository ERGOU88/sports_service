package mshop

import "sports_service/server/models"

type PlaceOrderReq struct {
	UserId         string            `json:"user_id"`
	ClientIp       string            `json:"client_ip"`
	UserAddrId     int               `json:"user_addr_id"`                        // 选择的地址
	Products       []*ProductParam   `json:"products" binding:"required"`
	ReqType        int               `json:"req_type" binding:"required"`         // 1 查询数据 2 详情页下单 3 购物车下单
}

type ProductParam struct {
	SkuId            int             `json:"sku_id" binding:"required"`
	ProductId        int             `json:"product_id" binding:"required"`
	Count            int             `json:"count" binding:"required"`
	CartId           int             `json:"cart_id"`
}

// 下单返回数据
type OrderResp struct {
	ClientIp         string          `json:"client_ip"`              // ip地址
	UserId           string          `json:"user_id"`
	OrderId          string          `json:"order_id"`               // 订单id
	IsEnough         bool            `json:"is_enough"`              // 库存标识 是否足够 false 库存不足
	Total            int             `json:"total"`                  // 总件数
	PayAmount        int             `json:"pay_amount"`             // 应付金额
	DiscountAmount   int             `json:"discount_amount"`        // 优惠金额
	DeliveryAmount   int             `json:"delivery_amount"`        // 配送费用
	OrderAmount      int             `json:"order_amount"`           // 合计金额
	ProductAmount    int             `json:"product_amount"`         // 商品总金额
	Products         []*Product      `json:"products"`               // 商品sku列表
	CreateTm         string          `json:"create_tm"`              // 订单创建时间
	CreateAt         int             `json:"create_at"`
	PayDuration      int             `json:"pay_duration"`           // 支付时长
	UserAddr         *models.UserAddress  `json:"user_addr"`         // 用户收获地址
}

type Product struct {
	UserId         string `json:"user_id"`
	OrderId        string `json:"order_id"`
	SkuId          int    `json:"sku_id"`
	ProductId      int    `json:"product_id"`
	Count          int    `json:"count"`
	IsEnough       bool   `json:"is_enough" xorm:"-"`
	CartId         int    `json:"cart_id" xorm:"-"`
	SkuImage       string `json:"sku_image"`
	SkuName        string `json:"sku_name"`
	ProductName    string `json:"product_name"`
	SkuNo          string `json:"sku_no"`
	CurPrice       int    `json:"cur_price" xorm:"-"`
	MarketPrice    int    `json:"market_price" xorm:"-"`
	IsFreeShip     int    `json:"is_free_ship" xorm:"-"`
	DiscountPrice  int    `json:"discount_price" xorm:"-"`
	StartTime      int    `json:"start_time" xorm:"-"`
	EndTime        int    `json:"end_time" xorm:"-"`
	RemainDuration int    `json:"remain_duration" xorm:"-"` // 活动剩余时长
	HasActivities  int32  `json:"has_activities" xorm:"-"`  // 1 有活动
	SkuSpec          []OwnSpec       `json:"own_spec"`                 // 商品实体的特有规格参数
	Indexes          string          `json:"indexes" xorm:"-"`         // 特有规格属性在商品属性模板中的对应下标组合
	Stock            int             `json:"stock" xorm:"-"`           // 库存
	MaxBuy           int             `json:"max_buy" xorm:"-"`         // 限购 0 表示无限制
	MinBuy           int             `json:"min_buy" xorm:"-"`         // 起购数
	CanBuy           bool            `json:"can_buy" xorm:"-"`         // 当前数量是否可购买
	
	PayAmount        int             `json:"pay_amount"`             // 应付金额
	DiscountAmount   int             `json:"discount_amount"`        // 优惠金额
	DeliveryAmount   int             `json:"delivery_amount"`        // 配送费用
	OrderAmount      int             `json:"order_amount"`           // 合计金额
	ProductAmount    int             `json:"product_amount"`         // 商品总金额
}

func (m *ShopModel) AddOrder(order *models.Orders) (int64, error) {
	return m.Engine.InsertOne(order)
}

func (m *ShopModel) AddOrderProduct(list []*Product) (int64, error) {
	return m.Engine.Table("order_product").InsertMulti(list)
}

func (m *ShopModel) AddBuyerDeliveryInfo(info *models.BuyerDeliveryInfo) (int64, error) {
	return m.Engine.InsertOne(info)
}

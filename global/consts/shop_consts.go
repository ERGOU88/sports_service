package consts

const (
	// 推荐的商品 默认最多取三条
	RECOMMEND_DEFAULT_LIMIT = 3
)

// 0 待支付
// 1 订单超时/未支付
// 2 已支付
const (
	SHOP_ORDER_TYPE_WAIT = iota
	SHOP_ORDER_TYPE_UNPAID
	SHOP_ORDER_TYPE_PAID
)

const (
	// 商城订单 可支付时长 24小时
	SHOP_PAYMENT_DURATION = 60 * 60 * 24
)

// 1 商品详情页下单
// 2 商品购物车下单
const (
	ORDER_ACTION_TYPE_DETAIL = 1
	ORDER_ACTION_TYPE_CART   = 2
)

// 配送状态： 0 未配送 1 已配送 2 已签收
const (
	NOT_DELIVERED = iota
	HAS_DELIVERED
	HAS_SIGNED
)

// TODO: 暂时写死
const (
	DELIVERY_NAME      = "顺丰速运"
	DELIVERY_TELEPHONE = "95338"
)

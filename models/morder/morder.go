package morder

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/rdskey"
	"sports_service/server/models"
)

// 支付请求参数
type PayReqParam struct {
	PayType   int     `binding:"required" json:"pay_type"`     // 1 支付宝 2 微信 3 钱包 4 苹果内购
	OrderId   string  `binding:"required" json:"order_id"`     // 订单id
	UserId    string
}

// 订单信息
type OrderInfo struct {
	CreatAt            string      `json:"creat_at"`            // 订单创建时间
	OrderType          int32       `json:"order_type"`          // 订单商品类型 1001 场馆预约 2001 购买月卡 2002 购买季卡 2003 购买年卡 2004 体验券 3001 私教（教练）订单 3002 课程订单 4001 充值订单
	OrderStatus        int32       `json:"order_status"`        // 订单状态 0 待支付 1 订单超时/未支付/已取消 2 已支付 3 已完成  4 退款中 5 已退款 6 退款失败
	Title              string      `json:"title"`               // 标题
	Amount             string      `json:"amount"`              // 金额
	TotalAmount        int         `json:"total_amount"`        // 总金额
	Duration           int64       `json:"duration"`            // 剩余支付时长
	UserId             string      `json:"user_id"`
	OrderId            string      `json:"order_id"`            // 订单id
	Count              int         `json:"count"`
	HasEvaluate        bool        `json:"has_evaluate"`        // 是否评价
}

// 订单退款/删除订单/取消订单
type ChangeOrder struct {
	OrderId    string  `binding:"required" json:"order_id"`     // 订单id
	UserId     string  `json:"user_id"`
}

// 券码信息
type CouponCodeInfo struct {
	Code        string `json:"code"`
	VenueName   string `json:"venue_name"`
	Subject     string `json:"subject"`
	TotalAmount int    `json:"total_amount"`
	Count       int    `json:"count"`
	ExpireTm    string `json:"expire_tm"`
	QrCodeInfo  string `json:"qr_code_info"`
	QrCodeExpireDuration int64 `json:"qr_code_expire_duration"`
}

type OrderModel struct {
	Engine         *xorm.Session
	Order          *models.VenuePayOrders
	OrderProduct   *models.VenueOrderProductInfo
	Record         *models.VenueAppointmentRecord
	Notify         *models.VenuePayNotify
}

func NewOrderModel(engine *xorm.Session) *OrderModel {
	return &OrderModel{
		Engine: engine,
		Order: new(models.VenuePayOrders),
		OrderProduct: new(models.VenueOrderProductInfo),
		Record: new(models.VenueAppointmentRecord),
		Notify: new(models.VenuePayNotify),
	}
}

// 添加订单
func (m *OrderModel) AddOrder() (int64, error) {
	return m.Engine.InsertOne(m.Order)
}

// 获取订单
func (m *OrderModel) GetOrder(orderId string) (bool, error) {
	m.Order = new(models.VenuePayOrders)
	return m.Engine.Where("pay_order_id=?", orderId).Get(m.Order)
}

// 查看订单商品流水表 获取商品销量
func (m *OrderModel) GetSalesByProduct() (int64, error) {
	return m.Engine.Where("product_id=? AND product_type=? AND status=2", m.OrderProduct.ProductId,
		m.OrderProduct.ProductType).SumInt(m.OrderProduct, "count")
}

// 批量添加订单商品流水
func (m *OrderModel) AddMultiOrderProduct(list []*models.VenueOrderProductInfo) (int64, error) {
	return m.Engine.InsertMulti(list)
}

// 添加订单商品流水
func (m *OrderModel) AddOrderProduct() (int64, error) {
	return m.Engine.InsertOne(m.OrderProduct)
}

// 更新订单信息
func (m *OrderModel) UpdateOrderStatus(orderId string, status int) (int64, error) {
	return m.Engine.Where("pay_order_id=? AND status=?", orderId, status).Cols("update_at",
		"status", "is_callback", "pay_time", "transaction", "refund_amount", "refund_fee").Update(m.Order)
}

// 通过订单id 获取订单流水信息
func (m *OrderModel) GetOrderProductsById(orderId string) (bool, error) {
	m.OrderProduct = new(models.VenueOrderProductInfo)
	return m.Engine.Where("pay_order_id=?", orderId).Get(m.OrderProduct)
}

// 更新订单商品状态
func (m *OrderModel) UpdateOrderProductStatus(orderId string, status int) (int64, error) {
	return m.Engine.Where("pay_order_id=? AND status=?", orderId, status).Cols("update_at", "status").Update(m.OrderProduct)
}

// 记录需处理超时的订单号
func (m *OrderModel) RecordOrderId(orderId string) (int, error) {
	rds := dao.NewRedisDao()
	return rds.SADD(rdskey.ORDER_EXPIRE_INFO, orderId)
}

// 记录订单回调通知
func (m *OrderModel) AddOrderPayNotify() (int64, error) {
	return m.Engine.InsertOne(m.Notify)
}

// 更新订单信息
func (m *OrderModel) UpdateOrderInfo(cols string) (int64, error) {
	return m.Engine.Where("pay_order_id=?", m.Order.PayOrderId).Cols(cols).Update(m.Order)
}

// 通过订单状态获取订单列表
// 订单状态：
// 0 待支付
// 1 订单超时/未支付
// 2 已支付
// ......
func (m *OrderModel) GetOrderListByStatus(condition string, offset, size int) ([]*models.VenuePayOrders, error) {
	var list []*models.VenuePayOrders
	if err := m.Engine.Where(condition).Desc("id").Limit(size, offset).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 按执行顺序 获取退款规则
func (m *OrderModel) GetRefundRules() ([]*models.VenueRefundRules, error) {
	var list []*models.VenueRefundRules
	if err := m.Engine.Where("status=0").Asc("rule_order").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

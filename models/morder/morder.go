package morder

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

// 支付请求参数
type PayReqParam struct {
	PayType   int     `binding:"required" json:"pay_type"`     // 1 支付宝 2 微信 3 钱包 4 苹果内购
	OrderId   string  `binding:"required" json:"order_id"`     // 订单id
	UserId    string
}

type OrderModel struct {
	Engine         *xorm.Session
	Order          *models.VenuePayOrders
	OrderProduct   *models.VenueOrderProductInfo
	Record         *models.AppointmentRecord
}

func NewOrderModel(engine *xorm.Session) *OrderModel {
	return &OrderModel{
		Engine: engine,
		Order: new(models.VenuePayOrders),
		OrderProduct: new(models.VenueOrderProductInfo),
		Record: new(models.AppointmentRecord),
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
	return m.Engine.Where("product_id=? AND order_type=? AND status=2", m.OrderProduct.ProductId,
		m.OrderProduct.OrderType).SumInt(m.OrderProduct, "count")
}

// 批量添加订单商品流水
func (m *OrderModel) AddMultiOrderProduct(list []*models.VenueOrderProductInfo) (int64, error) {
	return m.Engine.InsertMulti(list)
}

// 订单超时 更新订单状态
func (m *OrderModel) UpdateOrderStatus(orderId string, status int) (int64, error) {
	return m.Engine.Where("pay_order_id=? AND status=?", orderId, status).Cols("update_at", "status").Update(m.Order)
}

// 通过订单id 获取订单流水信息
func (m *OrderModel) GetOrderProductsById(orderId string, status int) (bool, error) {
	m.OrderProduct = new(models.VenueOrderProductInfo)
	return m.Engine.Where("pay_order_id=? AND status=?", orderId, status).Get(m.OrderProduct)
}

// 更新订单商品状态
func (m *OrderModel) UpdateOrderProductStatus(orderId string, status int) (int64, error) {
	return m.Engine.Where("pay_order_id=? AND status=?", orderId, status).Cols("update_at", "status").Update(m.OrderProduct)
}

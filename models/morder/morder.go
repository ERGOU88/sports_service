package morder

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

type OrderModel struct {
	Engine         *xorm.Session
	Order          *models.VenuePayOrders
	OrderProduct   *models.VenueOrderProductInfo
}

func NewOrderModel(engine *xorm.Session) *OrderModel {
	return &OrderModel{
		Engine: engine,
		Order: new(models.VenuePayOrders),
		OrderProduct: new(models.VenueOrderProductInfo),
	}
}

// 添加订单
func (m *OrderModel) AddOrder() (int64, error) {
	return m.Engine.InsertOne(m.Order)
}

// 查看订单商品流水表 获取商品销量
func (m *OrderModel) GetSalesByProduct() (int64, error) {
	return m.Engine.Where("product_id=? AND order_type=? AND status=2", m.OrderProduct.ProductId,
		m.OrderProduct.OrderType).SumInt(m.OrderProduct, "count")
}

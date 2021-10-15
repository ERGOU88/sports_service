package morder

import (
	"sports_service/server/models"
)

type RefundRateParam struct {
	Id       int     `json:"id"`
	Rate     int     `json:"rate"`
}

type OrderRecord struct {
	Id              int64     `json:"id"`
	VenueName       string    `json:"venue_name"`
	MobileNum       string    `json:"mobile_num"`
	PayOrderId      string    `json:"pay_order_id"`
	CreateAt        string    `json:"create_at"`
	ProductType     string    `json:"product_type"`
	OriginalAmount  string    `json:"original_amount"`
	Amount          string    `json:"amount"`
	Detail          string    `json:"detail"`                // 订单详情 例如：月卡 * 3
	Status          string    `json:"status"`
	PayChannel      string    `json:"pay_channel"`
}

// 获取场馆销售总额 已付款的
func (m *OrderModel) GetTotalSalesByVenue(venueId string) (int64, error) {
	return m.Engine.Where("venue_id=?", venueId).In("status", []int{2, 3, 4, 5}).SumInt(m.Order, "amount")
}

// 获取场馆订单数量（已付款的订单）
func (m *OrderModel) GetOrderCountByVenue(venueId string) (int64, error) {
	return m.Engine.Where("venue_id=?", venueId).In("status", []int{2, 3, 4, 5}).Count(&models.VenuePayOrders{})
}

// 获取场馆退款总额 退款中/已退款
func (m *OrderModel) GetTotalRefundByVenue(venueId string) (int64, error) {
	return m.Engine.Where("venue_id=?", venueId).In("status", []int{4, 5}).SumInt(m.Order, "refund_amount")
}

// 更新退款费率
func (m *OrderModel) UpdateRefundRate(id, rate int) (int64, error) {
	rules := new(models.VenueRefundRules)
	rules.FeeRate = rate
	return m.Engine.Where("id=?", id).Update(rules)
}

// 获取订单数量（所有场馆 已付款的订单）
func (m *OrderModel) GetOrderCount() (int64, error) {
	return m.Engine.In("status", []int{2, 3, 4, 5}).Count(&models.VenuePayOrders{})
}

// 获取订单列表
func (m *OrderModel) GetOrderList(offset, size int) ([]*models.VenuePayOrders, error) {
	var list []*models.VenuePayOrders
	if err := m.Engine.Limit(size, offset).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

package morder

import (
	"sports_service/server/models"
	"fmt"
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
	Status          int       `json:"status"`
	StatusCn        string    `json:"status_cn"`
	PayChannel      string    `json:"pay_channel"`
}

// 财务模块 订单统计数据
type OrderStat struct {
	TopInfo        map[string]interface{}      `json:"top_info"`         // 顶部统计数据
}


// 获取场馆销售总额 已付款的
func (m *OrderModel) GetTotalSalesByVenue(venueId string) (int64, error) {
	return m.Engine.Where("venue_id=?", venueId).In("status", []int{2, 3, 4, 5, 6}).SumInt(m.Order, "amount")
}

// 获取场馆订单数量（已付款的订单）
func (m *OrderModel) GetOrderCountByVenue(venueId string) (int64, error) {
	return m.Engine.Where("venue_id=?", venueId).In("status", []int{2, 3, 4, 5, 6}).Count(&models.VenuePayOrders{})
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
	if m.Order.UserId != "" {
		m.Engine.Where("user_id=?", m.Order.UserId)
	}
	return m.Engine.In("status", []int{2, 3, 4, 5, 6}).Count(&models.VenuePayOrders{})
}

// 获取订单列表
func (m *OrderModel) GetOrderList(offset, size int) ([]*models.VenuePayOrders, error) {
	var list []*models.VenuePayOrders
	if err := m.Engine.Desc("id").Limit(size, offset).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}


type RefundInfo struct {
	Id                 int64      `json:"id"`
	RefundChannelId    int        `json:"refund_channel_id"`
	Remark             string     `json:"remark"`
	CreateAt           int64      `json:"create_at"`
	OrderType          int        `json:"order_type"`
	RefundAmount       int        `json:"refund_amount"`
	RefundFee          int        `json:"refund_fee"`
	Status             int32      `json:"status"`
	Extra              string     `json:"extra"`
	PayOrderId         string     `json:"pay_order_id"`
	ProductType        int        `json:"product_type"`
	Amount             int        `json:"amount"`
	Detail             string     `json:"detail"`
	AmountCn           string     `json:"amount_cn"`
	RefundAmountCn     string     `json:"refund_amount_cn"`
	RefundFeeCn        string     `json:"refund_fee_cn"`
	OrderTypeCn        string     `json:"order_type_cn"`
	RefundChannelCn    string     `json:"refund_channel_cn"`
	CreateAtCn         string     `json:"create_at_cn"`
	MobileNum          string     `json:"mobile_num"`
	VenueName          string     `json:"venue_name"`
}

// 获取退款列表
func (m *OrderModel) GetRefundRecordList(orderId string, offset, size int) ([]*RefundInfo, error) {
	sql :=  "SELECT vrc.id, vrc.refund_channel_id,vrc.remark,vrc.create_at,o.order_type, o.refund_amount, " +
		"o.refund_fee, o.status, o.extra,o.pay_order_id,o.product_type,o.amount FROM venue_pay_orders AS o " +
		"LEFT JOIN venue_refund_record AS vrc ON vrc.pay_order_id = o.pay_order_id WHERE o.refund_amount >0 "

	if orderId != "" {
		sql += fmt.Sprintf("AND o.pay_order_id=%s ", orderId)
	}

	sql += "ORDER BY id DESC LIMIT ?,?"
	var list []*RefundInfo
	if err := m.Engine.SQL(sql, offset, size).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取订单收益流水[已付款/已退款]
func (m *OrderModel) GetRevenueFlow(minDate, maxDate, orderId string, offset, size int) ([]*models.VenuePayOrders, error) {
	var list []*models.VenuePayOrders
	m.Engine.Where("date(from_unixtime(create_at)) >= ? AND date(from_unixtime(create_at)) <= ?",
		minDate, maxDate).In("status", []int{2, 3, 4, 5, 6}).Limit(size, offset)
	if orderId != "" {
		m.Engine.Where("pay_order_id=?", orderId)
	}

	if err := m.Engine.Desc("id").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取总收入
func (m *OrderModel) GetTotalRevenue(minDate, maxDate string) (int64, error) {
	if minDate != "" && maxDate != "" {
		m.Engine.Where("date(from_unixtime(create_at)) >= ? AND date(from_unixtime(create_at)) <= ?",
			minDate, maxDate)
	}

	return m.Engine.In("status", []int{2, 3, 4, 5, 6}).SumInt(m.Order, "amount")
}

// 获取退款总金额
func (m *OrderModel) GetTotalRefund(minDate, maxDate string) (int64, error) {
	return m.Engine.Where("date(from_unixtime(create_at)) >= ? AND date(from_unixtime(create_at)) <= ?",
		minDate, maxDate).In("status", []int{4, 5}).SumInt(m.Order, "refund_amount")
}

// 通过日期新增会员用户数 / 总会员数
func (m *OrderModel) GetVipUserCount(minDate, maxDate string) (int64, error) {
	if minDate != "" && maxDate != "" {
		m.Engine.Where("date(from_unixtime(create_at))>=? AND date(from_unixtime(create_at))<=?", minDate, maxDate)
	}

	return m.Engine.Count(&models.VenueVipInfo{})
}

// 通过日期获取成功订单数[成功支付]
func (m *OrderModel) GetOrderNum(minDate, maxDate string) (int64, error) {
	return m.Engine.Where("date(from_unixtime(create_at)) >= ? AND date(from_unixtime(create_at)) <=?", minDate, maxDate).In("status", []int{2,3,4,5,6}).Count(m.Order)
}

// 通过日期获取场馆新增用户
func (m *OrderModel) GetDailyNewUsers() {

}

const (
	GET_VENUE_TOTAL_USER = "SELECT count(distinct o.user_id) AS count FROM venue_order_product_info AS vop " +
		"LEFT JOIN venue_pay_orders AS o ON vop.pay_order_id=o.pay_order_id WHERE o.status in(2,3,4,5,6) " +
		"AND vop.product_category in(1000,2000)"
)
// 获取所有场馆总用户数[会员/课程/私教/次卡]
func (m *OrderModel) GetVenueTotalUser() int64 {
	type stat struct {
		Count   int64
	}

	tmp := stat{}
	ok, err := m.Engine.SQL(GET_VENUE_TOTAL_USER).Get(&tmp)
	if !ok || err != nil {
		return 0
	}

	return tmp.Count
}

type SalesDetail struct {
	Count          int64    `json:"count"`
	Avg            float64  `json:"avg"`
	ProductType    int      `json:"product_type"`
	ProductName    string   `json:"product_name"`
	Dt             string   `json:"dt"`
	TotalSales     int      `json:"total_sales"`
	Alipay         int      `json:"alipay"`
	Wxpay          int      `json:"wxpay"`
	Cash           int      `json:"cash"`
	RefundAmount   int      `json:"refund_amount"`
	RefundCount    int      `json:"refund_count"`
}

type ResultList struct {
	Title string                 `json:"title"`
	List  map[string]interface{} `json:"list"`
}

type Result struct {
	Stat  float64          `json:"stat"`
	Title string           `json:"title"`
	List  []ResultList     `json:"list"`
}

// 获取销售明细
func (m *OrderModel) GetSalesDetail(queryType int, minDate, maxDate string) ([]*SalesDetail, error) {
	sql := "SELECT count(1) AS count, avg(amount) AS avg, product_type,date(from_unixtime(create_at)) AS dt, sum(amount) AS total_sales, " +
		"sum(if(pay_channel_id=1, amount, 0)) AS alipay, sum(if(pay_channel_id=2, amount, 0)) AS wxpay, " +
		"sum(if(status=5, refund_amount, 0)) AS refund_amount, sum(if(status=5, 1, 0)) AS refund_count, " +
		"sum(if(pay_channel_id=3,amount,0)) AS cash FROM venue_pay_orders WHERE status in(2,3,4,5,6) "
	if minDate != "" && maxDate != "" {
		sql += fmt.Sprintf("AND date(FROM_UNIXTIME(create_at)) >= '%s' AND date(FROM_UNIXTIME(create_at)) <= '%s'" +
			" ", minDate, maxDate)
	}

	switch queryType {
	// 根据商品分组
	case 1:
		sql += "GROUP BY product_type"
	// 根据日期分组
	case 2:
		sql += "GROUP BY dt"
	// 根据日期+商品分组
	case 3:
		sql += "GROUP BY dt,product_type"
	}

	var list []*SalesDetail
	if err := m.Engine.SQL(sql).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取用户消费总额
func (m *OrderModel) GetTotalSalesByUser(userId string) (int64, error) {
	return m.Engine.Where("user_id=?", userId).In("status", []int{2, 3, 4, 5, 6}).SumInt(m.Order, "amount")
}

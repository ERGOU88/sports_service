package cfinance

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/morder"
	"sports_service/server/models/mpay"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvenue"
	"sports_service/server/util"
	"fmt"
	"time"
)

// todo:
type FinanceModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	order       *morder.OrderModel
	venue       *mvenue.VenueModel
	pay         *mpay.PayModel
}

func New(c *gin.Context) FinanceModule {
	appSocket := dao.AppEngine.NewSession()
	defer appSocket.Close()

	venueSocket := dao.VenueEngine.NewSession()
	defer venueSocket.Close()
	return FinanceModule{
		context: c,
		user: muser.NewUserModel(appSocket),
		order: morder.NewOrderModel(venueSocket),
		venue: mvenue.NewVenueModel(venueSocket),
		pay: mpay.NewPayModel(venueSocket),
		engine: venueSocket,
	}
}

// 获取订单流水列表
func (svc *FinanceModule) GetOrderList(page, size int) (int, []*morder.OrderRecord) {
	offset := (page - 1) * size
	list, err := svc.order.GetOrderList(offset, size)
	if err != nil {
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*morder.OrderRecord{}
	}


	res := make([]*morder.OrderRecord, len(list))
	for index, item := range list {
		info := &morder.OrderRecord{
			Id:  item.Id,
			PayOrderId: item.PayOrderId,
			OriginalAmount: fmt.Sprintf("%.2f", float64(item.OriginalAmount)/100),
			CreateAt: time.Unix(int64(item.CreateAt), 0).Format(consts.FORMAT_TM),
			Amount: fmt.Sprintf("%.2f", float64(item.Amount)/100),
			Status: item.Status,
		}

		extra := mappointment.OrderResp{}
		if err := util.JsonFast.UnmarshalFromString(item.Extra, &extra); err == nil {
			info.MobileNum = extra.MobileNum
		}

		ok, err := svc.venue.GetVenueInfoById(fmt.Sprint(item.Id))
		if ok && err == nil {
			info.VenueName = svc.venue.Venue.VenueName
		}

		productName := svc.GetProductName(item.ProductType)
		info.Detail = fmt.Sprintf("%s * %d", productName, extra.Count)
		ok, err = svc.pay.GetPaymentChannel(svc.order.Order.PayChannelId)
		if ok && err == nil {
			info.PayChannel = svc.pay.PayChannel.Title
		}

		res[index] = info
	}

	return errdef.SUCCESS, res
}

// 获取退款流水列表
func (svc *FinanceModule) GetRefundList(orderId string, page, size int) (int, []*morder.RefundInfo) {
	offset := (page - 1) * size
	if orderId != "" {
		offset = 0
		size = 1
	}

	list, err := svc.order.GetRefundRecordList(orderId, offset, size)
	if err != nil {
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*morder.RefundInfo{}
	}

	for _, item := range list {
		info := &morder.RefundInfo{
			Id:  item.Id,
			PayOrderId: item.PayOrderId,
			AmountCn: fmt.Sprintf("%.2f", float64(item.Amount)/100),
			CreateAtCn: time.Unix(item.CreateAt, 0).Format(consts.FORMAT_TM),
			RefundAmountCn: fmt.Sprintf("%.2f", float64(item.RefundAmount)/100),
			Status: item.Status,
			RefundFeeCn: fmt.Sprintf("%.2f", float64(item.RefundFee)/100),
		}

		extra := mappointment.OrderResp{}
		if err := util.JsonFast.UnmarshalFromString(item.Extra, &extra); err == nil {
			info.MobileNum = extra.MobileNum
		}

		ok, err := svc.venue.GetVenueInfoById(fmt.Sprint(item.Id))
		if ok && err == nil {
			info.VenueName = svc.venue.Venue.VenueName
		}

		info.Detail = svc.GetProductName(item.ProductType)

		if info.OrderType == 1001 {
			info.OrderTypeCn = "线上退单"
		}

		if info.OrderType == 1002 {
			info.OrderTypeCn = "线下退单"
		}

		ok, err = svc.pay.GetPaymentChannel(svc.order.Order.PayChannelId)
		if ok && err == nil {
			info.RefundChannelCn = svc.pay.PayChannel.Title
		}

	}

	return errdef.SUCCESS, list
}

// 获取订单收益流水[已成功/已付款]
func (svc *FinanceModule) GetRevenueFlow(queryMinDate, queryMaxDate, orderId string, page, size int) (int, int64, []*morder.OrderRecord) {
	minDate := time.Now().Format(consts. FORMAT_DATE)
	maxDate := time.Now().Format(consts. FORMAT_DATE)
	if queryMinDate != "" && queryMaxDate != "" {
		minDate = queryMinDate
		maxDate = queryMaxDate
	}

	offset := (page - 1) * size
	list, err := svc.order.GetRevenueFlow(minDate, maxDate, orderId, offset, size)
	if err != nil {
		return errdef.ERROR, 0, []*morder.OrderRecord{}
	}

	if len(list) == 0 {
		return errdef.SUCCESS, 0, []*morder.OrderRecord{}
	}

	res := make([]*morder.OrderRecord, 0)
	for _, item := range list {
		statusCn, amountCn := svc.GetOrderStatusCn(item)
		info := &morder.OrderRecord{
			Id:    item.Id,
			PayOrderId: item.PayOrderId,
			CreateAt: time.Unix(int64(item.CreateAt), 0).Format(consts.FORMAT_TM),
			Amount: amountCn,
		    StatusCn: statusCn,
		}

		ok, err := svc.venue.GetVenueInfoById(fmt.Sprint(item.Id))
		if ok && err == nil {
			info.VenueName = svc.venue.Venue.VenueName
		}

		ok, err = svc.pay.GetPaymentChannel(svc.order.Order.PayChannelId)
		if ok && err == nil {
			info.PayChannel = svc.pay.PayChannel.Title
		}

		extra := mappointment.OrderResp{}
		if err := util.JsonFast.UnmarshalFromString(item.Extra, &extra); err == nil {
			info.MobileNum = extra.MobileNum
		}

		info.Detail = svc.GetProductName(item.ProductType)
		res = append(res, info)

		// 如果是退款流水 则应同时有一条购买流水
		if info.Status == consts.ORDER_TYPE_REFUND_SUCCESS || item.Status == consts.ORDER_TYPE_REFUND_WAIT {
			item.Status = consts.ORDER_TYPE_PAID
			info.StatusCn, info.Amount = svc.GetOrderStatusCn(item)
			res = append(res, info)
		}
	}

	// 累计收入= 销售总金额 - 退款总金额
	total := svc.GetTotalRevenue(minDate, maxDate) - svc.GetTotalRefund(minDate, maxDate)

	return errdef.SUCCESS, total, res
}

func (svc *FinanceModule) GetOrderStatusCn(item *models.VenuePayOrders) (string, string) {
	var statusCn, amountCn string
	switch item.Status {
	// 已支付/已完成/已过期 展示购买
	case consts.ORDER_TYPE_PAID, consts.ORDER_TYPE_COMPLETED, consts.ORDER_TYPE_EXPIRE:
		statusCn = "购买"
		amountCn = fmt.Sprintf("%.2f", float64(item.Amount)/100)
	// 退款中/已退款 展示退单
	case consts.ORDER_TYPE_REFUND_WAIT, consts.ORDER_TYPE_REFUND_SUCCESS:
		statusCn = "退单"
		amountCn = fmt.Sprintf("%.2f", float64(item.RefundAmount)/100)
	}

	return statusCn, amountCn

}

func (svc *FinanceModule) GetProductName(productType int) string {
	var productName string
	switch productType {
	case consts.ORDER_TYPE_APPOINTMENT_VENUE:
		productName = "预约场馆"
	case consts.ORDER_TYPE_APPOINTMENT_COACH:
		productName = "预约私教"
	case consts.ORDER_TYPE_APPOINTMENT_COURSE:
		productName = "预约课程"
	case consts.ORDER_TYPE_EXPERIENCE_CARD:
		productName = "次卡"
	case consts.ORDER_TYPE_MONTH_CARD:
		productName = "月卡"
	case consts.ORDER_TYPE_SEANSON_CARD:
		productName = "季卡"
	case consts.ORDER_TYPE_HALF_YEAR_CARD:
		productName = "半年卡"
	case consts.ORDER_TYPE_YEAR_CARD:
		productName = "年卡"
	default:
	   productName = "实物商品"
	}

	return productName
}

// 获取退款总额
func (svc *FinanceModule) GetTotalRefund(minDate, maxDate string) int64 {
	total, err := svc.order.GetTotalRefund(minDate, maxDate)
	if err != nil {
		return 0
	}

	return total
}

// 获取总收益
func (svc *FinanceModule) GetTotalRevenue(minDate, maxDate string) int64 {
	total, err := svc.order.GetTotalRevenue(minDate, maxDate)
	if err != nil {
		log.Log.Errorf("finance_trace: get total revenue fail, err:%s", err)
		return 0
	}

	return total
}

// 财务统计
func (svc *FinanceModule) HomePageStat(minDate, maxDate string) (int, *morder.OrderStat){
	// 总销售额
	totalSales, err := svc.order.GetTotalRevenue("", "")
	if err != nil {
		return errdef.ERROR, nil
	}

	orderStat := &morder.OrderStat{
		TopInfo: make(map[string]interface{}, 0),
	}
	orderStat.TopInfo["total_sales"] = totalSales

	today := time.Now().Format(consts. FORMAT_DATE)
	todaySales, err := svc.order.GetTotalRevenue(today, today)
	if err != nil {
		return errdef.ERROR, orderStat
	}
	orderStat.TopInfo["today_sales"] = todaySales


	yesterday := time.Now().AddDate(0, 0, -1).Format(consts.FORMAT_DATE)
	yesterdaySales, err := svc.order.GetTotalRevenue(yesterday, yesterday)
	if err != nil {
		return errdef.ERROR, orderStat
	}
	orderStat.TopInfo["yesterday_sales"] = yesterdaySales

	if todaySales > 0 && yesterdaySales > 0 || todaySales == 0 && yesterdaySales > 0 {
		orderStat.TopInfo["sales_rate"] = fmt.Sprintf("%.0f%s", (float64(todaySales)/float64(yesterdaySales)-1)*100, "%")
	}

	if todaySales > 0 && yesterdaySales == 0 || todaySales == 0 && yesterdaySales == 0 {
		orderStat.TopInfo["sales_rate"] = "--%"
	}

	orderStat.TopInfo["total_user"] = svc.order.GetVenueTotalUser()
	totalVip, err := svc.order.GetVipUserCount("")
	if err != nil {
		return errdef.ERROR, orderStat
	}

	orderStat.TopInfo["total_vip"] = totalVip

	todayOrder, err := svc.order.GetOrderNum(today, today)
	if err != nil {
		return errdef.ERROR, orderStat
	}
	orderStat.TopInfo["today_order"] = todayOrder

	week := time.Now().AddDate(0, 0, -6).Format(consts.FORMAT_DATE)
	weekOrder, err := svc.order.GetOrderNum(week, today)
	if err != nil {
		return errdef.ERROR, orderStat
	}
	orderStat.TopInfo["week_order"] = weekOrder

	return errdef.SUCCESS, orderStat
}

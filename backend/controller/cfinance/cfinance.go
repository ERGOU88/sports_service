package cfinance

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/dao"
	"sports_service/global/backend/errdef"
	"sports_service/global/backend/log"
	"sports_service/global/consts"
	"sports_service/models"
	"sports_service/models/mappointment"
	"sports_service/models/morder"
	"sports_service/models/mpay"
	"sports_service/models/muser"
	"sports_service/models/mvenue"
	"sports_service/util"
	"strings"
	"time"
)

// todo:
type FinanceModule struct {
	context *gin.Context
	engine  *xorm.Session
	user    *muser.UserModel
	order   *morder.OrderModel
	venue   *mvenue.VenueModel
	pay     *mpay.PayModel
}

func New(c *gin.Context) FinanceModule {
	appSocket := dao.AppEngine.NewSession()
	defer appSocket.Close()

	venueSocket := dao.VenueEngine.NewSession()
	defer venueSocket.Close()
	return FinanceModule{
		context: c,
		user:    muser.NewUserModel(appSocket),
		order:   morder.NewOrderModel(venueSocket),
		venue:   mvenue.NewVenueModel(venueSocket),
		pay:     mpay.NewPayModel(venueSocket),
		engine:  venueSocket,
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
			Id:             item.Id,
			PayOrderId:     item.PayOrderId,
			OriginalAmount: fmt.Sprintf("%.2f", float64(item.OriginalAmount)/100),
			CreateAt:       time.Unix(int64(item.CreateAt), 0).Format(consts.FORMAT_TM),
			Amount:         fmt.Sprintf("%.2f", float64(item.Amount)/100),
			Status:         item.Status,
		}

		extra := &mappointment.OrderResp{}
		if err := util.JsonFast.UnmarshalFromString(item.Extra, extra); err != nil {
			log.Log.Errorf("finance_trace: unmarshal order resp fail, orderId:%s, err:%s", item.PayOrderId, err)
		}

		if user := svc.user.FindUserByUserid(item.UserId); user != nil {
			info.MobileNum = util.HideMobileNum(fmt.Sprint(svc.user.User.MobileNum))
		}

		ok, err := svc.venue.GetVenueInfoById(fmt.Sprint(item.VenueId))
		if ok && err == nil {
			info.VenueName = svc.venue.Venue.VenueName
		}

		productName := svc.GetProductName(item.ProductType)
		if extra.Id > 0 {
			info.Detail = productName
			if extra.Count > 0 {
				info.Detail = fmt.Sprintf("%s * %d", productName, extra.Count)
			}
		}

		info.ProductName = productName
		ok, err = svc.pay.GetPaymentChannel(item.PayChannelId)
		if ok && err == nil {
			info.PayChannel = svc.pay.PayChannel.Title
		}

		res[index] = info
	}

	return errdef.SUCCESS, res
}

// 订单总数
func (svc *FinanceModule) GetOrderTotal() int64 {
	count, err := svc.order.GetOrderCount(nil)
	if err != nil {
		return 0
	}

	return count
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
		log.Log.Errorf("finance_trace: get refund list fail, err:%s", err)
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*morder.RefundInfo{}
	}

	for _, item := range list {
		item.AmountCn = fmt.Sprintf("%.2f", float64(item.Amount)/100)
		item.CreateAtCn = time.Unix(item.CreateAt, 0).Format(consts.FORMAT_TM)
		item.RefundAmountCn = fmt.Sprintf("%.2f", float64(item.RefundAmount)/100)
		item.RefundFeeCn = fmt.Sprintf("%.2f", float64(item.RefundFee)/100)

		if user := svc.user.FindUserByUserid(item.UserId); user != nil {
			item.MobileNum = util.HideMobileNum(fmt.Sprint(svc.user.User.MobileNum))
		}

		ok, err := svc.venue.GetVenueInfoById(fmt.Sprint(item.Id))
		if ok && err == nil {
			item.VenueName = svc.venue.Venue.VenueName
		}

		item.Detail = svc.GetProductName(item.ProductType)

		if item.OrderType == 1001 {
			item.OrderTypeCn = "线上退单"
		}

		if item.OrderType == 1002 {
			item.OrderTypeCn = "线下退单"
		}

		item.ProductName = svc.GetProductName(item.ProductType)

		ok, err = svc.pay.GetPaymentChannel(item.RefundChannelId)
		if ok && err == nil {
			item.RefundChannelCn = svc.pay.PayChannel.Title
		}

	}

	return errdef.SUCCESS, list
}

// 退款记录总数
func (svc *FinanceModule) GetRefundRecordTotal() int64 {
	count, err := svc.order.GetRefundRecordTotal()
	if err != nil {
		log.Log.Errorf("finance_trace: get refund total fail, err:%s", err)
		return 0
	}

	return count
}

// 获取订单收益流水[已成功/已付款]
func (svc *FinanceModule) GetRevenueFlow(queryMinDate, queryMaxDate, orderId string, page, size int) (int, int64, []morder.OrderRecord) {
	minDate := time.Now().Format(consts.FORMAT_DATE)
	maxDate := time.Now().Format(consts.FORMAT_DATE)
	if queryMinDate != "" && queryMaxDate != "" {
		minDate = queryMinDate
		maxDate = queryMaxDate
	}

	offset := (page - 1) * size
	list, err := svc.order.GetRevenueFlow(minDate, maxDate, orderId, offset, size)
	if err != nil {
		return errdef.ERROR, 0, []morder.OrderRecord{}
	}

	if len(list) == 0 {
		return errdef.SUCCESS, 0, []morder.OrderRecord{}
	}

	res := make([]morder.OrderRecord, 0)
	for _, item := range list {
		statusCn, amountCn := svc.GetOrderStatusCn(item)
		info := morder.OrderRecord{
			Id:         item.Id,
			PayOrderId: item.PayOrderId,
			CreateAt:   time.Unix(int64(item.CreateAt), 0).Format(consts.FORMAT_TM),
			Amount:     amountCn,
			StatusCn:   statusCn,
			Status:     item.Status,
		}

		if item.ProductType == 5102 {
			statusCn = "结算"
		}

		ok, err := svc.venue.GetVenueInfoById(fmt.Sprint(item.VenueId))
		if ok && err == nil {
			info.VenueName = svc.venue.Venue.VenueName
		}

		ok, err = svc.pay.GetPaymentChannel(item.PayChannelId)
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
		if item.Status == consts.ORDER_TYPE_REFUND_SUCCESS || item.Status == consts.ORDER_TYPE_REFUND_WAIT {
			item.Status = consts.ORDER_TYPE_PAID
			info.Status = consts.ORDER_TYPE_PAID
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
	case consts.ORDER_TYPE_SETTLEMENT:
		productName = "线下结算"
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

// 财务首页 顶部统计
func (svc *FinanceModule) TopStat() (int, *morder.OrderStat) {
	// 总销售额
	totalSales, err := svc.order.GetTotalRevenue("", "")
	if err != nil {
		return errdef.ERROR, nil
	}

	orderStat := &morder.OrderStat{
		TopInfo: make(map[string]interface{}, 0),
	}
	orderStat.TopInfo["total_sales"] = totalSales

	today := time.Now().Format(consts.FORMAT_DATE)
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

	orderStat.TopInfo["sales_rate"] = svc.GetRingRatio(todaySales, yesterdaySales)

	orderStat.TopInfo["total_user"] = svc.order.GetVenueTotalUser()
	totalVip, err := svc.order.GetVipUserCount("", "")
	if err != nil {
		return errdef.ERROR, orderStat
	}

	orderStat.TopInfo["total_vip"] = totalVip
	todayVip, err := svc.order.GetVipUserCount(today, today)
	if err != nil {
		return errdef.ERROR, orderStat
	}
	orderStat.TopInfo["today_vip"] = todayVip

	yesterdayVip, err := svc.order.GetVipUserCount(yesterday, yesterday)
	if err != nil {
		return errdef.ERROR, orderStat
	}
	orderStat.TopInfo["yesterday_vip"] = yesterdayVip
	orderStat.TopInfo["vip_rate"] = svc.GetRingRatio(todayVip, yesterdayVip)

	todayOrder, err := svc.order.GetOrderNum(today, today)
	if err != nil {
		return errdef.ERROR, orderStat
	}
	orderStat.TopInfo["today_order"] = todayOrder

	yesterdayOrder, err := svc.order.GetOrderNum(yesterday, yesterday)
	if err != nil {
		return errdef.ERROR, orderStat
	}
	orderStat.TopInfo["yesterday_order"] = yesterdayOrder
	orderStat.TopInfo["order_rate"] = svc.GetRingRatio(todayOrder, yesterdayOrder)

	week := time.Now().AddDate(0, 0, -6).Format(consts.FORMAT_DATE)
	weekOrder, err := svc.order.GetOrderNum(week, today)
	if err != nil {
		return errdef.ERROR, orderStat
	}
	orderStat.TopInfo["week_order"] = weekOrder

	weekDay := time.Now().AddDate(0, 0, -6)

	// 本周新增
	var weekNewUser int
	orderStat.TopInfo["week_new_users"] = 0
	weekUser, err := svc.order.GetVenueNewUsers(weekDay.Format(consts.FORMAT_DATE), today, "")
	if err == nil && weekUser != nil {
		weekNewUser = len(weekUser)
		orderStat.TopInfo["week_new_users"] = fmt.Sprintf("%.2f", float64(weekNewUser)/7)
		// 之前是否已成为场馆用户
		beforeUser, err := svc.order.GetVenueNewUsers("",
			weekDay.AddDate(0, 0, -1).Format(consts.FORMAT_DATE), strings.Join(weekUser, ","))
		if err == nil && beforeUser != nil {
			weekNewUser = len(weekUser) - len(beforeUser)
			orderStat.TopInfo["week_new_users"] = fmt.Sprintf("%.2f", float64(weekNewUser)/7)
		}
	}

	// 上周新增场馆用户
	var lastWeekNewUser int
	orderStat.TopInfo["last_week_new_users"] = 0
	lastWeekDay := weekDay.AddDate(0, 0, -1)
	lastWeekUser, err := svc.order.GetVenueNewUsers(lastWeekDay.AddDate(0, 0, -6).Format(consts.FORMAT_DATE),
		lastWeekDay.Format(consts.FORMAT_DATE), "")
	if err == nil && lastWeekUser != nil {
		lastWeekNewUser = len(lastWeekUser)
		orderStat.TopInfo["last_week_new_users"] = fmt.Sprintf("%.2f", float64(lastWeekNewUser)/7)
		// 之前是否已成为场馆用户
		beforeUser, err := svc.order.GetVenueNewUsers("",
			lastWeekDay.AddDate(0, 0, -8).Format(consts.FORMAT_DATE), strings.Join(lastWeekUser, ","))
		log.Log.Error("#######err:%s", err)
		if err == nil && beforeUser != nil {
			lastWeekNewUser = len(lastWeekUser) - len(beforeUser)
			// 7日平均
			orderStat.TopInfo["last_week_new_users"] = fmt.Sprintf("%.2f", float64(lastWeekNewUser)/7)
		}
	}

	orderStat.TopInfo["user_rate"] = svc.GetRingRatio(int64(weekNewUser), int64(lastWeekNewUser))

	return errdef.SUCCESS, orderStat
}

// 获取图表统计数据
func (svc *FinanceModule) GetChartStat(queryMinDate, queryMaxDate string) (int, map[string]interface{}) {
	days := 6
	minDate := time.Now().AddDate(0, 0, -days).Format(consts.FORMAT_DATE)
	maxDate := time.Now().Format(consts.FORMAT_DATE)
	if queryMinDate != "" && queryMaxDate != "" {
		minDate = queryMinDate
		maxDate = queryMaxDate
		min, err := time.Parse(consts.FORMAT_DATE, queryMinDate)
		if err != nil {
			log.Log.Errorf("finance_trace: time.Parse fail, err:%s", err)
			return errdef.ERROR, nil
		}

		max, err := time.Parse(consts.FORMAT_DATE, queryMaxDate)
		if err != nil {
			log.Log.Errorf("finance_trace: time.Parse fail, err:%s", err)
			return errdef.ERROR, nil
		}

		days = util.GetDiffDays(max, min)
	}

	// 商品分组
	list, err := svc.order.GetSalesDetail(1, minDate, maxDate)
	if err != nil {
		log.Log.Errorf("finance_trace: get sales detail fail, err:%s", err)
		return errdef.ERROR, nil
	}

	if len(list) > 0 {
		for _, item := range list {
			item.ProductName = svc.GetProductName(item.ProductType)
		}

	}

	// 合计
	totalInfo, err := svc.order.GetSalesDetail(0, minDate, maxDate)
	if err != nil {
		log.Log.Errorf("finance_trace: get sales detail fail, err:%s", err)
		return errdef.ERROR, nil
	}

	if len(totalInfo) > 0 {
		for _, item := range totalInfo {
			item.ProductType = 0
			item.ProductName = "合计"
		}
	}

	// 销售明细
	list = append(list, totalInfo...)
	result := make(map[string]interface{}, 0)
	result["sales_detail"] = list

	// 日期 + 商品分组
	salesByProduct, err := svc.order.GetSalesDetail(3, minDate, maxDate)
	if err != nil {
		log.Log.Errorf("finance_trace: get sales detail fail, err:%s", err)
		return errdef.ERROR, nil
	}

	// n天的销量总额 日期分组
	salesByDate, err := svc.order.GetSalesDetail(2, minDate, maxDate)
	if err != nil {
		log.Log.Errorf("finance_trace: get sales detail fail, err:%s", err)
		return errdef.ERROR, nil
	}
	if len(salesByDate) > 0 {
		for _, item := range salesByDate {
			item.ProductType = 0
			item.ProductName = "合计"
		}
	}

	var salesResultList []morder.ResultList
	salesResultList = append(salesResultList, morder.ResultList{
		Title: "销售总额",
		List:  svc.ResultInfoByDate(salesByDate, days, 0, 0, maxDate),
	}, morder.ResultList{
		Title: "课程预约",
		List:  svc.ResultInfoByDate(salesByProduct, days, consts.ORDER_TYPE_APPOINTMENT_COURSE, 1, maxDate),
	}, morder.ResultList{
		Title: "私教预约",
		List:  svc.ResultInfoByDate(salesByProduct, days, consts.ORDER_TYPE_APPOINTMENT_COACH, 1, maxDate),
	}, morder.ResultList{
		Title: "实体商品",
		List:  svc.ResultInfoByDate(salesByProduct, days, consts.ORDER_TYPE_PHYSICAL_GOODS, 1, maxDate),
	}, morder.ResultList{
		Title: "卡类商品[次卡/月卡/季卡/半年卡/年卡]",
		List:  svc.ResultInfoByDate(salesByProduct, days, 2000, 1, maxDate),
	})

	result["sales_result_list"] = salesResultList

	var payChannelResultList []morder.ResultList
	payChannelResultList = append(payChannelResultList, morder.ResultList{
		Title: "收款总额",
		List:  svc.ResultInfoByDate(salesByDate, days, 0, 0, maxDate),
	}, morder.ResultList{
		Title: "支付宝支付",
		List:  svc.ResultInfoByDate(salesByDate, days, 1, 2, maxDate),
	}, morder.ResultList{
		Title: "微信支付",
		List:  svc.ResultInfoByDate(salesByDate, days, 2, 2, maxDate),
	}, morder.ResultList{
		Title: "现金支付",
		List:  svc.ResultInfoByDate(salesByDate, days, 3, 2, maxDate),
	})
	result["pay_result_list"] = payChannelResultList

	var orderResultList []morder.ResultList
	orderResultList = append(orderResultList, morder.ResultList{
		Title: "订单均价",
		List:  svc.ResultInfoByDate(salesByDate, days, 0, 3, maxDate),
	})

	result["order_result_list"] = orderResultList

	return errdef.SUCCESS, result
}

// 获取环比
func (svc *FinanceModule) GetRingRatio(current, before int64) string {
	if current > 0 && before > 0 || current == 0 && before > 0 {
		return fmt.Sprintf("%.0f%s", (float64(current)/float64(before)-1)*100, "%")
	}

	if current > 0 && before == 0 || current == 0 && before == 0 {
		return "--%"
	}

	return ""
}

func (svc *FinanceModule) ResultInfoByDate(data []*morder.SalesDetail, days, condition, queryType int, maxDate string) map[string]interface{} {
	mapList := make(map[string]interface{})
	max, err := time.Parse(consts.FORMAT_DATE, maxDate)
	if err != nil {
		return nil
	}

	for i := 0; i <= days; i++ {
		date := max.AddDate(0, 0, -i).Format("2006-01-02")
		for _, v := range data {
			if v.Dt != date {
				continue
			}

			switch queryType {
			// 收款总额
			case 0:
				mapList[date] = v.TotalSales

			// 通过商品类型查询
			case 1:
				// 卡类 包含 次卡/月卡/季卡/年卡 需叠加
				if v.ProductType >= 2000 && v.ProductType < 3000 && condition == 2000 {
					if val, ok := mapList[date]; ok {
						mapList[date] = val.(int) + v.TotalSales
						continue
					}

					mapList[date] = v.TotalSales
					continue
				}

				if v.ProductType == condition {
					mapList[date] = v.TotalSales
				}

			// 通过支付渠道查询
			case 2:
				switch condition {
				case 0:
					mapList[date] = v.TotalSales
				// 支付宝
				case 1:
					mapList[date] = v.Alipay
				// 微信
				case 2:
					mapList[date] = v.Wxpay
				// 现金
				case 3:
					mapList[date] = v.Cash
				}
			// 订单均价
			case 3:
				mapList[date] = fmt.Sprintf("%.2f", v.Avg)
			}
		}

		if _, ok := mapList[date]; !ok {
			mapList[date] = 0
		}
	}

	return mapList
}

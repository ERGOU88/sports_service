package cappointment

import (
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/morder"
	"sports_service/server/util"
	"strings"
	"time"
)

type base struct {
	Engine      *xorm.Session
	appointment *mappointment.AppointmentModel
	order       *morder.OrderModel
	DateId      int
	Extra       *mappointment.AppointmentResp
	// 预约流水map
	recordMp    map[int64]*models.AppointmentRecord
	// 订单商品流水map
	orderMp     map[int64]*models.VenueOrderProductInfo
}

func New(socket *xorm.Session) *base {
	return &base{
		Engine: socket,
		appointment: mappointment.NewAppointmentModel(socket),
		order: morder.NewOrderModel(socket),
		Extra:  &mappointment.AppointmentResp{TimeNodeInfo: make([]*mappointment.TimeNodeInfo, 0)},
		recordMp: make(map[int64]*models.AppointmentRecord),
		orderMp: make(map[int64]*models.VenueOrderProductInfo),
	}
}

func (svc *base) AppointmentDateInfo(days, appointmentType int) interface{} {
	list := svc.GetAppointmentDate(days)
	res := make([]*mappointment.WeekInfo, len(list))
	for index, v := range list {
		info := &mappointment.WeekInfo{
			Id: v.Id,
			Week: v.Week,
			Date: v.Date,
			WeekCn: v.WeekCn,
		}

		svc.appointment.AppointmentInfo.WeekNum = v.Week
		svc.appointment.AppointmentInfo.AppointmentType = appointmentType
		svc.appointment.AppointmentInfo.CurAmount = 0

		if err := svc.appointment.GetMinPriceByWeek(); err != nil {
			log.Log.Errorf("venue_trace: get min price fail, err:%s", err)
		}

		info.MinPrice = svc.appointment.AppointmentInfo.CurAmount
		info.PriceCn = fmt.Sprintf("¥%.2f", float64(info.MinPrice)/100)
		res[index] = info
	}

	return res
}

// 获取预约的日期信息（从当天开始推算）
// days 天数
func (svc *base) GetAppointmentDate(days int) []util.DateInfo {
	curTime := time.Now()
	// 今天
	today := curTime.Format("2006-01-02")
	// 往后推6天 总共7天
	afterDay := curTime.AddDate(0, 0, days).Format("2006-01-02")
	dateInfo := util.GetBetweenDates(today, afterDay)

	return dateInfo
}

// 预约 场馆/私教 选项
func (svc *base) GetAppointmentOptions() ([]*models.VenueAppointmentInfo, error) {
	list, err := svc.appointment.GetOptionsByWeek()
	if err != nil {
		return nil, err
	}

	if list == nil {
		return []*models.VenueAppointmentInfo{}, nil
	}

	return list, nil
}

func (svc *base) SetWeek(week int) {
	svc.appointment.AppointmentInfo.WeekNum = week
}

func (svc *base) SetRelatedId(relatedId int) {
	svc.appointment.AppointmentInfo.RelatedId = int64(relatedId)
}

func (svc *base) SetAppointmentType(appointmentType int) {
	svc.appointment.AppointmentInfo.AppointmentType = appointmentType
}

func (svc *base) SetStockRelatedId(relatedId int64) {
	svc.appointment.Stock.RelatedId = relatedId
}

func (svc *base) SetStockDate(date string) {
	svc.appointment.Stock.Date = date
}

func (svc *base) SetStockTimeNode(timeNode string) {
	svc.appointment.Stock.TimeNode = timeNode
}

func (svc *base) SetStockAppointmentType(appointmentType int) {
	svc.appointment.Stock.AppointmentType = appointmentType
}

// 日期id
func (svc *base) SetDateId(id int) {
	svc.DateId = id
}

// 通过id获取日期
func (svc *base) GetDateById(id int, formatType string) string {
	var date string
	curTime := time.Now()

	if id >= 1 {
		date = curTime.AddDate(0, 0, id - 1).Format(formatType)
	}

	return date
}

func (svc *base) SetAppointmentOptionsRes(date string, item *models.VenueAppointmentInfo) *mappointment.OptionsInfo {
	if svc.DateId == 1 {
		nodes := strings.Split(item.TimeNode, "-")
		if len(nodes) ==  2 {
			now := time.Now().Unix()
			// 获取 预约时间配置的开始时间
			start := fmt.Sprintf("%s %s", date, nodes[0])
			ts := new(util.TimeS)
			startTm := ts.GetTimeStrOrStamp(start, "YmdHi")
			// 如果当前时间 > 配置中的开始时间 过滤该配置项
			if now > startTm.(int64) {
				log.Log.Errorf("过滤id：%d, date:%s", item.Id, date)
				return nil
			}
		}
	}

	info := &mappointment.OptionsInfo{
		RelatedId: item.RelatedId,
		CurAmount: item.CurAmount,
		TimeNode: item.TimeNode,
		DurationCn: util.ResolveTime(item.Duration),
		Duration: item.Duration,
		RealAmount: item.RealAmount,
		QuotaNum: item.QuotaNum,
		RecommendType: item.RecommendType,
		AppointmentType: item.AppointmentType,
		WeekNum: item.WeekNum,
		AmountCn: fmt.Sprintf("¥%.2f", float64(item.CurAmount)/100),
		Id: item.Id,
	}

	// 售价 < 定价 表示有优惠
	if item.CurAmount < item.RealAmount {
		info.HasDiscount = 1
		info.DiscountRate = item.DiscountRate
		info.DiscountAmount = item.DiscountAmount
		if item.DiscountRate % 10 == 0 {
			info.RateCn = fmt.Sprintf("%d", item.DiscountRate/10)
		} else {
			info.RateCn = fmt.Sprintf("%.1f", float64(item.DiscountRate)/10)
		}
	}

	ok, err := svc.QueryStockInfo(item.AppointmentType, item.RelatedId, date, item.TimeNode)
	if err != nil {
		log.Log.Errorf("venue_trace: get purchase num fail, err:%s", err)
	}

	if ok {
		log.Log.Debugf("stock:%+v", svc.appointment.Stock)
		info.PurchasedNum = svc.appointment.Stock.PurchasedNum
		info.QuotaNum = svc.appointment.Stock.QuotaNum
		// 库存 <= 冻结库存 表示已满场
		if info.QuotaNum <= info.PurchasedNum {
			info.IsFull = 1
		}
	}

	return info
}

func (svc *base) GetAppointmentConf(id int64) error  {
	ok, err := svc.appointment.GetAppointmentConfById(fmt.Sprint(id))
	if err != nil || !ok {
		return errors.New("get appointment fail")
	}

	return nil
}

func (svc *base) GetAppointmentConfByIds(ids []interface{}) ([]*models.VenueAppointmentInfo, error) {
	return svc.appointment.GetAppointmentConfByIds(ids)
}

// 预约类型、关联ID、日期、时间节点查询库存信息
func (svc *base) QueryStockInfo(appointmentType int, relatedId int64, date, timeNode string) (bool, error) {
	//svc.SetAppointmentType(appointmentType)
	//svc.SetStockRelatedId(relatedId)
	//svc.SetStockDate(date)
	//svc.SetStockTimeNode(timeNode)
	return svc.appointment.GetStockInfo(appointmentType, relatedId, date, timeNode)
}

func (svc *base) SetOrderProductInfo(orderId string, now, count int, relatedId int64) *models.VenueOrderProductInfo {
	return &models.VenueOrderProductInfo{
		ProductId:   svc.appointment.AppointmentInfo.Id,
		OrderType:   consts.ORDER_TYPE_APPOINTMENT_VENUE,
		Count:       count,
		RealAmount:  svc.appointment.AppointmentInfo.RealAmount,
		CurAmount:   svc.appointment.AppointmentInfo.CurAmount,
		DiscountRate: svc.appointment.AppointmentInfo.DiscountRate,
		DiscountAmount: svc.appointment.AppointmentInfo.DiscountAmount,
		Amount: svc.appointment.AppointmentInfo.CurAmount * count,
		Duration: svc.appointment.AppointmentInfo.Duration * count,
		CreateAt: now,
		UpdateAt: now,
		RelatedId: relatedId,
		PayOrderId: orderId,
	}
}

func (svc *base) SetAppointmentRecordInfo(userId, date, orderId string, now, count int, seatInfo []*mappointment.SeatInfo, relatedId int64) *models.AppointmentRecord {
	info, _ := util.JsonFast.MarshalToString(seatInfo)
	return &models.AppointmentRecord{
		UserId: userId,
		RelatedId: relatedId,
		AppointmentType: svc.appointment.AppointmentInfo.AppointmentType,
		TimeNode: svc.appointment.AppointmentInfo.TimeNode,
		Date: date,
		PayOrderId: orderId,
		CreateAt: now,
		UpdateAt: now,
		PurchasedNum: count,
		SeatInfo: info,
	}
}

// 添加库存
func (svc *base) AddStock(date string, now, count int) (bool, error) {
	// 库存数据不存在 写入
	data := &models.VenueAppointmentStock{
		Date:            date,
		TimeNode:        svc.appointment.AppointmentInfo.TimeNode,
		QuotaNum:        svc.appointment.AppointmentInfo.QuotaNum,
		PurchasedNum:    count,
		AppointmentType: svc.appointment.AppointmentInfo.AppointmentType,
		RelatedId:       svc.appointment.AppointmentInfo.RelatedId,
		CreateAt:        now,
		UpdateAt:        now,
	}

	_, err := svc.appointment.AddStockInfo(data)
	if err != nil {
		var myErr *mysql.MySQLError
		// 是否为mysql错误
		if errors.As(err, &myErr) {
			// 是唯一索引约束错误 则返回true
			if myErr.Number == 1062 {
				return true, err
			}
		}

		return false, err
	}

	return false, nil
}

// 更新库存
func (svc *base) UpdateStock(date string, count, now int) (int64, error) {
	// 更新库存
	return svc.appointment.UpdateStockInfo(svc.appointment.AppointmentInfo.TimeNode, date, count, now,
		svc.appointment.AppointmentInfo.AppointmentType, int(svc.appointment.AppointmentInfo.RelatedId))
}

// 添加预约流水
func (svc *base) AddAppointmentRecord() error {
	var rlst []*models.AppointmentRecord
	for _, val := range svc.recordMp {
		rlst = append(rlst, val)
	}

	// 添加预约记录流水
	affected, err := svc.appointment.AddMultiAppointmentRecord(rlst)
	if err != nil {
		return err
	}

	if affected != int64(len(rlst)) {
		return errors.New("add record fail, count not match~")
	}

	return nil
}

// 添加订单商品流水记录
func (svc *base) AddOrderProducts() error {
	var olst []*models.VenueOrderProductInfo
	for _, val := range svc.orderMp {
		olst = append(olst, val)
	}

	// 添加订单商品流水
	affected, err := svc.order.AddMultiOrderProduct(olst)
	if err != nil {
		return err
	}

	if affected != int64(len(olst)) {
		return errors.New("add order product fail, count not match~")
	}

	return nil
}

// 添加订单
func (svc *base) AddOrder(orderId, userId string, now, total int) error {
	extra, _ := util.JsonFast.MarshalToString(svc.Extra)
	svc.order.Order.Extra = extra
	svc.order.Order.PayOrderId = orderId
	svc.order.Order.UserId = userId
	svc.order.Order.OrderType = 1001
	svc.order.Order.CreateAt = now
	svc.order.Order.UpdateAt = now
	svc.order.Order.Amount = total
	affected, err := svc.order.AddOrder()
	if err != nil {
		return err
	}

	if affected != 1 {
		return errors.New("add order fail, affected not 1")
	}

	return nil
}

// 设置最新库存数据（返回使用）
func (svc *base) SetLatestInventoryResp(date string) (*mappointment.TimeNodeInfo, error) {
	info := &mappointment.TimeNodeInfo{}
	info.Id = svc.appointment.AppointmentInfo.Id
	info.TimeNode = svc.appointment.AppointmentInfo.TimeNode
	info.Date = date
	info.Amount = svc.appointment.AppointmentInfo.CurAmount
	info.Discount = svc.appointment.AppointmentInfo.DiscountAmount
	// 查看最新的库存 并返回 tips：读快照无问题 购买时保证数据一致性即可
	ok, err := svc.QueryStockInfo(svc.appointment.AppointmentInfo.AppointmentType,
		svc.appointment.AppointmentInfo.RelatedId, date, svc.appointment.AppointmentInfo.TimeNode)
	if !ok || err != nil {
		log.Log.Errorf("venue_trace: get stock info fail, err:%s", err)
		//return nil, errors.New("get stock info fail")
	} else {
		info.Count = svc.appointment.Stock.QuotaNum - svc.appointment.Stock.PurchasedNum
	}

	return info, nil
}

// 预约流程
func (svc *base) AppointmentProcess(userId, orderId string, relatedId int64, infos []*mappointment.AppointmentInfo) error {
	svc.Extra.IsEnough = true
	now := int(time.Now().Unix())
	//stockInfo := make([]*mappointment.StockInfoResp, 0)
	//// 订单商品流水map
	//orderMp := make(map[int64]*models.VenueOrderProductInfo)
	//// 预约流水map
	//recordMp := make(map[int64]*models.AppointmentRecord)
	for _, item := range infos {
		err := svc.GetAppointmentConf(item.Id)
		if svc.appointment.AppointmentInfo == nil || err != nil {
			log.Log.Errorf("venue_trace: get appointment conf fail, id:%d", item.Id)
			return errors.New("get appointment conf fail")
		}

		if len(item.SeatInfos) != item.Count {
			log.Log.Errorf("venue_trace: seat num not match, count:%d, len:%d", item.Count, len(item.SeatInfos))
			return errors.New("seat num not match")
		}

		log.Log.Errorf("appointment_info:%+v, id:%d", svc.appointment.AppointmentInfo, svc.appointment.AppointmentInfo.Id)

		if item.Count <= 0 || item.Count > svc.appointment.AppointmentInfo.QuotaNum {
			log.Log.Errorf("venue_trace: invalid count, count:%d, quotaNum:%d", item.Count, svc.appointment.AppointmentInfo.QuotaNum)
			return errors.New("invalid count")
		}

		date := svc.GetDateById(item.DateId, consts.FORMAT_DATE)
		if date == "" {
			log.Log.Errorf("venue_trace: invalid date, dateId:%d", item.DateId)
			return errors.New("invalid date")
		}

		svc.Extra.TotalTm += svc.appointment.AppointmentInfo.Duration * item.Count
		// 数量 * 售价
		svc.Extra.TotalAmount += item.Count * svc.appointment.AppointmentInfo.CurAmount
		svc.Extra.TotalDiscount += item.Count * svc.appointment.AppointmentInfo.DiscountAmount
		svc.orderMp[item.Id] = svc.SetOrderProductInfo(orderId, now, item.Count, relatedId)
		svc.recordMp[item.Id] = svc.SetAppointmentRecordInfo(userId, date, orderId, now, item.Count, item.SeatInfos, relatedId)

		ok, err := svc.appointment.HasExistsStockInfo(svc.appointment.AppointmentInfo.AppointmentType,
			svc.appointment.AppointmentInfo.RelatedId, date, svc.appointment.AppointmentInfo.TimeNode)
		if err != nil {
			log.Log.Errorf("venue_trace: query stock info fail, err:%s", err)
			return err
		}


		var b bool
		if !ok {
			b, err = svc.AddStock(date, now, item.Count)
			log.Log.Errorf("id:%d, b:%v, err:%s", svc.appointment.AppointmentInfo.Id, b, err)
			// false 表示存在错误 但不是唯一索引约束错误
			if !b && err != nil {
				log.Log.Errorf("venue_trace: add stock fail, err:%s", err)
				return err
			}
		}

		// 数据存在 或 插入时唯一索引约束错误[并发时 数据已写入]
		if ok || b {
			// 更新库存
			affected, err := svc.UpdateStock(date, item.Count, now)
			if err != nil {
				log.Log.Errorf("venue_trace: update stock fail, err:%s", err)
				return err
			}

			log.Log.Errorf("venue_trace: affected:%d", affected)

			// 更新未成功 库存不够 || 当前预约库存足够 但已出现库存不足的预约时间点 则需返回最新各预约节点的剩余库存
			if affected == 0 || affected == 1 && !svc.Extra.IsEnough {
				svc.Extra.IsEnough = false
				// 查看最新的库存 并返回 tips：读快照无问题 购买时保证数据一致性即可
				resp, err := svc.SetLatestInventoryResp(date)
				if err != nil {
					log.Log.Errorf("venue_trace:set inventory resp fail, err:%s", err)
					return err
				}

				svc.Extra.TimeNodeInfo = append(svc.Extra.TimeNodeInfo, resp)
				continue
			}
		}

		// 是否遍历完所有预约节点
		//if index < len(infos) - 1 {
		//	continue
		//}
	}

	//// 库存不足 返回最新库存
	//if !isEnough {
	//	return stockInfo, totalAmount, errors.New("not enough stock")
	//}

	return nil
}


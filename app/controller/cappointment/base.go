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
	"sports_service/server/models/mvenue"
	"sports_service/server/util"
	"strings"
	"time"
)

type base struct {
	Engine      *xorm.Session
	appointment *mappointment.AppointmentModel
	order       *morder.OrderModel
	DateId      int
	Extra       *mappointment.OrderResp
	// 预约流水map
	recordMp    map[int64]*models.VenueAppointmentRecord
	// 订单商品流水map
	orderMp     map[int64]*models.VenueOrderProductInfo
	venue       *mvenue.VenueModel
}

func New(socket *xorm.Session) *base {
	return &base{
		Engine: socket,
		appointment: mappointment.NewAppointmentModel(socket),
		order: morder.NewOrderModel(socket),
		Extra:  &mappointment.OrderResp{TimeNodeInfo: make([]*mappointment.TimeNodeInfo, 0)},
		recordMp: make(map[int64]*models.VenueAppointmentRecord),
		orderMp: make(map[int64]*models.VenueOrderProductInfo),
		venue:   mvenue.NewVenueModel(socket),
	}
}

func (svc *base) AppointmentDateInfo(days, appointmentType int) interface{} {
	list := svc.GetAppointmentDate(days)
	res := make([]*mappointment.WeekInfo, len(list))
	id := 0
	week := -1
	date := ""
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
		svc.appointment.AppointmentInfo.TimeNode = ""

		if err := svc.GetMinPriceByWeek(appointmentType); err != nil {
			log.Log.Errorf("venue_trace: get min price fail, err:%s", err)
		}

		info.MinPrice = svc.appointment.AppointmentInfo.CurAmount
		info.PriceCn = fmt.Sprintf("¥%.2f", float64(info.MinPrice)/100)

		condition, err := svc.GetQueryCondition()
		if err != nil {
			log.Log.Errorf("venue_trace: get query condition fail, err:%s", err)
			return nil
		}

		// 预约大课的时间节点选项 过期的节点会在服务端直接过滤掉 所以日期配置展示最低价格时 也应排除过期的节点
		if appointmentType == consts.APPOINTMENT_COURSE && v.Id == 1 {
			date := svc.GetDateById(v.Id, consts.FORMAT_DATE)
			svc.SetWeek(v.Week)
			svc.SetAppointmentType(appointmentType)
			list, err := svc.GetAppointmentOptions(condition)
			if len(list) > 0 && err == nil {
				for _, opts := range list {
					_, _, hasExpire := svc.TimeNodeHasExpire(date, opts.TimeNode)
					// 已过期
					if hasExpire {
						info.MinPrice = 0
						info.PriceCn = "¥0"
					} else {
						if info.MinPrice > opts.CurAmount || info.MinPrice == 0 {
							info.MinPrice = opts.CurAmount
							info.PriceCn = fmt.Sprintf("¥%.2f", float64(info.MinPrice)/100)
						}
					}
				}
			}
		}

		total, err := svc.appointment.GetTotalNodeByWeek(condition)
		if err != nil {
			log.Log.Errorf("venue_trace: get total node by week fail, err:%s", err)
		}
		info.Total = total

		if id == 0 && info.Total > 0 {
			id = v.Id
			date = v.Date
		}

		if week == -1 && info.Total > 0 {
			week = v.Week
		}

		res[index] = info
	}

	dateInfo := &mappointment.DateInfo{
		List: res,
		Id: id,
		Week: week,
		WeekCn: util.GetWeekCn(week),
		DateCn: date,
	}

	return dateInfo
}

// 根据预约类型 获取最低价格
func (svc *base) GetMinPriceByWeek(appointmentType int) error {
	switch appointmentType {
	case consts.APPOINTMENT_VENUE:
		return svc.appointment.GetVenueMinPriceByWeek()
	case consts.APPOINTMENT_COACH:
		return svc.appointment.GetCoachCourseMinPriceByWeek()
	case consts.APPOINTMENT_COURSE:
		return svc.appointment.GetCourseMinPriceByWeek()
	}

	return errors.New("unsupported appointment type")
}

// 获取查询条件
func (svc *base) GetQueryCondition() (string, error) {
	var condition string
	switch svc.appointment.AppointmentInfo.AppointmentType {
	case consts.APPOINTMENT_VENUE:
		condition = fmt.Sprintf("venue_id=%d AND week_num=%d AND appointment_type=%d AND status=0",
			svc.appointment.AppointmentInfo.VenueId, svc.appointment.AppointmentInfo.WeekNum,
			svc.appointment.AppointmentInfo.AppointmentType)
	case consts.APPOINTMENT_COACH:
		condition = fmt.Sprintf("course_id=%d AND coach_id=%d AND week_num=%d AND appointment_type=%d AND status=0",
			svc.appointment.AppointmentInfo.CourseId, svc.appointment.AppointmentInfo.CoachId, svc.appointment.AppointmentInfo.WeekNum,
			svc.appointment.AppointmentInfo.AppointmentType)
	case consts.APPOINTMENT_COURSE:
		condition = fmt.Sprintf("course_id=%d AND week_num=%d AND appointment_type=%d AND status=0",
			svc.appointment.AppointmentInfo.CourseId, svc.appointment.AppointmentInfo.WeekNum,
			svc.appointment.AppointmentInfo.AppointmentType)
	default:
		return "", errors.New("invalid appointmentType")
	}

	return condition, nil
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

// 预约 场馆/私教/大课 选项
func (svc *base) GetAppointmentOptions(condition string) ([]*models.VenueAppointmentInfo, error) {
	//var condition string
	//switch svc.appointment.AppointmentInfo.AppointmentType {
	//case consts.APPOINTMENT_VENUE:
	//	condition = fmt.Sprintf("venue_id=%d AND week_num=%d AND appointment_type=%d AND status=0",
	//		svc.appointment.AppointmentInfo.VenueId, svc.appointment.AppointmentInfo.WeekNum,
	//		svc.appointment.AppointmentInfo.AppointmentType)
	//case consts.APPOINTMENT_COACH:
	//	condition = fmt.Sprintf("course_id=%d AND coach_id=%d AND week_num=%d AND appointment_type=%d AND status=0",
	//		svc.appointment.AppointmentInfo.CourseId, svc.appointment.AppointmentInfo.CoachId, svc.appointment.AppointmentInfo.WeekNum,
	//		svc.appointment.AppointmentInfo.AppointmentType)
	//case consts.APPOINTMENT_COURSE:
	//	condition = fmt.Sprintf("course_id=%d AND week_num=%d AND appointment_type=%d AND status=0",
	//		svc.appointment.AppointmentInfo.CourseId, svc.appointment.AppointmentInfo.WeekNum,
	//		svc.appointment.AppointmentInfo.AppointmentType)
	//default:
	//	return nil, errors.New("invalid appointmentType")
	//}
	//
	list, err := svc.appointment.GetOptionsByWeek(condition)
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

func (svc *base) SetVenueId(venueId int) {
	svc.appointment.AppointmentInfo.VenueId = int64(venueId)
}

func (svc *base) SetCoachId(coachId int) {
	svc.appointment.AppointmentInfo.CoachId = int64(coachId)
}

func (svc *base) SetCourseId(courseId int) {
	svc.appointment.AppointmentInfo.CourseId = int64(courseId)
}

func (svc *base) SetAppointmentType(appointmentType int) {
	svc.appointment.AppointmentInfo.AppointmentType = appointmentType
}

//func (svc *base) SetStockRelatedId(relatedId int64) {
//	svc.appointment.Stock.RelatedId = relatedId
//}

func (svc *base) SetStockCoachId(coachId int64) {
	svc.appointment.Stock.CoachId = coachId
}

func (svc *base) SetStockCourseId(courseId int64) {
	svc.appointment.Stock.CourseId = courseId
}

func (svc *base) SetStockVenueId(venueId int64) {
	svc.appointment.Stock.VenueId = venueId
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
	var isExpire bool
	var start, end int64
	if svc.DateId == 1 {
		startTm, endTm, hasExpire := svc.TimeNodeHasExpire(date, item.TimeNode)
		start = startTm
		end = endTm
		if hasExpire {
			// 预约大课 如果当前时间 > 配置中的开始时间 过滤该配置项
			if item.AppointmentType == consts.APPOINTMENT_COURSE {
				log.Log.Errorf("过滤id：%d, date:%s", item.Id, date)
				return nil
			}

			// 预约场馆/ 预约私教 则给标示
			isExpire = true
		}

		//nodes := strings.Split(item.TimeNode, "-")
		//if len(nodes) ==  2 {
		//	now := time.Now().Unix()
		//	// 获取 预约时间配置的开始时间
		//	start := fmt.Sprintf("%s %s", date, nodes[0])
		//	ts := new(util.TimeS)
		//	startTm := ts.GetTimeStrOrStamp(start, "YmdHi")
		//	tm = startTm.(int64)
		//	if now > startTm.(int64) {
		//		// 预约大课 如果当前时间 > 配置中的开始时间 过滤该配置项
		//		if item.AppointmentType == consts.APPOINTMENT_COURSE {
		//			log.Log.Errorf("过滤id：%d, date:%s", item.Id, date)
		//			return nil
		//		}
		//
		//		// 预约场馆/ 预约私教 则给标示
		//		isExpire = true
		//	}
		//}
	}

	info := &mappointment.OptionsInfo{
		CurAmount: item.CurAmount,
		TimeNode: item.TimeNode,
		DurationCn: util.ResolveTime(item.Duration),
		Duration: item.Duration,
		RealAmount: item.RealAmount,
		QuotaNum: item.QuotaNum,
		RecommendType: item.RecommendType,
		AppointmentType: item.AppointmentType,
		WeekNum: item.WeekNum,
		AmountCn: fmt.Sprintf("%.2f", float64(item.CurAmount)/100),
		Id: item.Id,
		IsExpire: isExpire,
		StartTm: start,
		EndTm: end,
		Date: fmt.Sprintf("%s %s", date, item.TimeNode),
		CoachId: item.CoachId,
	}

	switch info.AppointmentType {
	case consts.APPOINTMENT_VENUE:
		info.RelatedId = item.VenueId
	case consts.APPOINTMENT_COACH, consts.APPOINTMENT_COURSE:
		info.RelatedId = item.CourseId
	}

	// 售价 < 定价 表示有优惠
	if item.CurAmount < item.RealAmount {
		info.DiscountRate = item.DiscountRate
		info.DiscountAmount = item.DiscountAmount
		if item.DiscountAmount > 0 {
			info.HasDiscount = 1
			if item.DiscountRate % 10 == 0 {
				info.RateCn = fmt.Sprintf("%d折券", item.DiscountRate/10)
			} else {
				info.RateCn = fmt.Sprintf("%.1f折券", float64(item.DiscountRate)/10)
			}
		}
	}

	svc.appointment.AppointmentInfo = item
	ok, err := svc.QueryStockInfo(date)
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
func (svc *base) QueryStockInfo(date string) (bool, error) {
	//svc.SetAppointmentType(appointmentType)
	//svc.SetStockRelatedId(relatedId)
	//svc.SetStockDate(date)
	//svc.SetStockTimeNode(timeNode)
	condition, err := svc.GetQueryStockCondition(date)
	if err != nil {
		return false, err
	}

	return svc.appointment.GetStockInfo(condition)
}

func (svc *base) SetOrderProductInfo(orderId string, now, count int, productId int64) *models.VenueOrderProductInfo {
	return &models.VenueOrderProductInfo{
		ProductId:   productId,
		ProductType: svc.Extra.OrderType,
		ProductCategory: consts.PRODUCT_CATEGORY_APPOINTMENT,
		Count:       count,
		RealAmount:  svc.appointment.AppointmentInfo.RealAmount,
		CurAmount:   svc.appointment.AppointmentInfo.CurAmount,
		DiscountRate: svc.appointment.AppointmentInfo.DiscountRate,
		DiscountAmount: svc.appointment.AppointmentInfo.DiscountAmount,
		VenueId: svc.appointment.AppointmentInfo.VenueId,
		Amount: svc.appointment.AppointmentInfo.CurAmount * count,
		CreateAt: now,
		UpdateAt: now,
		PayOrderId: orderId,
	}
}

func (svc *base) SetAppointmentRecordInfo(userId, date, orderId string, now, count int, seatInfo []*mappointment.SeatInfo,
	startTm, endTm int64) *models.VenueAppointmentRecord {
	info, _ := util.JsonFast.MarshalToString(seatInfo)
	return &models.VenueAppointmentRecord{
		UserId: userId,
		VenueId: svc.appointment.AppointmentInfo.VenueId,
		CourseId: svc.appointment.AppointmentInfo.CourseId,
		AppointmentType: svc.appointment.AppointmentInfo.AppointmentType,
		TimeNode: svc.appointment.AppointmentInfo.TimeNode,
		Date: date,
		PayOrderId: orderId,
		CreateAt: now,
		UpdateAt: now,
		PurchasedNum: count,
		SeatInfo: info,
		CoachId: svc.Extra.CoachId,
		SingleDuration: svc.appointment.AppointmentInfo.Duration,
		Duration: svc.appointment.AppointmentInfo.Duration * count,
		UnitDuration: svc.appointment.AppointmentInfo.UnitDuration,
		UnitPrice: svc.appointment.AppointmentInfo.UnitPrice,
		StartTm: int(startTm),
		EndTm: int(endTm),
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
		VenueId:         svc.appointment.AppointmentInfo.VenueId,
		CourseId:        svc.appointment.AppointmentInfo.CourseId,
		CoachId:         svc.appointment.AppointmentInfo.CoachId,
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
	switch svc.appointment.AppointmentInfo.AppointmentType {
	case consts.APPOINTMENT_VENUE:
		// 更新库存
		return svc.appointment.UpdateVenueStockInfo(svc.appointment.AppointmentInfo.TimeNode, date,
			svc.appointment.AppointmentInfo.QuotaNum, count, now, svc.appointment.AppointmentInfo.AppointmentType,
			int(svc.appointment.AppointmentInfo.VenueId))
	case consts.APPOINTMENT_COACH, consts.APPOINTMENT_COURSE:
		return svc.appointment.UpdateCourseStockInfo(svc.appointment.AppointmentInfo.TimeNode, date,
			svc.appointment.AppointmentInfo.QuotaNum, count, now, svc.appointment.AppointmentInfo.AppointmentType,
			int(svc.appointment.AppointmentInfo.VenueId), int(svc.appointment.AppointmentInfo.CoachId),
			int(svc.appointment.AppointmentInfo.CourseId))
	}

	return 0, errors.New("invalid appointmentType")
}

// 添加预约流水
func (svc *base) AddAppointmentRecord() error {
	for _, val := range svc.recordMp {
		// 实付金额为0 表示使用时长抵扣 或 活动免费 直接置为可用
		//if svc.order.Order.Amount == 0 {
		//	val.Status = 1
		//}

		affected, err := svc.appointment.AddAppointmentRecord(val)
		if err != nil {
			return err
		}

		if affected != 1 {
			return errors.New("add record fail, count not match~")
		}
	}

	return nil
}

// 添加订单商品流水记录
func (svc *base) AddOrderProducts() error {
	var olst []*models.VenueOrderProductInfo
	for _, val := range svc.orderMp {
		// 实付金额为0 表示使用时长抵扣 或 活动免费 直接置为可用
		if svc.order.Order.Amount == 0 {
			val.Status = consts.ORDER_TYPE_PAID
		}

		val.SnapshotId = svc.recordMp[val.ProductId].Id
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
func (svc *base) AddOrder(orderId, userId, subject string, now, productType int) error {
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	svc.Extra.CreateAt = time.Unix(int64(now), 0).In(cstSh).Format(consts.FORMAT_TM)
	extra, _ := util.JsonFast.MarshalToString(svc.Extra)
	svc.order.Order.Extra = extra
	svc.order.Order.PayOrderId = orderId
	svc.order.Order.UserId = userId
	svc.order.Order.OrderType = 1001
	svc.order.Order.CreateAt = now
	svc.order.Order.UpdateAt = now
	svc.order.Order.Amount = svc.Extra.TotalAmount
	svc.order.Order.OriginalAmount = svc.Extra.OriginalAmount
	svc.order.Order.ChannelId = svc.Extra.Channel
	svc.order.Order.Subject = subject
	svc.order.Order.ProductType = productType
	svc.order.Order.WriteOffCode = svc.Extra.WriteOffCode
	// todo: 暂时只有一个场馆
	svc.order.Order.VenueId = 1
	// 实付金额为0 表示使用时长抵扣 或 活动免费  订单直接置为成功
	if svc.order.Order.Amount == 0 {
		svc.order.Order.Status = consts.ORDER_TYPE_PAID
	}

	affected, err := svc.order.AddOrder()
	if err != nil {
		return err
	}

	if affected != 1 {
		log.Log.Errorf("venue_trace: add order fail, err:%s", err)
		return errors.New("add order fail, affected not 1")
	}

	return nil
}

// 预约的时间节点是否过期 true 表示节点已过期
func (svc *base) TimeNodeHasExpire(date, timeNode string) (int64, int64, bool) {
	nodes := strings.Split(timeNode, "-")
	if len(nodes) == 2 {
		now := time.Now().Unix()
		// 获取 预约时间配置的开始时间
		start := fmt.Sprintf("%s %s", date, nodes[0])
		end := fmt.Sprintf("%s %s", date, nodes[1])
		ts := new(util.TimeS)
		startTm := ts.GetTimeStrOrStamp(start, "YmdHi")
		endTm := ts.GetTimeStrOrStamp(end, "YmdHi")
		// 如果当前时间 > 配置中的开始时间 走库存不足流程 count置为0
		if now > startTm.(int64) {
			return startTm.(int64), endTm.(int64), true
		}

		return startTm.(int64), endTm.(int64), false
	}

	return 0, 0, false
}

// 设置最新库存数据（返回使用）
func (svc *base) SetLatestInventoryResp(date string, count int, isEnough bool) (*mappointment.TimeNodeInfo, error) {
	info := &mappointment.TimeNodeInfo{}
	info.Id = svc.appointment.AppointmentInfo.Id
	info.TimeNode = svc.appointment.AppointmentInfo.TimeNode
	info.Date = date
	info.Amount = svc.appointment.AppointmentInfo.CurAmount
	info.Discount = svc.appointment.AppointmentInfo.DiscountAmount

	// 节点是否已过期
	startTm, endTm, hasExpire := svc.TimeNodeHasExpire(info.Date, info.TimeNode)
	info.StartTm = startTm
	info.EndTm = endTm
	if hasExpire {
		// 节点过期 走库存不足流程 同时正常返回购物车数据 客户端清除库存不足的节点即可
		info.Count = 0
		info.IsEnough = false
		svc.Extra.IsEnough = false
		return info, nil
	}

	// 查看最新的库存 并返回 tips：读快照无问题 购买时保证数据一致性即可
	ok, err := svc.QueryStockInfo(date)
	if !ok || err != nil {
		log.Log.Errorf("venue_trace: get stock info fail, err:%s", err)
		//return nil, errors.New("get stock info fail")
		// 如果查询失败 返回当前数量
		info.Count = count
	} else {
		log.Log.Infof("venue_trace: info:%+v, quotaNum:%d, purchasedNum:%d, count:%d", info,
			svc.appointment.Stock.QuotaNum, svc.appointment.Stock.PurchasedNum, count)
		// 默认当前节点库存足够
		info.Count = count
		// 如果不够 更新库存失败了
		if !isEnough {
			// 可购数量 = 总库存 - 已购买
			info.Count = svc.appointment.Stock.QuotaNum - svc.appointment.Stock.PurchasedNum
		}
	}

	return info, nil
}

// 预约流程
func (svc *base) AppointmentProcess(userId, orderId string, relatedId int64, weekNum int, labelIds []interface{},
	infos []*mappointment.AppointmentInfo) error {

	if len(infos) == 0 {
		log.Log.Errorf("venue_trace: invalid params, infos len:%d", len(infos))
		return errors.New("invalid params")
	}

	svc.Extra.IsEnough = true
	now := int(time.Now().Unix())
	//stockInfo := make([]*mappointment.StockInfoResp, 0)
	//// 订单商品流水map
	//orderMp := make(map[int64]*models.VenueOrderProductInfo)
	//// 预约流水map
	//recordMp := make(map[int64]*models.VenueAppointmentRecord)
	for _, item := range infos {
		log.Log.Infof("params:%+v", item)
		err := svc.GetAppointmentConf(item.Id)
		if svc.appointment.AppointmentInfo == nil || err != nil {
			log.Log.Errorf("venue_trace: get appointment conf fail, id:%d", item.Id)
			return errors.New("get appointment conf fail")
		}

		if svc.appointment.AppointmentInfo.AppointmentType == consts.APPOINTMENT_VENUE {
			if len(item.SeatInfos) != item.Count {
				log.Log.Errorf("venue_trace: seat num not match, count:%d, len:%d", item.Count, len(item.SeatInfos))
				return errors.New("seat num not match")
			}
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

		tmInfo, err := time.Parse(consts.FORMAT_DATE, date)
		if err != nil {
			log.Log.Errorf("venue_trace: time.Parse fail, err:%s", err)
			return err
		}

		week := int(tmInfo.Weekday())
		if week != weekNum {
			log.Log.Errorf("venue_trace: week not match, week:%d, clientWeek:%d", week, weekNum)
			return errors.New("week not match")
		}

		if svc.Extra.Date == "" {
			svc.Extra.Date = date
		}

		if svc.appointment.AppointmentInfo.AppointmentType == 0 {
			var list []*models.VenuePersonalLabelConf
			labelType := 0
			if len(labelIds) > 0 {
				list, err = svc.appointment.GetLabelsByIds(labelIds)
				if err != nil || list == nil {
					log.Log.Errorf("venue_trace: get labels by ids fail, err:%s", err)
				}
			}

			if len(list) == 0 || len(labelIds) == 0 {
				list, err = svc.appointment.GetLabelsByRand()
				if err != nil || list == nil {
					log.Log.Errorf("venue_trace: get labels by rand fail, err:%s", err)
					return errors.New("get labels by rand fail")
				}

				labelType = 1

			}

			if len(list) > 0 {
				if err := svc.AddLabels(list, date, userId, orderId, relatedId, labelType); err != nil {
					log.Log.Errorf("venue_trace: add labels fail, err:%s", err)
					return err
				}
			}
		}

		svc.Extra.OrderType, err = svc.GetOrderType(svc.appointment.AppointmentInfo.AppointmentType)
		if err != nil {
			return err
		}

		svc.Extra.TotalTm += svc.appointment.AppointmentInfo.Duration * item.Count
		// 数量 * 售价
		svc.Extra.TotalAmount += item.Count * svc.appointment.AppointmentInfo.CurAmount
		svc.Extra.TotalDiscount += item.Count * svc.appointment.AppointmentInfo.DiscountAmount
		svc.Extra.OriginalAmount += item.Count * svc.appointment.AppointmentInfo.RealAmount
		// 购买总数量
		svc.Extra.Count += item.Count
		svc.orderMp[item.Id] = svc.SetOrderProductInfo(orderId, now, item.Count, item.Id)
		//svc.recordMp[item.Id] = svc.SetAppointmentRecordInfo(userId, date, orderId, now, item.Count, item.SeatInfos, relatedId)

		condition, err := svc.GetQueryStockCondition(date)
		if err != nil {
			log.Log.Errorf("venue_trace: get query stock condition fail, err:%s", err)
			return err
		}

		ok, err := svc.appointment.HasExistsStockInfo(condition)
		if err != nil {
			log.Log.Errorf("venue_trace: query stock info fail, err:%s", err)
			return err
		}


		var b bool
		// 默认该节点库存足够
		item.IsEnough = true
		if !ok {
			b, err = svc.AddStock(date, now, item.Count)
			log.Log.Errorf("venue_trace: id:%d, b:%v, err:%s", svc.appointment.AppointmentInfo.Id, b, err)
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

			log.Log.Errorf("venue_trace: affected:%d, id:%d, node:%s", affected, item.Id,
				svc.appointment.AppointmentInfo.TimeNode)

			// 更新未成功 库存不够 || 当前预约库存足够 但已出现库存不足的预约时间点 则需返回最新各预约节点的剩余库存
			//if affected == 0 || affected == 1 && !svc.Extra.IsEnough {
			//	svc.Extra.IsEnough = false
			//	//continue
			//}

			if affected == 0 {
				svc.Extra.IsEnough = false
				// 该节点库存不够
				item.IsEnough = false
			}

		}

		// 查看最新的库存 并返回 tips：读快照无问题 购买时保证数据一致性即可
		resp, err := svc.SetLatestInventoryResp(date, item.Count, item.IsEnough)
		if err != nil {
			log.Log.Errorf("venue_trace:set inventory resp fail, err:%s", err)
			return err
		}

		// 设置预约快照信息
		svc.recordMp[item.Id] = svc.SetAppointmentRecordInfo(userId, date, orderId, now, item.Count, item.SeatInfos,
			resp.StartTm, resp.EndTm)
		svc.Extra.TimeNodeInfo = append(svc.Extra.TimeNodeInfo, resp)

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

func (svc *base) GetQueryStockCondition(date string) (string, error) {
	var condition string
	switch svc.appointment.AppointmentInfo.AppointmentType {
	case consts.APPOINTMENT_VENUE:
		condition = fmt.Sprintf("appointment_type=%d AND venue_id=%d AND date='%s' AND time_node='%s'",
			svc.appointment.AppointmentInfo.AppointmentType, svc.appointment.AppointmentInfo.VenueId, date,
			svc.appointment.AppointmentInfo.TimeNode)
	case consts.APPOINTMENT_COACH,consts.APPOINTMENT_COURSE:
		condition = fmt.Sprintf("venue_id=%d AND appointment_type=%d AND course_id=%d AND coach_id=%d AND date='%s' AND time_node='%s'",
			svc.appointment.AppointmentInfo.VenueId, svc.appointment.AppointmentInfo.AppointmentType, svc.appointment.AppointmentInfo.CourseId,
			svc.appointment.AppointmentInfo.CoachId, date, svc.appointment.AppointmentInfo.TimeNode)
	default:
		return "", errors.New("invalid appointmentType")
	}

	log.Log.Infof("venue_trace: condition:%s", condition)
	return condition, nil
}

// labelType 0 表示用户添加 1 表示系统添加
func (svc *base) AddLabels(list []*models.VenuePersonalLabelConf, date, userId, orderId string, relatedId int64, labelType int) error {
	labels := make([]*models.VenueUserLabel, len(list))
	for k, v := range list {
		label := &models.VenueUserLabel{
			Date: date,
			TimeNode: svc.appointment.AppointmentInfo.TimeNode,
			UserId: userId,
			LabelId: v.Id,
			LabelName: v.LabelName,
			VenueId: relatedId,
			LabelType: labelType,
			PayOrderId: orderId,
		}

		labels[k] = label
	}

	affected, err := svc.appointment.AddLabels(labels)
	if err != nil || affected != int64(len(list)) {
		log.Log.Errorf("venue_trace: add labels failm err:%s, affected:%d, len:%d", err, affected, len(list))
		return errors.New("add labels fail")
	}

	return nil
}

// 获取订单类型
func (svc *base) GetOrderType(appointmentType int) (int, error) {
	switch appointmentType {
	case 0:
		return consts.ORDER_TYPE_APPOINTMENT_VENUE, nil
	case 1:
		return consts.ORDER_TYPE_APPOINTMENT_COACH, nil
	case 2:
		return consts.ORDER_TYPE_APPOINTMENT_COURSE, nil
	}

	return -1, errors.New("invalid appointmentType")
}

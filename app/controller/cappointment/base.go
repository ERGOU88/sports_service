package cappointment

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mappointment"
	"sports_service/server/util"
	"strconv"
	"time"
)

type base struct {
	Engine  *xorm.Session
	appointment *mappointment.AppointmentModel
	DateId  int
}

func New(socket *xorm.Session) *base {
	return &base{
		Engine: socket,
		appointment: mappointment.NewAppointmentModel(socket),
	}
}

func (svc *base) AppointmentDateInfo(days int) interface{} {
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
		svc.appointment.AppointmentInfo.WeekNum = v.Week
		svc.appointment.AppointmentInfo.AppointmentType = 0
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

func (svc *base) SetStockDate(date int) {
	svc.appointment.Stock.Date = int64(date)
}

func (svc *base) SetStockTimeNode(timeNode string) {
	svc.appointment.Stock.TimeNode = timeNode
}

// 日期id
func (svc *base) SetDateId(id int) {
	svc.DateId = id
}

// 通过id获取日期
func (svc *base) GetDateById(id int) int {
	var date string
	curTime := time.Now()

	if id >= 1 {
		date = curTime.AddDate(0, 0, id - 1).Format(consts.FORMAT_DATE)
	}

	res, err := strconv.Atoi(date)
	if err != nil {
		return 0
	}

	return res
}

func (svc *base) SetAppointmentOptionsRes(date int, item *models.VenueAppointmentInfo) *mappointment.OptionsInfo {
	info := &mappointment.OptionsInfo{
		RelatedId: item.RelatedId,
		CurAmount: item.CurAmount,
		TimeNode: item.TimeNode,
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
	}

	svc.SetStockRelatedId(item.RelatedId)
	svc.SetStockDate(date)
	svc.SetStockTimeNode(item.TimeNode)
	ok, err := svc.appointment.GetPurchaseNum()
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

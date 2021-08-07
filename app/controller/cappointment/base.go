package cappointment

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
	"sports_service/server/models/mappointment"
	"sports_service/server/util"
	"time"
	"sports_service/server/global/app/log"
	"fmt"
)

type base struct {
	Engine  *xorm.Session
	appointment *mappointment.AppointmentModel
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

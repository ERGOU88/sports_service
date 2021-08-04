package mappointment

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
	"sports_service/server/util"
	"time"
)

type AppointmentModel struct {
	AppointmentInfo  *models.VenueAppointmentInfo
	Engine           *xorm.Session
}

func NewAppointmentModel(engine *xorm.Session) *AppointmentModel {
	return &AppointmentModel{
		AppointmentInfo: new(models.VenueAppointmentInfo),
		Engine: engine,
	}
}

// 获取预约的日期信息（从当天开始推算）
// days 天数
func (m *AppointmentModel) GetAppointmentDate(days int) []*util.DateInfo {
	curTime := time.Now()
	// 今天
	today := curTime.Format("2006-01-02")
	// 往后推6天 总共7天
	afterDay := curTime.AddDate(0, 0, days).Format("2006-01-02")
	dateInfo := util.GetBetweenDates(today, afterDay)

	return dateInfo
}

// 根据星期 及 预约类型 获取最低价格
func (m *AppointmentModel) GetMinPriceByWeek() {
}

// 通过星期 及 预约类型 获取可预约选项
func (m *AppointmentModel) GetOptionsByWeek() ([]*models.VenueAppointmentInfo, error) {
	var list []*models.VenueAppointmentInfo
	if err := m.Engine.Where("week_num=? AND appointment_type=? AND status=0", m.AppointmentInfo.WeekNum,
		m.AppointmentInfo.AppointmentType).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

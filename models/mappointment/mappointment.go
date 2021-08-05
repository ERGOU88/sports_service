package mappointment

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

type AppointmentModel struct {
	AppointmentInfo  *models.VenueAppointmentInfo
	Engine           *xorm.Session
}

type WeekInfo struct {
	Date      string    `json:"date"`
	Week      int       `json:"week"`
	Id        int32     `json:"id"`
	WeekCn    string    `json:"week_cn"`
	MinPrice  int       `json:"min_price"`
	PriceCn   string    `json:"price_cn"`
}

func NewAppointmentModel(engine *xorm.Session) *AppointmentModel {
	return &AppointmentModel{
		AppointmentInfo: new(models.VenueAppointmentInfo),
		Engine: engine,
	}
}

const (
	QUERY_MIN_PRICE = "SELECT min(cur_amount) as cur_amount FROM venue_appointment_info WHERE week_num=? AND appointment_type=? AND status=0"
)
// 根据星期 及 预约类型 获取最低价格
func (m *AppointmentModel) GetMinPriceByWeek() error {
	ok, err := m.Engine.SQL(QUERY_MIN_PRICE, m.AppointmentInfo.WeekNum, m.AppointmentInfo.AppointmentType).Get(m.AppointmentInfo)
	if !ok || err != nil {
		return err
	}

	return nil
}

// 通过场馆id、私教id、星期 及 预约类型 获取可预约选项
func (m *AppointmentModel) GetOptionsByWeek() ([]*models.VenueAppointmentInfo, error) {
	var list []*models.VenueAppointmentInfo
	if err := m.Engine.Where("related_id=？ AND week_num=? AND appointment_type=? AND status=0", m.AppointmentInfo.RelatedId,
		m.AppointmentInfo.WeekNum, m.AppointmentInfo.AppointmentType).Asc("id").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

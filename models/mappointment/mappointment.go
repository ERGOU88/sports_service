package mappointment

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

type AppointmentModel struct {
	AppointmentInfo  *models.VenueAppointmentInfo
	Engine           *xorm.Session
	Stock            *models.VenueAppointmentStock
}

type WeekInfo struct {
	Date      string    `json:"date"`
	Week      int       `json:"week"`
	Id        int32     `json:"id"`
	WeekCn    string    `json:"week_cn"`
	MinPrice  int       `json:"min_price"`
	PriceCn   string    `json:"price_cn"`
}

type OptionsInfo struct {
	Id              int64  `json:"id"`
	TimeNode        string `json:"time_node"`
	Duration        int    `json:"duration"`
	RealAmount      int    `json:"real_amount"`
	CurAmount       int    `json:"cur_amount"`
	DiscountRate    int    `json:"discount_rate"`
	DiscountAmount  int    `json:"discount_amount"`
	QuotaNum        int    `json:"quota_num"`
	RelatedId       int64  `json:"related_id"`
	RecommendType   int    `json:"recommend_type"`
	AppointmentType int    `json:"appointment_type"`
	WeekNum         int    `json:"week_num"`

	HasDiscount     int    `json:"has_discount,omitempty"`       // 是否有优惠 0无 1有
	AmountCn        string `json:"amount_cn,omitempty"`          // 中文价格
	IsFull          int    `json:"is_full"`                      // 是否满场
	PurchasedNum    int    `json:"purchased_num,omitempty"`      // 已购买人数 包含[成功购买及已下单]
    Name            string `json:"name,omitempty"`               // 场馆名称
    Avatar          string `json:"avatar,omitempty"`             // 大课老师头像
    Labels          []*LabelInfo `json:"labels,omitempty"`       // 标签列表
}

type LabelInfo struct {
	UserId      string       `json:"user_id"`
	NickName    string       `json:"nick_name"`
	Avatar      string       `json:"avatar"`
	LabelId     int64        `json:"label_id"`
	LabelName   string       `json:"label_name"`
}

type Options struct {
	Id       int64        `json:"id"`
	Name     string       `json:"name"`
}

func NewAppointmentModel(engine *xorm.Session) *AppointmentModel {
	return &AppointmentModel{
		AppointmentInfo: new(models.VenueAppointmentInfo),
		Stock: new(models.VenueAppointmentStock),
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

// 通过场馆id、课程id、星期 及 预约类型 获取可预约选项
func (m *AppointmentModel) GetOptionsByWeek() ([]*models.VenueAppointmentInfo, error) {
	var list []*models.VenueAppointmentInfo
	if err := m.Engine.Where("related_id=? AND week_num=? AND appointment_type=? AND status=0", m.AppointmentInfo.RelatedId,
		m.AppointmentInfo.WeekNum, m.AppointmentInfo.AppointmentType).Asc("id").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取某时间点 场馆预约人数 包含已成功及已下单且订单未超时
func (m *AppointmentModel) GetPurchaseNum() (bool, error) {
	return m.Engine.Get(m.Stock)
}

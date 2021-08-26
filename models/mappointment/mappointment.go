package mappointment

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

type AppointmentModel struct {
	AppointmentInfo  *models.VenueAppointmentInfo
	Engine           *xorm.Session
	Stock            *models.VenueAppointmentStock
	Record           *models.VenueAppointmentRecord
	Vip              *models.VenueVipInfo
	Labels           *models.VenueUserLabel
	LabelConf        *models.VenuePersonalLabelConf
}

type WeekInfo struct {
	Date      string    `json:"date"`
	Week      int       `json:"week"`
	Id        int       `json:"id"`
	WeekCn    string    `json:"week_cn"`
	MinPrice  int       `json:"min_price"`
	PriceCn   string    `json:"price_cn"`
	Total     int64     `json:"total"`
}

// 日期信息
type DateInfo struct {
	List   []*WeekInfo   `json:"week_info"`
	Id      int          `json:"id"`
	Week    int          `json:"week"`
	WeekCn  string       `json:"week_cn"`
	DateCn  string       `json:"date_cn"`
}

type OptionsInfo struct {
	Id              int64  `json:"id"`
	TimeNode        string `json:"time_node"`
	Duration        int    `json:"duration"`
	DurationCn      string `json:"duration_cn"`
	RealAmount      int    `json:"real_amount"`
	RealAmountCn    string `json:"real_amount_cn"`
	CurAmount       int    `json:"cur_amount"`
	DiscountRate    int    `json:"discount_rate"`
	RateCn          string `json:"rate_cn"`
	DiscountAmount  int    `json:"discount_amount"`
	QuotaNum        int    `json:"quota_num"`
	RelatedId       int64  `json:"related_id"`
	RecommendType   int    `json:"recommend_type"`
	RecommendName   string `json:"recommend_name"`
	AppointmentType int    `json:"appointment_type"`
	WeekNum         int    `json:"week_num"`

	HasDiscount     int    `json:"has_discount,omitempty"`       // 是否有优惠 0无 1有
	AmountCn        string `json:"amount_cn,omitempty"`          // 中文价格
	IsFull          int    `json:"is_full"`                      // 是否满场
	PurchasedNum    int    `json:"purchased_num,omitempty"`      // 已购买人数 包含[成功购买及已下单]
	Name            string `json:"name,omitempty"`               // 场馆名称
	Avatar          string `json:"avatar,omitempty"`             // 大课老师头像
	Address         string `json:"address,omitempty"`            // 上课地点
	Labels          []*LabelInfo     `json:"labels,omitempty"`   // 标签列表
	ReservedUsers   []*SeatInfo      `json:"reserved_users"`     // 已预约人数
	IsExpire        bool   `json:"is_expire"`                    // 是否过期
	Date            string `json:"date"`                         // 年月日
	StartTm         int64  `json:"start_tm"`                     // 开始时间戳
	EndTm           int64  `json:"end_tm"`                       // 结束时间戳
	CoachId         int64  `json:"coach_id"`                     // 教练id
}

// 已预约人数
type ReservedUsers struct {
	UserId      string       `json:"user_id"`
	NickName    string       `json:"nick_name"`
	Avatar      string       `json:"avatar"`
}

type LabelInfo struct {
	UserId      string       `json:"user_id"`
	//NickName    string       `json:"nick_name"`
	//Avatar      string       `json:"avatar"`
	LabelId     int64        `json:"label_id"`
	LabelName   string       `json:"label_name"`
}

type Options struct {
	Id                 int64        `json:"id"`
	Name               string       `json:"name"`
	Title              string       `json:"title"`
	Avatar             string       `json:"avatar"`
	Describe           string       `json:"describe,omitempty"`
	CostDescription    string       `json:"cost_description,omitempty"`     // 费用说明
	Instructions       string       `json:"instructions,omitempty"`         // 购买须知
}

type AppointmentReq struct {
	AppointmentType int                 `json:"appointment_type"` // 0场馆 1私教课 2大课
	UserId          string              `json:"user_id"`          // 用户id
	Infos           []*AppointmentInfo  `json:"infos"`            // 预约请求数据 1个时间点对应1条数据
	Ids             []interface{}       `json:"ids"`              // 时间配置ids
	LabelIds        []interface{}       `json:"label_ids"`        // 用户选择的标签id列表
	RelatedId       int64               `json:"related_id"`       // 场馆id/私教课程id/大课id
	CoachId         int64               `json:"coach_id"`         // 老师id
	WeekNum         int                 `json:"week_num"`         // 星期几
	ReqType         int                 `json:"req_type"`         // 1 查询购物车数据 2 下单
	IsDiscount      int32               `json:"is_discount"`      // 是否抵扣时长 1 抵扣 0 不抵扣
	Channel         int                 `json:"channel"`          // 1001 安卓 1002 ios
}

// 预约请求数据
type AppointmentInfo struct {
	DateId      int            `json:"date_id"`    // 1为今天
	//TimeNode  string         `json:"time_node"`  // 预约时间节点
	Count       int            `json:"count"`      // 数量
	Id          int64          `json:"id"`         // 时间配置id
	SeatInfos   []*SeatInfo    `json:"seat_info"`  // 座位号
	IsEnough    bool
}

// 座位信息
type SeatInfo struct {
	UserId      string     `json:"user_id"`
	Avatar      string     `json:"avatar"`
	NickName    string     `json:"nick_name"`
	SeatNo      int        `json:"seat_no"`       // 座位号
}

// 预约请求返回库存数据
type StockInfoResp struct {
	Id        int64      `json:"id"`               // 时间配置id
	Stock     int        `json:"stock"`            // 剩余数量
}

// 预约订单返回数据
type OrderResp struct {
	Id       int64      `json:"id"`                            // 场馆id/商品id
	Name     string     `json:"name,omitempty"`                // 场馆名称
	Date     string     `json:"date,omitempty"`                // 预定日期
	WeekCn   string     `json:"week_cn,omitempty"`             // 星期几
	TotalTm  int        `json:"total_tm,omitempty"`            // 预约总时长
	TmCn     string     `json:"tm_cn,omitempty"`               // 总时长 中文
	IsEnough bool       `json:"is_enough"`                     // 库存是否足够 false 不足 true 足够
	IsDeduct bool       `json:"is_deduct"`                     // 是否可扣除会员时长
	TotalDeductionTm int `json:"total_deduction_tm"`           // 抵扣总时长
	TotalAmount int     `json:"total_amount"`                  // 总金额 真实支付金额
	TotalSalesPrice int  `json:"total_sales_price"`            // 总售价
	MobileNum string    `json:"mobile_num,omitempty"`          // 手机号
	TotalDiscount int   `json:"total_discount"`                // 总优惠
	TimeNodeInfo  []*TimeNodeInfo `json:"node_info,omitempty"` // 多时间节点预约数据
	OrderId       string `json:"order_id,omitempty"`           // 订单ID
	OrderType     int    `json:"order_type"`                   // 订单类型
	PayDuration      int64  `json:"pay_duration"`              // 支付时长
	Address          string `json:"address,omitempty"`         // 上课地点
	CoachName        string `json:"coach_name,omitempty"`      // 老师名称
	CoachId          int64  `json:"coach_id,omitempty"`        // 老师id
	CourseName       string `json:"course_name,omitempty"`     // 课程名称
	CourseId         int64  `json:"course_id,omitempty"`       // 课程id
	Channel          int    `json:"channel,omitempty"`         // 1001 安卓 1002 ios
	WriteOffCode     string `json:"write_off_code,omitempty"`  // 核销码
	CreateAt         string `json:"create_at,omitempty"`       // 下单时间
	Count            int    `json:"count,omitempty"`           // 次卡/月卡/季卡/年卡 数量
	ExpireDuration   int    `json:"expire_duration,omitempty"` // 次卡/月卡/季卡/年卡 过期时长[单个]
	VenueId          int64  `json:"venue_id,omitempty"`
	OrderStatus      int32  `json:"order_status"`                 // 订单状态
	OrderDescription string `json:"order_description,omitempty"` // 订单须知
	CanRefund        bool   `json:"can_refund"`                   // 是否可退款
	HasEvaluate      bool   `json:"has_evaluate,omitempty"`       // 是否已评价
	VenueName        string `json:"venue_name,omitempty"`         // 场馆名称
}

// 单时间节点预约数据
type TimeNodeInfo struct {
	Date         string       `json:"date"`         // 预约的日期
	TimeNode     string       `json:"time_node"`    // 预约时间节点
	Count        int          `json:"count"`        // 数量
	Id           int64        `json:"id"`           // 时间配置id
	Amount       int          `json:"amount"`       // 单价
	DeductionTm  int64        `json:"deduction_tm"` // 抵扣会员时长
	Discount     int          `json:"discount"`     // 优惠的金额
	IsEnough     bool         `json:"is_enough"`    // 当前节点库存是否足够
	StartTm      int64        `json:"start_tm"`     // 预约开始时间戳
	EndTm        int64        `json:"end_tm"`
}

func NewAppointmentModel(engine *xorm.Session) *AppointmentModel {
	return &AppointmentModel{
		AppointmentInfo: new(models.VenueAppointmentInfo),
		Stock: new(models.VenueAppointmentStock),
		Record: new(models.VenueAppointmentRecord),
		Labels: new(models.VenueUserLabel),
		LabelConf: new(models.VenuePersonalLabelConf),
		Vip: new(models.VenueVipInfo),
		Engine: engine,
	}
}

const (
	QUERY_MIN_PRICE = "SELECT min(cur_amount) as cur_amount FROM venue_appointment_info WHERE week_num=? " +
		"AND related_id=? AND appointment_type=? AND status=0"
)
// 根据星期 及 预约类型 获取最低价格
func (m *AppointmentModel) GetMinPriceByWeek() error {
	ok, err := m.Engine.SQL(QUERY_MIN_PRICE, m.AppointmentInfo.WeekNum, m.AppointmentInfo.RelatedId, m.AppointmentInfo.AppointmentType).Get(m.AppointmentInfo)
	if !ok || err != nil {
		return err
	}

	return nil
}

// 获取可预约的时间节点总数
func (m *AppointmentModel) GetTotalNodeByWeek() (int64, error) {
	return m.Engine.Where("week_num=? AND related_id=? AND appointment_type=? AND status=0",
		m.AppointmentInfo.WeekNum, m.AppointmentInfo.RelatedId, m.AppointmentInfo.AppointmentType).
		Count(&models.VenueAppointmentInfo{})
}

// 通过id获取预约配置
func (m *AppointmentModel) GetAppointmentConfById(id string) (bool, error) {
	m.AppointmentInfo = new(models.VenueAppointmentInfo)
	return m.Engine.ID(id).Get(m.AppointmentInfo)
}

// 通过ids[多个id]获取预约配置列表
func (m *AppointmentModel) GetAppointmentConfByIds(ids []interface{}) ([]*models.VenueAppointmentInfo, error) {
	var list []*models.VenueAppointmentInfo
	if err := m.Engine.In("id", ids...).OrderBy("cur_amount DESC").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
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

// 查库存表 获取某时间点 场馆预约人数 包含已成功及已下单且订单未超时
func (m *AppointmentModel) GetStockInfo(appointmentType int, relatedId int64, date, timeNode string) (bool, error) {
	m.Stock = new(models.VenueAppointmentStock)
	return m.Engine.Where("appointment_type=? AND related_id=? AND date=? AND time_node=?", appointmentType,
		relatedId, date, timeNode).Get(m.Stock)
}

// 是否存在库存信息 todo: 前期并发不高 可以读取快照
func (m *AppointmentModel) HasExistsStockInfo(appointmentType int, relatedId int64, date, timeNode string) (bool, error) {
	m.Stock = new(models.VenueAppointmentStock)
	return m.Engine.Where("appointment_type=? AND related_id=? AND date=? AND time_node=?", appointmentType,
		relatedId, date, timeNode).Exist(m.Stock)
}

// 添加预约库存数据[多条]
func (m *AppointmentModel) AddMultiStockInfo(stock []*models.VenueAppointmentStock) (int64, error) {
	return m.Engine.InsertMulti(stock)
}

// 添加预约库存数据[单条]
func (m *AppointmentModel) AddStockInfo(stock *models.VenueAppointmentStock) (int64, error) {
	return m.Engine.Insert(stock)
}

const (
	UPDATE_STOCK_INFO = "UPDATE `venue_appointment_stock` SET `quota_num`=?, `purchased_num`= `purchased_num`+ ?, `update_at`=? WHERE date=? AND " +
		"time_node=? AND appointment_type=? AND related_id=? AND ? >= `purchased_num`+ ? AND `purchased_num` >= 0 LIMIT 1"
)
// 需求：允许动态增加库存 需注意：不可动态减少
func (m *AppointmentModel) UpdateStockInfo(timeNode, date string, quotaNum, count, now, appointmentType, relatedId int) (int64, error) {
	res, err := m.Engine.Exec(UPDATE_STOCK_INFO, quotaNum, count, now, date, timeNode, appointmentType, relatedId, quotaNum, count)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

const (
	REVERT_STOCK_NUM = "UPDATE `venue_appointment_stock` SET `purchased_num`= `purchased_num`+ ?, `update_at`=? WHERE date=? AND " +
		"time_node=? AND appointment_type=? AND related_id=? AND `quota_num` >= `purchased_num`+ ? AND `purchased_num` + ? >= 0 LIMIT 1"
)
// 恢复冻结的库存
func (m *AppointmentModel) RevertStockNum(timeNode, date string, count, now, appointmentType, relatedId int) (int64, error) {
	res, err := m.Engine.Exec(REVERT_STOCK_NUM, count, now, date, timeNode, appointmentType, relatedId, count, count)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

// 获取成功预约的记录[包含已付款和支付中及已完成]
func (m *AppointmentModel) GetAppointmentRecord() ([]*models.VenueAppointmentRecord, error) {
	var list []*models.VenueAppointmentRecord
	sql := "SELECT * FROM venue_appointment_record WHERE status in(0, 2, 3) AND appointment_type=? AND related_id=? AND time_node=? " +
		"AND date=? ORDER BY id ASC"
	if err := m.Engine.SQL(sql, m.Record.AppointmentType, m.Record.RelatedId, m.Record.TimeNode, m.Record.Date).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取场馆会员数据
func (m *AppointmentModel) GetVenueVipInfo(userId string) (*models.VenueVipInfo, error) {
	ok, err := m.Engine.Where("user_id=?", userId).Get(m.Vip)
	if !ok || err != nil {
		return nil, err
	}

	return m.Vip, nil
}

const (
	UPDATE_VENUE_VIP_INFO = "UPDATE `venue_vip_info` SET `duration` = `duration` + ?, `update_at`=? WHERE `user_id`=? " +
		"AND `duration` + ? >= 0 LIMIT 1"
)
// 更新场馆会员数据 duration < 0 减少 duration > 0 增加
func (m *AppointmentModel) UpdateVenueVipInfo(duration int, userId string) (int64, error) {
	res, err := m.Engine.Exec(UPDATE_VENUE_VIP_INFO, duration, userId)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

// 批量添加预约记录
func (m *AppointmentModel) AddMultiAppointmentRecord(list []*models.VenueAppointmentRecord) (int64, error) {
	return m.Engine.InsertMulti(list)
}

const (
	GET_VENUE_USER_LABELS = "SELECT distinct(label_name),label_id FROM venue_user_label WHERE date=? AND time_node=? " +
		"AND status=0 AND venue_id=? LIMIT 3"
)
// 获取场馆用户标签[去重 最多取3条]
func (m *AppointmentModel) GetVenueUserLabels() ([]*models.VenueUserLabel, error) {
	var list []*models.VenueUserLabel
	if err := m.Engine.Where("date=? AND time_node=? AND status=0 AND venue_id=?", m.Labels.Date,
		m.Labels.TimeNode, m.Labels.VenueId).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取所有标签配置
func (m *AppointmentModel) GetUserLabelConf() ([]*models.VenuePersonalLabelConf, error) {
	var list []*models.VenuePersonalLabelConf
	if err := m.Engine.Where("status=0").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 通过id列表获取标签配置
func (m *AppointmentModel) GetLabelsByIds(ids []interface{}) ([]*models.VenuePersonalLabelConf, error) {
	var list []*models.VenuePersonalLabelConf
	if err := m.Engine.Where("status=0").In("id", ids...).OrderBy("id DESC").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

const (
	GET_LABELS_BY_RAND = "SELECT * FROM `venue_personal_label_conf` AS t1 JOIN (SELECT ROUND(RAND() * " +
		"((SELECT MAX(id) FROM `venue_personal_label_conf`)-(SELECT MIN(id) FROM `venue_personal_label_conf`)) " +
		"+ (SELECT MIN(id) FROM `venue_personal_label_conf`)) AS id) AS t2 WHERE t1.id >= t2.id ORDER BY t1.id LIMIT 3;"
)
// 随机获取标签
func (m *AppointmentModel) GetLabelsByRand() ([]*models.VenuePersonalLabelConf, error) {
	var list []*models.VenuePersonalLabelConf
	if err := m.Engine.SQL(GET_LABELS_BY_RAND).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 订单超时 更新标签数据状态
func (m *AppointmentModel) UpdateLabelsStatus(orderId string, status int) (int64, error) {
	return m.Engine.Where("pay_order_id=? AND status=?", orderId, status).Update(m.Labels)
}

// 添加场馆用户标签
func (m *AppointmentModel) AddLabels(labels []*models.VenueUserLabel) (int64, error) {
	return m.Engine.InsertMulti(labels)
}

// 通过订单id获取预约流水
func (m *AppointmentModel) GetAppointmentRecordByOrderId(orderId string, status int) ([]*models.VenueAppointmentRecord, error) {
	var list []*models.VenueAppointmentRecord
	if err := m.Engine.Where("pay_order_id=? AND status=?", orderId, status).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

const (
	UPDATE_RECORD_STATUS = "UPDATE `venue_appointment_record` SET `update_at`=?, `status`=? " +
		"WHERE `pay_order_id`=? AND status=?"
)
// 更新预约记录状态
func (m *AppointmentModel) UpdateAppointmentRecordStatus(orderId string, now, newStatus, status int) error {
	if _, err := m.Engine.Exec(UPDATE_RECORD_STATUS, now, newStatus, orderId, status); err != nil {
		return err
	}

	return nil
}

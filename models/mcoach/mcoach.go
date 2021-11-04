package mcoach

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
	"fmt"
)

type CoachModel struct {
	Coach       *models.VenueCoachDetail
	Engine      *xorm.Session
	Labels      *models.VenueCoachLabelConfig
	CoachScore  *models.VenueCoachScore
	Evaluate    *models.VenueUserEvaluateRecord
}

type CoachInfo struct {
	Id           int64      `json:"id"`
	Cover        string     `json:"cover"`
	Name         string     `json:"name"`
	Designation  string     `json:"designation"`
}

type CoachDetail struct {
	Id                int64   `json:"id"`
	Title             string  `json:"title"`
	Name              string  `json:"name"`
	Address           string  `json:"address"`
	Designation       string  `json:"designation"`
	Describe          string  `json:"describe"`
	AreasOfExpertise  string  `json:"areas_of_expertise"`
	Cover             string  `json:"cover"`
	Avatar            string  `json:"avatar"`
	Score             string  `json:"score"`
	TotalNum          int     `json:"total_num"`
	Courses           []*CourseInfo    `json:"courses"`
}

type CourseInfo struct {
	Id             int64  `json:"id"`
	CoachId        int64  `json:"coach_id"`
	Subhead        string `json:"subhead"`
	ClassPeriod    int    `json:"class_period"`
	Title          string `json:"title"`
	Describe       string `json:"describe"`
	Price          int    `json:"price"`
	PromotionPic   string `json:"promotion_pic"`
	Icon           string `json:"icon"`
	CourseType     int    `json:"course_type"`
	PeriodNum      int    `json:"period_num"`
	VenueId        int    `json:"venue_id"`

	CostDescription    string       `json:"cost_description"`     // 费用说明
	Instructions       string       `json:"instructions"`         // 购买须知
}

// 评价列表返回数据
type EvaluateResp struct {
	List     []*EvaluateInfo   `json:"list"`

}

// 评价信息
type EvaluateInfo struct {
	Id        int64        `json:"id"`
	//UserId    string       `json:"user_id"`
	//NickName  string       `json:"nick_name"`
	Avatar    string       `json:"avatar"`
	CoachId   int64       `json:"coach_id"`
	Star      int          `json:"star"`
	Content   string       `json:"content"`
	Labels    []*LabelInfo `json:"labels"`
}

type LabelInfo struct {
	Id          int64     `json:"id"`
	LabelName   string    `json:"label_name"`
}

type PubEvaluateParam struct {
	CoachId     int64          `binding:"required" json:"coach_id"`
	OrderId     string         `binding:"required" json:"order_id"`
	Star        int            `binding:"required" json:"star"`
	Content     string         `json:"content"`                        // 暂时无评价内容
	LabelIds    []interface{}  `json:"label_ids"`
}

func NewCoachModel(engine *xorm.Session) *CoachModel {
	return &CoachModel{
		Coach: new(models.VenueCoachDetail),
		Labels: new(models.VenueCoachLabelConfig),
		CoachScore: new(models.VenueCoachScore),
		Evaluate: new(models.VenueUserEvaluateRecord),
		Engine: engine,
	}
}

// 通过私教id 获取私教信息
func (m *CoachModel) GetCoachInfoById(id string) (bool, error) {
	m.Coach = new(models.VenueCoachDetail)
	return m.Engine.Where("id=?", id).Get(m.Coach)
}

const (
	GET_COACH_LIST = "SELECT va.coach_id, vc.* FROM venue_appointment_info " +
		"AS va LEFT JOIN venue_coach_detail AS vc ON va.coach_id=vc.id WHERE va.appointment_type=1 GROUP BY " +
		"va.coach_id ORDER BY va.id DESC LIMIT ?, ?"
)
// 通过课程id、私教类型 获取老师列表
func (m *CoachModel) GetCoachList(offset, size int) ([]*models.VenueCoachDetail, error) {
	var list []*models.VenueCoachDetail
	if err := m.Engine.SQL(GET_COACH_LIST, offset, size).Find(&list); err != nil {
		return nil , err
	}

	return list, nil
}

// 获取私教的评价列表
func (m *CoachModel) GetEvaluateListByCoach(coachId string, offset, size int) ([]*models.VenueUserEvaluateRecord, error) {
	var list []*models.VenueUserEvaluateRecord
	if err := m.Engine.Where("status=0 AND coach_id=?", coachId).Limit(size, offset).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取评价配置
func (m *CoachModel) GetEvaluateConfig() ([]*models.VenueCoachLabelConfig, error) {
	var list []*models.VenueCoachLabelConfig
	if err := m.Engine.Where("status=0").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

const (
	RECORD_COACH_SCORE_INFO = "INSERT INTO venue_coach_score(`coach_id`, `total_score`, `total_num`, `total%d_star`, `create_at`, `update_at`) " +
		"VALUES(?, ?, 1, 1, ?, ?) " +
		"ON DUPLICATE KEY UPDATE " +
		"total%d_star = total%d_star + 1, " +
		"total_num = total_num + 1, " +
		"total_score = total_score + ?, " +
		"update_at = ?"
)

// 记录私教评价总计
// 1星 = 1分
func (m *CoachModel) RecordCoachScoreInfo(coachId int64, starNum, now int) (int64, error) {
	sql := fmt.Sprintf(RECORD_COACH_SCORE_INFO, starNum, starNum, starNum)
	res, err := m.Engine.Exec(sql, coachId, starNum, now, now, starNum, now)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

// 获取私教评价总计
func (m *CoachModel) GetCoachScoreInfo(coachId string) (bool, error) {
	return m.Engine.Where("coach_id=? AND status=0", coachId).Get(m.CoachScore)
}

// 添加私教评价
func (m *CoachModel) AddCoachEvaluate() (int64, error) {
	return m.Engine.InsertOne(m.Evaluate)
}

// 用户是否已评价
func (m *CoachModel) HasEvaluateByUserId(userId, orderId string) (bool, error) {
	return m.Engine.Where("user_id=? AND order_id=? AND status=0", userId, orderId).Exist(m.Evaluate)
}

// 通过ids[多个id]获取标签配置列表
func (m *CoachModel) GetCoachLabelByIds(ids []interface{}) ([]*models.VenueCoachLabelConfig, error) {
	var list []*models.VenueCoachLabelConfig
	if err := m.Engine.In("id", ids...).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

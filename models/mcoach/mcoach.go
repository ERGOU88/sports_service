package mcoach

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

type CoachModel struct {
	Coach     *models.VenueCoachDetail
	Engine    *xorm.Session
	Labels    *models.VenueUserLabel
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
	Courses           []*CourseInfo    `json:"courses"`
}

type CourseInfo struct {
	Id             int64  `json:"id"`
	CoachId        int64  `json:"coach_id"`
	ClassPeriod    int    `json:"class_period"`
	Title          string `json:"title"`
	Describe       string `json:"describe"`
	Price          int    `json:"price"`
	PromotionPic   string `json:"promotion_pic"`
	Icon           string `json:"icon"`
	CourseType     int    `json:"course_type"`
	PeriodNum      int    `json:"period_num"`
}

// 评价信息
type EvaluateInfo struct {
	Id        int64        `json:"id"`
	//UserId    string       `json:"user_id"`
	//NickName  string       `json:"nick_name"`
	//Avatar    string       `json:"avatar"`
	CoachId   string       `json:"coach_id"`
	Star      int          `json:"star"`
	Content   string       `json:"content"`
	Labels    []*LabelInfo `json:"labels"`
}

type LabelInfo struct {
	Id     int64     `json:"id"`
	Name   string    `json:"name"`
}


func NewCoachModel(engine *xorm.Session) *CoachModel {
	return &CoachModel{
		Coach: new(models.VenueCoachDetail),
		Engine: engine,
	}
}

// 通过私教id 获取私教信息
func (m *CoachModel) GetCoachInfoById(id string) (bool, error) {
	m.Coach = new(models.VenueCoachDetail)
	return m.Engine.Where("id=?", id).Get(m.Coach)
}

// 通过课程id、私教类型 获取老师列表
func (m *CoachModel) GetCoachList(offset, size int) ([]*models.VenueCoachDetail, error) {
	var list []*models.VenueCoachDetail
	if err := m.Engine.Where("status=0 AND course_id=? AND coach_type=?", m.Coach.CourseId, m.Coach.CoachType).Limit(size, offset).Find(&list); err != nil {
		return nil, err
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

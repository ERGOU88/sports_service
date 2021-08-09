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


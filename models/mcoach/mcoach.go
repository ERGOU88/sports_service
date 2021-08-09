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

func NewCoachModel(engine *xorm.Session) *CoachModel {
	return &CoachModel{
		Coach: new(models.VenueCoachDetail),
		Engine: engine,
	}
}

// 通过私教id 获取私教信息
func (m *CoachModel) GetCoachInfoById(id int64) (bool, error) {
	m.Coach = new(models.VenueCoachDetail)
	return m.Engine.Where("id=?", id).Get(m.Coach)
}

// 通过课程id、私教类型 获取老师列表
func (m *CoachModel) GetCoachList() ([]*models.VenueCoachDetail, error) {
	var list []*models.VenueCoachDetail
	if err := m.Engine.Where("status=0 AND course_id=? AND coach_type=?", m.Coach.CourseId, m.Coach.CoachType).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}


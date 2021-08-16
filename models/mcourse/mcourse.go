package mcourse

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

type CourseModel struct {
	Course    *models.VenueCourseDetail
	Engine    *xorm.Session
	Labels    *models.VenueUserLabel
}

func NewCourseModel(engine *xorm.Session) *CourseModel {
	return &CourseModel{
		Course: new(models.VenueCourseDetail),
		Engine: engine,
	}
}

// 通过课程id 获取课程信息
func (m *CourseModel) GetCourseInfoById(id string) (bool, error) {
	m.Course = new(models.VenueCourseDetail)
	return m.Engine.Where("id=?", id).Get(m.Course)
}

// 通过私教id、课程类型 获取课程列表
func (m *CourseModel) GetCourseList() ([]*models.VenueCourseDetail, error) {
	var list []*models.VenueCourseDetail
	if err := m.Engine.Where("status=0 AND coach_id=? AND course_type=?", m.Course.CoachId, m.Course.CourseType).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 通过私教id获取 私教的课程
func (m *CourseModel) GetCourseByCoachId(coachId string) ([]*models.VenueCourseDetail, error) {
	var list []*models.VenueCourseDetail
	if err := m.Engine.Where("status=0 AND coach_id=? AND course_type=1", coachId).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}


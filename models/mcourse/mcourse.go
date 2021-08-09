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
func (m *CourseModel) GetCourseInfoById() (bool, error) {
	return m.Engine.Get(m.Course)
}

// 通过私教id、私教类型 获取课程列表
func (m *CourseModel) GetCourseList() ([]*models.VenueCourseDetail, error) {
	var list []*models.VenueCourseDetail
	if err := m.Engine.Where("status=0 AND coach_id=? AND course_type=?", m.Course.CoachId, m.Course.CourseType).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}


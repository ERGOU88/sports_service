package mcourse

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
	"sports_service/server/models/mcoach"
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
//func (m *CourseModel) GetCourseList() ([]*models.VenueCourseDetail, error) {
//	var list []*models.VenueCourseDetail
//	if err := m.Engine.Where("status=0 AND coach_id=? AND course_type=?", m.Course.CoachId, m.Course.CourseType).Find(&list); err != nil {
//		return nil, err
//	}
//
//	return list, nil
//}

const (
	GET_COURSE_LIST = "SELECT va.course_id, va.venue_id, va.duration AS class_period, vc.* FROM venue_appointment_info " +
		"AS va LEFT JOIN venue_course_detail AS vc ON va.course_id=vc.id WHERE va.appointment_type=2 GROUP BY " +
		"va.course_id ORDER BY va.id DESC"
)
// 获取大课列表
func (m *CourseModel) GetCourseList() ([]*mcoach.CourseInfo, error) {
	var list []*mcoach.CourseInfo
	if err := m.Engine.SQL(GET_COURSE_LIST).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 通过私教id获取 私教的课程
//func (m *CourseModel) GetCourseByCoachId(coachId string) ([]*models.VenueCourseDetail, error) {
//	var list []*models.VenueCourseDetail
//	if err := m.Engine.Where("status=0 AND coach_id=? AND course_type=1", coachId).Find(&list); err != nil {
//		return nil, err
//	}
//
//	return list, nil
//}


const (
	GET_COURSE_BY_COACH_ID = "SELECT va.course_id, va.coach_id, va.cur_amount AS price, va.venue_id, va.period_num, va.duration AS " +
		"class_period, vc.* FROM venue_appointment_info AS va LEFT JOIN venue_course_detail AS vc ON va.course_id=vc.id " +
		"WHERE va.appointment_type=1 AND va.coach_id=? GROUP BY va.course_id ORDER BY va.id DESC"
)
// 通过私教id获取 私教的课程
func (m *CourseModel) GetCourseByCoachId(coachId string) ([]*mcoach.CourseInfo, error) {
	var list []*mcoach.CourseInfo
	if err := m.Engine.SQL(GET_COURSE_BY_COACH_ID, coachId).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

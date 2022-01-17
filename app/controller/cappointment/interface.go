package cappointment

import (
	"sports_service/server/global/consts"
	"sports_service/server/models/mappointment"
)

type IAppointment interface {
	// 选项 [场馆、私课、大课]
	Options(relatedId int64) (int, interface{})
	// 进行预约
	Appointment(*mappointment.AppointmentReq) (int, interface{})
	// 预约时间选项
	AppointmentOptions() (int, interface{})
	// 预约详情
	AppointmentDetail() (int, interface{})
	// 预约日期
	AppointmentDate() (int, interface{})
	// 设置星期
	SetWeek(week int)
	// 设置场馆id
	SetVenueId(venueId int)
	// 设置教练id
	SetCoachId(coachId int)
	// 设置课程id
	SetCourseId(courseId int)
	// 设置预约类型
	SetAppointmentType(appointmentType int)
	// 设置预约配置id
	SetAppointmentId(id int64)
	// 设置日期id
	SetDateId(id int)
}

func GetOptions(i IAppointment, relatedId int64) (int, interface{}) {
	return i.Options(relatedId)
}

func UserAppointment(i IAppointment, param *mappointment.AppointmentReq) (int, interface{}) {
	return i.Appointment(param)
}

func GetAppointmentTimeOptions(i IAppointment, week, queryType, relatedId, id, coachId int) (int, interface{}) {
	i.SetWeek(week)
	i.SetAppointmentType(queryType)
	switch queryType {
	case consts.APPOINTMENT_VENUE:
		i.SetVenueId(relatedId)
	case consts.APPOINTMENT_COACH:
		i.SetCoachId(coachId)
		i.SetCourseId(relatedId)
	case consts.APPOINTMENT_COURSE:
		i.SetCourseId(relatedId)
	}

	i.SetDateId(id)
	return i.AppointmentOptions()
}

func GetAppointmentDetail(i IAppointment, dateId, id int) (int, interface{}) {
	i.SetDateId(dateId)
	i.SetAppointmentId(int64(id))
	return i.AppointmentDetail()
}

// 预约日期
func GetAppointmentDate(i IAppointment, queryType, relatedId, coachId int) (int, interface{}) {
	switch queryType {
	case consts.APPOINTMENT_VENUE:
		i.SetVenueId(relatedId)
	case consts.APPOINTMENT_COACH:
		i.SetCoachId(coachId)
		i.SetCourseId(relatedId)
	case consts.APPOINTMENT_COURSE:
		i.SetCourseId(relatedId)
	}

	return i.AppointmentDate()
}




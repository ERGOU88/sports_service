package cappointment

import (
	"sports_service/server/models/mappointment"
)

type IAppointment interface {
	// 选项 [场馆、私课、大课]
	Options(relatedId int64) (int, interface{})
	// 进行预约
	Appointment(*mappointment.AppointmentReq) (int, interface{})
	// 取消预约
	AppointmentCancel() int
	// 预约时间选项
	AppointmentOptions() (int, interface{})
	// 预约详情
	AppointmentDetail() (int, interface{})
	// 预约日期
	AppointmentDate() (int, interface{})
	// 设置星期
	SetWeek(week int)
	// 设置关联id
	SetRelatedId(relatedId int)
	// 设置预约类型
	SetAppointmentType(appointmentType int)
	// 设置日期id
	SetDateId(id int)
}

func GetOptions(i IAppointment, relatedId int64) (int, interface{}) {
	return i.Options(relatedId)
}

func UserAppointment(i IAppointment, param *mappointment.AppointmentReq) (int, interface{}) {
	return i.Appointment(param)
}

func UserAppointmentCancel(i IAppointment) int {
	return i.AppointmentCancel()
}

func GetAppointmentTimeOptions(i IAppointment, week, queryType, relatedId, id int) (int, interface{}) {
	i.SetWeek(week)
	i.SetAppointmentType(queryType)
	i.SetRelatedId(relatedId)
	i.SetDateId(id)
	return i.AppointmentOptions()
}

func GetAppointmentDetail(i IAppointment) (int, interface{}) {
	return i.AppointmentDetail()
}

// 预约日期
func GetAppointmentDate(i IAppointment) (int, interface{}) {
	return i.AppointmentDate()
}




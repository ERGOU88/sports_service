package cappointment

type IAppointment interface {
	// 进行预约
	Appointment() (int, interface{})
	// 取消预约
	AppointmentCancel() int
	// 预约选项
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

func UserAppointment(i IAppointment) (int, interface{}) {
	return i.Appointment()
}

func UserAppointmentCancel(i IAppointment) int {
	return i.AppointmentCancel()
}

func GetAppointmentOptions(i IAppointment) (int, interface{}) {
	return i.AppointmentOptions()
}

func GetAppointmentDetail(i IAppointment) (int, interface{}) {
	return i.AppointmentDetail()
}

// 预约日期
func GetAppointmentDate(i IAppointment) (int, interface{}) {
	return i.AppointmentDate()
}




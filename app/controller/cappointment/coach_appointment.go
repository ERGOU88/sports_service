package cappointment

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/mcourse"
	"sports_service/server/models/muser"
)

type CoachAppointmentModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	course      *mcourse.CourseModel
	*base
}

func NewCoach(c *gin.Context) *CoachAppointmentModule {
	venueSocket := dao.VenueEngine.NewSession()
	defer venueSocket.Close()
	appSocket := dao.AppEngine.NewSession()
	defer appSocket.Close()

	return &CoachAppointmentModule{
		context: c,
		user:    muser.NewUserModel(appSocket),
		course:  mcourse.NewCourseModel(venueSocket),
		engine:  venueSocket,
		base:    New(venueSocket),
	}
}

// 私教课程选项
func (svc *CoachAppointmentModule) Options(relatedId int64) (int, interface{}) {
	svc.course.Course.CoachId = relatedId
	svc.course.Course.CourseType = 1
	list, err := svc.course.GetCourseList()
	if err != nil {
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []interface{}{}
	}

	res := make([]*mappointment.Options, len(list))
	for index, item := range list {
		info := &mappointment.Options{
			Id: item.Id,
			Name: item.Title,
		}

		res[index] = info
	}

	return errdef.SUCCESS, res
}

// 预约私教
func (svc *CoachAppointmentModule) Appointment() (int, interface{}) {
	return 5000, nil
}

// 取消预约
func (svc *CoachAppointmentModule) AppointmentCancel() int {
	return 6000
}

// 获取某天的预约选项
func (svc *CoachAppointmentModule) AppointmentOptions() (int, interface{}) {
	date := svc.GetDateById(svc.DateId)
	if date <= 0 {
		return errdef.ERROR, nil
	}

	list, err := svc.GetAppointmentOptions()
	if err != nil {
		log.Log.Errorf("venue_trace: get options fail, err:%s", err)
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []interface{}{}
	}

	res := make([]*mappointment.OptionsInfo, len(list))
	for index, item := range list {
		info := svc.SetAppointmentOptionsRes(date, item)
		res[index] = info
	}


	return errdef.SUCCESS, res
}

func (svc *CoachAppointmentModule) AppointmentDetail() (int, interface{}) {
	return 8000, nil
}

// 预约私教日期配置
func (svc *CoachAppointmentModule) AppointmentDate() (int, interface{}) {
	return errdef.SUCCESS, svc.AppointmentDateInfo(6)
}



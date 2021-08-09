package cappointment

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/mcoach"
	"sports_service/server/models/mcourse"
	"sports_service/server/models/muser"
	"sports_service/server/global/app/log"
	"fmt"
)

type CourseAppointmentModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	course      *mcourse.CourseModel
	coach       *mcoach.CoachModel
	*base
}

func NewCourse(c *gin.Context) *CourseAppointmentModule {
	venueSocket := dao.VenueEngine.NewSession()
	defer venueSocket.Close()
	appSocket := dao.AppEngine.NewSession()
	defer appSocket.Close()

	return &CourseAppointmentModule{
		context: c,
		user:    muser.NewUserModel(appSocket),
		course:  mcourse.NewCourseModel(venueSocket),
		coach:   mcoach.NewCoachModel(venueSocket),
		engine:  venueSocket,
		base:    New(venueSocket),
	}
}

// 大课选项
func (svc *CourseAppointmentModule) Options(relatedId int64) (int, interface{}) {
	svc.course.Course.CoachId = 0
	svc.course.Course.CourseType = 2
	list, err := svc.course.GetCourseList()
	if err != nil {
		log.Log.Errorf("")
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
			Describe: item.Describe,
			CostDescription: "费用须知",
			Instructions: "购买说明",
		}

		res[index] = info
	}

	return errdef.SUCCESS, res
}


// 预约大课
func (svc *CourseAppointmentModule) Appointment() (int, interface{}) {
	return 5000, nil
}

// 取消预约
func (svc *CourseAppointmentModule) AppointmentCancel() int {
	return 6000
}

// 获取某天的预约选项
func (svc *CourseAppointmentModule) AppointmentOptions() (int, interface{}) {
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
		ok, err := svc.coach.GetCoachInfoById(fmt.Sprint(item.CoachId))
		if err != nil {
			log.Log.Errorf("venue_trace: get venue info by id fail, err:%s", err)
		}

		if ok {
			info.Name = svc.coach.Coach.Name
			info.Avatar = svc.coach.Coach.Avatar
		}


		res[index] = info
	}


	return errdef.SUCCESS, res
}

func (svc *CourseAppointmentModule) AppointmentDetail() (int, interface{}) {
	return 8000, nil
}

// 预约大课日期配置
func (svc *CourseAppointmentModule) AppointmentDate() (int, interface{}) {
	return errdef.SUCCESS, svc.AppointmentDateInfo(6)
}



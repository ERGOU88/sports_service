package cappointment

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/muser"
	"sports_service/server/global/app/log"
)

type CoachAppointmentModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
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
		engine:  venueSocket,
		base:    New(venueSocket),
	}
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



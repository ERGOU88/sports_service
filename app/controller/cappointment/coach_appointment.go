package cappointment

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/muser"
)

type CoachAppointmentModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	appointment *mappointment.AppointmentModel
}

func NewCoach(c *gin.Context) *CoachAppointmentModule {
	socket := dao.Engine.NewSession()
	defer socket.Close()
	return &CoachAppointmentModule{
		context: c,
		user: muser.NewUserModel(socket),
		appointment: mappointment.NewAppointmentModel(socket),
		engine:  socket,
	}
}

func (svc *CoachAppointmentModule) Appointment() (int, interface{}) {
	return 5000, nil
}

func (svc *CoachAppointmentModule) AppointmentCancel() int {
	return 6000
}

// 获取某天的预约选项
func (svc *CoachAppointmentModule) AppointmentOptions() (int, interface{}) {
	return 7000, nil
}

func (svc *CoachAppointmentModule) AppointmentDetail() (int, interface{}) {
	return 8000, nil
}

// 预约私教日期配置
func (svc *CoachAppointmentModule) AppointmentDate() (int, interface{}) {
	return errdef.SUCCESS, svc.appointment.GetAppointmentDate(6)
}



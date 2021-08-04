package cappointment

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/models/muser"
)

type CoachAppointmentModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
}

func NewCoach(c *gin.Context) *CoachAppointmentModule {
	socket := dao.Engine.NewSession()
	defer socket.Close()
	return &CoachAppointmentModule{
		context: c,
		user: muser.NewUserModel(socket),
		engine:  socket,
	}
}

func (service *CoachAppointmentModule) Appointment() (int, []interface{}) {
	return 5000, nil
}

func (service *CoachAppointmentModule) AppointmentCancel() int {
	return 6000
}

func (service *CoachAppointmentModule) AppointmentOptions() (int, map[string]interface{}) {
	return 7000, nil
}

func (service *CoachAppointmentModule) AppointmentDetail() (int, interface{}) {
	return 8000, nil
}



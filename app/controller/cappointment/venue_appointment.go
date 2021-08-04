package cappointment

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/models/mattention"
	"sports_service/server/models/muser"
)

type VenueAppointmentModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	attention   *mattention.AttentionModel
}

func NewVenue(c *gin.Context) *VenueAppointmentModule {
	socket := dao.Engine.NewSession()
	defer socket.Close()
	return &VenueAppointmentModule{
		context: c,
		user: muser.NewUserModel(socket),
		engine:  socket,
	}
}

// 预约场馆
func (service *VenueAppointmentModule) Appointment() (int, []interface{}) {
	return 1000, nil
}

func (service *VenueAppointmentModule) AppointmentCancel() int {
	return 2000
}

func (service *VenueAppointmentModule) AppointmentOptions() (int, map[string]interface{}) {
	return 3000, nil
}

func (service *VenueAppointmentModule) AppointmentDetail() (int, interface{}) {
	return 4000, nil
}


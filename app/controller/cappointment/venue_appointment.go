package cappointment

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/models"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/muser"
)

type VenueAppointmentModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	appointment *mappointment.AppointmentModel
}

func NewVenue(c *gin.Context) *VenueAppointmentModule {
	socket := dao.Engine.NewSession()
	defer socket.Close()
	return &VenueAppointmentModule{
		context: c,
		user: muser.NewUserModel(socket),
		appointment: mappointment.NewAppointmentModel(socket),
		engine:  socket,
	}
}

// 预约场馆
func (svc *VenueAppointmentModule) Appointment() (int, interface{}) {
	return 0, nil
}

func (svc *VenueAppointmentModule) AppointmentCancel() int {
	return 2000
}

// 预约场馆选项
func (svc *VenueAppointmentModule) AppointmentOptions() (int, interface{}) {
	list, err := svc.appointment.GetOptionsByWeek()
	if err != nil {
		return errdef.ERROR, nil
	}

	if list == nil {
		return errdef.SUCCESS, []*models.VenueAppointmentInfo{}
	}

	return errdef.SUCCESS, list
}

func (svc *VenueAppointmentModule) AppointmentDetail() (int, interface{}) {
	return 4000, nil
}

// 预约日期配置
func (svc *VenueAppointmentModule) AppointmentDate() (int, interface{}) {
	// todo: 查看当天最低价
	return errdef.SUCCESS, svc.appointment.GetAppointmentDate(6)
}


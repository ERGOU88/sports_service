package cappointment

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/models/muser"
)

type VenueAppointmentModule struct {
	context         *gin.Context
	engine          *xorm.Session
	user            *muser.UserModel
	*base
}

func NewVenue(c *gin.Context) *VenueAppointmentModule {
	socket := dao.Engine.NewSession()
	defer socket.Close()
	return &VenueAppointmentModule{
		context: c,
		user: muser.NewUserModel(socket),
		engine:  socket,
		base: New(socket),
	}
}

// 预约场馆
func (svc *VenueAppointmentModule) Appointment() (int, interface{}) {
	return 0, nil
}

// 取消预约
func (svc *VenueAppointmentModule) AppointmentCancel() int {
	return 2000
}

// 预约场馆选项
func (svc *VenueAppointmentModule) AppointmentOptions() (int, interface{}) {
	list, err := svc.GetAppointmentOptions()
	if err != nil {
		return errdef.ERROR, list
	}

	return errdef.SUCCESS, list
}

func (svc *VenueAppointmentModule) AppointmentDetail() (int, interface{}) {
	return 4000, nil
}

// 场馆预约日期配置
func (svc *VenueAppointmentModule) AppointmentDate() (int, interface{}) {
	return errdef.SUCCESS, svc.AppointmentDateInfo(6)
}

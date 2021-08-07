package cappointment

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/muser"
)

type VenueAppointmentModule struct {
	context         *gin.Context
	engine          *xorm.Session
	user            *muser.UserModel
	*base
}

func NewVenue(c *gin.Context) *VenueAppointmentModule {
	venueSocket := dao.VenueEngine.NewSession()
	defer venueSocket.Close()
	appSocket := dao.AppEngine.NewSession()
	defer appSocket.Close()
	return &VenueAppointmentModule{
		context: c,
		user:    muser.NewUserModel(appSocket),
		engine:  venueSocket,
		base:    New(venueSocket),
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
		log.Log.Errorf("venue_trace: get options fail, err:%s", err)
		return errdef.ERROR, list
	}

	if len(list) == 0 {
		return errdef.SUCCESS, list
	}

	res := make([]*mappointment.OptionsInfo, len(list))
	for _, item := range list {
		info := &mappointment.OptionsInfo{
			RelatedId: item.RelatedId,
			CurAmount: item.CurAmount,
			TimeNode: item.TimeNode,
			Duration: item.Duration,
			RealAmount: item.RealAmount,
			DiscountRate: item.DiscountRate,
			DiscountAmount: item.DiscountAmount,
			QuotaNum: item.QuotaNum,
			RecommendType: item.RecommendType,
			AppointmentType: item.AppointmentType,
			WeekNum: item.WeekNum,

		}
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

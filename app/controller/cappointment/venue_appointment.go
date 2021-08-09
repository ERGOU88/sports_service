package cappointment

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvenue"
)

type VenueAppointmentModule struct {
	context         *gin.Context
	engine          *xorm.Session
	user            *muser.UserModel
	venue           *mvenue.VenueModel
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
		venue:   mvenue.NewVenueModel(venueSocket),
		engine:  venueSocket,
		base:    New(venueSocket),
	}
}

// 场馆选项 todo：暂时只有一个场馆
func (svc *VenueAppointmentModule) Options(relatedId int64) (int, interface{}) {
	list, err := svc.venue.GetVenueList()
	if err != nil {
		return errdef.ERROR, nil
	}

	if list == nil {
		return errdef.SUCCESS, []interface{}{}
	}

	res := make([]*mappointment.Options, len(list))
	for index, item := range list {
		info := &mappointment.Options{
			Id: item.Id,
			Name: item.VenueName,
		}

		res[index] = info
	}

	return errdef.SUCCESS, res
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

		svc.venue.Venue.Id = item.RelatedId
		ok, err := svc.venue.GetVenueInfoById()
		if err != nil {
			log.Log.Errorf("venue_trace: get venue info by id fail, err:%s", err)
		}

		if ok {
			info.Name = svc.venue.Venue.VenueName
		}

		svc.venue.Labels.TimeNode = item.TimeNode
		svc.venue.Labels.Date = int64(date)
		svc.venue.Labels.VenueId = item.RelatedId
		labels, err := svc.venue.GetVenueUserLabels()
		if err != nil {
			log.Log.Errorf("venue_trace: get venue user lables fail, err:%s", err)
		}

		info.Labels = make([]*mappointment.LabelInfo, len(labels))
		for key, val := range labels {
			label := &mappointment.LabelInfo{
				UserId: val.UserId,
				LabelId: val.LabelId,
				LabelName: val.LabelName,
			}

			user := svc.user.FindUserByUserid(val.UserId)
			if user != nil {
				label.NickName = user.NickName
				label.Avatar = user.Avatar
			}

			info.Labels[key] = label
		}

		res[index] = info
	}


	return errdef.SUCCESS, res
}

func (svc *VenueAppointmentModule) AppointmentDetail() (int, interface{}) {
	return 4000, nil
}

// 场馆预约日期配置
func (svc *VenueAppointmentModule) AppointmentDate() (int, interface{}) {
	return errdef.SUCCESS, svc.AppointmentDateInfo(6)
}

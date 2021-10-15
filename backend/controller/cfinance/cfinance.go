package cfinance

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/morder"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvenue"
	"sports_service/server/util"
	"fmt"
	"time"
)

// todo:
type FinanceModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	order       *morder.OrderModel
	venue       *mvenue.VenueModel
}

func New(c *gin.Context) FinanceModule {
	appSocket := dao.AppEngine.NewSession()
	defer appSocket.Close()

	venueSocket := dao.VenueEngine.NewSession()
	defer venueSocket.Close()
	return FinanceModule{
		context: c,
		user: muser.NewUserModel(appSocket),
		order: morder.NewOrderModel(venueSocket),
		venue: mvenue.NewVenueModel(venueSocket),
		engine: venueSocket,
	}
}

// 获取订单流水列表
func (svc *FinanceModule) GetOrderList(page, size int) (int, []*morder.OrderRecord) {
	offset := (page - 1) * size
	list, err := svc.order.GetOrderList(offset, size)
	if err != nil {
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*morder.OrderRecord{}
	}


	var res []*morder.OrderRecord
	for index, item := range list {
		info := &morder.OrderRecord{
			Id:  item.Id,
			PayOrderId: item.PayOrderId,
			OriginalAmount: fmt.Sprintf("%.2f", float64(item.OriginalAmount)/100),
			CreateAt: time.Unix(int64(item.CreateAt), 0).Format(consts.FORMAT_TM),
			Amount: fmt.Sprintf("%.2f", float64(item.Amount)/100),
		}

		extra := mappointment.OrderResp{}
		if err := util.JsonFast.UnmarshalFromString(item.Extra, &extra); err != nil {
			continue
		}

		info.MobileNum = extra.MobileNum

		ok, err := svc.venue.GetVenueInfoById(fmt.Sprint(item.Id))
		if ok && err == nil {
			info.VenueName = svc.venue.Venue.VenueName
		}

		switch item.ProductType {
		case consts.ORDER_TYPE_APPOINTMENT_VENUE:
			info.Detail = fmt.Sprintf("预约场馆 * %d", extra.Count)
		case consts.ORDER_TYPE_APPOINTMENT_COACH:
			info.Detail = fmt.Sprintf("预约私教 %s", extra.CoachName)
		case consts.ORDER_TYPE_APPOINTMENT_COURSE:
			info.Detail = fmt.Sprintf("预约课程 %s", extra.CourseName)
		case consts.ORDER_TYPE_EXPERIENCE_CARD:
			info.Detail = fmt.Sprintf("次卡 * %d", extra.Count)
		case consts.ORDER_TYPE_MONTH_CARD:
			info.Detail = fmt.Sprintf("月卡 * %d", extra.Count)
		case consts.ORDER_TYPE_SEANSON_CARD:
			info.Detail = fmt.Sprintf("季卡 * %d", extra.Count)
		case consts.ORDER_TYPE_HALF_YEAR_CARD:
			info.Detail = fmt.Sprintf("半年卡 * %d", extra.Count)
		case consts.ORDER_TYPE_YEAR_CARD:
			info.Detail = fmt.Sprintf("年卡 * %d", extra.Count)
		}


	}
}

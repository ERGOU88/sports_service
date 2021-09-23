package cappointment

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/mcoach"
	"sports_service/server/models/mcourse"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvenue"
	"sports_service/server/util"
	"time"
)

type CoachAppointmentModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	course      *mcourse.CourseModel
	coach       *mcoach.CoachModel
	venue       *mvenue.VenueModel
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
		course:  mcourse.NewCourseModel(venueSocket),
		coach:   mcoach.NewCoachModel(venueSocket),
		venue:   mvenue.NewVenueModel(venueSocket),
		engine:  venueSocket,
		base:    New(venueSocket),
	}
}

// 私教课程选项
func (svc *CoachAppointmentModule) Options(relatedId int64) (int, interface{}) {
	//svc.course.Course.CoachId = relatedId
	//svc.course.Course.CourseType = 1
	list, err := svc.course.GetCourseByCoachId(fmt.Sprint(relatedId))
	if err != nil {
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
			Instructions: "预约须知",
		}

		res[index] = info
	}

	return errdef.SUCCESS, res
}

// 预约私教
func (svc *CoachAppointmentModule) Appointment(params *mappointment.AppointmentReq) (int, interface{}) {
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("venue_trace: session begin fail, err:%s", err)
		return errdef.ERROR, nil
	}

	if len(params.Infos) == 0 {
		svc.engine.Rollback()
		return errdef.APPOINTMENT_INVALID_INFO, nil
	}

	user := svc.user.FindUserByUserid(params.UserId)
	if user == nil {
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS, nil
	}

	list, err := svc.GetAppointmentConfByIds(params.Ids)
	if err != nil {
		svc.engine.Rollback()
		return errdef.APPOINTMENT_QUERY_NODE_FAIL, nil
	}

	if len(list) != len(params.Infos) {
		svc.engine.Rollback()
		return errdef.APPOINTMENT_INVALID_NODE_ID, nil
	}

	// 获取私教课程信息
	ok, err := svc.course.GetCourseInfoById(fmt.Sprint(params.RelatedId))
	if !ok || err != nil {
		log.Log.Errorf("venue_trace: get course info fail, err:%s", err)
		svc.engine.Rollback()
		return errdef.COURSE_NOT_EXISTS, nil
	}

	//if svc.course.Course.CoachId != params.CoachId {
	//	log.Log.Error("venue_trace: coach id not match, coachId:%d, curCoachId:%d",
	//	svc.course.Course.CoachId, params.CoachId)
	//	svc.engine.Rollback()
	//	return errdef.COACH_ID_NOT_MATCH, nil
	//}

	//ok, err = svc.coach.GetCoachInfoById(fmt.Sprint(svc.course.Course.CoachId))
	//if !ok || err != nil {
	//	log.Log.Errorf("venue_trace: get coach by id fail, err:%s", err)
	//	svc.engine.Rollback()
	//	return errdef.COACH_NOT_EXISTS, nil
	//}

	//if svc.coach.Coach.CoachType != 1 {
	//	log.Log.Errorf("venue_trace: coach type fail, coachType:%d", svc.coach.Coach.CoachType)
	//	return errdef.COACH_TYPE_FAIL, nil
	//}

	ok, err = svc.venue.GetVenueInfoById(fmt.Sprint(svc.course.Course.VenueId))
	if !ok || err != nil {
		log.Log.Errorf("venue_trace: get venue info by id fail, venueId:%d, err:%s", svc.course.Course.VenueId, err)
		svc.engine.Rollback()
		return errdef.VENUE_NOT_EXISTS, nil
	}

	svc.Extra.CoachId = svc.coach.Coach.Id
	svc.Extra.CoachName = svc.coach.Coach.Name
	svc.Extra.Address = svc.venue.Venue.Address
	svc.Extra.CourseId = svc.course.Course.Id
	svc.Extra.CourseName = svc.course.Course.Title
	svc.Extra.ProductImg = svc.course.Course.PromotionPic

	orderId := util.NewOrderId()
	now := int(time.Now().Unix())

	if err := svc.AppointmentProcess(user.UserId, orderId, params.RelatedId, params.WeekNum, params.LabelIds, params.Infos); err != nil {
		log.Log.Errorf("venue_trace: appointment fail, err:%s", err)
		svc.engine.Rollback()
		return errdef.APPOINTMENT_PROCESS_FAIL, nil
	}

	svc.Extra.Id = params.RelatedId
	svc.Extra.Name = svc.course.Course.Title
	svc.Extra.WeekCn = util.GetWeekCn(params.WeekNum)
	svc.Extra.MobileNum = util.HideMobileNum(fmt.Sprint(user.MobileNum))
	svc.Extra.TmCn = util.ResolveTime(svc.Extra.TotalTm)

	// 库存不足 返回最新数据 事务回滚
	if !svc.Extra.IsEnough {
		log.Log.Errorf("venue_trace: rollback, isEnough:%v, reqType:%d", svc.Extra.IsEnough, params.ReqType)
		svc.engine.Rollback()
		return errdef.APPOINTMENT_NOT_ENOUGH_STOCK, svc.Extra
	}

	// 查询数据 则返回200
	if params.ReqType != 2 {
		svc.engine.Rollback()
		return errdef.SUCCESS, svc.Extra
	}

	// 添加订单
	if err := svc.AddOrder(orderId, user.UserId, "预约私教", now, consts.ORDER_TYPE_APPOINTMENT_COACH); err != nil {
		log.Log.Errorf("venue_trace: add order fail, err:%s", err)
		svc.engine.Rollback()
		return errdef.ORDER_ADD_FAIL, nil
	}

	// 添加预约记录流水
	if err := svc.AddAppointmentRecord(); err != nil {
		log.Log.Errorf("venue_trace: add appointment record fail, err:%s", err)
		svc.engine.Rollback()
		return errdef.APPOINTMENT_ADD_RECORD_FAIL, nil
	}

	// 添加订单商品流水
	if err := svc.AddOrderProducts(); err != nil {
		log.Log.Errorf("venue_trace: add order products fail, err:%s", err)
		svc.engine.Rollback()
		return errdef.ORDER_PRODUCT_ADD_FAIL, nil
	}

	// 记录需处理支付超时的订单
	if _, err := svc.order.RecordOrderId(orderId); err != nil {
		log.Log.Errorf("venue_trace: record orderId fail, err:%s", err)
		svc.engine.Rollback()
		return errdef.APPOINTMENT_RECORD_ORDER_FAIL, nil
	}

	svc.engine.Commit()

	svc.Extra.OrderId = orderId
	svc.Extra.PayDuration = consts.PAYMENT_DURATION
	// 超时
	//redismq.PushOrderEventMsg(redismq.NewOrderEvent(user.UserId, svc.Extra.OrderId, int64(svc.order.Order.CreateAt) + svc.Extra.PayDuration,
	//	consts.ORDER_EVENT_COACH_TIME_OUT))
	return errdef.SUCCESS, svc.Extra
}

// 获取某天的预约选项
func (svc *CoachAppointmentModule) AppointmentOptions() (int, interface{}) {
	date := svc.GetDateById(svc.DateId, consts.FORMAT_DATE)
	if date == "" {
		return errdef.ERROR, nil
	}

	condition, err := svc.GetQueryCondition()
	if err != nil {
		log.Log.Errorf("venue_trace: get query condition fail, err:%s", err)
		return errdef.ERROR, nil
	}

	list, err := svc.GetAppointmentOptions(condition)
	if err != nil {
		log.Log.Errorf("venue_trace: get options fail, err:%s", err)
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []interface{}{}
	}

	res := make([]*mappointment.OptionsInfo, 0)
	for _, item := range list {
		info := svc.SetAppointmentOptionsRes(date, item)
		if info == nil {
			continue
		}

		res = append(res, info)
	}


	return errdef.SUCCESS, res
}

func (svc *CoachAppointmentModule) AppointmentDetail() (int, interface{}) {
	return 8000, nil
}

// 预约私教课程日期配置
func (svc *CoachAppointmentModule) AppointmentDate() (int, interface{}) {
	return errdef.SUCCESS, svc.AppointmentDateInfo(6, consts.APPOINTMENT_COACH)
}



package coach

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/dao"
	"sports_service/global/app/errdef"
	"sports_service/global/app/log"
	"sports_service/global/consts"
	"sports_service/models"
	"sports_service/models/mappointment"
	"sports_service/models/mcoach"
	"sports_service/models/mcourse"
	"sports_service/models/morder"
	"sports_service/models/muser"
	"sports_service/models/mvenue"
	"sports_service/tools/tencentCloud"
	"sports_service/util"
	"time"
)

type CoachModule struct {
	context *gin.Context
	engine  *xorm.Session
	coach   *mcoach.CoachModel
	course  *mcourse.CourseModel
	user    *muser.UserModel
	order   *morder.OrderModel
	venue   *mvenue.VenueModel
}

func New(c *gin.Context) *CoachModule {
	venueSocket := dao.VenueEngine.NewSession()
	defer venueSocket.Close()

	appSocket := dao.AppEngine.NewSession()
	defer appSocket.Close()
	return &CoachModule{
		context: c,
		coach:   mcoach.NewCoachModel(venueSocket),
		course:  mcourse.NewCourseModel(venueSocket),
		user:    muser.NewUserModel(appSocket),
		order:   morder.NewOrderModel(venueSocket),
		venue:   mvenue.NewVenueModel(venueSocket),
		engine:  venueSocket,
	}
}

// 获取私教列表
func (svc *CoachModule) GetCoachList(page, size int) (int, []*mcoach.CoachInfo) {
	offset := (page - 1) * size
	list, err := svc.coach.GetCoachList(offset, size)
	if err != nil {
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*mcoach.CoachInfo{}
	}

	res := make([]*mcoach.CoachInfo, len(list))
	for index, item := range list {
		info := &mcoach.CoachInfo{
			Id:          item.Id,
			Cover:       tencentCloud.BucketURI(item.Cover),
			Name:        item.Name,
			Designation: item.Designation,
		}

		res[index] = info
	}

	return errdef.SUCCESS, res
}

// 获取私教详细信息
func (svc *CoachModule) GetCoachDetail(coachId string) (int, *mcoach.CoachDetail) {
	ok, err := svc.coach.GetCoachInfoById(coachId)
	if err != nil || !ok {
		return errdef.ERROR, nil
	}

	res := &mcoach.CoachDetail{
		Id:               svc.coach.Coach.Id,
		Title:            svc.coach.Coach.Title,
		Name:             svc.coach.Coach.Name,
		Designation:      svc.coach.Coach.Designation,
		Describe:         svc.coach.Coach.Describe,
		AreasOfExpertise: svc.coach.Coach.AreasOfExpertise,
		Cover:            tencentCloud.BucketURI(svc.coach.Coach.Cover),
		Avatar:           tencentCloud.BucketURI(svc.coach.Coach.Avatar),
	}

	ok, err = svc.venue.GetVenueInfoById(fmt.Sprint(svc.course.Course.VenueId))
	if !ok || err != nil {
		log.Log.Errorf("venue_trace: get venue info by id fail, venueId:%d, err:%s", svc.course.Course.VenueId, err)
	} else {
		res.Address = svc.venue.Venue.Address
	}

	courses, err := svc.course.GetCourseByCoachId(coachId)
	if err != nil {
		log.Log.Errorf("coach_trace: get course by id fail, err:%s, coachId:%s", err, coachId)
	}

	if len(courses) > 0 {
		res.Courses = courses
	} else {
		res.Courses = make([]*mcoach.CourseInfo, 0)
	}

	return errdef.SUCCESS, res
}

// 获取评价列表
func (svc *CoachModule) GetEvaluateList(coachId string, page, size int) (int, []*mcoach.EvaluateInfo) {
	ok, err := svc.coach.GetCoachInfoById(coachId)
	if err != nil || !ok {
		return errdef.ERROR, nil
	}

	offset := (page - 1) * size
	list, err := svc.coach.GetEvaluateListByCoach(coachId, offset, size)
	if err != nil {
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*mcoach.EvaluateInfo{}
	}

	res := make([]*mcoach.EvaluateInfo, len(list))
	for index, item := range list {
		info := &mcoach.EvaluateInfo{
			Id: item.Id,
			//UserId: item.UserId,
			CoachId: item.CoachId,
			Star:    item.Star,
			Content: item.Content,
			Avatar:  consts.EVALUATE_DEFAULT_AVATAR,
		}

		//if user := svc.user.FindUserByUserid(item.UserId); user != nil {
		//	info.Avatar = user.Avatar
		//	info.NickName = user.NickName
		//}

		if err = util.JsonFast.UnmarshalFromString(item.LabelInfo, &info.Labels); err != nil {
			log.Log.Errorf("coach_trace: json unmarshal fail, coachId:%s", item.CoachId)
		}

		res[index] = info
	}

	return errdef.SUCCESS, res
}

// 获取评价配置
func (svc *CoachModule) GetEvaluateConfig() (int, []*models.VenueCoachLabelConfig) {
	list, err := svc.coach.GetEvaluateConfig()
	if err != nil {
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*models.VenueCoachLabelConfig{}
	}

	return errdef.SUCCESS, list
}

// 发布评价
func (svc *CoachModule) PubEvaluate(userId string, param *mcoach.PubEvaluateParam) int {
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("coach_trace: begin session fail, err:%s", err)
		return errdef.ERROR
	}

	if param.Star <= 0 || param.Star > 5 {
		log.Log.Errorf("coach_trace: invalid star, star:%d", param.Star)
		svc.engine.Rollback()
		return errdef.INVALID_PARAMS
	}

	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("coach_trace: user not found, userId:%s", userId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS
	}

	ok, err := svc.order.GetOrder(param.OrderId)
	if !ok || err != nil {
		log.Log.Errorf("coach_trace: coach order not found, err:%s", err)
		svc.engine.Rollback()
		return errdef.COACH_ORDER_NOT_EXISTS
	}

	products, err := svc.order.GetOrderProductsById(svc.order.Order.PayOrderId)
	if len(products) == 0 || err != nil {
		log.Log.Errorf("coach_trace: get order product by id fail, orderId:%s, err:%s", svc.order.Order.PayOrderId, err)
		svc.engine.Rollback()
		return errdef.COACH_ORDER_NOT_SUCCESS
	}

	for _, item := range products {
		// 如果订单类型 不是 预约私教
		if item.ProductType != consts.ORDER_TYPE_APPOINTMENT_COACH {
			log.Log.Errorf("coach_trace: invalid order product type, orderId:%s, productType:%d",
				svc.order.Order.PayOrderId, item.ProductType)
			svc.engine.Rollback()
			return errdef.COACH_ORDER_NOT_EXISTS
		}
	}

	if svc.order.Order.Status != consts.ORDER_TYPE_COMPLETED {
		log.Log.Errorf("coach_trace: coach order not success, status:%d", svc.order.Order.Status)
		svc.engine.Rollback()
		return errdef.COACH_ORDER_NOT_SUCCESS
	}

	extra := mappointment.OrderResp{}
	if err := util.JsonFast.UnmarshalFromString(svc.order.Order.Extra, &extra); err != nil {
		log.Log.Errorf("coach_trace: unmarshal extra fail, orderId:%s, err:%s", svc.order.Order.PayOrderId, err)
		svc.engine.Rollback()
		return errdef.COACH_ORDER_NOT_EXISTS
	}

	ok, err = svc.coach.GetCoachInfoById(fmt.Sprint(extra.CoachId))
	if !ok || err != nil {
		log.Log.Errorf("coach_trace: coach not found, coachId:%d", extra.CoachId)
		svc.engine.Rollback()
		return errdef.COACH_NOT_EXISTS
	}

	ok, err = svc.coach.HasEvaluateByUserId(userId, svc.order.Order.PayOrderId)
	if err != nil {
		log.Log.Errorf("coach_trace: get evaluate by userId fail, userId:%s, orderId:%s", userId, svc.order.Order.PayOrderId)
		svc.engine.Rollback()
		return errdef.COACH_PUB_EVALUATE_FAIL
	}

	if ok {
		log.Log.Errorf("coach_trace: coach already evaluate, userId:%s, orderId:%s", userId, svc.order.Order.PayOrderId)
		svc.engine.Rollback()
		return errdef.COACH_ALREADY_EVALUATE
	}

	var labels string
	if len(param.LabelIds) > 0 {
		list, err := svc.coach.GetCoachLabelByIds(param.LabelIds)
		if err != nil {
			log.Log.Errorf("coach_trace: get coach label fail, err:%s", err)
			svc.engine.Rollback()
			return errdef.COACH_GET_LABEL_FAIL
		}

		labels, _ = util.JsonFast.MarshalToString(list)
	}

	now := int(time.Now().Unix())
	svc.coach.Evaluate.CoachId = extra.CoachId
	svc.coach.Evaluate.UserId = userId
	svc.coach.Evaluate.Star = param.Star
	svc.coach.Evaluate.CreateAt = now
	svc.coach.Evaluate.UpdateAt = now
	svc.coach.Evaluate.LabelInfo = labels
	svc.coach.Evaluate.OrderType = 1
	svc.coach.Evaluate.OrderId = param.OrderId
	svc.coach.Evaluate.Content = param.Content

	// 添加私教评价
	if _, err := svc.coach.AddCoachEvaluate(); err != nil {
		log.Log.Errorf("coach_trace: add coach evaluate fail, err:%s", err)
		svc.engine.Rollback()
		return errdef.COACH_PUB_EVALUATE_FAIL
	}

	// 记录评价总计
	if _, err := svc.coach.RecordCoachScoreInfo(extra.CoachId, param.Star, now); err != nil {
		log.Log.Errorf("coach_trace: record coach score fail, err:%s", err)
		svc.engine.Rollback()
		return errdef.COACH_PUB_EVALUATE_FAIL
	}

	svc.engine.Commit()
	return errdef.SUCCESS
}

// 获取私教评分信息
func (svc *CoachModule) GetCoachScore(coachId string) (totalNum int, score float64) {
	ok, err := svc.coach.GetCoachScoreInfo(fmt.Sprint(coachId))
	if !ok || err != nil {
		log.Log.Errorf("coach_trace: get coach score info fail, err:%s", err)
		totalNum = 0
		score = 0
	} else {
		totalNum = svc.coach.CoachScore.TotalNum
		score = util.TruncFloat(float64(svc.coach.CoachScore.TotalScore)/float64(svc.coach.CoachScore.TotalNum), 1)
	}

	return
}

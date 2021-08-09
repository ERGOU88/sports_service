package coach

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/models/mcoach"
	"sports_service/server/models/mcourse"
	"sports_service/server/models/muser"
	"sports_service/server/util"
)

type CoachModule struct {
	context     *gin.Context
	engine      *xorm.Session
	coach       *mcoach.CoachModel
	course      *mcourse.CourseModel
	user        *muser.UserModel
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
		engine:  venueSocket,
	}
}

// 获取私教列表
func (svc *CoachModule) GetCoachList(page, size int) (int, []*mcoach.CoachInfo) {
	offset := (page - 1) * size
	svc.coach.Coach.CourseId = 0
	svc.coach.Coach.CoachType = 1
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
			Id:   item.Id,
			Cover: item.Cover,
			Name: item.Name,
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
		Id: svc.coach.Coach.Id,
		Title: svc.coach.Coach.Title,
		Name: svc.coach.Coach.Name,
		Address: svc.coach.Coach.Address,
		Designation: svc.coach.Coach.Designation,
		Describe: svc.coach.Coach.Describe,
		AreasOfExpertise: svc.coach.Coach.AreasOfExpertise,
		Cover: svc.coach.Coach.Cover,
		Avatar: svc.coach.Coach.Avatar,
	}

	courses, err := svc.course.GetCourseByCoachId(coachId)
	if err != nil {
		log.Log.Errorf("coach_trace: get course by id fail, err:%s, coachId:%s", err, coachId)
	}

	if len(courses) > 0 {
		res.Courses = make([]*mcoach.CourseInfo, len(courses))
		for key, val := range courses {
			course := &mcoach.CourseInfo{
				Id: val.Id,
				CourseType: val.CourseType,
				PeriodNum: val.PeriodNum,
				Price: val.Price,
				PromotionPic: val.PromotionPic,
				Icon: val.Icon,
				Title: val.Title,
				Describe: val.Describe,
				CoachId: val.CoachId,
				ClassPeriod: val.ClassPeriod,
			}

			res.Courses[key] = course
		}
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
			Star: item.Star,
			Content: item.Content,
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

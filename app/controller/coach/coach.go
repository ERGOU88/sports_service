package coach

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/models/mcoach"
)

type CoachModule struct {
	context     *gin.Context
	engine      *xorm.Session
	coach       *mcoach.CoachModel
}

func New(c *gin.Context) *CoachModule {
	venueSocket := dao.VenueEngine.NewSession()
	defer venueSocket.Close()
	return &CoachModule{
		context: c,
		coach:   mcoach.NewCoachModel(venueSocket),
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

package contest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mbanner"
	"sports_service/server/models/mcontest"
	"sports_service/server/models/muser"
	"sports_service/server/util"
	"time"
)

type ContestModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	banner      *mbanner.BannerModel
	contest     *mcontest.ContestModel
}

func New(c *gin.Context) ContestModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	return ContestModule{
		context: c,
		user: muser.NewUserModel(socket),
		banner: mbanner.NewBannerMolde(socket),
		contest: mcontest.NewContestModel(socket),
		engine: socket,
	}
}

// 获取赛事首页banner
func (svc *ContestModule) GetBanner() []*models.Banner {
	banners := svc.banner.GetRecommendBanners(int32(consts.CONTEST_BANNERS), time.Now().Unix(), 0, 10)
	if len(banners) == 0 {
		return []*models.Banner{}
	}

	return banners
}

// 获取直播列表 默认取两条最新的
func (svc *ContestModule) GetLiveList(page, size int) (int, []*mcontest.ContestLiveInfo) {
	offset := (page - 1) * size
	now := time.Now().Unix()
	list, err := svc.contest.GetLiveList(now, offset, size)
	if err != nil {
		return errdef.CONTEST_GET_LIVE_FAIL, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*mcontest.ContestLiveInfo{}
	}

	resp := make([]*mcontest.ContestLiveInfo, len(list))
	for index, item := range list {
		tm := time.Unix(int64(item.PlayTime), 0)
		live := &mcontest.ContestLiveInfo{
			Id: item.Id,
			UserId: item.UserId,
			RoomId: item.RoomId,
			GroupId: item.GroupId,
			Cover: item.Cover,
			RtmpAddr: item.RtmpAddr,
			HlsAddr: item.HlsAddr,
			FlvAddr: item.FlvAddr,
			PlayTime: item.PlayTime,
			Title: item.Title,
			HighLights: item.HighLights,
			Describe: item.Describe,
			Tags: item.Tags,
			LiveType: item.LiveType,
			Date: tm.Format("1月2日"),
			Week: util.GetWeekCn(int(tm.Weekday())),
		}

		user := svc.user.FindUserByUserid(item.UserId)
		if user != nil {
			live.NickName = user.NickName
			live.Avatar = user.Avatar
		}

		resp[index] = live
	}

	return errdef.SUCCESS, resp
}

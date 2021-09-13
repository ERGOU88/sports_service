package contest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mbanner"
	"sports_service/server/models/mcontest"
	"sports_service/server/models/muser"
	"sports_service/server/util"
	"time"
	"fmt"
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

// 获取直播列表 默认取两条最新的 todo:暂时只有一个赛事
func (svc *ContestModule) GetLiveList(status string, page, size int) (int, []*mcontest.ContestLiveInfo) {
	offset := (page - 1) * size
	now := time.Now().Unix()
	list, err := svc.contest.GetLiveList(now, offset, size, "1", status)
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
			Status: item.Status,
			HasReplay: 2,
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

// 获取赛程信息
func (svc *ContestModule) GetScheduleInfo() (int, *mcontest.ContestInfo) {
	if err := svc.GetContestInfo(); err != nil {
		log.Log.Errorf("contest_trace: get contest info fail, err:%s", err)
		return errdef.CONTEST_INFO_FAIL, nil
	}

	schedule, err := svc.contest.GetScheduleInfoByContestId(fmt.Sprint(svc.contest.Contest.Id))
	if err != nil {
		log.Log.Errorf("contest_trace: get schedule info fail, contestId:%d, err:%s", svc.contest.Contest.Id, err)
		return errdef.CONTEST_SCHEDULE_FAIL, nil
	}

	if schedule == nil {
		return errdef.SUCCESS, nil
	}

	resp := &mcontest.ContestInfo{ScheduleList: make([]*mcontest.ScheduleInfo, len(schedule))}
	resp.ContestId = svc.contest.Contest.Id
	resp.ContestName = svc.contest.Contest.ContestName
	for index, item := range schedule {
		info := &mcontest.ScheduleInfo{
			ScheduleId: item.Id,
			ScheduleName: item.ScheduleName,
			Description: item.ScheduleDescription,
			ShowType: item.ShowType,
		}

		resp.ScheduleList[index] = info
	}

	return errdef.SUCCESS, resp
}

// 获取最新的一个赛事
func (svc *ContestModule) GetContestInfo() error {
	ok, err := svc.contest.GetContestInfo(time.Now().Unix())
	if !ok || err != nil {
		return err
	}

	return nil
}

// 获取赛程晋级信息
func (svc *ContestModule) GetPromotionInfo(contestId, scheduleId string) (int, interface{}) {
	ok, err := svc.contest.GetScheduleInfoById(scheduleId)
	if !ok || err != nil {
		return errdef.CONTEST_SCHEDULE_FAIL, nil
	}

	switch svc.contest.Schedule.ShowType {
	// 列表展示
	case 1:
		list, err := svc.contest.GetScheduleDetailInfo(1, contestId, scheduleId)
		if err != nil {
			log.Log.Errorf("contest_trace: get promotion info fail, scheduleId:%s, err", scheduleId, err)
			return errdef.CONTEST_PROMOTION_INFO_FAIL, nil
		}

		mp := make(map[int64]*mcontest.ScheduleListDetailResp)
		ranking := 0
		for _, item := range list {
			// key 选手id
			if _, ok :=  mp[item.PlayerId]; !ok {
				detail := &mcontest.ScheduleListDetailResp{}
				detail.PlayerId = item.PlayerId
				detail.PlayerName = item.PlayerName
				detail.ContestId = item.ContestId
				detail.ScheduleId = item.ScheduleId
				//detail.IsWin = item.IsWin
				detail.Photo = item.Photo
				detail.BestScore = fmt.Sprintf("%.3f", float64(item.Score)/1000)
				if item.Rounds == 1 {
					detail.RoundOneScore = fmt.Sprintf("%.3f", float64(item.Score)/1000)
				}

				if item.Rounds == 2 {
					detail.RoundTwoScore = fmt.Sprintf("%.3f", float64(item.Score)/1000)
				}

				if item.Rounds == 3 {
					detail.RoundThreeScore = fmt.Sprintf("%.3f", float64(item.Score)/1000)
				}

				detail.Ranking = ranking
				ranking++
				mp[item.PlayerId] = detail
			} else {
				if item.Rounds == 1 {
					mp[item.PlayerId].RoundOneScore = fmt.Sprintf("%.3f", float64(item.Score)/1000)
				}

				if item.Rounds == 2 {
					mp[item.PlayerId].RoundTwoScore = fmt.Sprintf("%.3f", float64(item.Score)/1000)
				}

				if item.Rounds == 3 {
					mp[item.PlayerId].RoundThreeScore = fmt.Sprintf("%.3f", float64(item.Score)/1000)
				}
			}
		}

		// 防止数组越界
		if ranking > len(mp) {
			return errdef.CONTEST_PROMOTION_INFO_FAIL, nil
		}

		log.Log.Errorf("##########:len(map)", len(mp))
		resp := make([]*mcontest.ScheduleListDetailResp, len(mp))
		for _, val := range mp {
			log.Log.Infof("#######val:%+v", val)
			resp[val.Ranking] = val
		}

		return errdef.SUCCESS, resp

	// 分组展示
	case 2:
		list, err := svc.contest.GetScheduleDetailInfo(2, contestId, scheduleId)
		if err != nil {
			log.Log.Errorf("contest_trace: get promotion info fail, scheduleId:%s, err", scheduleId, err)
			return errdef.CONTEST_PROMOTION_INFO_FAIL, nil
		}

		index := 0
		mp := make(map[int]*mcontest.ScheduleGroupDetailResp, 0)
		for _, item := range list {
			var detail *mcontest.ScheduleGroupDetailResp
			// key 分组id
			if _, ok :=  mp[item.GroupNum]; !ok {
				detail = &mcontest.ScheduleGroupDetailResp{Player:  make([]mcontest.PlayerInfoResp, 0),
					Winner: make([]mcontest.PlayerInfoResp, 0)}
				detail.GroupNum = item.GroupNum
				detail.ScheduleId = item.ScheduleId
				detail.GroupName = item.GroupName
				detail.ContestId = item.ContestId
				detail.Index = index
				detail.BeginTm = time.Unix(int64(item.BeginTm), 0).Format(consts.FORMAT_DATE_STR)
				mp[item.GroupNum] = detail
				index++
			}

			player := mcontest.PlayerInfoResp{
				Id: item.Id,
				PlayerId: item.PlayerId,
				PlayerName: item.PlayerName,
				Photo: item.Photo,
				IsWin: item.IsWin,
				Score: fmt.Sprintf("%.3f", float64(item.Score)/1000),
				NumInGroup: item.NumInGroup,
			}

			mp[item.GroupNum].Player = append(mp[item.GroupNum].Player, player)
			// 1 表示胜出
			if item.IsWin == 1 {
				mp[item.GroupNum].Winner = append(mp[item.GroupNum].Winner, player)
			}
		}

		resp := make([]*mcontest.ScheduleGroupDetailResp, len(mp))
		for _, val := range mp {
			resp[val.Index] = val
		}

		return errdef.SUCCESS, resp

	}

	return errdef.ERROR, nil
}


// 获取赛事选手积分排行
func (svc *ContestModule) GetIntegralRanking(contestId string, page, size int) (int, []*mcontest.IntegralRanking) {
	offset := (page - 1) * size
	list, err := svc.contest.GetIntegralRankingByContestId(contestId, offset, size)
	if err != nil {
		return errdef.CONTEST_RANKING_FAIL, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*mcontest.IntegralRanking{}
	}

	for index, item := range list {
		item.Ranking = offset + index
		item.TotalIntegralStr = fmt.Sprintf("%.3f", float64(item.TotalIntegral) / 1000)
		item.BestScoreStr = fmt.Sprintf("%.3f", float64(item.BestScore) / 1000)
		item.TotalIntegral = 0
		item.BestScore = 0
	}

	return errdef.SUCCESS, list
}

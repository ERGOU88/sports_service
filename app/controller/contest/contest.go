package contest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sort"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mbanner"
	"sports_service/server/models/mcommunity"
	"sports_service/server/models/mcontest"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"sports_service/server/tools/tencentCloud"
	"sports_service/server/util"
	"time"
)

type ContestModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	banner      *mbanner.BannerModel
	contest     *mcontest.ContestModel
	community   *mcommunity.CommunityModel
	video       *mvideo.VideoModel
}

func New(c *gin.Context) ContestModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	return ContestModule{
		context: c,
		user: muser.NewUserModel(socket),
		banner: mbanner.NewBannerMolde(socket),
		contest: mcontest.NewContestModel(socket),
		community: mcommunity.NewCommunityModel(socket),
		video: mvideo.NewVideoModel(socket),
		engine: socket,
	}
}

// 获取赛事首页banner
func (svc *ContestModule) GetBanner() []*mbanner.Banner {
	banners := svc.banner.GetRecommendBanners(int32(consts.CONTEST_BANNERS), time.Now().Unix(), 0, 10)
	if len(banners) == 0 {
		return []*mbanner.Banner{}
	}

	return banners
}

// 获取推荐的直播列表 默认取2条同一天内最新的未开播/直播中的数据 todo:暂时只有一个赛事
// pullType 拉取类型  1 上拉加载 今天及未来赛事数据 [通过开播时间作为查询条件进行拉取] 2 下拉加载 历史赛事数据 [通过开播时间作为查询条件进行拉取] 默认上拉加载
// queryType 1 首页列表 [查询最近同一天内的 未开播/直播中的数据]
func (svc *ContestModule) GetLiveList(queryType, pullType, ts string, page, size int) (int, []*mcontest.ContestLiveInfo, int, int) {
	if page > 1 {
		page = 1
	}

	if err := svc.GetContestInfo(); err != nil {
		log.Log.Errorf("contest_trace: get contest info fail, err:%s", err)
		return errdef.CONTEST_INFO_FAIL, nil, 0, 0
	}

	offset := (page - 1) * size
	limitTm := ts
	if queryType == "1" || limitTm == "" {
		limitTm = fmt.Sprint(time.Now().Unix())
	}

	log.Log.Errorf("limitTm:%d", limitTm)
	list, err := svc.contest.GetLiveList(offset, size, fmt.Sprint(svc.contest.Contest.Id), limitTm, queryType, pullType)
	if err != nil {
		log.Log.Errorf("contest_trace: get live list fail, err:%s", err)
		return errdef.CONTEST_GET_LIVE_FAIL, nil, 0, 0
	}

	if len(list) == 0 {
		if pullType != "1" || ts != "" {
			return errdef.SUCCESS, []*mcontest.ContestLiveInfo{}, 0, 0
		}
		// 如果是取今天及未来赛事数据
		// 取历史数据
		pullType = "2"
		limitTm = fmt.Sprint(time.Now().Unix())
		list, err = svc.contest.GetLiveList(offset, size, fmt.Sprint(svc.contest.Contest.Id), limitTm, queryType, pullType)
		if err != nil {
			log.Log.Errorf("contest_trace: get live list fail, err:%s", err)
			return errdef.CONTEST_GET_LIVE_FAIL, nil, 0, 0
		}

		if len(list) == 0 {
			return errdef.SUCCESS, []*mcontest.ContestLiveInfo{}, 0, 0
		}
	}


	mp := make(map[string]*mcontest.ContestLiveInfo)
	var (
		pullUpTm, pullDownTm, index int
	)
	for _, item := range list {
		var detail *mcontest.ContestLiveInfo
		tm := time.Unix(int64(item.PlayTime), 0)
		date := tm.Format("1月2日")
		if _, ok := mp[date]; !ok {
			// queryType 为 1  只取最近同一天内的赛事
			if index == 1 && queryType == "1" {
				continue
			}

			detail = &mcontest.ContestLiveInfo{
				Date: date,
				Week: util.GetWeekCn(int(tm.Weekday())),
				Index: index,
				LiveInfo: make([]*mcontest.LiveInfo, 0),
			}

			if time.Now().Format("1月2日") == date {
				detail.IsToday = true
			}

			mp[date] = detail
			index++
		}

		live := &mcontest.LiveInfo{
			Id: item.Id,
			UserId: item.UserId,
			RoomId: item.RoomId,
			GroupId: item.GroupId,
			Cover: tencentCloud.BucketURI(item.Cover),
			RtmpAddr: item.RtmpAddr,
			HlsAddr: item.HlsAddr,
			FlvAddr: item.FlvAddr,
			PlayTime: item.PlayTime,
			Title: item.Title,
			Subhead: item.Subhead,
			HighLights: item.HighLights,
			Describe: item.Describe,
			Tags: item.Tags,
			LiveType: item.LiveType,
			Status: item.Status,
			// 默认无回放
			HasReplay: 2,
		}

		// 上拉加载 使用最大时间 拉取未来数据
		if pullUpTm == 0 || item.PlayTime > pullUpTm {
			pullUpTm = item.PlayTime
		}

		// 下拉加载 使用最小时间 拉取历史数据
 		if pullDownTm == 0 ||  item.PlayTime < pullDownTm {
			pullDownTm = item.PlayTime
		}

		// 如果直播已结束
		if item.Status == consts.LIVE_STATUS_END {
			live.HasReplay = 1
			// 获取回放数据
			svc.GetLiveReplayInfo(item.Id, live)
		}

		user := svc.user.FindUserByUserid(item.UserId)
		if user != nil {
			live.NickName = user.NickName
			live.Avatar = tencentCloud.BucketURI(user.Avatar)
		}

		mp[date].LiveInfo = append(mp[date].LiveInfo, live)
	}

	resp := make([]*mcontest.ContestLiveInfo, len(mp))
	for _, val := range mp {
		resp[val.Index] = val
	}

	if pullType == "2" {
		// 统一为正序返回
		sort.Sort(SortContestLive(resp))
	}

	return errdef.SUCCESS, resp, pullUpTm, pullDownTm
}

// 获取直播数量
func (svc *ContestModule) GetLiveCount() int64 {
	count, err := svc.contest.GetVideoLiveCount()
	if err != nil {
		log.Log.Errorf("contest_trace: get live count fail, err:%s", err)
		return 0
	}

	return count
}

// 获取直播回放信息
func (svc *ContestModule) GetLiveReplayInfo(id int64, live *mcontest.LiveInfo) {
	ok, err := svc.contest.GetVideoLiveReply(fmt.Sprint(id))
	if !ok || err != nil {
		log.Log.Errorf("contest_trace: get video live reply fail, liveId:%d, ok:%v, err:%s", id, ok, err)
	}

	// 存在回放数据
	if ok {
		live.HasReplay = 1
		replay := &mcontest.LiveReplayInfo{
			Id: svc.contest.VideoLiveReplay.Id,
			Size: svc.contest.VideoLiveReplay.Size,
			LiveId: svc.contest.VideoLiveReplay.LiveId,
			Describe: svc.contest.VideoLiveReplay.Describe,
			Duration: svc.contest.VideoLiveReplay.Duration,
			CreateAt: svc.contest.VideoLiveReplay.CreateAt,
			HistoryAddr: svc.video.AntiStealingLink(svc.contest.VideoLiveReplay.HistoryAddr),
			Title: svc.contest.VideoLiveReplay.Title,
			Subhead: svc.contest.VideoLiveReplay.Subhead,
			PlayNum: svc.contest.VideoLiveReplay.PlayNum,
		}

		if svc.contest.VideoLiveReplay.PlayInfo != "" {
			if err = util.JsonFast.UnmarshalFromString(svc.contest.VideoLiveReplay.PlayInfo, &replay.PlayInfo); err != nil {
				log.Log.Errorf("contest_trace: unmarshal playInfo fail, id:%d, err:%s", svc.contest.VideoLiveReplay.Id, err)
			}
		}

		if len(replay.PlayInfo) == 0 {
			replay.PlayInfo = []*mcontest.PlayInfo{}
		} else {
			for _, v := range replay.PlayInfo {
				// 添加防盗链
				v.Url = svc.video.AntiStealingLink(v.Url)
			}
		}

		live.LiveReplayInfo = replay
	}
}

// 获取赛程信息
func (svc *ContestModule) GetScheduleInfo() (int, *mcontest.ContestInfo) {
	if err := svc.GetContestInfo(); err != nil {
		log.Log.Errorf("contest_trace: get contest info fail, err:%s", err)
		return errdef.CONTEST_INFO_FAIL, nil
	}

	schedule, err := svc.contest.GetScheduleInfoByContestId(time.Now().Unix(), fmt.Sprint(svc.contest.Contest.Id))
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
		index := 0
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
				detail.BestScore = util.ResolveTimeByMilliSecond(item.Score)
				if item.Rounds == 1 {
					detail.RoundOneScore = util.ResolveTimeByMilliSecond(item.Score)
				}

				if item.Rounds == 2 {
					detail.RoundTwoScore = util.ResolveTimeByMilliSecond(item.Score)
				}

				if item.Rounds == 3 {
					detail.RoundThreeScore = util.ResolveTimeByMilliSecond(item.Score)
				}

				if item.Ranking > 0 {
					detail.Ranking = item.Ranking - 1
				} else {
					detail.Ranking = index
				}
				
				detail.Index = index
				index++
				mp[item.PlayerId] = detail
			} else {
				if item.Rounds == 1 {
					mp[item.PlayerId].RoundOneScore = util.ResolveTimeByMilliSecond(item.Score)
				}

				if item.Rounds == 2 {
					mp[item.PlayerId].RoundTwoScore = util.ResolveTimeByMilliSecond(item.Score)
				}

				if item.Rounds == 3 {
					mp[item.PlayerId].RoundThreeScore = util.ResolveTimeByMilliSecond(item.Score)
				}
			}
		}

		// 防止数组越界
		if index > len(mp) {
			return errdef.CONTEST_PROMOTION_INFO_FAIL, nil
		}

		log.Log.Errorf("##########:len(map)", len(mp))
		resp := make([]*mcontest.ScheduleListDetailResp, len(mp))
		for _, val := range mp {
			log.Log.Infof("#######val:%+v", val)
			resp[val.Index] = val
		}

		return errdef.SUCCESS, resp

	// 分组竞技/总决赛 展现形式
	case 2,3:
		playerMp := make(map[int64]bool, 0)
		list, err := svc.contest.GetScheduleDetailInfo(2, contestId, scheduleId)
		if err != nil {
			log.Log.Errorf("contest_trace: get promotion info fail, scheduleId:%s, err", scheduleId, err)
			return errdef.CONTEST_PROMOTION_INFO_FAIL, nil
		}

		index := 0
		mp := make(map[int]*mcontest.ScheduleGroupDetailResp, 0)
		for _, item := range list {
			if _, ok := playerMp[item.PlayerId]; ok {
				continue
			} else {
				playerMp[item.PlayerId] = true
			}

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
				Score: util.ResolveTimeByMilliSecond(item.Score),
				NumInGroup: item.NumInGroup,
				Integral: "0",
			}

			ok, err := svc.contest.GetTotalIntegralByPlayer(fmt.Sprint(item.ContestId), fmt.Sprint(item.PlayerId))
			if ok && err == nil {
				player.Integral = fmt.Sprintf("%.3f", float64(svc.contest.IntegralRanking.TotalIntegral) / 1000)
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
		item.BestScoreStr = util.ResolveTimeByMilliSecond(item.BestScore)
		item.TotalIntegral = 0
		item.BestScore = 0
	}

	return errdef.SUCCESS, list
}

// 获取赛事的社区板块
func (svc *ContestModule) GetContestSection() (int, int) {
	ok, err := svc.community.GetSectionByName("赛事")
	if !ok || err != nil {
		log.Log.Errorf("contest_trace: get section by name fail, err:%s", err)
		return errdef.ERROR, 0
	}

	return errdef.SUCCESS, svc.community.CommunitySection.Id
}

// 获取赛程直播 选手竞赛数据
func (svc *ContestModule) GetLiveScheduleData(liveId string, page, size int) (int, []*mcontest.LiveSchedulePlayerData) {
	if liveId == "" {
		return errdef.INVALID_PARAMS, nil
	}

	offset := (page - 1) * size
	list, err := svc.contest.GetLiveSchedulePlayerData(liveId, offset, size)
	if err != nil {
		log.Log.Errorf("contest_trace: get live schedule player data fail, err:%s", err)
		return errdef.CONTEST_LIVE_SCHEDULE_DATA, nil
	}

	if list == nil {
		return errdef.SUCCESS, []*mcontest.LiveSchedulePlayerData{}
	}

	resp := make([]*mcontest.LiveSchedulePlayerData, len(list))
	for index, item := range list {
		info := &mcontest.LiveSchedulePlayerData{
			Id: item.Id,
			ContestId: item.ContestId,
			ScheduleId: item.ScheduleId,
			PlayerId: item.PlayerId,
			LiveId: item.LiveId,
			RoundsNum: item.RoundsNum,
			Ranking: index+1,
			IntervalDuration: fmt.Sprintf("%.3f", float64(item.IntervalDuration)/1000),
			TopSpeed: fmt.Sprintf("%.3f", float64(item.IntervalDuration)/1000),
			ReceiveIntegral: fmt.Sprintf("%.3f", float64(item.ReceiveIntegral)/1000),
		}

		ok, err := svc.contest.GetPlayerInfoById(fmt.Sprint(info.PlayerId))
		if !ok || err != nil {
			log.Log.Errorf("contest_trace: get player info by id fail, playerId:%d, err:%s", info.PlayerId, err)
			continue
		}

		info.PlayerName = svc.contest.PlayerInfo.Name
		info.Photo = tencentCloud.BucketURI(svc.contest.PlayerInfo.Photo)

		resp[index] = info
	}

	return errdef.SUCCESS, resp
}

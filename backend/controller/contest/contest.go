package contest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/models"
	"sports_service/server/models/mcontest"
	"sports_service/server/global/backend/log"
	"sports_service/server/tools/im"
	"sports_service/server/tools/live"
	"sports_service/server/util"
	"time"
	"fmt"
)

type ContestModule struct {
	context     *gin.Context
	engine      *xorm.Session
	contest     *mcontest.ContestModel
}

func New(c *gin.Context) ContestModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	return ContestModule{
		context: c,
		contest: mcontest.NewContestModel(socket),
		engine: socket,
	}
}

func (svc *ContestModule) AddPlayer(player *models.FpvContestPlayerInformation) int {
	// 最新赛事
	ok, err := svc.contest.GetContestInfo(time.Now().Unix())
	if ok && err == nil {
		player.ContestId = svc.contest.Contest.Id
	}

	if _, err := svc.contest.AddPlayer(player); err != nil {
		return errdef.ERROR
	}

	return errdef.SUCCESS
}

func (svc *ContestModule) UpdatePlayer(player *models.FpvContestPlayerInformation) int {
	if _, err := svc.contest.UpdatePlayer(player); err != nil {
		return errdef.ERROR
	}

	return errdef.SUCCESS
}

func (svc *ContestModule) GetPlayerList(page, size int) (int, []*models.FpvContestPlayerInformation) {
	offset := (page - 1) * size
	list, err := svc.contest.GetPlayerList(offset, size)
	if err != nil {
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*models.FpvContestPlayerInformation{}
	}

	return errdef.SUCCESS, list
}

// 添加组别配置
func (svc *ContestModule) AddContestGroup(group *models.FpvContestScheduleGroup) int {
	// 最新赛事
	ok, err := svc.contest.GetContestInfo(time.Now().Unix())
	if ok && err == nil {
		group.ContestId = svc.contest.Contest.Id
	}

	if _, err := svc.contest.AddContestGroup(group); err != nil {
		log.Log.Errorf("contest_trace: add contest group fail, err:%s", err)
		return errdef.ERROR
	}

	return errdef.SUCCESS
}

// 编辑组别配置
func (svc *ContestModule) EditContestGroup(group *models.FpvContestScheduleGroup) int {
	// 最新赛事
	ok, err := svc.contest.GetContestInfo(time.Now().Unix())
	if ok && err == nil {
		group.ContestId = svc.contest.Contest.Id
	}

	if _, err := svc.contest.UpdateContestGroup(group); err != nil {
		log.Log.Errorf("contest_trace: add contest group fail, err:%s", err)
		return errdef.ERROR
	}

	return errdef.SUCCESS
}

// 获取赛事分组配置列表
func (svc *ContestModule) GetContestGroupList(page, size int, scheduleId, contestId string) (int, []*models.FpvContestScheduleGroup) {
	offset := (page - 1) * size
	list, err := svc.contest.GetContestGroupList(offset, size, scheduleId, contestId)
	if err != nil {
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*models.FpvContestScheduleGroup{}
	}

	return errdef.SUCCESS, list
}

// 获取赛程信息
func (svc *ContestModule) GetContestScheduleInfo() (int, []*models.FpvContestSchedule) {
	list, err := svc.contest.GetScheduleInfo()
	if err != nil {
		log.Log.Errorf("contest_trace: get schedule info fail, err:%s", err)
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*models.FpvContestSchedule{}
	}

	return errdef.SUCCESS, list
}

// 添加赛程详情
func (svc *ContestModule) AddContestScheduleDetail(param *mcontest.AddScheduleDetail) int {
	// 最新赛事
	ok, err := svc.contest.GetContestInfo(time.Now().Unix())
	if ok && err == nil {
		param.ContestId = svc.contest.Contest.Id
	}

	ok, err = svc.contest.GetScheduleInfoById(fmt.Sprint(param.ScheduleId))
	if !ok || err != nil {
		log.Log.Errorf("contest_trace: get schedule info by id fail, scheduleId:%d, err:%s", param.ScheduleId, err)
		return errdef.ERROR
	}

	if svc.contest.Schedule.RoundsNum > 0 {
		now := int(time.Now().Unix())
		list := make([]*models.FpvContestScheduleDetail, svc.contest.Schedule.RoundsNum)
		for round := 0; round < svc.contest.Schedule.RoundsNum; round++ {
			detail := &models.FpvContestScheduleDetail{
				Rounds: round+1,
				ScheduleId: param.ScheduleId,
				GroupNum: param.GroupNum,
				GroupName: param.GroupName,
				NumInGroup: param.NumInGroup,
				PlayerId: param.PlayerId,
				BeginTm: param.BeginTm,
				EndTm: param.EndTm,
				IsWin: param.IsWin,
				CreateAt: now,
				UpdateAt: now,
				Ranking: param.Ranking,
			}

			if round == 0 {
				detail.Score = param.RoundOneScore
				detail.ReceiveIntegral = param.RoundOneIntegral
			}

			if round == 1 {
				detail.Score = param.RoundTwoScore
				detail.ReceiveIntegral = param.RoundTwoIntegral
			}

			if round == 2 {
				detail.Score = param.RoundThreeScore
			}

			list[round] = detail
		}

		affected, err := svc.contest.AddContestScheduleDetail(list)
		if int(affected) != svc.contest.Schedule.RoundsNum || err != nil {
			log.Log.Errorf("contest_trace: add contest schedule detail fail, err:%s", err)
			return errdef.ERROR
		}
	}

	return errdef.SUCCESS
}

// 获取赛事赛程详情列表
func (svc *ContestModule) GetContestScheduleDetailList(scheduleId string) (int, []*mcontest.ScheduleListDetailResp) {
	var contestId int
	// 最新赛事
	ok, err := svc.contest.GetContestInfo(time.Now().Unix())
	if ok && err == nil {
		contestId = svc.contest.Contest.Id
	}

	list, err := svc.contest.GetScheduleDetailInfo(3, fmt.Sprint(contestId), scheduleId)
	if err != nil {
		log.Log.Errorf("contest_trace: get promotion info fail, scheduleId:%s, err", scheduleId, err)
		return errdef.ERROR, nil
	}

	mp := make(map[int64]*mcontest.ScheduleListDetailResp)
	ranking := 0
	for _, item := range list {
		// key 选手id
		if _, ok :=  mp[item.PlayerId]; !ok {
			detail := &mcontest.ScheduleListDetailResp{}
			detail.Id = item.Id
			detail.PlayerId = item.PlayerId
			detail.PlayerName = item.PlayerName
			detail.ContestId = item.ContestId
			detail.ScheduleId = item.ScheduleId
			detail.IsWin = item.IsWin
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
		return errdef.ERROR, nil
	}

	log.Log.Errorf("##########:len(map)", len(mp))
	resp := make([]*mcontest.ScheduleListDetailResp, len(mp))
	for _, val := range mp {
		log.Log.Infof("#######val:%+v", val)
		resp[val.Ranking] = val
	}

	return errdef.SUCCESS, resp
}

// 设置赛事积分榜
func (svc *ContestModule) SetIntegralRanking(info *models.FpvContestPlayerIntegralRanking) int {
	// 最新赛事
	ok, err := svc.contest.GetContestInfo(time.Now().Unix())
	if ok && err == nil {
		info.ContestId = svc.contest.Contest.Id
	}

	if _, err := svc.contest.SetIntegralRanking(info); err != nil {
		return errdef.ERROR
	}

	return errdef.SUCCESS
}


func (svc *ContestModule) UpdateIntegralRanking(info *models.FpvContestPlayerIntegralRanking) int {
	// 最新赛事
	ok, err := svc.contest.GetContestInfo(time.Now().Unix())
	if ok && err == nil {
		info.ContestId = svc.contest.Contest.Id
	}

	if _, err := svc.contest.UpdateIntegralRanking(info); err != nil {
		return errdef.ERROR
	}

	return errdef.SUCCESS
}

func (svc *ContestModule) GetIntegralRankingList(page, size int) (int, []*mcontest.IntegralRanking) {
	// 最新赛事
	ok, err := svc.contest.GetContestInfo(time.Now().Unix())
	if !ok || err != nil {
		return errdef.ERROR, nil
	}

	offset := (page - 1) * size
	list, err := svc.contest.GetIntegralRankingByContestId(fmt.Sprint(svc.contest.Contest.Id), offset, size)
	if err != nil {
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*mcontest.IntegralRanking{}
	}

	for _, item := range list {
		item.TotalIntegralStr = fmt.Sprintf("%.3f", float64(item.TotalIntegral) / 1000)
		item.BestScoreStr = fmt.Sprintf("%.3f", float64(item.BestScore) / 1000)

		item.TotalIntegralStr = util.ResolveTimeByMilliSecond(item.TotalIntegral)
		item.BestScoreStr = util.ResolveTimeByMilliSecond(item.BestScore)
		item.TotalIntegral = 0
		item.BestScore = 0
	}

	return errdef.SUCCESS, list
}

// 添加赛事直播
func (svc *ContestModule) AddContestLive(info *models.VideoLive) int {
	now := int(time.Now().Unix())
	if info.PlayTime < now {
		return errdef.INVALID_PARAMS
	}
	var err error
	info.GroupId, err = im.Im.CreateGroup("AVChatRoom", "", info.Title, info.Describe, "",
		"")
	if err != nil {
		log.Log.Errorf("contest_trace: create group fail, err:%s", err)
		return errdef.ERROR
	}

	// 最新赛事
	ok, err := svc.contest.GetContestInfo(time.Now().Unix())
	if ok && err == nil {
		info.ContestId = svc.contest.Contest.Id
	}

	info.CreateAt = now
	info.UpdateAt = now

	duration := info.PlayTime - now
	// todo: 过期时间待确认
	expireTm := int64(duration + 86400 * 30)
	roomId := fmt.Sprint(util.GetXID())
	info.RoomId = roomId
	info.PushStreamUrl, info.PushStreamKey = live.Live.GenPushStream(roomId, expireTm)
	streamInfo := live.Live.GenPullStream(roomId, expireTm)
	info.RtmpAddr = streamInfo.RtmpAddr
	info.HlsAddr = streamInfo.HlsAddr
	info.FlvAddr = streamInfo.FlvAddr

	if _, err := svc.contest.AddContestLive(info); err != nil {
		return errdef.ERROR
	}

	return errdef.SUCCESS
}

func (svc *ContestModule) UpdateContestLive(live *models.VideoLive) int {
	if _, err := svc.contest.UpdateContestLive(live); err != nil {
		return errdef.ERROR
	}

	return errdef.SUCCESS
}


func (svc *ContestModule) DelContestLive(id string) int {
	if _, err := svc.contest.DelContestLive(id); err != nil {
		return errdef.ERROR
	}

	return errdef.SUCCESS
}

func (svc *ContestModule) GetContestLiveList(page, size int) (int, []*models.VideoLive) {
	offset := (page - 1) * size
	list, err := svc.contest.GetContestLiveList(offset, size)
	if err != nil {
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*models.VideoLive{}
	}

	return errdef.SUCCESS, list
}

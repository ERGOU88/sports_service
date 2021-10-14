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

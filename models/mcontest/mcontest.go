package mcontest

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

// 赛事
type ContestModel struct {
	Engine          *xorm.Session
	VideoLive       *models.VideoLive
	Schedule        *models.FpvContestSchedule
	Contest         *models.FpvContestInfo
	IntegralRanking *models.FpvContestPlayerIntegralRanking
	PlayerInfo      *models.FpvContestPlayerInformation
}

// 赛事直播信息
type ContestLiveInfo struct {
	Id             int64  `json:"id"`
	UserId         string `json:"user_id"`
	RoomId         string `json:"room_id"`
	GroupId        string `json:"group_id"`
	Cover          string `json:"cover"`
	RtmpAddr       string `json:"rtmp_addr"`
	FlvAddr        string `json:"flv_addr"`
	HlsAddr        string `json:"hls_addr"`
	PlayTime       int    `json:"play_time"`
	Title          string `json:"title"`
	HighLights     string `json:"high_lights"`
	Describe       string `json:"describe"`
	Tags           string `json:"tags"`
	LiveType       int    `json:"live_type"`
	NickName       string `json:"nick_name"`
	Avatar         string `json:"avatar"`
	Date           string `json:"date"`
	Week           string `json:"week"`
	Status         int    `json:"int"`           // 状态 0未直播 1直播中 2 已结束
	HasReplay      int    `json:"has_replay"`    // 是否有回放 1 有 2 无
}

// 赛事信息[包含赛程]
type ContestInfo struct {
	ContestId      int       `json:"contest_id"`
	ContestName    string    `json:"contest_name"`

	ScheduleList   []*ScheduleInfo `json:"schedule_info"`          // 赛程列表
}

// 赛程信息
type ScheduleInfo struct {
	ScheduleId     int       `json:"schedule_id"`
	ScheduleName   string    `json:"schedule_name"`
	Description    string    `json:"description"`
}

// 实例
func NewContestModel(engine *xorm.Session) *ContestModel {
	return &ContestModel{
		Engine: engine,
		VideoLive: new(models.VideoLive),
		Schedule: new(models.FpvContestSchedule),
		Contest: new(models.FpvContestInfo),
		IntegralRanking: new(models.FpvContestPlayerIntegralRanking),
		PlayerInfo: new(models.FpvContestPlayerInformation),
	}
}

// 获取直播列表
func (m *ContestModel) GetLiveList(now int64, offset, size int, contestId, status string) ([]*models.VideoLive, error){
	var list []*models.VideoLive
	m.Engine.Where("play_time >= ? AND contest_id=?", now, contestId)
	if status == "1" {
		m.Engine.Where("status=1")
	}

	if err := m.Engine.Asc("play_time").Limit(size, offset).Find(&list); err != nil {
		return []*models.VideoLive{}, err
	}

	return list, nil
}

// 通过赛事id获取赛程信息
func (m *ContestModel) GetScheduleInfoByContestId(contestId string) ([]*models.FpvContestSchedule, error) {
	var list []*models.FpvContestSchedule
	if err := m.Engine.Where("contest_id=? AND status=0", contestId).Asc("order").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取最新一个赛事信息
func (m *ContestModel) GetContestInfo(now int64) (bool, error) {
	m.Contest = new(models.FpvContestInfo)
	return m.Engine.Where("start_tm<=? AND end_tm>=?", now, now).Desc("id").Limit(1).Get(m.Contest)
}

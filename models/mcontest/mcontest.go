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
	ScheduleDetail  *models.FpvContestScheduleDetail
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

// 赛程详情
type ScheduleDetail struct {
	Id         int64  `json:"id"`
	ScheduleId int    `json:"schedule_id"`
	Rounds     int    `json:"rounds"`
	GroupNum   int    `json:"group_num"`
	GroupName  string `json:"group_name"`
	PlayerId   int64  `json:"player_id"`
	PlayerName string `json:"name"`
	Photo      string `json:"photo"`
	Score      int    `json:"score"`
	IsWin      int    `json:"is_win"`
	NumInGroup int    `json:"num_in_group"`
	ContestId  int    `json:"contest_id"`
}


// 赛程详情返回数据
type ScheduleDetailResp struct {
	Id         int64  `json:"id,omitempty"`
	ScheduleId int    `json:"schedule_id"`
	PlayerId   int64  `json:"player_id"`
	PlayerName string `json:"player_name"`
	Photo      string `json:"photo"`
	//IsWin      int    `json:"is_win"`
	ContestId  int    `json:"contest_id"`
	Ranking    int    `json:"ranking"`

	BestScore          string    `json:"best_score"`
	RoundOneScore      string    `json:"round_one_score"`
	RoundTwoScore      string    `json:"round_two_score"`
	RoundThreeScore    string    `json:"round_three_score"`
}

type IntegralRanking struct {
	PlayerId   int64  `json:"player_id"`
	PlayerName string `json:"player_name"`
	Photo      string `json:"photo"`
	ContestId  int    `json:"contest_id"`

	Ranking            int       `json:"ranking"`
	TotalIntegral      int       `json:"integral,omitempty"`
	BestScore          int       `json:"score,omitempty"`

	TotalIntegralStr   string    `json:"total_integral"`
	BestScoreStr       string    `json:"best_score"`
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
		ScheduleDetail: new(models.FpvContestScheduleDetail),
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

// 通过id获取赛程信息
func (m *ContestModel) GetScheduleInfoById(scheduleId string) (bool, error) {
	m.Schedule = new(models.FpvContestSchedule)
	return m.Engine.Where("id=?", scheduleId).Get(m.Schedule)
}

// 获取最新一个赛事信息
func (m *ContestModel) GetContestInfo(now int64) (bool, error) {
	m.Contest = new(models.FpvContestInfo)
	return m.Engine.Where("start_tm<=? AND end_tm>=?", now, now).Desc("id").Limit(1).Get(m.Contest)
}

// 获取赛程详情
func (m *ContestModel) GetScheduleDetail(contestId, scheduleId string) ([]*models.FpvContestScheduleDetail, error) {
	var list []*models.FpvContestScheduleDetail
	if err := m.Engine.Where("contest_id=? AND schedule_id=? AND status=0", contestId, scheduleId).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

const (
	GET_SCHEDULE_DETAIL_BY_SCORE = "SELECT p.id AS player_id, p.name AS player_name, p.photo, cs.* FROM fpv_contest_player_information AS p " +
		"LEFT JOIN fpv_contest_schedule_detail AS cs ON p.id = cs.player_id AND cs.contest_id=? AND cs.schedule_id=? " +
		" WHERE p.status = 0 ORDER BY cs.score is null, cs.score ASC, p.id ASC"
)
// 获取赛程信息[成绩正序]
func (m *ContestModel) GetScheduleDetailByScore(contestId, scheduleId string) ([]*ScheduleDetail, error) {
	var list []*ScheduleDetail
	if err := m.Engine.SQL(GET_SCHEDULE_DETAIL_BY_SCORE, contestId, scheduleId).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取赛事参赛选手列表
func (m *ContestModel) GetPlayerByContestId(contestId string) ([]*models.FpvContestPlayerInformation, error) {
	var list []*models.FpvContestPlayerInformation
	if err := m.Engine.Where("contest_id=? AND status=0", contestId).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

const (
	GET_INTEGRAL_RANKING = "SELECT p.id AS player_id, p.name AS player_name, p.photo, rk.contest_id, " +
		"rk.total_integral, rk.best_score FROM fpv_contest_player_information AS p " +
		"LEFT JOIN fpv_contest_player_integral_ranking AS rk " +
		"ON p.id = rk.player_id AND rk.contest_id=? " +
		"WHERE p.status = 0 ORDER BY rk.total_integral DESC, p.id ASC LIMIT ?, ?"
)
// 通过赛事id 获取选手积分排行
func (m *ContestModel) GetIntegralRankingByContestId(contestId string, offset, size int) ([]*IntegralRanking, error) {
	var list []*IntegralRanking
	if err := m.Engine.SQL(GET_INTEGRAL_RANKING, contestId, offset,size).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

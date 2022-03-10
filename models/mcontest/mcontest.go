package mcontest

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
	"sports_service/server/tools/tencentCloud"
)

// 赛事
type ContestModel struct {
	Engine              *xorm.Session
	VideoLive           *models.VideoLive
	VideoLiveReplay     *models.VideoLiveReplay
	Schedule            *models.FpvContestSchedule
	Contest             *models.FpvContestInfo
	IntegralRanking     *models.FpvContestPlayerIntegralRanking
	PlayerInfo          *models.FpvContestPlayerInformation
	ScheduleDetail      *models.FpvContestScheduleDetail
	ScheduleLiveData    *models.FpvContestScheduleLiveData
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
	ShowType       int       `json:"show_type"`
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
	Photo      tencentCloud.BucketURI `json:"photo"`
	Score      int    `json:"score"`
	IsWin      int    `json:"is_win"`
	NumInGroup int    `json:"num_in_group"`
	ContestId  int    `json:"contest_id"`
	BeginTm    int    `json:"begin_tm"`
	EndTm      int    `json:"end_tm"`
	ReceiveIntegral int `json:"receive_integral"`
	Ranking   int     `json:"ranking"`
}


// 赛程列表详情返回数据
type ScheduleListDetailResp struct {
	Id         int64  `json:"id,omitempty"`
	ScheduleId int    `json:"schedule_id"`
	PlayerId   int64  `json:"player_id"`
	PlayerName string `json:"player_name"`
	Photo      tencentCloud.BucketURI `json:"photo"`
	IsWin      int    `json:"is_win"`
	ContestId  int    `json:"contest_id"`
	Ranking    int    `json:"ranking"`
	Ids        []int64 `json:"ids,omitempty"`
	Index      int     `json:"index"`

	BestScore          string    `json:"best_score"`
	RoundOneScore      string    `json:"round_one_score"`
	RoundTwoScore      string    `json:"round_two_score"`
	RoundThreeScore    string    `json:"round_three_score"`
	RoundOneIntegral   int       `json:"round_one_integral"`
	RoundTwoIntegral   int       `json:"round_two_integral"`
	GroupName          string    `json:"group_name"`
	GroupNum           int       `json:"group_num"`
	NumInGroup         int       `json:"num_in_group"`
	BeginTm            int       `json:"begin_tm"`
}

// 赛程分组详情返回数据
type ScheduleGroupDetailResp struct {
	GroupNum   int    `json:"group_num"`
	GroupName  string `json:"group_name"`
	ContestId  int    `json:"contest_id"`
	ScheduleId int    `json:"schedule_id"`
	Index      int    `json:"index"`
	BeginTm    string `json:"begin_tm"`
	//EndTm      string `json:"end_tm"`

	Player     []PlayerInfoResp  `json:"player"`
	Winner     []PlayerInfoResp  `json:"winner"`
}

// 赛程分组详情返回数据
type PlayerInfoResp struct {
	Id         int64  `json:"id,omitempty"`
	PlayerId   int64  `json:"player_id"`
	PlayerName string `json:"player_name"`
	Photo      tencentCloud.BucketURI `json:"photo"`
	IsWin      int    `json:"is_win"`
	Score      string `json:"score"`
	NumInGroup int    `json:"num_in_group"`
	Integral   string `json:"integral"`
}

type IntegralRanking struct {
	Id         int    `json:"id"`
	PlayerId   int64  `json:"player_id"`
	PlayerName string `json:"player_name"`
	Photo      tencentCloud.BucketURI `json:"photo"`
	ContestId  int    `json:"contest_id"`

	Ranking            int       `json:"ranking"`
	TotalIntegral      int       `json:"integral,omitempty"`
	BestScore          int       `json:"score,omitempty"`

	TotalIntegralStr   string    `json:"total_integral"`
	BestScoreStr       string    `json:"best_score"`
}

// 赛程直播 选手竞赛数据
type LiveSchedulePlayerData struct {
	Id               int64    `json:"id"`
	ContestId        int      `json:"contest_id"`
	ScheduleId       int      `json:"schedule_id"`
	PlayerId         int64    `json:"player_id"`
	PlayerName       string   `json:"player_name"`
	Photo            tencentCloud.BucketURI   `json:"photo"`
	LiveId           int64    `json:"live_id"`
	RoundsNum        int      `json:"rounds_num"`
	IntervalDuration string   `json:"interval_duration"`
	TopSpeed         string   `json:"top_speed"`
	ReceiveIntegral  string   `json:"receive_integral"`
	Ranking          int      `json:"ranking"`
}

// 实例
func NewContestModel(engine *xorm.Session) *ContestModel {
	return &ContestModel{
		Engine: engine,
		VideoLive: new(models.VideoLive),
		VideoLiveReplay: new(models.VideoLiveReplay),
		Schedule: new(models.FpvContestSchedule),
		Contest: new(models.FpvContestInfo),
		IntegralRanking: new(models.FpvContestPlayerIntegralRanking),
		PlayerInfo: new(models.FpvContestPlayerInformation),
		ScheduleDetail: new(models.FpvContestScheduleDetail),
		ScheduleLiveData: new(models.FpvContestScheduleLiveData),
	}
}

// 通过赛事id获取正在进行的赛程信息
func (m *ContestModel) GetScheduleInfoByContestId(now int64, contestId string) ([]*models.FpvContestSchedule, error) {
	var list []*models.FpvContestSchedule
	if err := m.Engine.Where("contest_id=? AND status=0 AND start_tm <= ? AND end_tm >= ?", contestId, now, now).
		Asc("order").Find(&list); err != nil {
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
		" WHERE p.status = 0 ORDER BY ISNULL(cs.ranking),cs.ranking ASC,cs.score is null, cs.score ASC,cs.id DESC, p.id ASC"

	GET_SCHEDULE_DETAIL_BY_GROUP = "SELECT cs.*, p.id AS player_id, p.name AS player_name, p.photo FROM " +
		"fpv_contest_schedule_detail AS cs LEFT JOIN fpv_contest_player_information AS p " +
		"ON cs.player_id=p.id WHERE cs.contest_id=? AND cs.schedule_id=? ORDER BY cs.group_num ASC, cs.score ASC, cs.id DESC, group_num ASC"

	GET_SCHEDULE_DETAIL_BY_BACKEND = "SELECT p.id AS player_id, p.name AS player_name, p.photo, cs.* FROM " +
		"fpv_contest_schedule_detail AS cs LEFT JOIN fpv_contest_player_information AS p ON p.id = cs.player_id " +
		"AND cs.contest_id=? AND cs.schedule_id=? WHERE p.status = 0 ORDER BY ISNULL(cs.ranking),cs.ranking ASC,cs.score is null, cs.score ASC, cs.id DESC, " +
		" p.id ASC"


)
// 获取赛程信息[成绩正序]
func (m *ContestModel) GetScheduleDetailInfo(showType int, contestId, scheduleId string) ([]*ScheduleDetail, error) {
	sql := GET_SCHEDULE_DETAIL_BY_SCORE
	if showType == 2 {
		sql = GET_SCHEDULE_DETAIL_BY_GROUP
	}
 
	if showType == 3 {
		sql = GET_SCHEDULE_DETAIL_BY_BACKEND
	}

	var list []*ScheduleDetail
	if err := m.Engine.SQL(sql, contestId, scheduleId).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 删除赛程详情
func (m *ContestModel) DelScheduleDetail(ids []int) (int64, error) {
	return m.Engine.In("id", ids).Delete(&models.FpvContestScheduleDetail{})
}

const (
	GET_INTEGRAL_RANKING = "SELECT p.id AS player_id, p.name AS player_name, p.photo, rk.id,rk.contest_id, " +
		"rk.total_integral, rk.best_score FROM fpv_contest_player_information AS p " +
		"LEFT JOIN fpv_contest_player_integral_ranking AS rk " +
		"ON p.id = rk.player_id AND rk.contest_id=? " +
		"WHERE p.status = 0 ORDER BY rk.total_integral DESC, ISNULL(rk.ranking),rk.ranking ASC, p.id ASC LIMIT ?, ?"
)
// 通过赛事id 获取选手积分排行
func (m *ContestModel) GetIntegralRankingByContestId(contestId string, offset, size int) ([]*IntegralRanking, error) {
	var list []*IntegralRanking
	if err := m.Engine.SQL(GET_INTEGRAL_RANKING, contestId, offset,size).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取排行榜总数
func (m *ContestModel) GetIntegralRankingTotal(contestId string) (int64, error) {
	return m.Engine.Where("status=0 AND contest_id=?", contestId).Count(&models.FpvContestPlayerInformation{})
}

// 获取选手总积分
func (m *ContestModel) GetTotalIntegralByPlayer(contestId, playerId string) (bool, error) {
	m.IntegralRanking = new(models.FpvContestPlayerIntegralRanking)
	return m.Engine.Where("status=0 AND contest_id=? AND player_id=?", contestId, playerId).Get(m.IntegralRanking)
}

const (
	GET_LIVE_SCHEDULE_PLAYER_DATA = "SELECT * FROM fpv_contest_schedule_live_data WHERE status=0 AND live_id=? " +
		"AND rounds_num=(select max(rounds_num) from fpv_contest_schedule_live_data) ORDER BY receive_integral DESC" +
		" LIMIT ?, ?"
)
// 赛程直播 选手竞赛数据
func (m *ContestModel) GetLiveSchedulePlayerData(liveId string, offset, size int) ([]*models.FpvContestScheduleLiveData, error) {
	var list []*models.FpvContestScheduleLiveData
	if err := m.Engine.SQL(GET_LIVE_SCHEDULE_PLAYER_DATA, liveId, offset, size).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 通过id 获取选手信息
func (m *ContestModel) GetPlayerInfoById(playerId string) (bool, error) {
	m.PlayerInfo = new(models.FpvContestPlayerInformation)
	return m.Engine.Where("id=?", playerId).Get(m.PlayerInfo)
}

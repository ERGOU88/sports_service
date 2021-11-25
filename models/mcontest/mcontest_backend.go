package mcontest

import (
	"sports_service/server/models"
)

// 添加赛程详情请求参数
type AddScheduleDetail struct {
	PlayerId         int64   `json:"player_id"`
	PlayerName       string  `json:"player_name"`
	GroupName        string  `json:"group_name"`
	GroupNum         int     `json:"group_num"`
	RoundOneScore    int     `json:"round_one_score"`
	RoundTwoScore    int     `json:"round_two_score"`
	RoundThreeScore  int     `json:"round_three_score"`
	RoundOneIntegral int     `json:"round_one_integral"`
	RoundTwoIntegral int     `json:"round_two_integral"`
	Ranking          int     `json:"ranking"`
	IsWin            int     `json:"is_win"`
	NumInGroup       int     `json:"num_in_group"`
	BeginTm          int     `json:"begin_tm"`
	EndTm            int     `json:"end_tm"`
	ScheduleId       int     `json:"schedule_id"`
	ContestId        int     `json:"contest_id"`
}

// 添加选手信息
func (m *ContestModel) AddPlayer(player *models.FpvContestPlayerInformation) (int64, error) {
	return m.Engine.InsertOne(player)
}

// 更新选手信息
func (m *ContestModel) UpdatePlayer(player *models.FpvContestPlayerInformation) (int64, error) {
	return m.Engine.Where("id=?", player.Id).Update(player)
}

// 获取选手列表
func (m *ContestModel) GetPlayerList(offset, size int) ([]*models.FpvContestPlayerInformation, error) {
	var list []*models.FpvContestPlayerInformation
	if err := m.Engine.Where("status=0").Limit(size, offset).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取选手总数
func (m *ContestModel) GetPlayerCount() int64 {
	count, err := m.Engine.Count(&models.FpvContestPlayerInformation{})
	if err != nil {
		return 0
	}
	
	return count
}


// 添加赛程组别
func (m *ContestModel) AddContestGroup(group *models.FpvContestScheduleGroup) (int64, error) {
	return m.Engine.InsertOne(group)
}

// 更新赛程组别信息
func (m *ContestModel) UpdateContestGroup(group *models.FpvContestScheduleGroup) (int64, error) {
	return m.Engine.Where("id=?", group.Id).Update(group)
}

// 获取赛事 赛程组别配置信息
func (m *ContestModel) GetContestGroupList(offset, size int, scheduleId, contestId string) ([]*models.FpvContestScheduleGroup, error) {
	var list []*models.FpvContestScheduleGroup
    m.Engine.Where("status=0")
	if scheduleId != "" {
		m.Engine.Where("schedule_id=?", scheduleId)
	}

	if contestId != "" {
		m.Engine.Where("contest_id=?", contestId)
	}

	if err := m.Engine.Asc("order").Limit(size, offset).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}


// 获取赛程信息
func (m *ContestModel) GetScheduleInfo() ([]*models.FpvContestSchedule, error) {
	var list []*models.FpvContestSchedule
	if err := m.Engine.Where("status=0").
		Asc("order").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 设置赛事积分排行
func (m *ContestModel) SetIntegralRanking(info *models.FpvContestPlayerIntegralRanking) (int64, error) {
	return m.Engine.InsertOne(info)
}

// 更新赛事积分排行信息
func (m *ContestModel) UpdateIntegralRanking(info *models.FpvContestPlayerIntegralRanking) (int64, error) {
	return m.Engine.Where("id=?", info.Id).Update(info)
}

// 添加赛事详情
func (m *ContestModel) AddContestScheduleDetail(list []*models.FpvContestScheduleDetail) (int64, error) {
	return m.Engine.InsertMulti(list)
}

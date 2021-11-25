package mcontest

import (
	"sports_service/server/models"
	"sports_service/server/tools/tencentCloud"
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

type FpvContestPlayerInformation struct {
	Id        int64  `json:"id" xorm:"pk autoincr comment('自增id') BIGINT(20)"`
	Name      string `json:"name" xorm:"not null comment('选手名称 例如：pdd') VARCHAR(60)"`
	Photo     tencentCloud.BucketURI `json:"photo" xorm:"not null default '' comment('选手照片') VARCHAR(512)"`
	Country   string `json:"country" xorm:"not null default '' comment('国家') VARCHAR(128)"`
	Province  string `json:"province" xorm:"not null default '' comment('省份') VARCHAR(128)"`
	City      string `json:"city" xorm:"not null default '' comment('城市') VARCHAR(128)"`
	Age       int    `json:"age" xorm:"not null default 0 comment('年龄') INT(3)"`
	Hobby     string `json:"hobby" xorm:"not null default '' comment('爱好') VARCHAR(255)"`
	ContestId int    `json:"contest_id" xorm:"not null default 0 comment('参加的赛事id') index INT(11)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('0 正常 1 废弃') TINYINT(1)"`
	IdCard    string `json:"id_card" xorm:"not null default '' comment('证件号码') VARCHAR(255)"`
	IdType    int    `json:"id_type" xorm:"not null default 0 comment('1 身份证 2 居住证 3 护照 4 港澳') TINYINT(2)"`
	Gender    int    `json:"gender" xorm:"not null default 0 comment('0 未知 1 男 2 女') TINYINT(1)"`
	Born      string `json:"born" xorm:"not null default '' comment('出生年月日') VARCHAR(128)"`
	MobileNum string `json:"mobile_num" xorm:"not null default '' comment('手机号码') VARCHAR(60)"`
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
func (m *ContestModel) GetPlayerList(offset, size int) ([]*FpvContestPlayerInformation, error) {
	var list []*FpvContestPlayerInformation
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

// 获取赛程组别总数
func (m *ContestModel) GetContestGroupCount(scheduleId, contestId string) int64 {
	m.Engine.Where("status=0")
	if scheduleId != "" {
		m.Engine.Where("schedule_id=?", scheduleId)
	}
	
	if contestId != "" {
		m.Engine.Where("contest_id=?", contestId)
	}
	
	count, err := m.Engine.Count(&models.FpvContestScheduleGroup{})
	if err != nil {
		return 0
	}
	
	return count
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

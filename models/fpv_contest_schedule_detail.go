package models

type FpvContestScheduleDetail struct {
	Id              int64  `json:"id" xorm:"pk autoincr comment('自增id') BIGINT(20)"`
	ScheduleId      int    `json:"schedule_id" xorm:"not null comment('赛程id') index INT(11)"`
	Rounds          int    `json:"rounds" xorm:"not null comment('轮次 第几轮') INT(8)"`
	GroupNum        int    `json:"group_num" xorm:"not null default 1 comment('第几组') INT(6)"`
	GroupName       string `json:"group_name" xorm:"not null default '' comment('组名 例如A组') VARCHAR(128)"`
	PlayerId        int64  `json:"player_id" xorm:"not null comment('选手id') BIGINT(20)"`
	BeginTm         int    `json:"begin_tm" xorm:"not null default 0 comment('开始时间') INT(11)"`
	EndTm           int    `json:"end_tm" xorm:"not null default 0 comment('结束时间') INT(11)"`
	Status          int    `json:"status" xorm:"not null default 0 comment('0 正常 1 废弃') TINYINT(1)"`
	Score           int    `json:"score" xorm:"not null default 0 comment('比赛成绩 (暂定 * 1000存储)') INT(11)"`
	IsWin           int    `json:"is_win" xorm:"not null default 0 comment('是否胜出 1 胜出') TINYINT(1)"`
	NumInGroup      int    `json:"num_in_group" xorm:"not null default 0 comment('分组内编号') INT(3)"`
	CreateAt        int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt        int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
	ContestId       int    `json:"contest_id" xorm:"not null comment('赛事id') INT(11)"`
	Ranking         int    `json:"ranking" xorm:"comment('排名') INT(8)"`
	ReceiveIntegral int    `json:"receive_integral" xorm:"not null default 0 comment('获得积分') INT(11)"`
}

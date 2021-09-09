package models

type FpvContestPlayerIntegralRecord struct {
	Id              int64 `json:"id" xorm:"pk autoincr comment('自增id') BIGINT(20)"`
	ContestId       int64 `json:"contest_id" xorm:"not null comment('所属赛事id') index(contest_player_schedule) BIGINT(20)"`
	ScheduleId      int64 `json:"schedule_id" xorm:"not null comment('赛程id') index(contest_player_schedule) BIGINT(20)"`
	PlayerId        int64 `json:"player_id" xorm:"not null comment('选手id') index(contest_player_schedule) BIGINT(20)"`
	Ranking         int   `json:"ranking" xorm:"not null default 0 comment('排名') INT(8)"`
	ReceiveIntegral int   `json:"receive_integral" xorm:"not null default 0 comment('获取积分数') INT(11)"`
	Status          int   `json:"status" xorm:"not null default 0 comment('0 正常 1 废弃') TINYINT(1)"`
	CreateAt        int   `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt        int   `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
}

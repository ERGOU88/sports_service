package models

type FpvContestScheduleGroup struct {
	Id             int    `json:"id" xorm:"not null pk autoincr comment('自增id') INT(11)"`
	ContestId      int    `json:"contest_id" xorm:"not null comment('所属赛事id') index INT(11)"`
	ScheduleId     int    `json:"schedule_id" xorm:"not null comment('赛程id') index INT(11)"`
	ScheduleName   string `json:"schedule_name" xorm:"not null comment('赛程名称 例如：测试赛、资格赛、32强赛、16强赛等') VARCHAR(128)"`
	GroupName      string `json:"group_name" xorm:"not null comment('组别名称') VARCHAR(128)"`
	Order          int    `json:"order" xorm:"not null default 0 comment('组别顺序 1 表示最先开始') INT(8)"`
	GroupPlayerNum int    `json:"group_player_num" xorm:"not null default 0 comment('组内选手数 4人一组 则 设为4') INT(8)"`
	Status         int    `json:"status" xorm:"not null default 0 comment('0表示正常 1表示废弃') TINYINT(1)"`
	CreateAt       int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt       int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
}

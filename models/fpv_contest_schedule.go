package models

type FpvContestSchedule struct {
	Id                  int64  `json:"id" xorm:"pk autoincr comment('赛程id') BIGINT(20)"`
	ContestId           int64  `json:"contest_id" xorm:"not null comment('所属赛事id') index BIGINT(20)"`
	ScheduleName        string `json:"schedule_name" xorm:"not null comment('赛程名称 例如：测试赛、资格赛、32强赛、16强赛等') VARCHAR(128)"`
	Order               int    `json:"order" xorm:"not null default 0 comment('赛程顺序 1 表示最先开始') INT(8)"`
	RoundsNum           int    `json:"rounds_num" xorm:"not null comment('轮数 当前赛事规则下 进行几轮角逐 ') INT(8)"`
	CanCompeteNum       int    `json:"can_compete_num" xorm:"not null default 0 comment('可参与的选手人数') INT(8)"`
	PromotionNumber     int    `json:"promotion_number" xorm:"not null default 0 comment('当前轮次可晋级的人数 0表示 非晋级竞赛') INT(8)"`
	TotalGroupNum       int    `json:"total_group_num" xorm:"not null default 1 comment('分组数 默认为1组  [分几组竞技]') INT(6)"`
	GroupPlayerNum      int    `json:"group_player_num" xorm:"not null default 0 comment('分组竞技选手数 例如：4人一组 则 设为4') INT(8)"`
	Status              int    `json:"status" xorm:"not null default 0 comment('0表示正常 1表示废弃') TINYINT(1)"`
	ScheduleDescription string `json:"schedule_description" xorm:"not null default '' comment('赛程说明 例： 32强赛：32人参赛，4人一组 分8组角逐16强') VARCHAR(1024)"`
	CreateAt            int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt            int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
}

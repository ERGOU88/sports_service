package models

type FpvContestScheduleLiveData struct {
	Id               int64  `json:"id" xorm:"pk autoincr comment('自增id') BIGINT(20)"`
	ContestId        int    `json:"contest_id" xorm:"not null comment('所属赛事id') INT(11)"`
	ScheduleId       int    `json:"schedule_id" xorm:"not null comment('赛程id') INT(11)"`
	PlayerId         int64  `json:"player_id" xorm:"not null comment('选手id') BIGINT(20)"`
	LiveId           int64  `json:"live_id" xorm:"not null default 0 comment('直播数据id') BIGINT(20)"`
	RoundsNum        int    `json:"rounds_num" xorm:"not null default 0 comment('圈数') INT(8)"`
	IntervalDuration int    `json:"interval_duration" xorm:"not null default 0 comment('间隔时长 (暂定 * 1000存储)') INT(11)"`
	TopSpeed         int    `json:"top_speed" xorm:"not null default 0 comment('最快圈速 (暂定 * 1000存储)') INT(11)"`
	ReceiveIntegral  int    `json:"receive_integral" xorm:"not null default 0 comment('获取积分数 (暂定 * 1000存储)') INT(11)"`
	Status           int    `json:"status" xorm:"not null default 0 comment('0 正常 1 废弃') TINYINT(1)"`
	CreateAt         int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt         int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
	Ranking          int    `json:"ranking" xorm:"not null default 0 comment('排名') INT(11)"`
	Name             string `json:"name" xorm:"not null default '' comment('选手名称') VARCHAR(60)"`
	Photo            string `json:"photo" xorm:"not null default '' comment('选手照片') VARCHAR(512)"`
	Gender           int    `json:"gender" xorm:"not null default 0 comment('0 未知 1 男 2 女') TINYINT(1)"`
}

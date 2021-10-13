package models

type FpvContestPlayerIntegralRanking struct {
	Id            int64 `json:"id" xorm:"pk autoincr comment('自增id') BIGINT(20)"`
	ContestId     int   `json:"contest_id" xorm:"not null comment('所属赛事id') INT(11)"`
	PlayerId      int64 `json:"player_id" xorm:"not null comment('选手id') BIGINT(20)"`
	TotalIntegral int64 `json:"total_integral" xorm:"not null default 0 comment('总积分(暂定 * 1000存储)') BIGINT(20)"`
	BestScore     int   `json:"best_score" xorm:"not null default 0 comment('最佳成绩 (暂定 * 1000存储)') INT(11)"`
	Status        int   `json:"status" xorm:"not null default 0 comment('0 正常 1 隐藏') TINYINT(1)"`
	CreateAt      int   `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt      int   `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
	Ranking       int   `json:"ranking" xorm:"comment('排名') INT(8)"`
}

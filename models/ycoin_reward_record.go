package models

type YcoinRewardRecord struct {
	Id          int64  `json:"id" xorm:"pk autoincr comment('自增主键') BIGINT(20)"`
	GiverId     string `json:"giver_id" xorm:"not null comment('赠予人uid') index VARCHAR(60)"`
	DonneId     string `json:"donne_id" xorm:"not null comment('受赠人uid') index VARCHAR(60)"`
	Status      int    `json:"status" xorm:"not null default 1 comment('1-正常展示, 2-不展示') TINYINT(1)"`
	RewardYcoin int    `json:"reward_ycoin" xorm:"not null comment('打赏的游币数量') INT(11)"`
	RecordId    string `json:"record_id" xorm:"not null default ' ' comment('打赏的视频记录id（暂时只有视频可以打赏）') VARCHAR(100)"`
	RewardType  int    `json:"reward_type" xorm:"not null comment('打赏类型 1 打赏视频') TINYINT(1)"`
	CreateAt    int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
}

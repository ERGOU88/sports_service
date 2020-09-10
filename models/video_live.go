package models

type VideoLive struct {
	Id          int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	AnchorId    string `json:"anchor_id" xorm:"not null comment('主播id') index VARCHAR(60)"`
	RoomId      string `json:"room_id" xorm:"not null comment('房间id') VARCHAR(60)"`
	Cover       string `json:"cover" xorm:"comment('直播封面') VARCHAR(200)"`
	RtmpAddr    string `json:"rtmp_addr" xorm:"not null default '' comment('rtmp地址') VARCHAR(512)"`
	FlvAddr     string `json:"flv_addr" xorm:"not null default '' comment('flv地址') VARCHAR(512)"`
	HlsAddr     string `json:"hls_addr" xorm:"not null default '' comment('hls地址') VARCHAR(512)"`
	StreamUrl   string `json:"stream_url" xorm:"not null default '' comment('推流url') VARCHAR(255)"`
	StreamKey   string `json:"stream_key" xorm:"comment('推流密钥') VARCHAR(255)"`
	PlayTime    int    `json:"play_time" xorm:"not null comment('开播时间') INT(11)"`
	EndTime     int    `json:"end_time" xorm:"not null comment('结束时间') INT(11)"`
	IncomeYcoin int    `json:"income_ycoin" xorm:"comment('本次直播收益（游币）') INT(11)"`
	Status      int    `json:"status" xorm:"default 0 comment('状态 0未直播 1直播中 2异常') TINYINT(1)"`
	Describe    string `json:"describe" xorm:"not null default '' comment('描述') VARCHAR(255)"`
	Tags        string `json:"tags" xorm:"not null default '' comment('直播标签') VARCHAR(255)"`
	CreateAt    int    `json:"create_at" xorm:"not null default 0 comment('记录创建时间') INT(11)"`
	UpdateAt    int    `json:"update_at" xorm:"not null default 0 comment('记录更新时间') INT(11)"`
	Duration    int64  `json:"duration" xorm:"default 0 comment('时长（秒）') BIGINT(20)"`
	LiveType    int    `json:"live_type" xorm:"not null default 0 comment('直播类型（0：管理员[sys_user]，1：用户[user]）') TINYINT(1)"`
	Manager     int    `json:"manager" xorm:"not null default 0 comment('后台操作用户') INT(11)"`
	Sequence    string `json:"sequence" xorm:"not null default '' comment('推流唯一标识') VARCHAR(50)"`
}

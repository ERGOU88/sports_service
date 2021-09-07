package models

type VideoLive struct {
	Id             int64  `json:"id" xorm:"pk autoincr comment('主键') BIGINT(20)"`
	UserId         string `json:"user_id" xorm:"not null default '' comment('主播id') VARCHAR(60)"`
	RoomId         string `json:"room_id" xorm:"not null default '' comment('房间id') index VARCHAR(60)"`
	Cover          string `json:"cover" xorm:"not null default '' comment('直播封面') VARCHAR(512)"`
	RtmpAddr       string `json:"rtmp_addr" xorm:"not null default '' comment('rtmp地址[拉流]') VARCHAR(512)"`
	FlvAddr        string `json:"flv_addr" xorm:"not null default '' comment('flv地址[拉流]') VARCHAR(512)"`
	HlsAddr        string `json:"hls_addr" xorm:"not null default '' comment('hls地址[拉流]') VARCHAR(512)"`
	PushStreamAddr string `json:"push_stream_addr" xorm:"not null default '' comment('推流地址') VARCHAR(512)"`
	PlayTime       int    `json:"play_time" xorm:"not null default 0 comment('开播时间') INT(11)"`
	EndTime        int    `json:"end_time" xorm:"not null default 0 comment('结束时间') INT(11)"`
	Status         int    `json:"status" xorm:"default 0 comment('状态 0未直播 1直播中 2异常') TINYINT(1)"`
	Describe       string `json:"describe" xorm:"not null default '' comment('描述') VARCHAR(512)"`
	Tags           string `json:"tags" xorm:"not null default '' comment('直播标签') VARCHAR(512)"`
	CreateAt       int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt       int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	Duration       int64  `json:"duration" xorm:"default 0 comment('时长（秒）') BIGINT(20)"`
	LiveType       int    `json:"live_type" xorm:"not null default 0 comment('直播类型（0：管理员[sys_user]，1：用户[user]）') TINYINT(1)"`
}

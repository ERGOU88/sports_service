package models

type VideoAlbumDetail struct {
	AlbumId   int64  `json:"album_id" xorm:"not null pk comment('专辑id') BIGINT(20)"`
	VideoId   int64  `json:"video_id" xorm:"not null pk comment('视频id') BIGINT(20)"`
	AlbumName string `json:"album_name" xorm:"comment('专辑名') VARCHAR(60)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('0为正常 1为废弃') TINYINT(1)"`
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
}

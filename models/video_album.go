package models

type VideoAlbum struct {
	Id        int64  `json:"id" xorm:"pk autoincr comment('专辑id') BIGINT(20)"`
	UserId    string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	AlbumName string `json:"album_name" xorm:"not null comment('专辑名称') VARCHAR(60)"`
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('0 正常 1 废弃') TINYINT(1)"`
}

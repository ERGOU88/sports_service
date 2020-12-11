package models

import (
	"time"
)

type Bili struct {
	Author      string    `json:"author" xorm:"VARCHAR(255)"`
	AuthorId    string    `json:"author_id" xorm:"VARCHAR(255)"`
	Created     time.Time `json:"created" xorm:"created DATETIME"`
	Description string    `json:"description" xorm:"TEXT"`
	Id          int64     `json:"id" xorm:"pk autoincr BIGINT(20)"`
	Pic         string    `json:"pic" xorm:"VARCHAR(500)"`
	Tags        string    `json:"tags" xorm:"VARCHAR(500)"`
	Title       string    `json:"title" xorm:"VARCHAR(255)"`
	Updated     time.Time `json:"Updated" xorm:"DATETIME"`
	VideoPage   string    `json:"video_page" xorm:"unique VARCHAR(255)"`
	VideoUrl    string    `json:"video_url" xorm:"default '' VARCHAR(255)"`
}

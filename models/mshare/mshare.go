package mshare

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

type ShareModel struct {
	Engine         *xorm.Session
	Share          *models.ShareRecord
}

// 分享请求参数
type ShareParams struct {
	SharePlatform      int    `binding:"required" json:"share_platform"`        // 分享平台 1 微信 2 微博 3 qq 4 app内
	ShareType          int    `binding:"required" json:"share_type"`            // 分享类型  1 分享视频 2 分享帖子
	ComposeId          int    `binding:"required" json:"compose_id"`            // 视频/帖子id

}

// 转发请求参数
type ForwardParams struct {
	ForwardPlatform     int     `json:"forward_platform"`      // 转发平台
}

func NewForwardModel(engine *xorm.Session) *ShareModel {
	return &ShareModel{
		Engine: engine,
		Share: new(models.ShareRecord),
	}
}

// 添加转发记录
func (m *ShareModel) AddForward() (int64, error) {
	return m.Engine.InsertOne(m.Share)
}

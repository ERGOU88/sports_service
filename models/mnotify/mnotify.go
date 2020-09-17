package mnotify

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

type NotifyModel struct {
	NofitySetting  *models.SystemNoticeSettings
	Engine         *xorm.Session
}

// 通知设置请求参数
type NotifySettingParams struct {
	CommentPushSet   int    `json:"commentPushSet"`         // 评论推送 0 接收 1 不接收
	ThumbUpPushSet   int    `json:"thumbUpPushSet"`         // 点赞推送 0 接收 1 不接收
	AttentionPushSet int    `json:"attentionPushSet"`       // 关注推送 0 接收 1 不接收
	SharePushSet     int    `json:"sharePushSet"`           // 分享推送 0 接收 1 不接收
	SlotPushSet      int    `json:"slotPushSet"`            // 投币推送 0 接收 1 不接收
}

// 实例
func NewNotifyModel(engine *xorm.Session) *NotifyModel {
	return &NotifyModel{
		NofitySetting: new(models.SystemNoticeSettings),
		Engine: engine,
	}
}

// 添加用户通知设置（默认接收）
func (m *NotifyModel) AddUserNotifySetting(userId string, now int) error {
	m.NofitySetting.UserId = userId
	m.NofitySetting.CreateAt = now
	m.NofitySetting.UpdateAt = now
	if _, err := m.Engine.InsertOne(m.NofitySetting); err != nil {
		return err
	}

	return nil
}

// 修改用户通知设置
func (m *NotifyModel) UpdateUserNotifySetting() error {
	if _, err := m.Engine.Update(m.NofitySetting); err != nil {
		return err
	}

	return nil
}





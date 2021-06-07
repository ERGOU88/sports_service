package mnotify

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/rdskey"
	"sports_service/server/models"
	"time"
	"fmt"
)

type NotifyModel struct {
	NofitySetting  *models.SystemNoticeSettings
	SystemNotify   *models.SystemMessage
	Engine         *xorm.Session
}

// 通知设置请求参数
type NotifySettingParams struct {
	CommentPushSet   int    `json:"comment_push_set" example:"0"`          // 评论推送 0 接收 1 不接收
	ThumbUpPushSet   int    `json:"thumb_up_push_set" example:"0"`         // 点赞推送 0 接收 1 不接收
	AttentionPushSet int    `json:"attention_push_set" example:"0"`        // 关注推送 0 接收 1 不接收
	SharePushSet     int    `json:"share_push_set" example:"0"`            // 分享推送 0 接收 1 不接收
	SlotPushSet      int    `json:"slot_push_set" example:"0"`             // 投币推送 0 接收 1 不接收
}

// 收到的@信息（评论/回复）
type ReceiveCommentAtInfo struct {
	ComposeId     int64                 `json:"compose_id" example:"1000000000"`      // 视频作品id
	Title         string                `json:"title" example:"视频标题"`               // 标题
	Describe      string                `json:"describe" example:"视频描述"`            // 描述
	Cover         string                `json:"cover" example:"视频封面地址"`           // 封面
	VideoAddr     string                `json:"video_addr" example:"视频地址"`         // 视频地址
	IsRecommend   int                   `json:"is_recommend" example:"0"`             // 是否推荐
	IsTop         int                   `json:"is_top" example:"0"`                   // 是否置顶
	VideoDuration int                   `json:"video_duration" example:"1000"`        // 视频时长
	VideoWidth    int64                 `json:"video_width" example:"1000"`           // 视频宽
	VideoHeight   int64                 `json:"video_height" example:"1000"`          // 视频高
	Status        int32                 `json:"status" example:"1"`                   // 审核状态
	CreateAt      int                   `json:"create_at" example:"1600000000"`       // 视频创建时间
	BarrageNum    int                   `json:"barrage_num" example:"1000"`           // 弹幕数
	BrowseNum     int                   `json:"browse_num" example:"1000"`            // 浏览数（播放数）
	UserId        string                `json:"user_id" example:"执行@的用户id"`        // 执行@的用户id
	Avatar        string                `json:"avatar" example:"执行@的用户头像"`        // 执行@的用户头像
	Nickname      string                `json:"nick_name" example:"执行@的用户昵称"`     // 执行@的用户昵称
	ToUserId      string                `json:"to_user_id" example:"被@的用户id"`       // 被@的用户id
	ToUserAvatar  string                `json:"to_user_avatar" example:"被@的用户头像"`  // 被@用户头像
	ToUserName    string                `json:"to_user_name" example:"被@的用户昵称"`    // 被@的用户昵称
	Content       string                `json:"content" example:"内容"`                 // 评论内容
	Reply         string                `json:"reply"  example:"回复的内容"`             // 回复的内容
	AtTime        int                   `json:"at_time" example:"1600000000"`          // 用户@的时间
	Type          int                   `json:"type" example:"1"`                      // 类型 1 视频 2 帖子 3 评论
	CommentType   int                   `json:"comment_type" example:"1"`              // 1 为评论 2 为回复
	IsAt          int                   `json:"is_at"`                                 // 1为回复的回复 0 为1级评论/1级评论的回复
	ParentComment string                `json:"parent_comment"`                        // 父级评论（1级评论）
	IsLike        int                   `json:"is_like" example:"0"`                   // 是否点赞
	CommentId     int64                 `json:"comment_id"`                            // 评论id（1级评论的id）

	ReplyCommentId int64                `json:"reply_comment_id"`                      // 回复评论所用id
	TotalLikeNum   int64                `json:"total_like_num"`                        // 总点赞数
}

// 首页通知
type HomePageNotify struct {
	UnBrowsedNum     int64       `json:"un_browsed_num"` // 未浏览的视频（关注模块 关注用户发布的视频） 客户端处理：当前用户有未浏览的 则展示红点）
	UnreadNum        int64       `json:"unread_num"`     // 首页未读消息总数
}

// 实例
func NewNotifyModel(engine *xorm.Session) *NotifyModel {
	return &NotifyModel{
		NofitySetting: new(models.SystemNoticeSettings),
		SystemNotify: new(models.SystemMessage),
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
	if _, err := m.Engine.Where("user_id=?", m.NofitySetting.UserId).Cols("attention_push_set, comment_push_set, thumb_up_push_set, share_push_set, slot_push_set, update_at").Update(m.NofitySetting); err != nil {
		return err
	}

	return nil
}

// 获取用户通知设置
func (m *NotifyModel) GetUserNotifySetting(userId string) *models.SystemNoticeSettings {
	setting := new(models.SystemNoticeSettings)
	ok, err := m.Engine.Where("user_id=?", userId).Get(setting)
	if !ok || err != nil {
		return nil
	}

	return setting
}

// 获取系统通知
func (m *NotifyModel) GetSystemNotify(userId string, offset, size int) []*models.SystemMessage {
	sql := "SELECT `system_id`, `cover`, `send_type`, `receive_id`, `system_topic`, `system_content`, `send_time`, `extra`, `status` FROM system_message "
	if userId != "" {
		sql = fmt.Sprintf("%s WHERE receive_id=%s AND send_status=0 OR send_default=1 AND send_status=0", sql, userId)
	} else {
		sql = fmt.Sprintf("%s WHERE receive_id='' AND send_default=1 AND send_status=0", sql)
	}

	sql = fmt.Sprintf("%s ORDER BY system_id DESC LIMIT %d, %d", sql, offset, size)

	var list []*models.SystemMessage
	if err := m.Engine.SQL(sql).Find(&list); err != nil {
		return nil
	}

	return list
}


// 通过id获取通知详情
func (m *NotifyModel) GetSystemNotifyById(systemId string) *models.SystemMessage {
	msg := new(models.SystemMessage)
	ok, err := m.Engine.Where("system_id=?", systemId).Get(msg)
	if !ok || err != nil {
		return nil
	}

	return msg
}

func (m *NotifyModel) UpdateSystemNotifyStatus(id int64) error {
	sql := fmt.Sprintf("UPDATE `system_message` SET `status`=1 WHERE system_id <= %d AND send_status=0", id)
	if _, err := m.Engine.Exec(sql); err != nil {
		return err
	}

	return nil
}

// 获取未读系统消息总数
func (m *NotifyModel) GetUnreadSystemMsgNum(userId string) int64 {
	count, err := m.Engine.Where("status=0 AND receive_id=? AND send_status=0", userId).Count(&models.SystemMessage{})
	if err != nil {
		return 0
	}

	return count
}

// 记录用户读取被点赞通知消息的时间
func (m *NotifyModel) RecordReadBeLikedTime(userId string) error {
	rds := dao.NewRedisDao()
	return rds.Set(rdskey.MakeKey(rdskey.USER_READ_BELIKED_NOTIFY, userId), time.Now().Unix())
}

// 获取用户读取被点赞通知消息的时间
func (m *NotifyModel) GetReadBeLikedTime(userId string) (string, error) {
	rds := dao.NewRedisDao()
	return rds.Get(rdskey.MakeKey(rdskey.USER_READ_BELIKED_NOTIFY, userId))
}

// 记录用户读取关注用户发布的视频的时间
func (m *NotifyModel) RecordReadAttentionPubVideo(userId string) error {
	rds := dao.NewRedisDao()
	return rds.Set(rdskey.MakeKey(rdskey.USER_READ_ATTENTION_VIDEO, userId), time.Now().Unix())
}

// 获取用户读取关注用户发布的视频的时间
func (m *NotifyModel) GetReadAttentionPubVideo(userId string) (string, error) {
	rds := dao.NewRedisDao()
	return rds.Get(rdskey.MakeKey(rdskey.USER_READ_ATTENTION_VIDEO, userId))
}

// 记录用户读取@通知消息的时间
func (m *NotifyModel) RecordReadAtTime(userId string) error {
	rds := dao.NewRedisDao()
	return rds.Set(rdskey.MakeKey(rdskey.USER_READ_AT_NOTIFY, userId), time.Now().Unix())
}

// 获取用户读取@通知消息的时间
func (m *NotifyModel) GetReadAtTime(userId string) (string, error) {
	rds := dao.NewRedisDao()
	return rds.Get(rdskey.MakeKey(rdskey.USER_READ_AT_NOTIFY, userId))
}




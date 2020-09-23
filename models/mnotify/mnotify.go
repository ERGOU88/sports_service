package mnotify

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/rdskey"
	"sports_service/server/models"
	"time"
)

type NotifyModel struct {
	NofitySetting  *models.SystemNoticeSettings
	Engine         *xorm.Session
}

// 通知设置请求参数
type NotifySettingParams struct {
	CommentPushSet   int    `json:"comment_push_set"`          // 评论推送 0 接收 1 不接收
	ThumbUpPushSet   int    `json:"thumb_up_push_set"`         // 点赞推送 0 接收 1 不接收
	AttentionPushSet int    `json:"attention_push_set"`        // 关注推送 0 接收 1 不接收
	SharePushSet     int    `json:"share_push_set"`            // 分享推送 0 接收 1 不接收
	SlotPushSet      int    `json:"slot_push_set"`             // 投币推送 0 接收 1 不接收
}

// 收到的@信息（评论/回复）
type ReceiveCommentAtInfo struct {
	ComposeId     int64                 `json:"compose_id"`      // 作品id
	Title         string                `json:"title"`           // 标题
	Describe      string                `json:"describe"`        // 描述
	Cover         string                `json:"cover"`           // 封面
	VideoAddr     string                `json:"video_addr"`      // 视频地址
	IsRecommend   int                   `json:"is_recommend"`    // 是否推荐
	IsTop         int                   `json:"is_top"`          // 是否置顶
	VideoDuration int                   `json:"video_duration"`  // 视频时长
	VideoWidth    int64                 `json:"video_width"`     // 视频宽
	VideoHeight   int64                 `json:"video_height"`    // 视频高
	Status        int32                 `json:"status"`          // 审核状态
	CreateAt      int                   `json:"create_at"`       // 视频创建时间
	BarrageNum    int                   `json:"barrage_num"`     // 弹幕数
	BrowseNum     int                   `json:"browse_num"`      // 浏览数（播放数）
	UserId        string                `json:"user_id"`         // 执行@的用户id
	Avatar        string                `json:"avatar"`          // 执行@的用户头像
	Nickname      string                `json:"nick_name"`       // 执行@的用户昵称
	ToUserId      string                `json:"to_user_id"`      // 被@的用户id
	ToUserAvatar  string                `json:"avatar"`          // 被@用户头像
	ToUserName    string                `json:"to_user_name"`    // 被@的用户昵称
	Content       string                `json:"content"`         // 评论内容
	Reply         string                `json:"reply"`           // 回复的内容
	AtTime        int                   `json:"at_time"`         // 用户@的时间
	Type          int                   `json:"type"`            // 类型 1 视频 2 帖子 3 评论
	CommentType   int                   `json:"comment_type"`    // 1 为评论 2 为回复
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

// 获取用户通知设置
func (m *NotifyModel) GetUserNotifySetting(userId string) *models.SystemNoticeSettings {
	setting := new(models.SystemNoticeSettings)
	ok, err := m.Engine.Where("user_id=?", userId).Get(setting)
	if !ok || err != nil {
		return nil
	}

	return setting
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




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
	SharePlatform      int                `binding:"required" json:"share_platform"`        // 分享平台 1 微信 2 微博 3 qq 4 app内
	ShareType          int                `binding:"required" json:"share_type"`            // 分享类型  1 分享视频 2 分享帖子
	ComposeId          int                `binding:"required" json:"compose_id"`            // 视频/帖子id
	SectionId         int                 `json:"section_id"`                               // 主模块id
	TopicIds          []string            `json:"topic_ids"`                                // 话题id （多个）
	Title             string              `json:"title"`                                    // 标题
	Describe          string              `json:"describe"`                                 // 描述
	//ForwardVideo      *ShareVideoInfo     `json:"forward_video"`                            // 转发的视频内容
	//ForwardPost       *SharePostInfo      `json:"forward_post"`                             // 转发的帖子内容
	UserId            string              `json:"user_id"`                                  // 用户id
}

// 分享的视频信息
type ShareVideoInfo struct {
	VideoId       int64                 `json:"video_id"`
	Title         string                `json:"title"  example:"标题"`                 // 标题
	Describe      string                `json:"describe"  example:"描述"`              // 描述
	Cover         string                `json:"cover"  example:"封面"`                 // 封面
	VideoAddr     string                `json:"video_addr"  example:"视频地址"`         // 视频地址
	VideoDuration int                   `json:"video_duration" example:"100000"`       // 视频时长
	CreateAt      int                   `json:"create_at" example:"1600000000"`        // 视频创建时间
	BrowseNum     int                   `json:"browse_num" example:"10"`               // 浏览数（播放数）
	UserId        string                `json:"user_id" example:"发布视频的用户id"`      // 发布视频的用户id
	Avatar        string                `json:"avatar" example:"头像"`                 // 头像
	Nickname      string                `json:"nick_name"  example:"昵称"`             // 昵称
	Size          int64                 `json:"size"`                                 // 视频总字节数
	BarrageNum    int                   `json:"barrage_num" example:"1"`              // 弹幕数
	Labels        []*models.VideoLabels `json:"label_names,omitempty"`                // 标签名称 多个用逗号分隔
}

// 分享的帖子信息 todo: 产品当前逻辑 转发的帖子皆为文本
type SharePostInfo struct {
	PostId        int64                       `json:"post_id"`                // 转发的帖子id
	PostingType   int                         `json:"posting_type"`           // 帖子类型  0 纯文本 1 图文 2 视频 + 文
	ContentType   int                         `json:"content_type"`           // 0 社区发布 1 转发视频 2 转发帖子
	Topics        []*models.CommunityTopic    `json:"topic_names"`            // 话题
	Title         string                      `json:"title"`                  // 标题
	Describe      string                      `json:"describe"`               // 描述
	Content       string                      `json:"content,omitempty"`       // 暂时使用不到
	BrowseNum     int                         `json:"browse_num" example:"10"` // 浏览数（播放数）
	CommentNum    int                         `json:"comment_num"`             // 评论数
	UserId        string                      `json:"user_id"`                 // up主id
	Nickname      string                      `json:"nick_name"`               // up主昵称
	Avatar        string                      `json:"avatar"`                  // up主头像
	ImagesAddr    []string                    `json:"images_addr,omitempty"`   // 图片地址
}

func NewShareModel(engine *xorm.Session) *ShareModel {
	return &ShareModel{
		Engine: engine,
		Share: new(models.ShareRecord),
	}
}

// 添加转发记录
func (m *ShareModel) AddShareRecord() (int64, error) {
	return m.Engine.InsertOne(m.Share)
}

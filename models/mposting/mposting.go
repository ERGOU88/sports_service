package mposting

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

// 帖子模块
type PostingModel struct {
	Engine            *xorm.Session
	Posting           *models.PostingInfo
	PostingSection    *models.PostingSection
	PostingTopic      *models.PostingTopic
	Browse            *models.UserBrowseRecord
}

// 发布帖子请求参数
type PostPublishParam struct {
	Title          string     `json:"title"`                              // 标题
	Describe       string     `json:"describe"`                           // 富文本内容
	Content        string     `json:"content" binding:"required"`         // 图片列表 或 转发的视频内容 todo: 结构体
	VideoId        string     `json:"video_id"`         // 转发的视频id
	//PostingType    int        `json:"posting_type" binding:"required"`    // 帖子类型  0 纯文本 1 图文 2 视频 + 文
	ImagesAddr     []string   `json:"images_addr"`                        // 图片地址
	SectionId      int        `json:"section_id"`                         // 主模块id
	TopicIds       []string   `json:"topic_ids"`                          // 话题id （多个）
	ContentType    int        `json:"content_type"`                       // 0 发布 1 转发
}

// 帖子详情数据
type PostDetailInfo struct {
	PostId        int64                  `json:"post_id"  example:"1000000000"`        // 帖子id
	Cover         string                 `json:"cover"  example:"封面"`                 // 封面
	VideoAddr     string                 `json:"video_addr"  example:"视频地址"`         // 视频地址
	VideoDuration int                    `json:"video_duration" example:"100000"`       // 视频时长
	Size          int                    `json:"size"`                                 // 视频大小
	Title         string                 `json:"title"  example:"标题"`                 // 标题
	Describe      string                 `json:"describe"  example:"描述"`              // 描述
	Images        []string               `json:"images,omitempty"`                     // 图片列表
	ForwardVideo  string                 `json:"forward_video,omitempty"`              // 分享的视频信息 todo: 结构体
	IsRecommend   int                    `json:"is_recommend" example:"0"`             // 是否推荐
	IsTop         int                    `json:"is_top"  example:"0"`                   // 是否置顶
	Status        int32                  `json:"status"  example:"1"`                   // 审核状态
	CreateAt      int                    `json:"create_at" example:"1600000000"`        // 创建时间
	FabulousNum   int                    `json:"fabulous_num" example:"10"`             // 点赞数
	CommentNum    int                    `json:"comment_num" example:"10"`              // 评论数
	ShareNum      int                    `json:"share_num" example:"10"`                // 分享数
	BrowseNum     int                    `json:"browse_num" example:"10"`               // 浏览数
	UserId        string                 `json:"user_id" example:"发布视频的用户id"`      // 发布者用户id
	Avatar        string                 `json:"avatar" example:"头像"`                 // 头像
	Nickname      string                 `json:"nick_name"  example:"昵称"`             // 昵称
	IsAttention   int                    `json:"is_attention" example:"1"`             // 是否关注 1 关注 0 未关注
	IsLike        int                    `json:"is_like" example:"1"`                  // 是否点赞
	FansNum       int64                  `json:"fans_num" example:"100"`               // 粉丝数
	Topics        []*models.PostingTopic `json:"topics"`                               // 所属话题
}

// 实栗
func NewPostingModel(engine *xorm.Session) *PostingModel {
	return &PostingModel{
		Engine: engine,
		Posting: new(models.PostingInfo),
		PostingSection: new(models.PostingSection),
		PostingTopic: new(models.PostingTopic),
	}
}

// 通过id获取帖子
func (m *PostingModel) GetPostById(id string) (*models.PostingInfo, error) {
	m.Posting = new(models.PostingInfo)
	ok, err := m.Engine.Where("id=?", id).Get(m.Posting)
	if !ok || err != nil {
		return nil, err
	}

	return m.Posting, nil
}

// 获取帖子所属板块 [1对1]
func (m *PostingModel) GetPostSection(postId int64) (*models.PostingSection, error) {
	m.PostingSection = new(models.PostingSection)
	if err := m.Engine.Where("posting_id=?", postId).Asc( "id").Find(m.PostingSection); err != nil {
		return nil, err
	}

	return m.PostingSection, nil
}

// 添加帖子所属板块
func (m *PostingModel) AddPostSection() (int64, error) {
	return m.Engine.InsertOne(m.PostingSection)
}

// 获取帖子所属话题 [1对多]
func (m *PostingModel) GetPostTopic(postId string) ([]*models.PostingTopic, error) {
	var list []*models.PostingTopic
	if err := m.Engine.Where("posting_id=?", postId).Asc("id").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 添加帖子
func (m *PostingModel) AddPost() (int64, error) {
	return m.Engine.InsertOne(m.Posting)
}

// 删除帖子
func (m *PostingModel) DelPost(postId int64) error {
	if _, err := m.Engine.Where("id=?", postId).Delete(&models.PostingInfo{}); err != nil {
		return err
	}

	return nil
}

// 添加帖子所属话题(一次插入多条)
func (m *PostingModel) AddPostingTopics(topics []*models.PostingTopic) (int64, error) {
	return m.Engine.InsertMulti(topics)
}

const (
	UPDATE_POST_BROWSE_NUM  = "UPDATE `posting_statistic` SET `browse_num` = `browse_num` + ?, `update_at`=? WHERE " +
		"`posting_id`=? AND `browse_num` + ? >= 0 LIMIT 1"
)
// 更新帖子浏览数
func (m *PostingModel) UpdatePostBrowseNum(postId int64, now, num int) error {
	if _, err := m.Engine.Exec(UPDATE_POST_BROWSE_NUM, num, now, postId, num); err != nil {
		return err
	}

	return nil
}

const (
	UPDATE_POST_LIKE_NUM  = "UPDATE `posting_statistic` SET `fabulous_num` = `fabulous_num` + ?, `update_at`=? WHERE `posting_id`=? AND `fabulous_num` + ? >= 0 LIMIT 1"
)
// 更新帖子点赞数
func (m *PostingModel) UpdatePostLikeNum(videoId int64, now, num int) error {
	if _, err := m.Engine.Exec(UPDATE_POST_LIKE_NUM, num, now, videoId, num); err != nil {
		return err
	}

	return nil
}

const (
	UPDATE_POST_COMMENT_NUM = "UPDATE `posting_statistic` SET `comment_num` = `comment_num` + ?, `update_at`=? WHERE `posting_id`=? AND `comment_num` + ? >= 0 LIMIT 1"
)
// 更新帖子评论数
func (m *PostingModel) UpdatePostCommentNum(postId int64, now, num int) error {
	if _, err := m.Engine.Exec(UPDATE_POST_COMMENT_NUM, num, now, postId, num); err != nil {
		return err
	}

	return nil
}


// 获取用户浏览过的帖子
func (m *PostingModel) GetUserBrowsePost(userId string, composeType int, composeId int64) *models.UserBrowseRecord {
	m.Browse = new(models.UserBrowseRecord)
	ok, err := m.Engine.Where("user_id=? AND compose_type=? AND compose_id=?", userId, composeType, composeId).Get(m.Browse)
	if !ok || err != nil {
		return nil
	}

	return m.Browse
}

// 记录用户浏览的帖子记录
func (m *PostingModel) RecordUserBrowsePost() error {
	if _, err := m.Engine.InsertOne(m.Browse); err != nil {
		return err
	}

	return nil
}

// 之前有浏览记录 更新浏览时间
func (m *PostingModel) UpdateUserBrowsePost(userId string,  composeType int, composeId int64) error {
	if _, err := m.Engine.Where("user_id=? AND compose_id=? AND compose_type=?", userId, composeId, composeType).Cols("create_at, update_at").Update(m.Browse); err != nil {
		return err
	}

	return nil
}

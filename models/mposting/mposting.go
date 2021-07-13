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
}

// 发布帖子请求参数
type PostPublishParam struct {
	Title          string     `json:"title"`                              // 标题
	Content        string     `json:"content" binding:"required"`         // 富文本内容
	Cover          string     `json:"cover"`                            // 视频封面
	VideoAddr      string     `json:"video_addr"`     // 视频地址
	VideoDuration  int        `json:"video_duration"` // 视频时长（秒）
	//VideoId        string     `json:"video_id"`         // 转发的视频id
	//PostingType    int        `json:"posting_type" binding:"required"`    // 帖子类型  0 纯文本 1 图文 2 视频 + 文
	ImagesAddr     []string   `json:"images_addr"`                        // 图片地址
	SectionId      int        `json:"section_id"`                         // 主模块id
	TopicIds       []string   `json:"topic_ids"`                          // 话题id （多个）
	ContentType    int        `json:"content_type"`                       // 0 发布 1 转发
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

// 获取帖子所属板块 [1对多]
func (m *PostingModel) GetPostSection(postId int64) ([]*models.PostingSection, error){
	var list []*models.PostingSection
	if err := m.Engine.Where("posting_id=?", postId).Asc( "id").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取帖子所属话题 [1对多]
func (m *PostingModel) GetPostTopic(postId int64) ([]*models.PostingTopic, error) {
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



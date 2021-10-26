package mcommunity

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

// 社区模块
type CommunityModel struct {
	Engine              *xorm.Session
	CommunitySection    *models.CommunitySection
	CommunityTopic      *models.CommunityTopic
}

// 社区话题信息
type CommunityTopicInfo struct {
	Id        int    `json:"id"`
	TopicName string `json:"topic_name"`
	IsHot     int    `json:"is_hot"`
	Cover     string `json:"cover,omitempty"`
	Describe  string `json:"describe,omitempty"`
	PostNum   int64  `json:"post_num"`      // 帖子数量
}

// 社区板块信息
type CommunitySectionInfo struct {
	Id          int    `json:"id"`
	SectionName string `json:"section_name"`
	PostNum     int64  `json:"post_num"`
}

// 社区实栗
func NewCommunityModel(engine *xorm.Session) *CommunityModel {
	return &CommunityModel{
		Engine: engine,
		CommunitySection: new(models.CommunitySection),
		CommunityTopic: new(models.CommunityTopic),
	}
}

// 通过id获取板块信息
func (m *CommunityModel) GetSectionInfo(id string) (*models.CommunitySection, error) {
	m.CommunitySection = new(models.CommunitySection)
	ok, err := m.Engine.Where("id=? AND status=1", id).Get(m.CommunitySection)
	if !ok || err != nil {
		return nil, err
	}

	return m.CommunitySection, nil
}

// 通过名称获取板块
func (m *CommunityModel) GetSectionByName(name string) (bool, error) {
	return m.Engine.Where("section_name=? AND status=1", name).Get(m.CommunitySection)
}

// 通过id获取社区话题信息
func (m *CommunityModel) GetTopicInfo(id string) (*models.CommunityTopic, error) {
	m.CommunityTopic = new(models.CommunityTopic)
	ok, err := m.Engine.Where("id =? AND status=1", id).Get(m.CommunityTopic)
	if !ok || err != nil {
		return nil, err
	}

	return m.CommunityTopic, nil
}

// 根据话题id获取多个话题信息
func (m *CommunityModel) GetTopicByIds(ids []string) ([]*models.CommunityTopic, error) {
	var list []*models.CommunityTopic
	if err := m.Engine.In("id", ids).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取社区所有板块
func (m *CommunityModel) GetAllSection() ([]*models.CommunitySection, error) {
	var list []*models.CommunitySection
	if err := m.Engine.Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取社区话题
func (m *CommunityModel) GetCommunityTopics(sectionId, isHot string, offset, size int) ([]*models.CommunityTopic, error) {
	var list []*models.CommunityTopic
	table := m.Engine.Where("status=1").Desc("sortorder")

	if sectionId != "" {
		table = table.Where("section_id=?", sectionId)
	}

	if isHot != "" {
		table = table.Where("is_hot=?", isHot)
	}

	if offset >= 0 && size > 0 {
		table = table.Limit(size, offset)
	}

	if err := table.Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

const (
	QUERY_TOPIC_LIST = "SELECT ct.*, post_num FROM community_topic AS ct LEFT JOIN (SELECT count(1) as post_num, " +
		"topic_id FROM posting_topic as pt WHERE pt.status=1 GROUP BY pt.`topic_id`) AS pt ON ct.id=pt.topic_id " +
		"ORDER BY pt.post_num DESC, ct.sortorder DESC, ct.id DESC LIMIT ?, ?"
)
// 获取话题列表 [按话题下的帖子数量排序]
func (m *CommunityModel) GetTopicListOrderByPostNum(offset, size int) ([]*CommunityTopicInfo, error) {
	var list []*CommunityTopicInfo
	if err := m.Engine.SQL(QUERY_TOPIC_LIST, offset, size).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取所有话题
func (m *CommunityModel) GetAllTopic() ([]*models.CommunityTopic, error) {
	var list []*models.CommunityTopic
	if err := m.Engine.Where("status=1").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

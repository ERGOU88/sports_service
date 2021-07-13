package mcommunity

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

// 社区模块
type CommunityModel struct {
	Engine              *xorm.Session
	Community           *models.CommunityInfo
	CommunitySection    *models.CommunitySection
	CommunityTopic      *models.CommunityTopic
}

// 社区实栗
func NewCommunityModel(engine *xorm.Session) *CommunityModel {
	return &CommunityModel{
		Engine: engine,
		Community: new(models.CommunityInfo),
		CommunitySection: new(models.CommunitySection),
		CommunityTopic: new(models.CommunityTopic),
	}
}

// 通过社区id获取社区信息
func (m *CommunityModel) GetCommunityInfo(id int) (*models.CommunityInfo, error) {
	m.Community = new(models.CommunityInfo)
	ok, err := m.Engine.Where("id=? AND status=1", id).Get(m.Community)
	if !ok || err != nil {
		return nil, err
	}

	return m.Community, nil
}

// 通过社区id获取板块信息
func (m *CommunityModel) GetSectionInfo(id int) (*models.CommunitySection, error) {
	m.CommunitySection = new(models.CommunitySection)
	ok, err := m.Engine.Where("id=? AND status=1", id).Get(m.CommunitySection)
	if !ok || err != nil {
		return nil, err
	}

	return m.CommunitySection, nil
}

// 通过社区id获取社区话题信息
func (m *CommunityModel) GetTopicInfo(id int) (*models.CommunityTopic, error) {
	m.CommunityTopic = new(models.CommunityTopic)
	ok, err := m.Engine.Where("id =? AND status=1", id).Get(m.CommunityTopic)
	if !ok || err != nil {
		return nil, err
	}

	return m.CommunityTopic, nil
}

// 获取多个社区话题信息
func (m *CommunityModel) GetTopicByIds(ids string) ([]*models.PostingTopic, error) {
	var list []*models.PostingTopic
	if err := m.Engine.In("id", ids).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

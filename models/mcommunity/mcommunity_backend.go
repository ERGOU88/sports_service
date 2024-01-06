package mcommunity

import "sports_service/models"

// 添加板块
type AddSection struct {
	Id          int    `json:"id"`
	SectionName string `json:"section_name"`
	Sortorder   int    `json:"sortorder"`
	Status      int    `json:"status"`
}

type DelSection struct {
	Id int `json:"id"`
}

type AddTopic struct {
	Id        int    `json:"id"`
	TopicName string `json:"topic_name"`
	Sortorder int    `json:"sortorder"`
	Cover     string `json:"cover"`
	Describe  string `json:"describe"`
	Status    int    `json:"status"`
	IsHot     int    `json:"is_hot"`
	SectionId int    `json:"section_id"`
}

type DelTopic struct {
	Id int `json:"id"`
}

func (m *CommunityModel) SectionTableName() string {
	return "community_section"
}

func (m *CommunityModel) TopicTableName() string {
	return "community_topic"
}

// 添加社区板块
func (m *CommunityModel) AddCommunitySection() (int64, error) {
	return m.Engine.InsertOne(m.CommunitySection)
}

// 删除社区板块
func (m *CommunityModel) DelCommunitySection(id int) (int64, error) {
	return m.Engine.Where("id=?", id).Delete(&models.CommunitySection{})
}

func (m *CommunityModel) UpdateSectionInfo(id int, mp map[string]interface{}) (int64, error) {
	return m.Engine.Table(m.SectionTableName()).Where("id=?", id).Update(mp)
}

func (m *CommunityModel) AddTopic() (int64, error) {
	return m.Engine.InsertOne(m.CommunityTopic)
}

// 修改话题
func (m *CommunityModel) UpdateTopicInfo(id int, mp map[string]interface{}) (int64, error) {
	return m.Engine.Table(m.TopicTableName()).Where("id=?", id).Update(mp)
}

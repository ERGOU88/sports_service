package mposting

import "sports_service/server/models"

// 添加申精记录
func (m *PostingModel) AddApplyCreamRecord() (int64, error) {
	return m.Engine.InsertOne(m.ApplyCream)
}

// 通过帖子id获取正在审批或审批通过的申精记录
func (m *PostingModel) GetApplyCreamRecord(postId string) (*models.PostingApplyCream, error) {
	m.ApplyCream = new(models.PostingApplyCream)
	ok, err := m.Engine.Where("post_id=? AND status in(0, 1)", postId).Get(m.ApplyCream)
	if !ok || err != nil {
		return nil, err
	}

	return m.ApplyCream, nil
}


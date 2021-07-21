package mvideo

import (
	"sports_service/server/models"
)

// 获取视频分区配置列表
func (m *VideoModel) GetSubAreaList() ([]*models.VideoSubarea, error) {
	var list []*models.VideoSubarea
	if err := m.Engine.Where("status=0").Desc("sortorder").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 通过分区id获取视频分区信息
func (m *VideoModel) GetSubAreaById(id string) (*models.VideoSubarea, error) {
	m.Subarea = new(models.VideoSubarea)
	ok, err := m.Engine.Where("id=?", id).Get(m.Subarea)
	if !ok || err != nil {
		return nil, err
	}

	return m.Subarea, nil
}

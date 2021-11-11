package mvideo

import (
	"sports_service/server/models"
)

type AddSubarea struct {
	Name      string    `json:"name"`
	SortOrder int       `json:"sortorder"`
	SysId     int       `json:"sys_id"`
	SysUser   string    `json:"sys_user"`
	Id        int       `json:"id"`
	Status    int       `json:"status"`
}

type DelSubarea struct {
	Id      int     `json:"id"`
}

func (m *VideoModel) AddSubArea() (int64, error) {
	return m.Engine.InsertOne(m.Subarea)
}

func (m *VideoModel) UpdateSubArea() (int64, error) {
	return m.Engine.Where("id=?", m.Subarea.Id).Cols("status, update_at, sortorder").Update(m.Subarea)
}

func (m *VideoModel) DelSubArea(id int) (int64, error) {
	return m.Engine.Where("id=?", id).Delete(m.Subarea)
}

// 获取视频分区配置列表
func (m *VideoModel) GetSubAreaList(status []int) ([]*models.VideoSubarea, error) {
	var list []*models.VideoSubarea
	if err := m.Engine.In("status", status).Desc("sortorder").Find(&list); err != nil {
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

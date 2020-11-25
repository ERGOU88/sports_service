package mvideo

import (
  "fmt"
  "sports_service/server/global/app/log"
  "sports_service/server/models"
)

// 添加视频标签(一次插入多条)
func (m *VideoModel) AddVideoLabels(labels []*models.VideoLabels) (int64, error) {
	return m.Engine.InsertMulti(labels)
}

// 获取视频标签
func (m *VideoModel) GetVideoLabels(videoId string) []*models.VideoLabels {
	var labels []*models.VideoLabels
	if err := m.Engine.Where("video_id=?", videoId).Find(&labels); err != nil {
		log.Log.Errorf("video_trace: get video labels err:%s", err)
		return nil
	}

	return labels
}

// 通过标签id 查询视频id列表
func (m *VideoModel) GetVideoIdsByLabelId(labelId string, offset, size int) []string {
	var videoIds []string
	if err := m.Engine.Table(&models.VideoLabels{}).Where("label_id=?", labelId).Cols("video_id").Limit(size, offset).Find(&videoIds); err != nil {
		log.Log.Errorf("video_trace: get videoIds by labelid err:%s", err)
		return nil
	}

	return videoIds
}

// 通过标签id列表 查询视频id列表
func (m *VideoModel) FindVideoIdsByLabelIds(labelIds string, offset, size int) []string {
	var videoIds []string
	sql := fmt.Sprintf("SELECT DISTINCT(video_id) FROM video_labels WHERE label_id IN(%s) LIMIT ?, ?", labelIds)
	if err := m.Engine.Table(&models.VideoLabels{}).SQL(sql, offset, size).Find(&videoIds); err != nil {
		log.Log.Errorf("video_trace: find videoIds by label ids err:%s", err)
		return nil
	}

	return videoIds
}

// 删除视频标签
func (m *VideoModel) DelVideoLabels(videoId string) error {
	if _, err := m.Engine.Where("video_id=?", videoId).Delete(m.Labels); err != nil {
		return err
	}

	return nil
}

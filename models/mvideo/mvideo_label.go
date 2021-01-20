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

// 更新视频标签信息
func (m *VideoModel) UpdateVideoLabelInfo(condition, cols string) (int64, error) {
  return m.Engine.Where(condition).Cols(cols).Update(m.Labels)
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

const (
  RANDOM_QUERY_VIDEO_BY_LABELS = "SELECT DISTINCT(vl.video_id) FROM `video_labels` AS vl " +
    "JOIN (SELECT RAND() * ((SELECT MAX(video_id) FROM `video_labels`)-(SELECT MIN(video_id) FROM `video_labels`)) " +
    "+ (SELECT MIN(video_id) FROM `video_labels`) AS video_id) AS vl2 " +
    "WHERE vl.video_id >= vl2.video_id-%d AND vl.video_id != %s AND vl.label_id in(%s) AND vl.status=1 " +
    "ORDER BY vl.video_id LIMIT %d"
)
// 通过标签列表 随机获取同标签类型的视频列表 (取20条)
func (m *VideoModel) RandomGetVideoIdByLabels(videoId, labelIds string, size int) []string {
  var videoIds []string
  sql := fmt.Sprintf(RANDOM_QUERY_VIDEO_BY_LABELS, size, videoId, labelIds, size)
  if err := m.Engine.Table(&models.VideoLabels{}).SQL(sql).Find(&videoIds); err != nil {
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

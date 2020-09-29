package mbarrage

import (
	"sports_service/server/global/app/log"
	"sports_service/server/models"
	"github.com/go-xorm/xorm"
)

type BarrageModel struct {
	Barrage      *models.VideoBarrage
	Engine       *xorm.Session
}

// 发送弹幕请求参数
type SendBarrageParams struct {
	Color            string `json:"color"`
	Content          string `binding:"required" json:"content"`
	Font             string `json:"font"`
	Location         int    `json:"location"`
	VideoCurDuration int    `binding:"required" json:"video_cur_duration"`
	VideoId          int64  `binding:"required" json:"video_id"`
}

// 实栗
func NewBarrageModel(engine *xorm.Session) *BarrageModel {
	return &BarrageModel{
		Engine: engine,
		Barrage: new(models.VideoBarrage),
	}
}

// 记录视频弹幕
func (m *BarrageModel) RecordVideoBarrage() error {
	if _, err := m.Engine.Insert(m.Barrage); err != nil {
		log.Log.Errorf("barrage_trace: record video barrage err:%s", err)
		return err
	}

	return nil
}

// 根据视频时长区间获取弹幕 todo:限制下 最多取最新的1000条 根据取的时间区间大小做调整
func (m *BarrageModel) GetBarrageByDuration(videoId, minDuration, maxDuration string, offset, limit int) []*models.VideoBarrage {
	var list []*models.VideoBarrage
	if err := m.Engine.Where(" video_id =? AND video_cur_duration >= ? AND video_cur_duration <= ?", videoId,
		minDuration, maxDuration).Desc("id").Limit(limit, offset).Find(&list); err != nil {
		return []*models.VideoBarrage{}
	}

	return list
}

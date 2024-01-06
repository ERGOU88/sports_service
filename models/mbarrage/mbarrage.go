package mbarrage

import (
	"github.com/go-xorm/xorm"
	"sports_service/global/app/log"
	"sports_service/models"
)

type BarrageModel struct {
	Barrage *models.VideoBarrage
	Engine  *xorm.Session
}

// 发送弹幕请求参数
type SendBarrageParams struct {
	Color            string `json:"color"`
	Content          string `binding:"required" json:"content"`
	Font             string `json:"font"`
	Location         int    `json:"location"`
	VideoCurDuration int    `json:"video_cur_duration"`
	VideoId          int64  `binding:"required" json:"video_id"` // 视频id/直播id
	BarrageType      int64  `json:"barrage_type"`                // 0首页视频弹幕 1直播/回放视频弹幕
}

// 实栗
func NewBarrageModel(engine *xorm.Session) *BarrageModel {
	return &BarrageModel{
		Engine:  engine,
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
func (m *BarrageModel) GetBarrageByDuration(videoId, barrageType, minDuration, maxDuration string, offset, limit int) []*models.VideoBarrage {
	var list []*models.VideoBarrage
	if err := m.Engine.Where(" video_id =? AND barrage_type=? AND video_cur_duration >= ? AND video_cur_duration <= ?", videoId,
		barrageType, minDuration, maxDuration).Desc("id").Limit(limit, offset).Find(&list); err != nil {
		return []*models.VideoBarrage{}
	}

	return list
}

// 获取用户视频总弹幕数
func (m *BarrageModel) GetUserTotalVideoBarrage(userId string) int64 {
	total, err := m.Engine.Where("user_id=?", userId).Count(&models.VideoBarrage{})
	if err != nil {
		log.Log.Errorf("barrage_trace: get user total barrage err:%s, uid:%s", err, userId)
		return 0
	}

	return total
}

type VideoBarrageInfo struct {
	Id               int64  `json:"id"  example:"10000000000"`
	VideoId          int64  `json:"video_id" example:"10000000000"`
	VideoCurDuration int    `json:"video_cur_duration" example:"100"`
	Content          string `json:"content" example:"弹幕内容"`
	UserId           string `json:"user_id" example:"10000000000"`
	Color            string `json:"color" example:"颜色"`
	Font             string `json:"font" example:"字体"`
	BarrageType      int    `json:"barrage_type" example:"0"`
	Location         int    `json:"location" example:"0"`
	SendTime         int64  `json:"send_time" example:"1600000000"`
	Title            string `json:"title,omitempty" example:"标题"`
	VideoAddr        string `json:"video_addr,omitempty" example:"视频地址"`
}

// 后台分页获取 视频弹幕 列表
func (m *BarrageModel) GetVideoBarrageList(offset, size int) []*VideoBarrageInfo {
	sql := "SELECT vb.*, v.title, v.video_addr FROM video_barrage AS vb LEFT JOIN videos AS v ON vb.video_id=v.video_id " +
		"WHERE vb.barrage_type=0 GROUP BY vb.id ORDER BY vb.id DESC LIMIT ?, ?"
	var list []*VideoBarrageInfo
	if err := m.Engine.SQL(sql, offset, size).Find(&list); err != nil {
		log.Log.Errorf("barrage_trace: get video barrage list by sort, err:%s", err)
		return []*VideoBarrageInfo{}
	}

	return list
}

// 后台分页获取 直播弹幕 列表
func (m *BarrageModel) GetLiveBarrageList(offset, size int) []*models.VideoBarrage {
	var list []*models.VideoBarrage
	if err := m.Engine.Where("barrage_type=1").Desc("id").Limit(size, offset).Find(&list); err != nil {
		return []*models.VideoBarrage{}
	}

	return list
}

// 后台删除弹幕请求参数
type DelBarrageParam struct {
	Id          string `binding:"required" json:"id"` // 弹幕id
	BarrageType int32  `json:"barrage_type"`          // 弹幕类型
}

const (
	DELETE_VIDEO_BARRAGE = "DELETE FROM `video_barrage` WHERE `id`=?"
)

// 删除弹幕
func (m *BarrageModel) DelVideoBarrage(id string) error {
	if _, err := m.Engine.Exec(DELETE_VIDEO_BARRAGE, id); err != nil {
		return err
	}

	return nil
}

// 获取弹幕总数
func (m *BarrageModel) GetVideoBarrageTotal(barrageType string) int64 {
	count, err := m.Engine.Where("barrage_type=?", barrageType).Count(&models.VideoBarrage{})
	if err != nil {
		log.Log.Errorf("comment_trace: get total barrage err:%s", err)
		return 0
	}

	return count
}

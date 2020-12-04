package mcollect

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/global/app/log"
	"sports_service/server/models"
	"time"
	"fmt"
)

type CollectModel struct {
	CollectRecord *models.CollectRecord
	Engine        *xorm.Session
}

// 添加收藏请求参数
type AddCollectParam struct {
	VideoId       int64     `binding:"required" json:"video_id" example:"10001"`       // 收藏的视频id
}

// 取消收藏请求参数
type CancelCollectParam struct {
	VideoId       int64     `binding:"required" json:"video_id" example:"10001"`       // 取消收藏的视频id
}

// 删除收藏记录请求参数
type DeleteCollectParam struct {
	ComposeIds        []string     `binding:"required" json:"compose_ids"`           // 作品id列表
}

// 实栗
func NewCollectModel(engine *xorm.Session) *CollectModel {
	return &CollectModel{
		CollectRecord: new(models.CollectRecord),
		Engine:        engine,
	}
}

// 添加视频收藏
func (m *CollectModel) AddCollectVideo(userId, toUserId string, videoId int64, status, composeType int) error {
	m.CollectRecord.UserId = userId
	m.CollectRecord.ToUserId = toUserId
	m.CollectRecord.ComposeId = videoId
	m.CollectRecord.ComposeType = composeType
	m.CollectRecord.UpdateAt = int(time.Now().Unix())
	m.CollectRecord.CreateAt = int(time.Now().Unix())
	m.CollectRecord.Status = status
	if _, err := m.Engine.InsertOne(m.CollectRecord); err != nil {
		return err
	}

	return nil
}

// 获取收藏的信息
func (m *CollectModel) GetCollectInfo(userId string, videoId int64, composeType int) *models.CollectRecord {
	m.CollectRecord = new(models.CollectRecord)
	ok, err := m.Engine.Where("user_id=? AND compose_id=? AND compose_type=?", userId, videoId, composeType).Get(m.CollectRecord)
	if !ok || err != nil {
		return nil
	}

	return m.CollectRecord
}

// 更新收藏状态 收藏/取消收藏
func (m *CollectModel) UpdateCollectStatus() error {
	if _, err := m.Engine.ID(m.CollectRecord.Id).
		Cols("status, update_at").
		Update(m.CollectRecord); err != nil {
		return err
	}

	return nil
}

type CollectVideosInfo struct {
	ComposeId int64 `json:"compose_id"`
	UpdateAt  int   `json:"update_at"`
}
// 获取收藏的作品id列表
func (m *CollectModel) GetCollectList(userId string, offset, size int) []*CollectVideosInfo {
	var list []*CollectVideosInfo
	if err := m.Engine.Table(&models.CollectRecord{}).Where("status=1 AND user_id=?", userId).
		Cols("compose_id, update_at").
		Desc("update_at", "id").
		Limit(size, offset).
		Find(&list); err != nil {
		log.Log.Errorf("collect_trace: get collect videos err:%s", err)
		return nil
	}

	return list
}

// 通过id列表删除收藏记录
func (m *CollectModel) DeleteCollectByIds(userId string, ids string) error {
  var sql string
  if ids == "-1" {
    // -1 删除所有收藏的视频
    sql = fmt.Sprintf("DELETE FROM `collect_record` WHERE `user_id`=%s AND compose_type = 0", userId)
  } else {
    sql = fmt.Sprintf("DELETE FROM `collect_record` WHERE `user_id`=%s AND compose_id in(%s)", userId,  ids)
  }

	if _, err := m.Engine.Exec(sql); err != nil {
		return err
	}

	return nil
}

// 获取用户收藏的作品总数
func (m *CollectModel) GetUserTotalCollect(userId string) int64 {
	total, err := m.Engine.Where("user_id=? AND status=1", userId).Count(m.CollectRecord)
	if err != nil {
		log.Log.Errorf("collect_trace: get collect total err:%s, uid:%s", err, userId)
		return 0
	}

	return total
}

// 获取用户被收藏的作品列表
func (m *CollectModel) GetCollectedList(toUserId string, offset, size int) []*models.CollectRecord {
	var list []*models.CollectRecord
	if err := m.Engine.Where("to_user_id=? AND status=1", toUserId).Desc("id").Limit(size, offset).Find(&list); err != nil {
		log.Log.Errorf("collect_trace: get collected list err:%s", err)
		return nil
	}

	return list
}

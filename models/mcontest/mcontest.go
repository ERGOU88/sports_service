package mcontest

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

// 赛事
type ContestModel struct {
	Engine          *xorm.Session
	VideoLive       *models.VideoLive
}

// 赛事直播信息
type ContestLiveInfo struct {
	Id             int64  `json:"id"`
	UserId         string `json:"user_id"`
	RoomId         string `json:"room_id"`
	GroupId        string `json:"group_id"`
	Cover          string `json:"cover"`
	RtmpAddr       string `json:"rtmp_addr"`
	FlvAddr        string `json:"flv_addr"`
	HlsAddr        string `json:"hls_addr"`
	PlayTime       int    `json:"play_time"`
	Title          string `json:"title"`
	HighLights     string `json:"high_lights"`
	Describe       string `json:"describe"`
	Tags           string `json:"tags"`
	LiveType       int    `json:"live_type"`
	NickName       string `json:"nick_name"`
	Avatar         string `json:"avatar"`
	Date           string `json:"date"`
	Week           string `json:"week"`
}

// 实例
func NewContestModel(engine *xorm.Session) *ContestModel {
	return &ContestModel{
		Engine: engine,
		VideoLive: new(models.VideoLive),
	}
}

// 获取直播列表
func (m *ContestModel) GetLiveList(now int64, offset, size int) ([]*models.VideoLive, error){
	var list []*models.VideoLive
	if err := m.Engine.Where("play_time >= ?", now).Asc("play_time").Limit(size, offset).Find(&list); err != nil {
		return []*models.VideoLive{}, err
	}

	return list, nil
}



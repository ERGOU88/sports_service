package mcontest

import "sports_service/server/models"

type baseInfo struct {
	EventType        int       `json:"event_type"`          // 1推流  0断流  100直播录制
	Appid            int       `json:"appid"`
	App              string    `json:"app"`
	Appname          string    `json:"appname"`
	StreamID         string    `json:"stream_id"`
	ChannelID        string    `json:"channel_id"`
	StreamParam      string    `json:"stream_param"`
	Sign             string    `json:"sign"`
	T                int       `json:"t"`
}

// 推流/断流/流录制 回调信息
type StreamCallbackInfo struct {
	baseInfo
	Errcode      int       `json:"errcode"`
	Errmsg       string    `json:"errmsg"`
	EventTime    int       `json:"event_time"`
	SetID        int       `json:"set_id"`
	Node         string    `json:"node"`
	Sequence     string    `json:"sequence"`
	UserIP       string    `json:"user_ip"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	PushDuration string   `json:"push_duration"`
	VideoURL     string    `json:"video_url"`
	FileID       string    `json:"file_id"`
	FileFormat   string    `json:"file_format"`
	TaskID       string    `json:"task_id"`
	StartTime    int       `json:"start_time"`         // 录制文件起始时间戳
	EndTime      int       `json:"end_time"`           // 录制文件结束时间戳
	Duration     int       `json:"duration"`           // 录制文件时长，单位秒
	FileSize     int       `json:"file_size"`          // 录制文件大小，单位字节
	Status       int                                   // 0未直播 1直播中 2 已结束
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
	Status         int    `json:"int"`           // 状态 0未直播 1直播中 2 已结束
	HasReplay      int    `json:"has_replay"`    // 是否有回放 1 有 2 无
}

// 获取直播列表
func (m *ContestModel) GetLiveList(now int64, offset, size int, contestId, status string) ([]*models.VideoLive, error){
	var list []*models.VideoLive
	m.Engine.Where("play_time >= ? AND contest_id=?", now, contestId)
	if status == "1" {
		m.Engine.Where("status=1")
	}

	if err := m.Engine.Asc("play_time").Limit(size, offset).Find(&list); err != nil {
		return []*models.VideoLive{}, err
	}

	return list, nil
}

// 通过房间id获取直播信息
func (m *ContestModel) GetLiveInfoByRoomId(roomId string) (bool, error) {
	m.VideoLive = new(models.VideoLive)
	return m.Engine.Where("room_id=?", roomId).Get(m.VideoLive)
}

// 更新直播数据
func (m *ContestModel) UpdateLiveInfo(cols string) (int64, error) {
	return m.Engine.ID(m.VideoLive.Id).Cols(cols).Update(m.VideoLive)
}

// 添加直播回放
func (m *ContestModel) AddVideoLiveReply() (int64, error) {
	return m.Engine.InsertOne(m.VideoLiveReplay)
}

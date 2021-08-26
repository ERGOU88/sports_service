package vod

import "time"

// 拉取事件通知
type PullEventNotify struct {
	Response Response `json:"Response"`
}

type SourceInfo struct {
	SourceType string `json:"SourceType"`
	SourceContext string `json:"SourceContext"`
}

type MediaBasicInfo struct {
	Name string `json:"VenueName"`
	Description string `json:"Description"`
	CreateTime time.Time `json:"CreateTime"`
	UpdateTime time.Time `json:"UpdateTime"`
	ExpireTime time.Time `json:"ExpireTime"`
	ClassID int `json:"ClassId"`
	ClassName string `json:"ClassName"`
	ClassPath string `json:"ClassPath"`
	CoverURL string `json:"CoverUrl"`
	Type string `json:"Type"`
	MediaURL string `json:"MediaUrl"`
	TagSet []interface{} `json:"TagSet"`
	StorageRegion string `json:"StorageRegion"`
	SourceInfo SourceInfo `json:"SourceInfo"`
	Vid string `json:"Vid"`
}

type FileUploadEvent struct {
	FileID string `json:"FileId"`
	MediaBasicInfo MediaBasicInfo `json:"MediaBasicInfo"`
	ProcedureTaskID string `json:"ProcedureTaskId"`
}

type EventSet struct {
	EventHandle string `json:"EventHandle"`
	EventType string `json:"EventType"`
	FileUploadEvent FileUploadEvent `json:"FileUploadEvent"`
}

type Response struct {
	EventSet []EventSet `json:"EventSet"`
	RequestID string `json:"RequestId"`
}

// 普通事件回调
type EventNotify struct {
	EventType string `json:"EventType"`
	FileUploadEvent FileUploadEvent `json:"FileUploadEvent"`
}

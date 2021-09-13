package live

import (
	"fmt"
)

type liveRealize struct {}

var Live ILive

type (
	BaseResponse struct {
		ActionStatus string `json:"ActionStatus"`    // OK 表示处理成功，FAIL 表示失败
		ErrorCode int `json:"ErrorCode"`
		ErrorInfo string `json:"ErrorInfo"`          // 0表示成功，非0表示失败
	}

	ResponseAddUser struct {
		BaseResponse
	}
)

const (
	EXPIRE_TM = 86400
	PUSH_STREAM_HOST = "livepush.bluetrans.cn"
	PULL_STREAM_HOST = "livepull.bluetrans.cn"
	LIVE_PUSH_KEY    = "21a40af5285c5cf9e9d6e1bb84cb3d56"
	LIVE_PULL_KEY    = "uVf1BiIOSh0KQsbH"
	LIVE_CALLBACK_KEY= "V0pQFDIdq4D5IBsQoS56QbSdkeDNlPmL"
)

func NewLiveRealize() *liveRealize {
	return &liveRealize{}
}

func (live *liveRealize) GenPushStream(roomId string, expireTm int64) string {
	txTime := BuildTxTime(expireTm)
	txSecret := BuildTxSecret(LIVE_PUSH_KEY, roomId, expireTm)
	stream := fmt.Sprintf("rtmp://%s/live/%s?txSecret=%s&txTime=%s", PUSH_STREAM_HOST, roomId, txSecret, txTime)
	return stream
}

type PullStreamInfo struct {
	RtmpAddr    string   `json:"rtmp_addr"`
	FlvAddr     string   `json:"flv_addr"`
	HlsAddr     string   `json:"hls_addr"`
}

func (live *liveRealize) GenPullStream(roomId string, expireTime int64) *PullStreamInfo {
	txTime := BuildTxTime(expireTime)
	txSecret := BuildTxSecret(LIVE_PULL_KEY, roomId, expireTime)
	streamInfo := &PullStreamInfo{
		RtmpAddr: fmt.Sprintf("rtmp://%s/live/%s?txSecret=%s&txTime=%s", PULL_STREAM_HOST, roomId, txSecret, txTime),
		FlvAddr: fmt.Sprintf("https://%s/live/%s.flv?txSecret=%s&txTime=%s", PULL_STREAM_HOST, roomId, txSecret, txTime),
		HlsAddr: fmt.Sprintf("https://%s/live/%s.m3u8?txSecret=%s&txTime=%s", PULL_STREAM_HOST, roomId, txSecret, txTime),
	}

	return streamInfo
}

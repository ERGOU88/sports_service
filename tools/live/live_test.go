package live

import (
	"testing"
	"fmt"
)

// 获取推流地址
func TestGenPushStream(t *testing.T) {
	txTime := BuildTxTime(1000)
	txSecret := BuildTxSecret(LIVE_PUSH_KEY,"123", 1000)
	stream := fmt.Sprintf("rtmp://%s/live/%s?txSecret=%s&txTime=%s", PUSH_STREAM_HOST,
		"123", txSecret, txTime)
	t.Logf("stream:%s", stream)
}

func TestGenPullStream(t *testing.T) {
	txTime := BuildTxTime(3600)
	txSecret := BuildTxSecret(LIVE_PULL_KEY, "123", 3600)
	streamInfo := &PullStreamInfo{
		RtmpAddr: fmt.Sprintf("rtmp://%s/live/%s?txSecret=%s&txTime=%s", PULL_STREAM_HOST, "123", txSecret, txTime),
		FlvAddr: fmt.Sprintf("https://%s/live/%s.flv?txSecret=%s&txTime=%s", PULL_STREAM_HOST, "123", txSecret, txTime),
		HlsAddr: fmt.Sprintf("https://%s/live/%s.m3u8?txSecret=%s&txTime=%s", PULL_STREAM_HOST, "123", txSecret, txTime),
	}

	t.Logf("streamInfo:%+v", streamInfo)
}


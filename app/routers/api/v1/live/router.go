package live

import (
	"github.com/gin-gonic/gin"
)

// 直播模块
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	live := api.Group("/live")
	{
		// 推流回调
		live.POST("/push/stream/callback", PushStreamCallback)
		// 断流回调
		live.POST("/disconnect/stream/callback", DisconnectStreamCallback)
		// 录制回调
		live.POST("/transcribe/stream/callback", TranscribeStreamCallback)
	}
}


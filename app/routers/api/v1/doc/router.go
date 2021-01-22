package doc

import (
	"github.com/gin-gonic/gin"
)

// api错误码路由
func Router(engine *gin.Engine) {
	// 错误码文档
	api := engine.Group("/api/v1")
	api.GET("/doc", ApiCode)
  // 推送通知文档
  api.GET("/notify/doc", NotifyDoc)
}

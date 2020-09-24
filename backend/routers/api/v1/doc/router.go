package doc

import (
	"github.com/gin-gonic/gin"
)

// api错误码路由
func Router(engine *gin.Engine) {
	// 错误码文档
	api := engine.Group("/backend/v1")
	api.GET("/doc", ApiCode)
}

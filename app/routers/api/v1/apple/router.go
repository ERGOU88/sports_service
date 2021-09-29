package apple

import (
	"github.com/gin-gonic/gin"
)

// 苹果相关路由
func Router(engine *gin.Engine) {
	engine.GET("/.well-known/apple-app-site-association", AppleLink)
	engine.GET("/qq_conn/101964013", AppleLink)
}

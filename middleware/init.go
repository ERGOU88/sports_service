package middleware

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/header"
)

// 初始化部分中间件
func InitMiddleware(engine *gin.Engine) {
	// 跨域处理
	engine.Use(header.Options)
}

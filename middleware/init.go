package middleware

import (
	"github.com/gin-gonic/gin"
	"sports_service/log"
	"sports_service/middleware/engineLog"
	"sports_service/middleware/header"
)

// 初始化部分中间件
func InitMiddleware(engine *gin.Engine, log log.ILogger, showColor bool) {
	// 跨域处理
	engine.Use(header.Cors)
	// 日志中间件
	engine.Use(engineLog.EngineLog(log, showColor))
}

package middleware

import (
  "github.com/gin-gonic/gin"
  "sports_service/server/log"
  "sports_service/server/middleware/engineLog"
  "sports_service/server/middleware/header"
)

// 初始化部分中间件
func InitMiddleware(engine *gin.Engine, log log.ILogger, showColor bool) {
  // 跨域处理
  engine.Use(header.Options)
  // 日志中间件
  engine.Use(engineLog.EngineLog(log, showColor))
}

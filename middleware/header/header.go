package header

import (
	"github.com/gin-gonic/gin"
  "net/http"
)

// 跨域处理
func Options(c *gin.Context) {
  if c.Request.Method == "OPTIONS" {
    c.Header("Access-Control-Allow-Origin", "*")
    c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE")
    c.Header("Access-Control-Allow-Headers", "x-requested-with")
    c.Header("Access-Control-Allow-Headers", "Cookie")
    c.Header("Access-Control-Allow-Headers", "Authorization")
    c.Header("Access-Control-Allow-Headers", "auth")
    c.Header("Access-Control-Allow-Headers", "Content-Type")
    // 允许请求带有验证信息
    c.Header("Access-Control-Allow-Credentials", "true")
    c.AbortWithStatus(http.StatusOK)
  } else {
    c.Header("Access-Control-Allow-Origin", "*")
    c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE")
    c.Header("Access-Control-Allow-Headers", "x-requested-with")
    c.Header("Access-Control-Allow-Headers", "Cookie")
    c.Header("Access-Control-Allow-Headers", "Authorization")
    c.Header("Access-Control-Allow-Headers", "auth")
    c.Header("Access-Control-Allow-Headers", "Content-Type")
    // 允许请求带有验证信息
    c.Header("Access-Control-Allow-Credentials", "true")
  }
}

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
    // 允许请求带有验证信息
    c.Header("Access-Control-Allow-Credentials", "true")
    c.Header("Access-Control-Allow-Headers", "Cookie, Authorization, x-requested-with, origin, Content-Type, auth")
    c.AbortWithStatus(http.StatusOK)
  } else {
    //c.Header("Access-Control-Allow-Origin", "http://fpv-web-qa.youzu.com")
    c.Header("Access-Control-Allow-Origin", "*")
    c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE")
    // 允许请求带有验证信息
    c.Header("Access-Control-Allow-Credentials", "true")
    c.Header("Access-Control-Allow-Headers", "Cookie, Authorization, x-requested-with, origin, Content-Type, auth")
  }
}

package header

import (
  "github.com/gin-gonic/gin"
)

// 跨域处理
func Cors(c *gin.Context) {
    //c.Header("Access-Control-Allow-Origin", "http://fpv-web-qa.youzu.com")
    c.Header("Access-Control-Allow-Origin", "*")
    c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE")
    // 允许请求带有验证信息
    c.Header("Access-Control-Allow-Credentials", "true")
    c.Header("Access-Control-Allow-Headers", "Cookie, Authorization, x-requested-with, origin, Content-Type, auth")
    //if c.Request.Method == "OPTIONS" {
    //  c.AbortWithStatus(http.StatusNoContent)
    //}
}

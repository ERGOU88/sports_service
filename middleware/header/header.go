package header

import (
	"github.com/gin-gonic/gin"
)

// 跨域处理
func Options(c *gin.Context) {
  //if c.Request.Method != "OPTIONS" {
  //  c.Next()
  //} else {
  //  c.Header("Access-Control-Allow-Origin", "*")
  //  c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
  //  c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept, app-id, signature, sq-id, timestamp")
  //  c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
  //  c.Header("Access-Control-Allow-Headers", "Cookie")
    //c.Header("Content-Type", "application/json")
    //c.AbortWithStatus(200)
  //}

  //if c.Request.Method != "OPTIONS" {
  // c.Next()
  //} else {
  // c.Header("Access-Control-Allow-Origin", "http://fpv-web-qa.youzu.com")
  // c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
  // c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept, app-id, signature, sq-id, timestamp")
  // c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
  // c.Header("Access-Control-Allow-Headers", "Cookie")
  // //c.Header("Content-Type", "application/json")
  // c.AbortWithStatus(200)
  //}
	c.Header("Access-Control-Allow-Origin", "*")
  c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "x-requested-with")
	c.Header("Access-Control-Allow-Headers", "Cookie")
	c.Header("Access-Control-Allow-Headers", "Authorization")
	c.Header("Access-Control-Allow-Headers", "auth")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	// 允许请求带有验证信息
	c.Header("Access-Control-Allow-Credentials", "true")
}

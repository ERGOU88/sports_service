package tencentCloud

import (
  "github.com/gin-gonic/gin"
)

// 腾讯云
func Router(engine *gin.Engine) {
  api := engine.Group("/api/v1")
  cloud := api.Group("/cloud")
  {
    // 获取腾讯cos通行证
    cloud.GET("/cos/access", CosTempAccess)
  }
}


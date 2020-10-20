package admin

import "github.com/gin-gonic/gin"

// 管理员模块后台路由
func Router(engine *gin.Engine) {
  api := engine.Group("/backend/v1")
  admin := api.Group("/admin")
  {
    // 注册后台用户（测试使用）
    admin.POST("/reg", RegAdminUser)
    // 后台管理员登陆
    admin.POST("/login", LoginByPassword)
  }
}

package admin

import (
  "github.com/gin-gonic/gin"
  "sports_service/server/middleware/jwt"
)

// 管理员模块后台路由
func Router(engine *gin.Engine) {
  api := engine.Group("/backend/v1")
  admin := api.Group("/admin")
  {
    // 注册后台用户（测试使用）
    admin.POST("/reg", RegAdminUser)
    // 域账号登陆
    admin.POST("/ad/login", AdLogin)
    // 管理员账号登录
    admin.POST("/login", AdminLogin)
    // 管理员上传
    admin.POST("/upload", jwt.JwtAuth(), UploadFile)
    // 添加管理员
    admin.POST("/add", jwt.JwtAuth(), AddAdmin)
    // 更新管理员
    admin.POST("/update", jwt.JwtAuth(), UpdateAdmin)
    // 删除管理员
    admin.DELETE("/del", jwt.JwtAuth(), DelAdmin)
  }
}

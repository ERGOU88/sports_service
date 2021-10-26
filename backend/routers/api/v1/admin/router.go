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
    // 禁用/解禁 管理员
    admin.POST("/forbid", jwt.JwtAuth(), ForbidAdmin)
    // 管理员详情
    admin.GET("/detail", AdminDetail)
    // 管理员列表
    admin.GET("/list", AdminList)
    // 角色列表
    admin.GET("/role/list", RoleList)
    // 添加角色
    admin.POST("/add/role", AddRole)
    // 设置角色可浏览的菜单
    admin.POST("/add/role/menu", AddRoleMenu)
    // 更新角色可浏览的菜单
    admin.POST("/update/role/menu", UpdateRoleMenu)
    // 获取角色可浏览的菜单
    admin.GET("/role/menu", GetRoleMenu)
    // 添加菜单
    admin.POST("/add/menu", AddMenu)
    // 更新菜单
    admin.POST("/update/menu", jwt.JwtAuth(), UpdateMenu)
    // 菜单详情
    admin.GET("/menu/detail", MenuDetail)
    // 菜单列表
    admin.GET("/menu/list", MenuList)
  }
}

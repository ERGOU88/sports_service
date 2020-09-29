package user

import "github.com/gin-gonic/gin"

func Router(engine *gin.Engine) {
	api := engine.Group("/backend/v1")
	user := api.Group("/user")
	{
		// 获取用户列表
		user.GET("/list", UserList)
		// 后台封禁用户
		user.POST("/forbid", ForbidUser)
		// 后台解封用户
		user.POST("/unforbid", UnForbidUser)
	}
}

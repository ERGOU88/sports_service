package attention

import (
	"github.com/gin-gonic/gin"
	"sports_service/middleware/sign"
	"sports_service/middleware/token"
)

// 关注模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	attention := api.Group("/attention")
	attention.Use(sign.CheckSign())
	{
		// 关注用户
		attention.POST("/user", token.TokenAuth(), AttentionUser)
		// 取消关注
		attention.POST("/cancel", token.TokenAuth(), CancelAttention)
		// 关注的用户列表
		attention.GET("/list", AttentionList)
		// 用户的粉丝列表
		attention.GET("/fans", FansList)
	}
}

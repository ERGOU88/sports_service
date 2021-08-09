package coach

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
)

// 私教模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	coach := api.Group("/coach")
	coach.Use(sign.CheckSign())
	{
		// 获取私教列表
		coach.GET("/list", CoachList)
	}
}

package coach

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
	"sports_service/server/middleware/token"
)

// 私教模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	coach := api.Group("/coach")
	coach.Use(sign.CheckSign())
	{
		// 获取私教列表
		coach.GET("/list", CoachList)
		// 私教详情
		coach.GET("/detail", CoachDetail)
        // 评价列表
		coach.GET("/evaluate/list", CoachEvaluate)
		// 评价配置
		coach.GET("/evaluate/conf", CoachEvaluateConf)
		// 发布对私教的评价
		coach.POST("/pub/evaluate", token.TokenAuth(), PubEvaluate)
	}
}

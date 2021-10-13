package contest

import "github.com/gin-gonic/gin"

func Router(engine *gin.Engine) {
	api := engine.Group("/backend/v1")
	contest := api.Group("/contest")
	{
		// 添加选手
		contest.POST("/add/player", AddPlayer)
	}
}

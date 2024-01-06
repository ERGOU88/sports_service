package stat

import (
	"github.com/gin-gonic/gin"
	"sports_service/middleware/jwt"
)

// 统计模块
func Router(engine *gin.Engine) {
	backend := engine.Group("/backend/v1")
	stat := backend.Group("/stat")
	stat.Use(jwt.JwtAuth())
	{
		// 首页统计数据
		stat.GET("/homepage", HomePageInfo)
		// 生态数据
		stat.GET("/ecological/info", EcologicalInfo)
	}
}

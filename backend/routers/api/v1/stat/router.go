package stat

import "github.com/gin-gonic/gin"

// 统计模块
func Router(engine *gin.Engine) {
	backend := engine.Group("/backend/v1")
	stat := backend.Group("/stat")
	{
		// 首页统计数据
		stat.GET("/homepage", HomePageInfo)
	}
}


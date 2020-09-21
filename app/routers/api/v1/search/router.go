package search

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/token"
)

// 搜索模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	search := api.Group("/search", token.TokenAuth())
	{
		// 搜索视频
		search.GET("/videos", VideoSearch)
		// 搜索用户
		search.GET("/users", UserSearch)
		// 综合搜索
		search.GET("/colligate", ColligateSearch)
		// 标签搜索
		search.GET("/label", LabelSearch)
		// 热门搜索
		search.GET("/hot", HotSearch)
	}
}


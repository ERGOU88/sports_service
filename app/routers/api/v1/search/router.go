package search

import (
  "github.com/gin-gonic/gin"
  "sports_service/server/middleware/token"
)

// 搜索模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	search := api.Group("/search")
	{
		// 搜索视频
		search.GET("/videos", VideoSearch)
		// 搜索用户
		search.GET("/users", UserSearch)
		// 综合搜索
		search.GET("/colligate", ColligateSearch)
		// 标签搜索
		search.GET("/label", LabelSearch)
		// 热门搜索 及 历史搜索记录
		search.GET("/hot", HotSearch)
		// 搜索关注的用户
		search.GET("/attention", token.TokenAuth(), AttentionSearch)
		// 搜索粉丝
		search.GET("/fans", token.TokenAuth(), FansSearch)
		// 清空用户搜索历史
		search.POST("/clean/history", token.TokenAuth(), CleanHistorySearch)
	}
}


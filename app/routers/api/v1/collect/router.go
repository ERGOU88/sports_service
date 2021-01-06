package collect

import (
  "github.com/gin-gonic/gin"
  "sports_service/server/middleware/sign"
  "sports_service/server/middleware/token"
)

// 收藏模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	collect := api.Group("/collect")
  collect.Use(sign.CheckSign())
	{
		// 收藏视频
		collect.POST("/video", token.TokenAuth(), CollectVideo)
		// 取消收藏
		collect.POST("/video/cancel", token.TokenAuth(), CancelCollect)
		// 用户收藏的视频列表
		collect.GET("/video/list", CollectVideoList)
    // 其他用户收藏的视频列表 todo:预留
		collect.GET("/other/list", OtherUserCollectVideoList)
		// 删除收藏
		collect.POST("/delete", token.TokenAuth(), DeleteCollect)
	}
}

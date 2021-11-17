package pub

import "github.com/gin-gonic/gin"

func Router(engine *gin.Engine) {
	backend := engine.Group("/backend/v1")
	pub := backend.Group("/pub")
	{
		// 发布视频
		pub.POST("/video", PubVideo)
		// 发布帖子
		pub.POST("/post", PubPost)
		// 腾讯云vod签名
		pub.GET("/upload/sign", UploadSign)
		// 发布资讯
		pub.POST("/information", PubInformation)
		// 首页板块
		pub.POST("/section", SectionInfo)
	}
}

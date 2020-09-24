package video

import "github.com/gin-gonic/gin"

// 视频点播模块后台路由
func Router(engine *gin.Engine) {
	api := engine.Group("/backend/v1")
	video := api.Group("/video")
	{
		// 视频审核 修改视频状态
		video.POST("/edit/status", EditVideoStatus)
		// 分页获取视频列表（已审核通过的）
		video.GET("/list", VideoList)
		// 编辑视频是否置顶
		video.POST("/edit/top", EditVideoTop)
		// 编辑视频是否推荐
		video.POST("/edit/recommend", EditVideoRecommend)
		// 审核/审核失败的视频列表
		video.GET("/review/list", VideoReviewList)
		// 获取视频标签列表
		video.GET("/label/list", VideoLabelList)
		// 添加视频标签
		video.POST("/add/label", AddVideoLabel)
		// 删除视频标签
		video.POST("/del/label", DelVideoLabel)
	}
}

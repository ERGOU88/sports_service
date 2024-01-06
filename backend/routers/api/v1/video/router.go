package video

import (
	"github.com/gin-gonic/gin"
	"sports_service/middleware/jwt"
	//"sports_service/middleware/jwt"
)

// 视频点播模块后台路由
func Router(engine *gin.Engine) {
	api := engine.Group("/backend/v1")
	video := api.Group("/video")
	video.Use(jwt.JwtAuth())
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
		// 编辑视频标签
		video.POST("/edit/label", EditVideoLabel)
		// 删除视频标签
		video.POST("/del/label", DelVideoLabel)
		// 添加分区配置
		video.POST("/add/subarea/conf", AddVideoSubareaConf)
		// 编辑分区配置
		video.POST("/edit/subarea/conf", EditVideoSubareaConf)
		// 删除分区配置
		video.POST("/del/subarea/conf", DelVideoSubareaConf)
		// 获取视频分区列表
		video.GET("/subarea/list", VideoSubareaList)
		// 批量编辑视频信息
		video.POST("/batch/edit", BatchEditVideoInfo)
		// 添加视频专辑
		video.POST("/add/album", AddAlbum)
		// 用户发布的专辑列表
		video.GET("/album/list", VideoAlbumList)
		// 首页板块信息
		video.GET("/section/info", SectionInfo)
	}
}

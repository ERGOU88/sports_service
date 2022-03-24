package course

import (
	"github.com/gin-gonic/gin"
)

// 课程模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	course := api.Group("/course")
	// todo: 先不校验签名
	//course.Use(sign.CheckSign())
	{
		// 课程分类配置
		course.GET("/category/config", CourseCategoryConfig)
		// 某分类下的课程列表
		course.GET("/category/list", CourseListByCategory)
		// 获取课程详情
		course.GET("/detail", CourseDetail)
		// 获取某一课时视频
		course.GET("/video", CourseVideo)
		// 客户端埋点 点击立即学习
		course.POST("/click/learn", ClickLearn)
		// 用户学习的课程记录
		course.GET("/user/learn/record", UserLearnRecord)
		// 客户端埋点 记录用户学习课程视频数据
		course.POST("/study/video/record", UserStudyVideoInfo)
	}
}

package course

import "github.com/gin-gonic/gin"

// 课程管理模块
func Router(engine *gin.Engine) {
	backend := engine.Group("/backend/v1")
	course := backend.Group("/course")
	{
		// 添加课程
		course.POST("/add", AddCourse)
		// 更新课程
		course.POST("/update", UpdateCourse)
		// 删除课程
		course.POST("/del", DelCourse)
		// 获取课程列表 or 搜索课程（名称、id）
		course.GET("/list", GetCourseList)
		// 获取上传签名（腾讯云）
		course.GET("/upload/sign", UploadSign)
		// 课程详情
		course.GET("/detail", CourseDetail)
		// 设置首页推荐（课程）
		course.POST("/set/homepage/recommend", SetHomePageRecommend)
		// 课程分类
		course.GET("/course/category", CourseCategory)
		// 添加课程分类
		course.POST("/add/course/category", AddCourseCategory)
		// 编辑课程分类
		course.POST("/edit/course/category", EditCourseCategory)
	}
}


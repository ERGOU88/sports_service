package course

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/backend/controller/course"
	"sports_service/global/backend/errdef"
	"sports_service/global/backend/log"
	"sports_service/models"
	"sports_service/models/medu"
	"sports_service/util"
)

// 添加课程
func AddCourse(c *gin.Context) {
	reply := errdef.New(c)
	args := new(medu.AddCourseArgs)
	if err := c.BindJSON(args); err != nil {
		log.Log.Errorf("course_trace: invalid params, args:%+v, err:%s", args, err)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := course.New(c)
	syscode := svc.AddCourse(args)
	reply.Response(http.StatusOK, syscode)
}

// 更新课程
func UpdateCourse(c *gin.Context) {
	reply := errdef.New(c)
	args := new(medu.UpdateCourseArgs)
	if err := c.BindJSON(args); err != nil {
		log.Log.Errorf("course_trace: invalid params, args:%+v, err:%s", args, err)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := course.New(c)
	syscode := svc.UpdateCourse(args)
	reply.Response(http.StatusOK, syscode)
}

// 获取上传签名（腾讯云）
func UploadSign(c *gin.Context) {
	reply := errdef.New(c)

	svc := course.New(c)
	syscode, sign, taskId := svc.GetUploadSign()
	reply.Data["sign"] = sign
	reply.Data["task_id"] = taskId

	log.Log.Errorf("#####taskId:%s", taskId)
	reply.Response(http.StatusOK, syscode)
}

// 删除课程
func DelCourse(c *gin.Context) {
	reply := errdef.New(c)
	arg := new(medu.DelCourseParam)
	if err := c.BindJSON(arg); err != nil {
		log.Log.Errorf("course_trace: invalid params, args:%+v, err:%s", arg, err)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := course.New(c)
	syscode := svc.DelCourse(arg.Id)
	reply.Response(http.StatusOK, syscode)
}

// 获取课程列表
func GetCourseList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	name := c.Query("name")

	svc := course.New(c)
	list, total := svc.GetCourseList(name, page, size)
	reply.Data["list"] = list
	reply.Data["total"] = total
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// 课程详情
func CourseDetail(c *gin.Context) {
	reply := errdef.New(c)
	id := c.Query("id")
	svc := course.New(c)
	detail := svc.GetCourseDetailInfo(id)
	reply.Data["detail"] = detail
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// 设置首页推荐(课程)
func SetHomePageRecommend(c *gin.Context) {
	reply := errdef.New(c)
	param := new(medu.SetHomePageRecommend)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("course_trace: invalid params, err:%s, param:%+v", err, param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := course.New(c)
	syscode := svc.SetHomePageRecommend(param)
	reply.Response(http.StatusOK, syscode)
}

// 课程分类列表
func CourseCategory(c *gin.Context) {
	reply := errdef.New(c)
	svc := course.New(c)
	list := svc.GetCourseCategoryList()
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

func AddCourseCategory(c *gin.Context) {
	reply := errdef.New(c)
	params := &models.CourseCategoryConfig{}
	if err := c.BindJSON(params); err != nil {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := course.New(c)
	reply.Response(http.StatusOK, svc.AddCourseCategory(params))
}

func EditCourseCategory(c *gin.Context) {
	reply := errdef.New(c)
	params := &models.CourseCategoryConfig{}
	if err := c.BindJSON(params); err != nil {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := course.New(c)
	reply.Response(http.StatusOK, svc.EditCourseCategory(params))
}

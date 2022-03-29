package course

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/course"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/medu"
	"sports_service/server/util"
)

// @Summary 获取某一分类下的课程列表 (ok)
// @Tags 课程模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Gid           header    string 	true  "调用/api/v1/client/init接口 服务端下发的gid"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	    header    string 	true  "版本" default(1.0.0)
// @Param   page	  	    query  	  string 	true  "页码 从1开始"
// @Param   size	  	    query  	  string 	true  "每页展示多少 最多50条"
// @Param   cate_id       query     string  true  "课程分类id"
// @Param   user_id       query     string  true  "用户id"
// @Success 200 {string}  json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string}  json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/course/category/list [get]
// 通过指定排序条件 获取某一分类下的课程列表
func CourseListByCategory(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	cateId := c.Query("cate_id")
	userId := c.Query("user_id")
	svc := course.New(c)
	// 通过指定分类 获取 课程列表
	list := svc.GetCourseListByCategory(userId, cateId, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 课程详情信息 (ok)
// @Tags 课程模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Gid           header    string 	true  "调用/api/v1/client/init接口 服务端下发的gid"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	    header    string 	true  "版本" default(1.0.0)
// @Param   course_id     query     string  true  "课程id"
// @Param   user_id       query     string  true  "用户id"
// @Success 200 {string}  json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string}  json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/course/detail [get]
// 课程详情信息
func CourseDetail(c *gin.Context) {
	reply := errdef.New(c)
	courseId := c.Query("course_id")
	userId := c.Query("user_id")
	
	svc := course.New(c)
	detail := svc.GetCourseDetailInfo(userId, courseId)
	reply.Data["detail"] = detail
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 获取某一课时视频 (ok)
// @Tags 课程模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Gid           header    string 	true  "调用/api/v1/client/init接口 服务端下发的gid"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	    header    string 	true  "版本" default(1.0.0)
// @Param   id            query     string  true  "课程视频id"
// @Param   course_id     query     string  true  "课程id"
// @Success 200 {string}  json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string}  json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/course/video [get]
// 课程视频
func CourseVideo(c *gin.Context) {
	reply := errdef.New(c)
	//userId, _ := c.Get(consts.USER_ID)
	userId := c.Query("user_id")
	courseId := c.Query("course_id")
	// 课程视频id
	id := c.Query("id")
	
	svc := course.New(c)
	syscode, videoInfo := svc.GetCourseVideoInfo(userId, courseId, id)
	if syscode != errdef.SUCCESS {
		reply.Response(http.StatusOK, syscode)
		return
	}
	
	reply.Data["video_info"] = videoInfo
	reply.Response(http.StatusOK, syscode)
}

// @Summary  客户端埋点 课程详情页 用户点击"立即学习"时调用 (ok)
// @Tags 课程模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Gid           header    string 	true  "调用/api/v1/client/init接口 服务端下发的gid"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   ClickLearnParam  body mcourse.ClickLearnParam true "用户点击立即学习 请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/course/click/learn [post]
// 客户端埋点 点击立即学习
func ClickLearn(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	param := new(medu.ClickLearnParam)
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("course_trace: click learn param err:%s, param:%+v", err, param)
		reply.Response(http.StatusBadRequest, errdef.SUCCESS)
		return
	}
	
	svc := course.New(c)
	syscode := svc.UserClickLearn(userId.(string), param.CourseId)
	reply.Response(http.StatusOK, syscode)
}

// @Summary 用户 "我的学习" (ok)
// @Tags 课程模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Gid           header    string 	true  "调用/api/v1/client/init接口 服务端下发的gid"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	    header    string 	true  "版本" default(1.0.0)
// @Param   Authorization header    string  true  "jwt"
// @Param   page          query     string  true  "页码 从1开始"
// @Param   size          query     string  true  "每页取多少"
// @Success 200 {string}  json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string}  json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/course/user/learn/record [get]
// 用户学习记录
func UserLearnRecord(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	
	svc := course.New(c)
	list := svc.GetUserLearnRecord(userId.(string), page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary  购买课程 (ok)
// @Tags 课程模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Gid           header    string 	true  "调用/api/v1/client/init接口 服务端下发的gid"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   PurchaseCourseParam  body mcourse.PurchaseCourseParam true "购买课程请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/course/purchase [post]
// 购买课程
//func PurchaseCourse(c *gin.Context) {
//	reply := errdef.New(c)
//	userId, _ := c.Get(consts.USER_ID)
//	params := new(medu.PurchaseCourseParam)
//	if err := c.BindJSON(params); err != nil {
//		log.Log.Errorf("course_trace: purchase course param err:%s", err)
//		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
//		return
//	}
//
//	svc := course.New(c)
//	syscode, info := svc.PurchaseCourse(userId.(string), params)
//	if syscode != errdef.SUCCESS {
//		reply.Response(http.StatusOK, syscode)
//		return
//	}
//
//	reply.Data["payment_info"] = info
//	reply.Response(http.StatusOK, syscode)
//}

// @Summary 客户端埋点 记录用户学习课程视频数据 (ok)
// @Tags 课程模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Gid           header    string 	true  "调用/api/v1/client/init接口 服务端下发的gid"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   RecordUserStudyInfo  body mcourse.RecordUserStudyInfo true "记录用户学习课程视频数据 请求参数"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/course/study/video/record [post]
// 客户端埋点 记录用户学习课程视频数据
func UserStudyVideoInfo(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	params := new(medu.RecordUserStudyInfo)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("course_trace: invalid params, params:%+v", params)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}
	
	svc := course.New(c)
	syscode := svc.RecordUserStudyVideoInfo(userId.(string), params)
	reply.Response(http.StatusOK, syscode)
}

func CourseCategoryConfig(c *gin.Context) {
	reply := errdef.New(c)
	
	svc := course.New(c)
	reply.Data["list"] = svc.GetCourseCategory()
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// 搜索课程
func CourseSearch(c *gin.Context) {
	reply := errdef.New(c)
	name := c.Query("name")
	userId := c.Query("user_id")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	
	svc := course.New(c)
	list := svc.CourseSearch(userId, name, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

func RecommendCourse(c *gin.Context) {
	reply := errdef.New(c)
	curCourseId := c.DefaultQuery("course_id", "0")
	userId := c.Query("user_id")
	
	svc := course.New(c)
	code, list := svc.RecommendCourse(userId, curCourseId)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

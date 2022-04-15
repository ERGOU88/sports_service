package course

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mconfigure"
	"sports_service/server/models/medu"
	"sports_service/server/models/morder"
	"sports_service/server/models/muser"
	cloud "sports_service/server/tools/tencentCloud"
	"sports_service/server/util"
	"sports_service/server/dao"
	"sports_service/server/global/backend/log"
	"strings"
	"time"
)

type CourseModule struct {
	context     *gin.Context
	engine      *xorm.Session
	course      *medu.EduModel
	user        *muser.UserModel
	configure   *mconfigure.ConfigModel
	order       *morder.OrderModel
}

func New(c *gin.Context) CourseModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	return CourseModule{
		context: c,
		course: medu.NewEduModel(socket),
		user: muser.NewUserModel(socket),
		configure: mconfigure.NewConfigModel(socket),
		order: morder.NewOrderModel(socket),
		engine: socket,
	}
}

// 添加课程
func (svc *CourseModule) AddCourse(args *medu.AddCourseArgs) int {
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("course_trace: session begin err:%s", err)
		return errdef.ERROR
	}
	
	//if args.SaiCoin == 0 {
	//  svc.engine.Rollback()
	//  return backend.COURSE_INVALID_SAI_COIN
	//}
	
	if args.EventEndTime != 0 || args.EventStartTime != 0 {
		if args.EventEndTime <= args.EventStartTime ||
			args.EventEndTime <= int(time.Now().Unix()) {
			log.Log.Errorf("course_trace: invalid activity time, endTime:%v", args.EventEndTime)
			svc.engine.Rollback()
			return errdef.ERROR
		}
	}
	
	videoNum := len(args.Videos)
	if videoNum  == 0 {
		log.Log.Errorf("course_trace: course video is empty, length:%d", len(args.Videos))
		svc.engine.Rollback()
		return errdef.ERROR
	}
	
	now := int(time.Now().Unix())
	
	svc.course.Course.AreasOfExpertise = args.AreasOfExpertise
	svc.course.Course.VipIsFree = args.VipIsFree
	svc.course.Course.IsTop = args.IsTop
	svc.course.Course.IsRecommend = args.IsRecommend
	svc.course.Course.Sortorder = args.Sortorder
	svc.course.Course.TeacherName = args.TeacherName
	svc.course.Course.TeacherPhoto = args.TeacherPhoto
	svc.course.Course.TeacherTitle = args.TeacherTitle
	svc.course.Course.PromotionPic = args.PromotionPic
	svc.course.Course.EventEndTime = args.EventEndTime
	svc.course.Course.EventStartTime = args.EventStartTime
	svc.course.Course.EventSaiCoin = args.EventSaiCoin
	svc.course.Course.SaiCoin = args.SaiCoin
	svc.course.Course.Describe = args.Describe
	svc.course.Course.Title = args.Title
	// 1 待发布 需等待腾讯云点播 事件回调 写入视频数据
	svc.course.Course.Status = 1
	svc.course.Course.CreateAt = now
	svc.course.Course.UpdateAt = now
	// 添加课程
	result, err := svc.course.AddCourse()
	if err != nil {
		log.Log.Errorf("course_trace: add course err:%s", err)
		svc.engine.Rollback()
		return errdef.ERROR
	}
	
	// 获取课程id
	courseId, err := result.LastInsertId()
	if err != nil {
		log.Log.Errorf("course_trace: get last insert id err:%s", err)
		svc.engine.Rollback()
		return errdef.ERROR
	}
	
	categorys := strings.Split(args.CourseCategory, ",")
	infos := make([]*models.CourseCategory, len(categorys))
	for index, id := range categorys {
		// 获取课程分类信息
		category := svc.course.GetCourseCategoryById(id)
		if category == nil {
			log.Log.Error("course_trace: course category not found")
			svc.engine.Rollback()
			return errdef.ERROR
		}
		
		info := new(models.CourseCategory)
		info.Name = category.Name
		info.CateId = category.Id
		info.CreateAt = now
		info.CourseId = fmt.Sprint(courseId)
		infos[index] = info
	}
	
	// 添加课程分类(多个)
	affected, err := svc.course.AddCourseCategorys(infos)
	if err != nil || affected != int64(len(categorys)) {
		log.Log.Errorf("course_trace: add course category err:%s, affected:%d", err, affected)
		svc.engine.Rollback()
		return errdef.ERROR
	}
	
	videos := make([]*models.CourseVideos, videoNum)
	for index, val := range args.Videos {
		video := new(models.CourseVideos)
		video.CreateAt = now
		video.UpdateAt = now
		video.Title = val.Title
		video.VideoAddr = val.VideoAddr
		video.FileOrder = val.FileOrder
		video.TxFileId = val.TxFileId
		// 1 待发布
		video.Status = 1
		video.CourseId = courseId
		// 记录任务id -> 腾讯文件id
		if err := svc.course.RecordUploadFileId(val.TxFileId, val.TaskId); err != nil {
			log.Log.Errorf("course_trace: record upload taskId err:%s, taskId:%d", err, val.TaskId)
			svc.engine.Rollback()
			return errdef.ERROR
		}
		
		videos[index] = video
	}
	
	// 添加课程视频
	affected, err = svc.course.AddCourseVideos(videos)
	if err != nil || int(affected) != videoNum {
		log.Log.Errorf("course_trace: add course videos err:%s, affected:%d, len:%d", err, affected, videoNum)
		svc.engine.Rollback()
		return errdef.ERROR
	}
	
	svc.engine.Commit()
	return errdef.SUCCESS
}

// 更新课程
func (svc *CourseModule) UpdateCourse(args *medu.UpdateCourseArgs) int {
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("course_trace: session begin err:%s", err)
		return errdef.ERROR
	}
	
	//if args.SaiCoin == 0 {
	//  svc.engine.Rollback()
	//  return backend.COURSE_INVALID_SAI_COIN
	//}
	
	if args.EventEndTime != 0 || args.EventStartTime != 0 {
		if args.EventEndTime <= args.EventStartTime ||
			args.EventEndTime <= int(time.Now().Unix()) {
			log.Log.Errorf("course_trace: invalid activity time, endTime:%v", args.EventEndTime)
			svc.engine.Rollback()
			return errdef.ERROR
		}
	}
	
	// 获取课程
	course := svc.course.GetCourseById(fmt.Sprint(args.Id))
	if course == nil {
		log.Log.Errorf("course_trace: course not found, courseId:%d", args.Id)
		svc.engine.Rollback()
		return errdef.ERROR
	}
	
	now := int(time.Now().Unix())
	svc.course.Course.AreasOfExpertise = args.AreasOfExpertise
	svc.course.Course.VipIsFree = args.VipIsFree
	svc.course.Course.IsTop = args.IsTop
	svc.course.Course.IsRecommend = args.IsRecommend
	svc.course.Course.Sortorder = args.Sortorder
	svc.course.Course.TeacherName = args.TeacherName
	svc.course.Course.TeacherPhoto = args.TeacherPhoto
	svc.course.Course.TeacherTitle = args.TeacherTitle
	svc.course.Course.PromotionPic = args.PromotionPic
	svc.course.Course.EventEndTime = args.EventEndTime
	svc.course.Course.EventStartTime = args.EventStartTime
	svc.course.Course.EventSaiCoin = args.EventSaiCoin
	svc.course.Course.SaiCoin = args.SaiCoin
	svc.course.Course.Describe = args.Describe
	svc.course.Course.Title = args.Title
	svc.course.Course.UpdateAt = now
	
	cols := "areas_of_expertise, vip_is_free, is_top, is_recommend, sortorder, teacher_name, teacher_photo, " +
		"teacher_title, promotion_pic, event_end_time, event_start_time, event_sai_coin, sai_coin, describe, title, update_at"
	condition := fmt.Sprintf("id=%d", course.Id)
	// 更新课程
	affected, err := svc.course.UpdateCourseInfo(condition, cols)
	if affected != 1 || err != nil {
		log.Log.Errorf("course_trace: update course err:%s", err)
		svc.engine.Rollback()
		return errdef.ERROR
	}
	
	// 获取课程的所有分类id
	ids := svc.course.GetCourseCategoryId(fmt.Sprint(course.Id))
	if ids != nil {
		all := strings.Join(ids, ",")
		// 不一致 表示分类有更新
		if strings.Compare(all, args.CourseCategory) != 0 {
			// 删除之前所属分类
			res, err := svc.course.DelCourseCategorys(all, fmt.Sprint(course.Id))
			if err != nil {
				log.Log.Errorf("course_trace: del course category err:%s", err)
				svc.engine.Rollback()
				return errdef.ERROR
			}
			
			affected, err = res.RowsAffected()
			if affected != int64(len(ids)) || err != nil {
				log.Log.Errorf("course_trace: del course category err:%s, affected:%d", err, affected)
				svc.engine.Rollback()
				return errdef.ERROR
			}
			
			ids := strings.Split(args.CourseCategory, ",")
			// 添加新的分类
			infos := make([]*models.CourseCategory, len(ids))
			// 课程分类
			for index, id := range ids {
				category := svc.course.GetCourseCategoryById(id)
				if category == nil {
					log.Log.Error("course_trace: course category not found")
					svc.engine.Rollback()
					return errdef.ERROR
				}
				
				info := new(models.CourseCategory)
				info.Name = category.Name
				info.CateId = category.Id
				info.CreateAt = now
				info.CourseId = fmt.Sprint(course.Id)
				infos[index] = info
			}
			
			// 添加课程分类(多个)
			affected, err := svc.course.AddCourseCategorys(infos)
			if err != nil || affected != int64(len(ids)) {
				log.Log.Errorf("course_trace: add course category err:%s, affected:%d", err, affected)
				svc.engine.Rollback()
				return errdef.ERROR
			}
		}
	}
	
	if len(args.Videos) > 0 {
		videos := make([]*models.CourseVideos, len(args.Videos))
		for index, val := range args.Videos {
			video := new(models.CourseVideos)
			video.CreateAt = now
			video.UpdateAt = now
			video.Title = val.Title
			video.VideoAddr = val.VideoAddr
			video.FileOrder = val.FileOrder
			video.TxFileId = val.TxFileId
			// 1 待发布
			video.Status = 1
			video.CourseId = course.Id
			// 记录任务id -> 腾讯文件id
			if err := svc.course.RecordUploadFileId(val.TxFileId, val.TaskId); err != nil {
				log.Log.Errorf("course_trace: record upload taskId err:%s, taskId:%d", err, val.TaskId)
				svc.engine.Rollback()
				return errdef.ERROR
			}
			
			videos[index] = video
		}
		
		// 添加新的课程视频
		affected, err = svc.course.AddCourseVideos(videos)
		if err != nil || int(affected) != len(args.Videos) {
			log.Log.Errorf("course_trace: add course videos err:%s, affected:%d, len:%d", err, affected, len(args.Videos))
			svc.engine.Rollback()
			return errdef.ERROR
		}
	}
	
	svc.engine.Commit()
	return errdef.SUCCESS
}

// 获取上传签名
func (svc *CourseModule) GetUploadSign() (int, string, int64) {
	// todo:暂时写死admin
	userId := "admin"
	client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
	taskId := util.GetXID()
	
	// 透传数据
	//sourceContext := &medu.SourceContext{
	//	EduUserId: userId,
	//	EduTaskId: taskId,
	//}
	
	//context, _ := util.JsonFast.Marshal(sourceContext)
	sign := client.GenerateSign(userId, consts.VOD_PROCEDURE_NAME, taskId)
	
	// 记录上传课程视频的后台用户
	if err := svc.course.RecordUploadUser(userId, taskId); err != nil {
		log.Log.Errorf("video_trace: record upload task id err:%s", err)
		return errdef.ERROR, "", 0
	}
	
	return errdef.SUCCESS, sign, taskId
}

// 删除课程
func (svc *CourseModule) DelCourse(id int64) int {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("course_trace: session begin err:%s", err)
		return errdef.ERROR
	}
	
	// 删除课程详情
	affected, err := svc.course.DelCourse(id)
	if affected != 1 || err != nil {
		svc.engine.Rollback()
		return errdef.ERROR
	}
	
	// 删除课程分类
	if _, err := svc.course.DelCourseCategoryById(id); err != nil {
		svc.engine.Rollback()
		return errdef.ERROR
	}
	
	// 删除该课程 用户们的学习记录
	if _, err := svc.course.DelCourseStudyRecord(id); err != nil {
		svc.engine.Rollback()
		return errdef.ERROR
	}
	
	// 删除该课程所有视频
	if _, err := svc.course.DelCourseVideosById(id); err != nil {
		svc.engine.Rollback()
		return errdef.ERROR
	}
	
	// 提交事务
	svc.engine.Commit()
	
	return errdef.SUCCESS
}

// 后台查找课程列表（课程标题、课程id）
func (svc *CourseModule) SearchCourseList(name string, offset, size int) []*models.CourseDetail {
	list := svc.course.SearchCourseList(name, offset, size)
	if len(list) == 0 {
		return []*models.CourseDetail{}
	}
	
	return list
}

// 获取搜索到的课程总数
func (svc *CourseModule) GetCourseTotalBySearch(name string) int64 {
	return svc.course.GetCourseTotalBySearch(name)
}

// 获取课程列表
func (svc *CourseModule) GetCourseList(name string, page, size int) ([]*medu.CourseInfoResp, int64) {
	var (
		list []*models.CourseDetail
		total int64
	)
	
	offset := (page - 1) * size
	if name != "" {
		list = svc.SearchCourseList(name, offset, size)
		total = svc.GetCourseTotalBySearch(name)
	} else {
		list = svc.course.GetAllCourse(offset, size)
		total = svc.GetCourseTotalCount()
	}
	
	if list == nil {
		return []*medu.CourseInfoResp{}, 0
	}
	
	res := make([]*medu.CourseInfoResp, len(list))
	for index, val := range list {
		info := new(medu.CourseInfoResp)
		info.Id = val.Id
		info.VipIsFree = val.VipIsFree
		info.Title = val.Title
		info.Sortorder = val.Sortorder
		info.SaiCoin = fmt.Sprintf("%.2f", float64(val.SaiCoin)/100)
		info.IsRecommend = val.IsRecommend
		info.IsFree = val.IsFree
		info.Status = val.Status
		info.CreateAt = val.CreateAt
		info.Describe = val.Describe
		info.Icon = val.Icon
		info.PromotionPic = val.PromotionPic
		info.TeacherName = val.TeacherName
		info.TeacherPhoto = val.TeacherPhoto
		info.TeacherTitle = val.TeacherTitle
		info.IsTop = val.IsTop
		info.EventSaiCoin = fmt.Sprintf("%.2f", float64(val.EventSaiCoin)/100)
		info.EventStartTime = val.EventStartTime
		info.EventEndTime = val.EventEndTime
		now := int(time.Now().Unix())
		// 是否在活动时间内 1为 已开启活动
		if now > val.EventStartTime && now < val.EventEndTime {
			info.HasActivity = 1
		}
		
		// 获取成功购买当前课程的用户数量(学习该课程的总人数)
		//info.StudyNum = svc.order.GetPurchaseCourseNum(fmt.Sprint(val.Id))
		// 学习该课程的用户总数
		info.StudyNum = svc.course.GetTotalStudyNumById(fmt.Sprint(val.Id))
		// 通过课程id 获取该课程下的视频总数
		videos := svc.course.GetCourseVideoSimpleInfoById(fmt.Sprint(val.Id))
		info.TotalChapters = len(videos)
		info.CategoryName = svc.course.GetCourseCategoryNameById(fmt.Sprint(val.Id))
		if info.CategoryName == nil {
			info.CategoryName = []*medu.CategoryInfo{}
		}
		
		res[index] = info
	}
	
	return res, total
}

// 获取课程总条目
func (svc *CourseModule) GetCourseTotalCount() int64 {
	return svc.course.GetCourseTotalCount()
}

// 获取课程详情信息
func (svc *CourseModule) GetCourseDetailInfo(courseId string) *medu.CourseDetail {
	// 查看课程是否存在
	course := svc.course.GetCourseById(courseId)
	if course == nil {
		return nil
	}
	
	info := new(medu.CourseDetail)
	info.Id = course.Id
	info.VipIsFree = course.VipIsFree
	info.Title = course.Title
	info.Sortorder = course.Sortorder
	info.SaiCoin = fmt.Sprintf("%.2f", float64(course.SaiCoin)/100)
	info.IsRecommend = course.IsRecommend
	info.Status = course.Status
	info.CreateAt = course.CreateAt
	info.Describe = course.Describe
	info.Icon = course.Icon
	info.PromotionPic = course.PromotionPic
	info.TeacherName = course.TeacherName
	info.TeacherPhoto = course.TeacherPhoto
	info.TeacherTitle = course.TeacherTitle
	info.IsTop = course.IsTop
	info.AreasOfExpertise = course.AreasOfExpertise
	info.EventStartTime = course.EventStartTime
	info.EventEndTime = course.EventEndTime
	info.EventSaiCoin = fmt.Sprintf("%.2f", float64(course.EventSaiCoin)/100)
	info.CategoryName = svc.course.GetCourseCategoryNameById(fmt.Sprint(courseId))
	if info.CategoryName == nil {
		info.CategoryName = []*medu.CategoryInfo{}
	}
	
	// 获取课程视频列表简单信息
	videos := svc.course.GetCourseVideoSimpleInfo(courseId)
	if len(videos) == 0 {
		videos = []*medu.CourseVideoSimple{}
	}
	
	info.CourseVideos = videos
	info.VideoTotal = len(videos)
	
	now := int(time.Now().Unix())
	// 是否在活动时间内 1为 已开启活动
	if now > course.EventStartTime && now < course.EventEndTime {
		info.HasActivity = 1
	}
	
	return info
}

// 设置首页推荐(课程)
func (svc *CourseModule) SetHomePageRecommend(param *medu.SetHomePageRecommend) int {
	svc.course.Course.IsRecommend = int(param.IsRecommend)
	svc.course.Course.UpdateAt = int(time.Now().Unix())
	svc.course.Course.Sortorder = int(param.Sortorder)
	condition := fmt.Sprintf("id=%d", param.Id)
	cols := "is_recommend, update_at, sortorder"
	affected, err := svc.course.UpdateCourseInfos(condition, cols)
	if affected != 1 || err != nil {
		log.Log.Errorf("course_trace: update course infos err:%s, affected:%d", err, affected)
		return errdef.ERROR
	}
	
	return errdef.SUCCESS
}

// 获取课程分类列表
func (svc *CourseModule) GetCourseCategoryList() []*models.CourseCategoryConfig {
	list := svc.course.GetCourseCategoryByLevel()
	if len(list) == 0 {
		return []*models.CourseCategoryConfig{}
	}
	
	return list
}

func (svc *CourseModule) AddCourseCategory(category *models.CourseCategoryConfig) int {
	svc.course.CourseCategory = category
	if err := svc.course.AddCourseCategory(); err != nil {
		return errdef.ERROR
	}
	
	return errdef.SUCCESS
}

func (svc *CourseModule) EditCourseCategory(category *models.CourseCategoryConfig) int {
	svc.course.CourseCategory = category
	if err := svc.course.UpdateCourseCategory(); err != nil {
		return errdef.ERROR
	}
	
	return errdef.SUCCESS
}

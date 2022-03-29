package course

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"math"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/models"
	"sports_service/server/models/mconfigure"
	"sports_service/server/models/medu"
	"sports_service/server/models/morder"
	"sports_service/server/models/muser"
	"sports_service/server/util"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"strconv"
	"time"
)

type CourseModule struct {
	context *gin.Context
	engine  *xorm.Session
	edu     *medu.EduModel
	conf    *mconfigure.ConfigModel
	order   *morder.OrderModel
	user    *muser.UserModel
}

func New(c *gin.Context) CourseModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	return CourseModule{
		context: c,
		engine:  socket,
		edu:     medu.NewEduModel(socket),
		conf:    mconfigure.NewConfigModel(socket),
		order:   morder.NewOrderModel(socket),
		user:    muser.NewUserModel(socket),
	}
}

// 获取首页推荐的课程 todo:默认先取4个 需求 有多少 取多少  100个足够了
func (svc *CourseModule) GetRecommendCourse(userId string) []*medu.CourseInfoResp {
	list := svc.edu.GetRecommendCourseList(0, 100)
	if len(list) == 0 {
		return []*medu.CourseInfoResp{}
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
		
		now := int(time.Now().Unix())
		// 是否在活动时间内 1为 已开启活动
		if now > val.EventStartTime && now < val.EventEndTime {
			info.HasActivity = consts.COURSE_HAS_ACTIVITY
			info.EventSaiCoin = fmt.Sprintf("%.2f", float64(val.EventSaiCoin)/100)
		}
		
		if userId != "" {
			// 当前用户是否已购买该课程
		}
		
		// 学习该课程的用户总数
		info.StudyNum = svc.edu.GetTotalStudyNumById(fmt.Sprint(val.Id))
		// 获取成功购买当前课程的用户数量
		//info.StudyNum = svc.order.GetPurchaseCourseNum(fmt.Sprint(val.Id))
		
		res[index] = info
	}
	
	return res
}

// 通过分类id获取课程列表
func (svc *CourseModule) GetCourseListByCategory(userId, cateId string, page, size int) []*medu.CourseInfoResp {
	var list []*models.CourseDetail
	offset := (page - 1) * size
	// 10表示限时免费
	if cateId == consts.LIMITED_FREE_COURSE {
		// 获取限时免费的课程
		list = svc.edu.GetLimitedFreeCourse(time.Now().Unix(), offset, size)
		if list == nil {
			log.Log.Error("course_trace: limited free edu not found")
			return []*medu.CourseInfoResp{}
		}
		
		log.Log.Errorf("list:%+v", list)
		
	} else {
		if info := svc.edu.GetCourseCategoryById(cateId); info == nil {
			log.Log.Errorf("course_trace: category not found, cateId:%s", cateId)
			return []*medu.CourseInfoResp{}
		}
		
		list = svc.edu.GetCourseByCategory(cateId, offset, size)
		if len(list) == 0 {
			log.Log.Errorf("list len: %d", len(list))
			return []*medu.CourseInfoResp{}
		}
		
		log.Log.Debugf("list:%+v", list)
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
		info.Describe = util.TrimHtml(val.Describe)
		info.Icon = val.Icon
		info.PromotionPic = val.PromotionPic
		info.TeacherName = val.TeacherName
		info.TeacherPhoto = val.TeacherPhoto
		info.TeacherTitle = val.TeacherTitle
		info.IsTop = val.IsTop
		
		now := int(time.Now().Unix())
		// 是否在活动时间内 1为 已开启活动
		if now > val.EventStartTime && now < val.EventEndTime {
			info.HasActivity = consts.COURSE_HAS_ACTIVITY
			info.EventSaiCoin = fmt.Sprintf("%.2f", float64(val.EventSaiCoin)/100)
			if val.EventSaiCoin == 0 {
				info.IsFree = consts.COURSE_IS_FREE
			}
		}
		
		if userId != "" {
			// 当前用户是否已购买该课程
		}
		
		// 获取学习当前课程的用户数量
		info.StudyNum = svc.edu.GetTotalStudyNumById(fmt.Sprint(val.Id))
		
		res[index] = info
	}
	
	return res
}

// 获取课程详情信息
func (svc *CourseModule) GetCourseDetailInfo(userId, courseId string) *medu.CourseDetailInfo {
	// 查看课程是否存在
	course := svc.edu.GetCourseById(courseId)
	if course == nil {
		return nil
	}
	
	info := new(medu.CourseDetailInfo)
	info.Id = course.Id
	info.VipIsFree = course.VipIsFree
	info.Title = course.Title
	info.Sortorder = course.Sortorder
	info.SaiCoin = fmt.Sprintf("%.2f", float64(course.SaiCoin)/100)
	info.IsRecommend = course.IsRecommend
	info.IsFree = course.IsFree
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
	
	// 获取课程视频列表简单信息
	videos := svc.edu.GetCourseVideoSimpleInfoById(courseId)
	if len(videos) == 0 {
		videos = []*medu.CourseVideoSimpleInfo{}
	}
	info.CourseVideos = videos
	info.VideoTotal = len(videos)
	
	now := int(time.Now().Unix())
	// 是否在活动时间内 1为 已开启活动
	if now > course.EventStartTime && now < course.EventEndTime {
		info.HasActivity = consts.COURSE_HAS_ACTIVITY
		info.EventSaiCoin = fmt.Sprintf("%.2f", float64(course.EventSaiCoin)/100)
		if course.EventSaiCoin == 0 {
			info.IsFree = consts.COURSE_IS_FREE
		}
	}
	
	if userId != "" {
		// 当前用户是否已购买该课程
		info.UserId = userId
		
		// 用户已购买 / 课程免费 / 用户是vip 且 课程vip免费 返回第一课时的数据
		if (info.HasPurchase == consts.COURSE_HAS_PURCHASED || info.IsFree == consts.COURSE_IS_FREE &&
			info.VipIsFree == consts.COURSE_IS_FREE_BY_VIP)  && info.VideoTotal > 0 {
			videoInfo := svc.edu.GetCourseVideoById(fmt.Sprint(info.CourseVideos[0].Id))
			//log.Log.Debugf("videoInfo:%+v", videoInfo)
			if videoInfo != nil {
				// 防盗链
				videoInfo.VideoAddr = svc.edu.AntiStealingLink(videoInfo.VideoAddr)
				// 满足播放条件 返回第一课时的视频数据（包含播放地址）
				info.FirstPeriod = videoInfo
				playRecord := svc.edu.GetUserPlayVideoRecordByFileId(userId, fmt.Sprint(info.CourseVideos[0].Id))
				if playRecord != nil {
					info.FirstPeriod.PlayDuration = playRecord.PlayDuration
				}
				
				if videoInfo.PlayInfo == nil {
					info.FirstPeriod.PlayInfo = []*medu.PlayInfo{}
				} else {
					for _, video := range info.FirstPeriod.PlayInfo {
						video.Url = svc.edu.AntiStealingLink(video.Url)
					}
				}
			}
			
		}
		
		// 用户是否学习过该课程
		if record := svc.edu.GetCourseStudyRecordByUser(userId, fmt.Sprint(course.Id)); record != nil {
			info.HasStudy = consts.HAS_STUDY_COURSE
		}
	}
	// 获取成功购买当前课程的用户数量
	info.StudyNum = svc.edu.GetTotalStudyNumById(fmt.Sprint(course.Id))
	
	return info
}

// 获取课程某一课时视频信息
func (svc *CourseModule) GetCourseVideoInfo(userId, courseId, id string) (int, *medu.CourseVideoInfo) {
	//if userId == "" {
	//	log.Log.Error("course_trace: need login")
	//	return errdef.USER_NO_LOGIN, nil
	//}
	
	// 查看课程是否存在
	course := svc.edu.GetCourseById(courseId)
	if course == nil {
		return errdef.EDU_COURSE_NOT_EXISTS, nil
	}
	
	// 查看课程视频是否存在
	videoInfo := svc.edu.GetCourseVideoById(id)
	if videoInfo == nil {
		return errdef.EDU_COURSE_VIDEO_NOT_EXISTS, nil
	}
	
	var (
		hasPurchase int32   // 用户是否已购买该课程
	)
	
	// 获取用户播放记录
	playRecord := svc.edu.GetUserPlayVideoRecordByFileId(userId, id)
	if playRecord != nil {
		// 最后播放时长
		videoInfo.PlayDuration = playRecord.PlayDuration
	}
	
	// 课程免费 / 该课程视频免费 / 用户为vip且课程针对vip免费 / 用户为vip且该课时针对vip免费 返回课程视频数据
	if course.IsFree == consts.COURSE_IS_FREE || videoInfo.IsFree == consts.COURSE_VIDEO_IS_FREE ||
		hasPurchase == consts.COURSE_HAS_PURCHASED &&
			videoInfo.VipIsFree == consts.COURSE_VIDEO_IS_FREE_BY_VIP && course.VipIsFree == consts.COURSE_IS_FREE_BY_VIP {
		videoInfo.VideoAddr = svc.edu.AntiStealingLink(videoInfo.VideoAddr)
		if videoInfo.PlayInfo != nil {
			for _, video := range videoInfo.PlayInfo {
				video.Url = svc.edu.AntiStealingLink(video.Url)
			}
		} else {
			videoInfo.PlayInfo = []*medu.PlayInfo{}
		}
		
		return errdef.SUCCESS, videoInfo
	}
	
	return errdef.EDU_COURSE_NOT_HAVE_ACCESS, nil
}

// 客户端埋点 用户点击 "立即学习"
func(svc *CourseModule) UserClickLearn(userId string, courseId int64) int {
	if userId == "" {
		log.Log.Errorf("course_trace: user not login, userId:%s", userId)
		return errdef.USER_NO_LOGIN
	}
	
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("course_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS
	}
	
	now := int(time.Now().Unix())
	// 用户是否学习过该课程 未学习过 添加记录
	if record := svc.edu.GetCourseStudyRecordByUser(userId, fmt.Sprint(courseId)); record == nil {
		svc.edu.CourseStudy.UserId = userId
		svc.edu.CourseStudy.CourseId = courseId
		svc.edu.CourseStudy.CreateAt = now
		svc.edu.CourseStudy.UpdateAt = now
		if err := svc.edu.SaveCourseStudyRecord(); err != nil {
			log.Log.Errorf("course_trace: save edu study record err:%s, userId:%s, courseId:%d", err, userId, courseId)
			return errdef.EDU_COURSE_SAVE_STUDY_RECORD
		}
		
	} else {
		// 已学习过 更新学习时间（错误无需返回）
		svc.edu.CourseStudy.UpdateAt = now
		// 购买状态不变
		svc.edu.CourseStudy.Status = record.Status
		if err := svc.edu.UpdateCourseStudyRecord(); err != nil {
			log.Log.Errorf("course_trace: update edu study record err:%s, userId:%s, courseId:%d", err, userId, courseId)
		}
	}
	
	return errdef.SUCCESS
}

// 获取用户学习课程的历史记录
func (svc *CourseModule) GetUserLearnRecord(userId string, page, size int) []*medu.UserCourseStudyInfo {
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("course_trace: user not found, userId:%s", userId)
		return []*medu.UserCourseStudyInfo{}
	}
	
	offset := (page - 1) * size
	now := int(time.Now().Unix())
	// 课程在活动期间 且 用户未购买过 且 活动价格为0 或 用户已购买过 在"我的学习"中展示
	list := svc.edu.FindStudyRecordsByUserId(userId, now, offset, size)
	if len(list) == 0 {
		log.Log.Error("course_trace: find edu study records by user fail")
		return []*medu.UserCourseStudyInfo{}
	}
	
	resp := make([]*medu.UserCourseStudyInfo, 0)
	for _, info := range list {
		var err error
		// 学习总时长
		info.StudyDuration, err = svc.edu.GetStudyTotalDurationByCourseId(userId, fmt.Sprint(info.Id))
		if err != nil {
			log.Log.Errorf("course_trace: get study total duration by edu id err:%s", err)
		}
		
		if info.StudyDuration > 0 {
			info.HasStudy = consts.HAS_STUDY_COURSE
		}
		
		// 通过课程id 获取 课程视频
		videos := svc.edu.GetCourseVideosById(fmt.Sprint(info.Id))
		// 未获取到视频 跳过
		if len(videos) == 0 {
			info.CourseVideos = []*medu.CourseVideoInfo{}
			info.VideoTotal = 0
			resp = append(resp, info)
			continue
		} else {
			info.CourseVideos = videos
			info.VideoTotal = len(videos)
		}
		
		// 获取用户最新的播放记录
		playRecord := svc.edu.GetUserPlayVideoRecord(userId, fmt.Sprint(info.Id))
		if playRecord != nil {
			// 通过视频id获取某一课时视频信息
			info.LastPlay = svc.edu.GetCourseVideoById(fmt.Sprint(playRecord.FileId))
			if info.LastPlay != nil {
				info.LastPlay.VideoAddr = svc.edu.AntiStealingLink(info.LastPlay.VideoAddr)
				if info.LastPlay.PlayInfo == nil {
					info.LastPlay.PlayInfo = []*medu.PlayInfo{}
				} else {
					for _, video := range info.LastPlay.PlayInfo {
						video.Url = svc.edu.AntiStealingLink(video.Url)
					}
				}
				
				// 当前已播放的时间节点
				info.LastPlay.PlayDuration = playRecord.PlayDuration
				// todo: 是否判断最近观看的视频 是否已学习完毕
				info.LastPlay.StudyStatus = consts.HAS_STUDY_COURSE
			}
		}
		
		for _, video := range info.CourseVideos {
			playRecord = svc.edu.GetUserPlayVideoRecordByFileId(userId, fmt.Sprint(video.Id))
			if playRecord != nil {
				video.StudyStatus = consts.HAS_STUDY_COURSE
				video.PlayDuration = playRecord.PlayDuration
			}
			
			video.VideoAddr = svc.edu.AntiStealingLink(video.VideoAddr)
			
			if video.PlayInfo == nil {
				video.PlayInfo = []*medu.PlayInfo{}
			} else {
				for _, playInfo := range video.PlayInfo {
					playInfo.Url = svc.edu.AntiStealingLink(playInfo.Url)
				}
			}
			
			// 课程总时长
			info.TotalDuration += int64(video.VideoDuration)/1000
			
			var maxProgress float64
			// 播放该视频的最长进度
			maxStr := svc.edu.GetMaxProgressOfPlayVideo(userId, fmt.Sprint(video.Id))
			if maxStr != "" {
				maxProgress, err = strconv.ParseFloat(maxStr, 64)
				if err != nil {
					log.Log.Errorf("course_trace: max progress parseFloat err:%s", err)
				}
			}
			
			// 获取用户总学习时长（针对某一课程的单一视频）
			duration, err := svc.edu.GetStudyTotalDurationByFileId(userId, fmt.Sprint(video.Id))
			if err != nil {
				log.Log.Errorf("course_trace: get study total duration by fileId err:%s", err)
			}
			
			// 记录已学习完毕的课程视频数
			var count int
			// 每个课程视频播放90%以上时间节点 且 观看时长 >= 视频总时长一半 则认为已学习
			if math.Ceil(maxProgress*100) > 90 && duration >= int64(float64(video.VideoDuration/1000/2) * 0.9) {
				video.StudyStatus = consts.END_STUDY_COURSE
				count++
			}
			
			// 每个视频占比 * 已学习的课程视频数 = 已学习百分比
			info.StudyRate += 100 / info.VideoTotal * count
			if info.StudyRate > 100 {
				info.StudyRate = 100
			}
		}
		
		resp = append(resp, info)
		
	}
	
	return resp
}

// 记录用户学习课程视频数据
func (svc *CourseModule) RecordUserStudyVideoInfo(userId string, param *medu.RecordUserStudyInfo) int {
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("course_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS
	}
	
	video := svc.edu.GetCourseVideoById(param.VideoId)
	if video == nil {
		log.Log.Errorf("course_trace: video not found, videoId:%s", param.VideoId)
		return errdef.EDU_COURSE_VIDEO_NOT_EXISTS
	}
	
	svc.edu.UserStudy.UserId = userId
	// 视频总时长存的毫秒
	svc.edu.UserStudy.TotalDuration = video.VideoDuration/1000
	svc.edu.UserStudy.StudyDuration = param.StudyDuration
	svc.edu.UserStudy.CourseId = video.CourseId
	svc.edu.UserStudy.FileId = video.Id
	svc.edu.UserStudy.PlayDuration = param.PlayDuration
	svc.edu.UserStudy.CreateAt = int(time.Now().Unix())
	svc.edu.UserStudy.CurProgress = fmt.Sprintf("%.2f", float64(param.PlayDuration) / float64(svc.edu.UserStudy.TotalDuration))
	// 添加用户学习课程视频记录
	if err := svc.edu.AddUserStudyRecord(); err != nil {
		log.Log.Errorf("course_trace: add user study record err:%s", err)
		return errdef.EDU_COURSE_RECORD_STUDY_VIDEO
	}
	
	return errdef.SUCCESS
}

// 课程查询
func (svc *CourseModule) CourseSearch(userId, name string, page, size int) []*medu.CourseInfoResp {
	if name == "" {
		log.Log.Errorf("search_trace: search consultant name can't empty, name:%s", name)
		return []*medu.CourseInfoResp{}
	}
	
	length := util.GetStrLen([]rune(name))
	if length > 20 {
		log.Log.Errorf("search_trace: invalid search name len, len:%s", length)
		return []*medu.CourseInfoResp{}
	}
	
	offset := (page - 1) * size
	list := svc.edu.SearchCourse(name, offset, size)
	if len(list) == 0 {
		return []*medu.CourseInfoResp{}
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
		
		now := int(time.Now().Unix())
		// 是否在活动时间内 1为 已开启活动
		if now > val.EventStartTime && now < val.EventEndTime {
			info.HasActivity = consts.COURSE_HAS_ACTIVITY
			info.EventSaiCoin = fmt.Sprintf("%.2f", float64(val.EventSaiCoin)/100)
		}
		
		if userId != "" {
			// 当前用户是否已购买该课程
			
		}
		// 获取成功购买当前课程的用户数量
		//info.StudyNum = svc.order.GetPurchaseCourseNum(fmt.Sprint(val.Id))
		// 学习该课程的用户总数
		info.StudyNum = svc.edu.GetTotalStudyNumById(fmt.Sprint(val.Id))
		
		res[index] = info
	}
	
	return res
}

func (svc *CourseModule) GetCourseCategory() []*models.CourseCategoryConfig {
	return svc.edu.GetCourseCategoryByLevel()
}

func (svc *CourseModule) RecommendCourse(userId, courseId string) (int, []*medu.CourseInfoResp) {
	list, err := svc.edu.GetRecommendCourse(courseId, 4)
	if err != nil {
		return errdef.ERROR, nil
	}
	
	if len(list) == 0 {
		return errdef.SUCCESS, []*medu.CourseInfoResp{}
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
		
		now := int(time.Now().Unix())
		// 是否在活动时间内 1为 已开启活动
		if now > val.EventStartTime && now < val.EventEndTime {
			info.HasActivity = consts.COURSE_HAS_ACTIVITY
			info.EventSaiCoin = fmt.Sprintf("%.2f", float64(val.EventSaiCoin)/100)
		}
		
		if userId != "" {
			// 当前用户是否已购买该课程
		}
		
		// 学习该课程的用户总数
		info.StudyNum = svc.edu.GetTotalStudyNumById(fmt.Sprint(val.Id))
		// 获取成功购买当前课程的用户数量
		//info.StudyNum = svc.order.GetPurchaseCourseNum(fmt.Sprint(val.Id))
		
		res[index] = info
	}
	
	return errdef.SUCCESS, res
}

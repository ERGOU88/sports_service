package medu

import (
	"github.com/go-xorm/xorm"
	"net/url"
	"sports_service/server/models"
	"sports_service/server/util"
	"strings"
	"time"
	"fmt"
)

// 教育模块
type EduModel struct {
	Engine                *xorm.Session
	Course                *models.CourseDetail
	Category              *models.CourseCategory
	Videos                *models.CourseVideos              // 课程视频
	CourseStudy           *models.CourseStudyRecord         // 用户课程学习记录
	UserStudy             *models.UserStudyRecord           // 用户观看视频的记录
	Events                *models.TencentCloudEvents
	CourseCategory        *models.CourseCategoryConfig
}

// 课程简单信息
type CourseInfoResp struct {
	CreateAt       int    `json:"create_at"`
	Describe       string `json:"describe"`            // 课程描述
	EventSaiCoin   string `json:"event_sai_coin"`      // 活动价格
	EventStartTime int    `json:"event_start_time,omitempty"`
	EventEndTime   int    `json:"event_end_time,omitempty"`
	Icon           string `json:"icon"`
	Id             int64  `json:"id"`
	IsFree         int    `json:"is_free"`
	IsRecommend    int    `json:"is_recommend"`
	IsTop          int    `json:"is_top"`
	PromotionPic   string `json:"promotion_pic"`
	SaiCoin        string `json:"sai_coin"`
	Sortorder      int    `json:"sortorder"`
	Status         int    `json:"status"`
	TeacherName    string `json:"teacher_name"`
	TeacherPhoto   string `json:"teacher_photo"`
	TeacherTitle   string `json:"teacher_title"`
	Title          string `json:"title"`
	VipIsFree      int    `json:"vip_is_free"`
	StudyNum       int64  `json:"study_num"`                    // 学习总人数
	HasActivity    int32  `json:"has_activity"`                 // 是否开启活动
	//ActivityPrice  string `json:"activity_price"`             // 活动价格
	HasPurchase    int32  `json:"has_purchase"`                 // 用户是否已购买
	TotalChapters  int    `json:"total_chapters,omitempty"`     // 总章节数(后台)
	
	CategoryName   []*CategoryInfo `json:"category_name,omitempty"`    // 分类名称
}

// 课程详情数据
type CourseDetailInfo struct {
	CreateAt       int    `json:"create_at"`
	Describe       string `json:"describe"`           // 课程描述
	EventSaiCoin   string `json:"event_sai_coin"`     // 活动价格
	Icon           string `json:"icon"`
	Id             int64  `json:"id"`                 // 课程id
	IsFree         int    `json:"is_free"`            // 是否免费 默认0 付费 1 免费
	IsRecommend    int    `json:"is_recommend"`
	IsTop          int    `json:"is_top"`
	PromotionPic   string `json:"promotion_pic"`
	SaiCoin        string `json:"sai_coin"`
	Sortorder      int    `json:"sortorder"`
	Status         int    `json:"status"`
	TeacherName    string `json:"teacher_name"`
	TeacherPhoto   string `json:"teacher_photo"`
	TeacherTitle   string `json:"teacher_title"`
	Title          string `json:"title"`
	VipIsFree      int    `json:"vip_is_free"`
	StudyNum       int64  `json:"study_num"`           // 学习总人数
	HasActivity    int32  `json:"has_activity"`        // 是否开启活动 1 有 0 无
	HasPurchase    int32  `json:"has_purchase"`        // 用户是否已购买
	IsVip          int32  `json:"is_vip"`              // 是否为vip
	UserId         string `json:"user_id"`             // 当前用户id
	VideoTotal     int    `json:"video_total"`         // 课时总数
	
	TeacherDescribe   string                   `json:"teacher_describe"`      // 老师简介
	FirstPeriod       *CourseVideoInfo         `json:"first_period"`          // 第一课时视频数据
	CourseVideos      []*CourseVideoSimpleInfo `json:"course_videos"`         // 课程视频信息
	HasStudy          int32                    `json:"has_study"`             // 用户是否学习过
	AreasOfExpertise  string                   `json:"areas_of_expertise"`    // 老师擅长领域
	
}

// 课程视频信息
type CourseVideoSimpleInfo struct {
	CourseId      int64  `json:"course_id"`           // 课程id
	FileOrder     int    `json:"file_order"`          // 文件序列
	Id            int64  `json:"id"`                  // 视频id
	IsFree        int    `json:"is_free"`             // 是否免费 0 收费 1 免费
	Title         string `json:"title"`               // 课程视频标题
	VideoDuration int    `json:"video_duration"`      // 课程视频时长（毫秒）
	VipIsFree     int    `json:"vip_is_free"`         // 会员是否免费 0 收费 1 免费
}

// 课程视频信息
type CourseVideoInfo struct {
	CourseId      int64       `json:"course_id"`
	Cover         string      `json:"cover"`
	Id            int64       `json:"id"`
	PlayInfo      []*PlayInfo `json:"play_info"`
	Size          int64       `json:"size"`
	Title         string      `json:"title"`
	VideoAddr     string      `json:"video_addr"`
	VideoDuration int         `json:"video_duration"`            // 视频时长
	PlayDuration  int         `json:"play_duration" xorm:"-"`    // 已播时长
	IsFree        int         `json:"is_free"`                   // 是否免费 0 收费 1 免费
	VipIsFree     int         `json:"vip_is_free"`               // 会员是否免费 0 收费 1 免费
	StudyStatus   int32       `json:"study_status" xorm:"-"`     // 0未学习 1 学习中 2 已学习完毕
	FileOrder     int         `json:"file_order"`                 // 文件序列
}

// 视频转码信息
type PlayInfo struct {
	Type     string   `json:"type" example:"1 流畅（FLU） 2 标清（SD）3 高清（HD）4 全高清（FHD）5 2K 6 4K"`    // 1 流畅（FLU） 2 标清（SD）3 高清（HD）4 全高清（FHD）5 2K 6 4K
	Url      string   `json:"url" example:"对应类型的视频地址"`
	Size     int64    `json:"size" example:"1000000000"`
	Duration int64    `json:"duration" example:"1000000000"`
}

// 用户点击 "立即学习" 请求参数
type ClickLearnParam struct {
	CourseId       int64     `binding:"required" json:"course_id"`     // 课程id
}

// 用户 "我的学习" 记录数据
type UserCourseStudyInfo struct {
	Id             int64  `json:"id"`                 // 课程id
	Describe       string `json:"describe"`           // 课程描述
	Icon           string `json:"icon"`
	PromotionPic   string `json:"promotion_pic"`
	TeacherName    string `json:"teacher_name"`
	TeacherPhoto   string `json:"teacher_photo"`
	TeacherTitle   string `json:"teacher_title"`
	Title          string `json:"title"`
	UserId         string `json:"user_id"`             // 当前用户id
	VideoTotal     int    `json:"video_total"`         // 课时总数
	TotalDuration  int64  `json:"total_duration"`      // 视频总时长
	StudyDuration  int64  `json:"study_duration"`      // 用户学习总时长
	StudyRate      int    `json:"study_rate"`          // 已学习的进度百分比
	
	
	LastPlay       *CourseVideoInfo         `json:"last_play"`     // 上一次播放的视频数据
	CourseVideos   []*CourseVideoInfo       `json:"course_videos"` // 课程视频信息
	HasStudy       int32                    `json:"has_study"`     // 用户是否学习过 0 开始学习 1 继续学习
}

// 购买课程请求参数
type PurchaseCourseParam struct {
	CourseId    string    `json:"course_id" binding:"required"`     // 课程id
}

// 购买课程信息
type PurchaseCourseInfo struct {
	Title         string    `json:"title"`         // 标题 例： 课程购买
	CourseId      int64     `json:"course_id"`     // 课程id
	CourseTitle   string    `json:"course_title"`  // 课程标题
	PromotionPic  string    `json:"promotion_pic"` // 课程宣传图
	TeacherName   string    `json:"teacher_name"`  // 老师
	TeacherTitle  string    `json:"teacher_title"` // 老师抬头
	TeacherPhoto  string    `json:"teacher_photo"` // 老师照片
	SaiCoin       string    `json:"sai_coin"`      // 课程价格（标价）
	RealSaiCoin   string    `json:"real_sai_coin"` // 需支付的真实价格
	Balance       string    `json:"balance"`       // 余额
	IsEnough      int32     `json:"is_enough"`     // 余额是否足够 0 不够 1 足够
	OrderId       string    `json:"order_id"`      // 订单id
}

// 记录用户学习课程视频数据（时长、进度等）
type RecordUserStudyInfo struct {
	StudyDuration    int     `json:"study_duration" binding:"required"`  // 学习时长 [秒]（真实停留时长）
	VideoId          string  `json:"video_id" binding:"required"`        // 视频文件id
	PlayDuration     int     `json:"play_duration"`                      // 视频已观看的时长 [秒]（用户可能直接拉取进度条）
}

func NewEduModel(engine *xorm.Session) *EduModel {
	return &EduModel{
		Engine: engine,
		Course: new(models.CourseDetail),
		Category: new(models.CourseCategory),
		Videos: new(models.CourseVideos),
		CourseStudy: new(models.CourseStudyRecord),
		UserStudy: new(models.UserStudyRecord),
		Events: new(models.TencentCloudEvents),
		CourseCategory: new(models.CourseCategoryConfig),
	}
}

// 通过id获取课程详情
func (m *EduModel) GetCourseById(courseId string) *models.CourseDetail {
	m.Course = new(models.CourseDetail)
	ok, err := m.Engine.Where("id=?", courseId).Get(m.Course)
	if !ok || err != nil {
		return nil
	}
	
	return m.Course
}

// 获取推荐的课程列表
func (m *EduModel) GetRecommendCourseList(offset, size int) []*models.CourseDetail {
	var list []*models.CourseDetail
	if err := m.Engine.Where("is_recommend=1").Desc("sortorder", "id").Limit(size, offset).Find(&list); err != nil {
		return nil
	}
	
	return list
}


// 获取限时免费的课程(活动开始时间 > 当前时间 且 活动结束时间 < 当前时间 且 活动价格为0)
func (m *EduModel) GetLimitedFreeCourse(tm int64, offset, size int) []*models.CourseDetail {
	var list []*models.CourseDetail
	if err := m.Engine.Where("event_start_time < ? AND event_end_time > ? AND event_sai_coin = 0", tm).Desc("id").Limit(size, offset).Find(&list); err != nil {
		return nil
	}
	
	return list
	
}

const (
	QUERY_COURSE_BY_CATEGORY = "SELECT course.*, category.* FROM `course_detail` AS course " +
		"LEFT JOIN `course_category` AS category ON course.`id`=category.`course_id` WHERE category.`cate_id` = ? " +
		"AND course.`status`=0 GROUP BY course.`id` ORDER BY course.is_top DESC, course.sortorder DESC, course.id DESC LIMIT ?, ?"
)
// 某一分类下的课程列表 [默认排序规则]
func (m *EduModel) GetCourseByCategory(cateId string, offset, size int) []*models.CourseDetail {
	var list []*models.CourseDetail
	if err := m.Engine.Table(&models.CourseDetail{}).SQL(QUERY_COURSE_BY_CATEGORY, cateId, offset, size).Find(&list); err != nil {
		return nil
	}
	
	return list
}

type CategoryInfo struct {
	Name      string    `json:"name"`    // 名称
	CateId    int32     `json:"cate_id"` // 分类id
}
// 通过课程id获取课程分类名称（多个）
func (m *EduModel) GetCourseCategoryNameById(courseId string) []*CategoryInfo {
	var info []*CategoryInfo
	if err := m.Engine.Table(&models.CourseCategory{}).Where("course_id=?", courseId).Cols("name, cate_id").Find(&info); err != nil {
		return nil
	}
	
	return info
}

// 课程标题 搜索课程
func (m *EduModel) SearchCourse(name string, offset, size int) []*models.CourseDetail {
	sql := "SELECT * FROM `course_detail` WHERE status=0 AND title like '%" + name + "%' ORDER BY id DESC LIMIT ?, ?"
	var list []*models.CourseDetail
	if err := m.Engine.SQL(sql, offset, size).Find(&list); err != nil {
		return nil
	}
	
	return list
}

// 通过课程id获取课程视频详情列表
func (m *EduModel) GetCourseVideosById(courseId string) []*CourseVideoInfo {
	var list []*CourseVideoInfo
	if err := m.Engine.Table(&models.CourseVideos{}).Where("course_id=? AND status=0", courseId).Asc("file_order").Find(&list); err != nil {
		return nil
	}
	
	return list
	
}

// 通过课程id获取课程视频简单信息列表
func (m *EduModel) GetCourseVideoSimpleInfoById(courseId string) []*CourseVideoSimpleInfo {
	var list []*CourseVideoSimpleInfo
	if err := m.Engine.Table("course_videos").Where("course_id=? AND status=0", courseId).Asc("file_order").Find(&list); err != nil {
		return nil
	}
	
	return list
}

// 通过课程id查询学习该课程的用户总数
func (m *EduModel) GetTotalStudyNumById(courseId string) int64 {
	total, err := m.Engine.Where("course_id=?", courseId).Count(&models.CourseStudyRecord{})
	if err != nil {
		return 0
	}
	
	return total
}

// 通过视频id获取某一课时视频信息
func (m *EduModel) GetCourseVideoById(videoId string) *CourseVideoInfo {
	info := new(CourseVideoInfo)
	ok, err := m.Engine.Table(&models.CourseVideos{}).Where("id=?", videoId).Get(info)
	if !ok || err != nil {
		return nil
	}
	
	return info
}

// 保存用户课程学习记录 todo: (用户购买过的课程 默认写入学习记录)
func (m *EduModel) SaveCourseStudyRecord() error {
	if _, err := m.Engine.InsertOne(m.CourseStudy); err != nil {
		return err
	}
	
	return nil
}

// 更新用户学习课程的时间 todo: (用户购买过的课程 默认写入学习记录)
func (m *EduModel) UpdateCourseStudyRecord() error {
	if _, err := m.Engine.ID(m.CourseStudy.Id).
		Cols("update_at, status").
		Update(m.CourseStudy); err != nil {
		return err
	}
	
	return nil
}

// 获取用户的课程学习记录（课程id及用户id获取）
func (m *EduModel) GetCourseStudyRecordByUser(userId, courseId string) *models.CourseStudyRecord {
	m.CourseStudy = new(models.CourseStudyRecord)
	ok, err := m.Engine.Where("user_id=? AND course_id=?", userId, courseId).Get(m.CourseStudy)
	if !ok || err != nil {
		return nil
	}
	
	return m.CourseStudy
}

// 分页 获取用户课程学习记录
func (m *EduModel) FindCourseStudyRecordsByUser(offset, size int, userId string) []*models.CourseStudyRecord {
	var list []*models.CourseStudyRecord
	if err := m.Engine.Where("user_id=?", userId).Desc("id").Limit(size, offset).Find(&list); err != nil {
		return nil
	}
	
	return list
}

const (
	FIND_STUDY_RECORDS = "SELECT csr.user_id, cd.id, cd.describe, cd.icon, cd.promotion_pic, cd.teacher_name, cd.teacher_photo, " +
		"cd.teacher_title, cd.title  FROM course_study_record AS csr LEFT JOIN course_detail AS cd " +
		"ON csr.course_id=cd.id WHERE csr.user_id=? AND cd.event_start_time < ? AND cd.event_end_time > ? AND csr.status=0 AND cd.event_sai_coin=0 " +
		"OR csr.user_id=? AND csr.status=1 ORDER BY csr.id DESC LIMIT ?, ?"
)
// 课程在活动期间 且 用户未购买过 且 活动价格为0 或 用户已购买过 在"我的学习"中展示
func (m *EduModel) FindStudyRecordsByUserId(userId string, now, offset, size int) []*UserCourseStudyInfo {
	var list []*UserCourseStudyInfo
	if err := m.Engine.SQL(FIND_STUDY_RECORDS, userId, now, now, userId, offset, size).Find(&list); err != nil {
		return nil
	}
	
	return list
}

// 添加用户学习课程视频数据记录 [记录播放时长、进度等]
func (m *EduModel) AddUserStudyRecord() error {
	if _, err := m.Engine.InsertOne(m.UserStudy); err != nil {
		return err
	}
	
	return nil
}

// 获取用户播放某课程视频的最新记录
func (m *EduModel) GetUserPlayVideoRecord(userId, courseId string) *models.UserStudyRecord {
	m.UserStudy = new(models.UserStudyRecord)
	ok, err := m.Engine.Where("user_id=? AND course_id=?", userId, courseId).Desc("id").Limit(1, 0).Get(m.UserStudy)
	if !ok || err != nil {
		return nil
	}
	
	return m.UserStudy
}

// 通过文件id获取用户播放课程视频的最新记录
func (m *EduModel) GetUserPlayVideoRecordByFileId(userId, fileId string) *models.UserStudyRecord {
	m.UserStudy = new(models.UserStudyRecord)
	ok, err := m.Engine.Where("user_id=? and file_id=?", userId, fileId).Desc("id").Limit(1, 0).Get(m.UserStudy)
	if !ok || err != nil {
		return nil
	}
	
	return m.UserStudy
}

// 获取用户总学习时长（针对某一个课程的所有视频）
func (m *EduModel) GetStudyTotalDurationByCourseId(userId, courseId string) (int64, error) {
	m.UserStudy = new(models.UserStudyRecord)
	return m.Engine.Where("user_id=? and course_id=?", userId, courseId).SumInt(m.UserStudy, "study_duration")
}

// 获取用户总学习时长（包含所有课程）
func (m *EduModel) GetStudyTotalDuration(userId string) (int64, error) {
	m.UserStudy = new(models.UserStudyRecord)
	return m.Engine.Where("user_id=?", userId).SumInt(m.UserStudy, "study_duration")
}

// 获取用户总学习时长（针对某一课程的单一视频）
func (m *EduModel) GetStudyTotalDurationByFileId(userId, fileId string) (int64, error) {
	m.UserStudy = new(models.UserStudyRecord)
	return m.Engine.Where("user_id=? and file_id=?", userId, fileId).SumInt(m.UserStudy, "study_duration")
}

// 获取用户观看某课程视频的最长进度
func (m *EduModel) GetMaxProgressOfPlayVideo(userId, fileId string) string {
	type TmpInfo struct {
		CurProgress    string   `json:"cur_progress"`  // 最大进度
	}
	
	tmp := new(TmpInfo)
	ok, err := m.Engine.Table(&models.UserStudyRecord{}).Where("user_id=? and file_id=?", userId, fileId).
		Cols("cur_progress").Desc("cur_progress").Limit(1).Get(tmp)
	if !ok || err != nil {
		return "0.00"
	}
	
	return tmp.CurProgress
}

// 获取用户学习过的课程总数
func (m *EduModel) GetStudyTotalCourse(userId string) int64 {
	total, err := m.Engine.Where("user_id=?", userId).Count(&models.CourseStudyRecord{})
	if err != nil {
		return 0
	}
	
	return total
}

const (
	LINK_KEY  = "MgbbCBVMhIwJCi9gnhhj"
	EXPIRE_TM = 3600 * 3
)
// 视频防盗链(有效时长3个小时)
func (m *EduModel) AntiStealingLink(videoUrl string) string {
	if videoUrl == "" {
		return ""
	}
	
	urls, err := url.Parse(videoUrl)
	if err != nil {
		return ""
	}
	
	path := urls.Path
	path = path[:strings.LastIndex(path, "/") + 1]
	rand := util.GenSecret(util.MIX_MODE, 10)
	tm := fmt.Sprintf("%x", time.Now().Unix() + EXPIRE_TM)
	signStr := fmt.Sprintf("%s%s%s%s", LINK_KEY, path, tm, rand)
	signStr = strings.Trim(signStr, " ")
	signStr = strings.Replace(signStr, "\n", "", -1)
	sign := util.Md5String(signStr)
	return fmt.Sprintf("%s?t=%s&us=%s&sign=%s", videoUrl, tm, rand, sign)
}



package medu

import (
	"database/sql"
	"fmt"
	"sports_service/dao"
	"sports_service/global/rdskey"
	"sports_service/models"
)

// 添加课程请求参数
type AddCourseArgs struct {
	Title            string `json:"title" binding:"required"` //
	Describe         string `json:"describe" binding:"required"`
	SaiCoin          int    `json:"sai_coin"`
	EventSaiCoin     int    `json:"event_sai_coin"`
	EventStartTime   int    `json:"event_start_time"`
	EventEndTime     int    `json:"event_end_time"`
	PromotionPic     string `json:"promotion_pic"`
	TeacherPhoto     string `json:"teacher_photo" binding:"required"`
	TeacherTitle     string `json:"teacher_title" binding:"required"`
	TeacherName      string `json:"teacher_name" binding:"required"`
	Sortorder        int    `json:"sortorder"`
	IsRecommend      int    `json:"is_recommend"`
	IsTop            int    `json:"is_top"`
	VipIsFree        int    `json:"vip_is_free"`
	AreasOfExpertise string `json:"areas_of_expertise" binding:"required"`
	CourseCategory   string `json:"course_category" binding:"required"`

	Videos []*CourseVideos `json:"videos" binding:"required"` // 课程视频列表
}

// 更新课程请求参数
type UpdateCourseArgs struct {
	Id               int64  `json:"id" binding:"required"`
	Title            string `json:"title" binding:"required"` //
	Describe         string `json:"describe" binding:"required"`
	SaiCoin          int    `json:"sai_coin"`
	EventSaiCoin     int    `json:"event_sai_coin"`
	EventStartTime   int    `json:"event_start_time"`
	EventEndTime     int    `json:"event_end_time"`
	PromotionPic     string `json:"promotion_pic"`
	TeacherPhoto     string `json:"teacher_photo" binding:"required"`
	TeacherTitle     string `json:"teacher_title" binding:"required"`
	TeacherName      string `json:"teacher_name" binding:"required"`
	Sortorder        int    `json:"sortorder"`
	IsRecommend      int    `json:"is_recommend"`
	IsTop            int    `json:"is_top"`
	VipIsFree        int    `json:"vip_is_free"`
	AreasOfExpertise string `json:"areas_of_expertise" binding:"required"`
	CourseCategory   string `json:"course_category" binding:"required"`

	Videos []*CourseVideos `json:"videos"` // 新加的课程视频
}

// 课程详情数据(管理后台)
type CourseDetail struct {
	CreateAt       int    `json:"create_at"`
	Describe       string `json:"describe"`       // 课程描述
	EventSaiCoin   string `json:"event_sai_coin"` // 活动价格
	Icon           string `json:"icon"`
	Id             int64  `json:"id"` // 课程id
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
	HasActivity    int32  `json:"has_activity"`     // 是否开启活动 1 有 0 无
	VideoTotal     int    `json:"video_total"`      // 课时总数
	EventStartTime int    `json:"event_start_time"` // 活动开始时间
	EventEndTime   int    `json:"event_end_time"`   // 活动结束时间

	CategoryName     []*CategoryInfo      `json:"category_name,omitempty"` // 分类名称
	TeacherDescribe  string               `json:"teacher_describe"`        // 老师简介
	CourseVideos     []*CourseVideoSimple `json:"course_videos"`           // 课程视频信息
	AreasOfExpertise string               `json:"areas_of_expertise"`      // 老师擅长领域
}

// 透传数据
type SourceContext struct {
	EduUserId string `json:"edu_user_id"` // 用户id
	EduTaskId int64  `json:"sai_task_id"` // 任务id
}

// 课程视频
type CourseVideos struct {
	TaskId    int64  `json:"task_id" binding:"required"`
	Title     string `json:"title" binding:"required"`
	VideoAddr string `json:"video_addr" binding:"required"`
	FileOrder int    `json:"file_order" binding:"required"`
	TxFileId  string `json:"tx_file_id" binding:"required"`
}

// 课程视频信息
type CourseVideoSimple struct {
	CourseId      int64  `json:"course_id"`      // 课程id
	FileOrder     int    `json:"file_order"`     // 文件序列
	Id            int64  `json:"id"`             // 视频id
	IsFree        int    `json:"is_free"`        // 是否免费 0 收费 1 免费
	Title         string `json:"title"`          // 课程视频标题
	VideoDuration int    `json:"video_duration"` // 课程视频时长（毫秒）
	VipIsFree     int    `json:"vip_is_free"`    // 会员是否免费 0 收费 1 免费
	VideoAddr     string `json:"video_addr"`     // 视频链接
}

// 设置首页推荐（课程）
type SetHomePageRecommend struct {
	Id          int64 `json:"id" binding:"required"` // 数据id
	IsRecommend int32 `json:"is_recommend"`          // 0 不推荐 1 推荐
	Sortorder   int64 `json:"sortorder"`             // 权重值
}

const (
	ADD_COURSE = "INSERT INTO `course_detail` (`title`, `describe`, `sai_coin`, `event_sai_coin`, `event_start_time`, " +
		"`event_end_time`, `promotion_pic`, `teacher_photo`, `teacher_title`, `teacher_name`, `sortorder`, `is_recommend`, " +
		"`is_top`, `vip_is_free`, `areas_of_expertise`, `status`, `create_at`, `update_at`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, " +
		"?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
)

// 添加课程
func (m *EduModel) AddCourse() (sql.Result, error) {
	return m.Engine.Exec(ADD_COURSE, m.Course.Title, m.Course.Describe, m.Course.SaiCoin, m.Course.EventSaiCoin, m.Course.EventStartTime,
		m.Course.EventEndTime, m.Course.PromotionPic, m.Course.TeacherPhoto, m.Course.TeacherTitle, m.Course.TeacherName,
		m.Course.Sortorder, m.Course.IsRecommend, m.Course.IsTop, m.Course.VipIsFree, m.Course.AreasOfExpertise, m.Course.Status,
		m.Course.CreateAt, m.Course.UpdateAt)
}

// 删除课程请求参数
type DelCourseParam struct {
	Id int64 `json:"id" binding:"required"` // 课程详情表 数据id
}

// 删除课程
func (m *EduModel) DelCourse(id int64) (int64, error) {
	return m.Engine.ID(id).Delete(&models.CourseDetail{})
}

// 通过课程id获取课程视频简单信息（管理后台）
func (m *EduModel) GetCourseVideoSimpleInfo(courseId string) []*CourseVideoSimple {
	var list []*CourseVideoSimple
	if err := m.Engine.Table("course_videos").Where("course_id=? AND status=0", courseId).Asc("file_order").Find(&list); err != nil {
		return nil
	}

	return list
}

// 删除课程所属分类
func (m *EduModel) DelCourseCategorys(ids, courseId string) (sql.Result, error) {
	return m.Engine.Exec(fmt.Sprintf("DELETE FROM `course_category` WHERE cate_id in(%s) AND course_id=%s", ids, courseId))
}

// 添加课程视频（多个）
func (m *EduModel) AddCourseVideos(videos []*models.CourseVideos) (int64, error) {
	return m.Engine.InsertMulti(videos)
}

// 添加课程分类（多个）
func (m *EduModel) AddCourseCategorys(infos []*models.CourseCategory) (int64, error) {
	return m.Engine.InsertMulti(infos)
}

// 记录任务id -> 腾讯文件id
func (m *EduModel) RecordUploadFileId(fileId string, taskId int64) error {
	rds := dao.NewRedisDao()
	return rds.SETEX(rdskey.MakeKey(rdskey.COURSE_VIDEO_UPLOAD_TASK, taskId), rdskey.KEY_EXPIRE_DAY*3, fileId)
}

// 通过任务id 获取 腾讯文件id
func (m *EduModel) GetUploadFileIdByTaskId(taskId int64) (string, error) {
	rds := dao.NewRedisDao()
	return rds.Get(rdskey.MakeKey(rdskey.COURSE_VIDEO_UPLOAD_TASK, taskId))
}

// 记录任务id -> 用户 (1天过期)
func (m *EduModel) RecordUploadUser(userId string, taskId int64) error {
	rds := dao.NewRedisDao()
	return rds.SETEX(rdskey.MakeKey(rdskey.COURSE_VIDEO_UPLOAD_USER, taskId), rdskey.KEY_EXPIRE_DAY, userId)
}

// 通过任务id 获取 用户
func (m *EduModel) GetUploadUserByTaskId(taskId int64) (string, error) {
	rds := dao.NewRedisDao()
	return rds.Get(rdskey.MakeKey(rdskey.COURSE_VIDEO_UPLOAD_USER, taskId))
}

// 删除记录腾讯文件id的key
func (m *EduModel) DelRecordFileIdKey(taskId int64) (int, error) {
	rds := dao.NewRedisDao()
	return rds.Del(rdskey.MakeKey(rdskey.COURSE_VIDEO_UPLOAD_TASK, taskId))
}

// 通过腾讯云返回的文件id查询课程视频
func (m *EduModel) GetCourseVideoByFileId(fileId string) *models.CourseVideos {
	m.Videos = new(models.CourseVideos)
	ok, err := m.Engine.Where("tx_file_id=?", fileId).Get(m.Videos)
	if !ok || err != nil {
		return nil
	}

	return m.Videos
}

// 更新课程视频转码数据
func (m *EduModel) UpdateCourseVideoPlayInfo(id string) error {
	if _, err := m.Engine.Where("id=?", id).Cols("play_info", "status").Update(m.Videos); err != nil {
		return err
	}

	return nil
}

// 更新课程视频数据（时长、封面、大小等）
func (m *EduModel) UpdateCourseVideoInfo(id string) (int64, error) {
	return m.Engine.Where("id=?", id).Cols("cover", "video_duration", "size", "video_width",
		"video_height", "update_at", "status").Update(m.Videos)
}

// 记录腾讯事件回调信息
func (m *EduModel) RecordTencentEvent() (int64, error) {
	return m.Engine.InsertOne(m.Events)
}

// 获取所有课程（分页）
func (m *EduModel) GetAllCourse(offset, size int) []*models.CourseDetail {
	var list []*models.CourseDetail
	if err := m.Engine.Desc("is_recommend", "sortorder", "id").Limit(size, offset).Find(&list); err != nil {
		return nil
	}

	return list
}

// 更新课程信息
func (m *EduModel) UpdateCourseInfo(condition, cols string) (int64, error) {
	return m.Engine.Where(condition).Cols(cols).Update(m.Course)
}

// 获取课程总数
func (m *EduModel) GetCourseTotalCount() int64 {
	count, err := m.Engine.Count(&models.CourseDetail{})
	if err != nil {
		return 0
	}

	return count
}

// 后台搜索课程列表
func (m *EduModel) SearchCourseList(name string, offset, size int) []*models.CourseDetail {
	sql := "SELECT * FROM `course_detail` WHERE title like '%" + name + "%' OR id like '%" + name + "%' ORDER BY is_recommend DESC, sortorder DESC, id DESC LIMIT ?, ?"
	var list []*models.CourseDetail
	if err := m.Engine.SQL(sql, offset, size).Find(&list); err != nil {
		return nil
	}

	return list
}

// 获取搜索到的课程总数
func (m *EduModel) GetCourseTotalBySearch(name string) int64 {
	sql := "SELECT count(1) as total FROM `course_detail` WHERE title like '%" + name + "%' OR id like '%" + name + "%'"
	type tmp struct {
		Total int64 `json:"total"`
	}

	info := new(tmp)
	ok, err := m.Engine.Table(&models.CourseDetail{}).SQL(sql).Get(info)
	if !ok || err != nil {
		return 0
	}

	return info.Total
}

// 更新课程信息
func (m *EduModel) UpdateCourseInfos(condition, cols string) (int64, error) {
	return m.Engine.Where(condition).Cols(cols).Update(m.Course)
}

// 通过课程id获取课程分类id（多个）
func (m *EduModel) GetCourseCategoryId(courseId string) []string {
	var ids []string
	if err := m.Engine.Table(&models.CourseCategory{}).Where("course_id=?", courseId).Cols("cate_id").Find(&ids); err != nil {
		return nil
	}

	return ids
}

// 通过课程id 删除所有分类
func (m *EduModel) DelCourseCategoryById(courseId int64) (int64, error) {
	return m.Engine.Where("course_id=?", courseId).Delete(&models.CourseCategory{})
}

// 通过课程id 删除课程所有视频
func (m *EduModel) DelCourseVideosById(courseId int64) (int64, error) {
	return m.Engine.Where("course_id=?", courseId).Delete(&models.CourseVideos{})
}

// 通过课程id 删除用户课程学习记录
func (m *EduModel) DelCourseStudyRecord(courseId int64) (int64, error) {
	return m.Engine.Where("course_id=?", courseId).Delete(&models.CourseStudyRecord{})
}

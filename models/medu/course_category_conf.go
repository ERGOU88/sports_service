package medu

import (
	"sports_service/server/models"
	"sync"
	"fmt"
	"reflect"
	"sports_service/server/global/app/log"
)

// 课程分类
type CourseCategory struct {
	CreateAt  int               `json:"create_at"`
	Icon      string            `json:"icon"`
	Id        int               `json:"id"`
	Name      string            `json:"name"`
	Pid       int               `json:"pid"`
	Sortorder int               `json:"sortorder"`
	Status    int               `json:"status"`
	UpdateAt  int               `json:"update_at"`
	Child     []*CourseCategory `json:"child" xorm:"-"` // 子类信息
}

// 添加课程分类请求参数
type AddCourseCategoryParam struct {
	Name         string    `json:"name"`           // 分类名称
	Icon         string    `json:"icon"`           // icon
	Sortorder    int       `json:"sortorder"`      // 权重
}

// 删除课程分类请求参数
type DelCourseCategoryParam struct {
	Id      string    `json:"id"`                  // 分类id
}

var courseCategory []*CourseCategory

// id -> CourseCategory
var courseCategoryMp map[string]*CourseCategory

var courseMutex sync.Mutex

func init() {
	courseCategoryMp = make(map[string]*CourseCategory)
}

// 删除课程分类
func (m *EduModel) DelCourseCategory(id string) error {
	if _, err := m.Engine.Where("id=?", id).Delete(&models.CourseCategoryConfig{}); err != nil {
		return err
	}
	
	return nil
}

// 添加课程分类
func (m *EduModel) AddCourseCategory() error {
	if _, err := m.Engine.Insert(m.CourseCategory); err != nil {
		return err
	}
	
	return nil
}

func (m *EduModel) UpdateCourseCategory() error {
	if _, err := m.Engine.Where("id=?", m.CourseCategory.Id).Update(m.CourseCategory); err != nil {
		return err
	}
	
	return nil
}

// 通过分类id获取课程分类信息
func (m *EduModel) GetCourseCategoryById(id string) *models.CourseCategoryConfig {
	m.CourseCategory = new(models.CourseCategoryConfig)
	ok, err := m.Engine.Where("id=?", id).Get(m.CourseCategory)
	if !ok || err != nil {
		log.Log.Errorf("configure_trace: get course category info by id err:%s", err)
		return nil
	}
	
	return m.CourseCategory
}

// 分类是否存在
func (m *EduModel) IsExistsCourseCategory(id string) bool {
	if _, ok := courseCategoryMp[id]; !ok {
		return false
	}
	
	return true
}

// 清空课程分类map
func (m *EduModel) DelAllCourseCategoryByMem() {
	courseCategoryMp = make(map[string]*CourseCategory)
}

// 通过id 获取课程分类信息
func (m *EduModel) GetCourseCategoryByMem(id string) *CourseCategory {
	category, ok := courseCategoryMp[id]
	if !ok {
		return nil
	}
	
	return category
}

// 通过id 获取课程分类名称
func (m *EduModel) GetCourseCategoryNameByMem(id string) string {
	category, ok := courseCategoryMp[id]
	if ok {
		return category.Name
	}
	
	return ""
}

// 删除课程分类信息（内存）
func (m *EduModel) DelCourseCategoryByMem(id string) {
	courseMutex.Lock()
	defer courseMutex.Unlock()
	if _, ok := courseCategoryMp[id]; ok {
		delete(courseCategoryMp, id)
	}
}

// 修改频率低的数据 直接清理内存数据 下次从数据库load内存
func (m *EduModel) CleanCourseCategoryByMem() {
	courseCategory = nil
}

// 添加课程分类信息（内存）
func (m *EduModel) AddCourseCategoryByMem() {
	courseMutex.Lock()
	defer courseMutex.Unlock()
	info := &CourseCategory{
		CreateAt: m.CourseCategory.CreateAt,
		Icon: m.CourseCategory.Icon,
		Id: m.CourseCategory.Id,
		Name: m.CourseCategory.Name,
		Pid: m.CourseCategory.Pid,
		Sortorder: m.CourseCategory.Sortorder,
		Status: m.CourseCategory.Status,
		UpdateAt: m.CourseCategory.UpdateAt,
	}
	
	courseCategoryMp[fmt.Sprint(m.CourseCategory.Id)] = info
}

// 从内存读取课程分类 （第一次请求 内存没有 则从数据库load到内存）
func (m *EduModel) GetCourseCategoryList() []*CourseCategory {
	if len(courseCategory) == 0 {
		log.Log.Debugf("load mysql")
		var err error
		err, courseCategory = m.LoadCourseCategoryByDb()
		if err != nil {
			return []*CourseCategory{}
		}
	}
	
	log.Log.Debugf("load mem")
	return courseCategory
}

// 获取课程二级分类列表
func (m *EduModel) GetCourseCategoryByLevel() []*models.CourseCategoryConfig {
	var list []*models.CourseCategoryConfig
	if err := m.Engine.Where("pid != 0").Desc("sortorder").Find(&list); err != nil {
		return nil
	}
	
	return list
}

const (
	QUERY_PARENT_COURSE_CATEGORY = "SELECT * FROM `course_category_config` WHERE `pid` = 0 ORDER BY sortorder DESC"
)

// 从数据库获取分类信息
func (m *EduModel) LoadCourseCategoryByDb() (error, []*CourseCategory) {
	// 定义指针切片用来存储所有分类
	var info []*CourseCategory
	// 定义指针切片返回控制器
	var res []*CourseCategory
	
	// 找出所有1级类别
	if err := m.Engine.Table(&models.CourseCategoryConfig{}).SQL(QUERY_PARENT_COURSE_CATEGORY).Find(&info); err != nil {
		log.Log.Errorf("configure_trace: get consultant category info by db err:%s", err)
		return err, nil
	}
	
	// 判断是否存在数据 存在 则进行树状图重构
	if reflect.ValueOf(info).IsValid() {
		res = m.courseTree(info)
	}
	
	return nil, res
}

const (
	QUERY_SUB_COURSE_CATEGORY = "SELECT * FROM `course_category_config` WHERE `pid`=?" // 查询某个父类下的所有子类别
)

// 通过父级分类id查询下属的子分类
func (m *EduModel) FindSubCourseCategoryByPid(category *CourseCategory) ([]*CourseCategory, error) {
	var child []*CourseCategory
	if err := m.Engine.Table(&models.CourseCategoryConfig{}).SQL(QUERY_SUB_COURSE_CATEGORY, category.Id).Find(&child); err != nil {
		log.Log.Errorf("configure_trace: get child category info err:%s", err)
		return []*CourseCategory{}, err
	}
	
	return child, nil
}

// 树状图重构
func (m *EduModel) courseTree(info []*CourseCategory) []*CourseCategory {
	if reflect.ValueOf(info).IsValid() {
		// 循环所有1级标签
		for k, v := range info {
			courseCategoryMp[fmt.Sprint(v.Id)] = v
			log.Log.Debugf("courseCategoryMp:%+v", courseCategoryMp)
			// 查询所有一级分类下的所有子分类
			child, err := m.FindSubCourseCategoryByPid(v)
			if err != nil || len(child) == 0 {
				continue
			}
			
			// 将子类别的数据循环赋值
			for k2, _ := range child {
				info[k].Child = append(info[k].Child, child[k2])
			}
			
			// 将刚刚查询出来的子类别进行递归 查询出三级,四级...子类
			m.courseTree(child)
		}
		
	}
	
	return info
}



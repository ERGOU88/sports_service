package mlabel

import (
	"sports_service/server/dao"
	"sports_service/server/global/app/log"
	"sports_service/server/models"
	"github.com/go-xorm/xorm"
	"reflect"
	"sync"
)

type LabelModel struct {
	Engine         *xorm.Session
	VideoLabels    *models.VideoLabelConfig
	PostLabels     *models.PostLabelConfig
}

// 标签列表信息
type VideoLabelList struct {
	CreateAt  int               `json:"create_at" `
	Icon      string            `json:"icon"`
	LabelId   int               `json:"label_id"`
	LabelName string            `json:"label_name"`
	Pid       int               `json:"pid"`
	Sortorder int               `json:"sortorder"`
	Status    int               `json:"status"`
	UpdateAt  int               `json:"update_at"`
	Child     []*VideoLabelList `json:"child" xorm:"-"`               // 子类标签信息
}

var videoLabels []*VideoLabelList

// labelId -> labelName
var labelMp map[int]string

var mutex sync.Mutex

func init() {
	labelMp = make(map[int]string)
}

// 实栗
func NewLabelModel(engine *xorm.Session) *LabelModel {
	return &LabelModel{
		Engine: engine,
		VideoLabels: new(models.VideoLabelConfig),
		PostLabels: new(models.PostLabelConfig),
	}
}

// 通过labeiid获取视频标签信息
func (m *LabelModel) GetVideoLabelInfoById(labelId string) *models.VideoLabelConfig {
	m.VideoLabels = new(models.VideoLabelConfig)
	ok, err := m.Engine.Where("label_id=?", labelId).Get(m.VideoLabels)
	if !ok || err != nil {
		log.Log.Errorf("label_trace: get label info by id err:%s", err)
		return nil
	}

	return m.VideoLabels
}

// 标签是否存在
func (m *LabelModel) IsExistsLabel(labelId int) bool {
	if _, ok := labelMp[labelId]; !ok {
		return false
	}

	return true
}

// 清空标签map
func (m *LabelModel) DelAllLabel() {
	labelMp = make(map[int]string)
}

// 更新标签信息
func (m *LabelModel) UpdateLabelInfo(labelId int) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(labelMp, labelId)
}

// 从内存读取视频标签 （第一次请求 内存没有 则从数据库load到内存）
func (m *LabelModel) GetVideoLabelList() []*VideoLabelList {
	if len(videoLabels) == 0 {
		var err error
		err, videoLabels = m.LoadLabelsInfoByDb()
		if err != nil {
			return nil
		}
	}

	return videoLabels
}

const (
	QUERY_PARENT_LABELS = "SELECT * FROM `video_label_config` WHERE `pid` = 0 ORDER BY sortorder DESC"
)

// 从数据库获取标签信息
func (m *LabelModel) LoadLabelsInfoByDb() (error, []*VideoLabelList) {
	// 定义指针切片用来存储所有标签
	var info []*VideoLabelList
	// 定义指针切片返回控制器
	var res []*VideoLabelList

	// 找出所有1级类别
	if err := dao.Engine.SQL(QUERY_PARENT_LABELS).Find(&info); err != nil {
		log.Log.Errorf("labels_trace: get labels info err:%s", err)
		return err, nil
	}

	// 判断是否存在数据 存在 则进行树状图重构
	if reflect.ValueOf(info).IsValid() {
		res = m.tree(info)
	}

	return nil, res
}

const (
	QUERY_SUB_LABELS = "SELECT * FROM `video_label_config` WHERE `pid`=?" // 查询某个父类下的所有子类别
)

// 通过父类标签id查询下属的子标签
func (m *LabelModel) FindSubLabelsByPid(label *VideoLabelList) ([]*VideoLabelList, error) {
	var child []*VideoLabelList
	if err := dao.Engine.SQL(QUERY_SUB_LABELS, label.LabelId).Find(&child); err != nil {
		log.Log.Errorf("label_trace: get child labels info err:%s", err)
		return nil, err
	}

	return child, nil
}

// 树状图重构
func (m *LabelModel) tree(info []*VideoLabelList) []*VideoLabelList {
	if reflect.ValueOf(info).IsValid() {
		// 循环所有1级标签
		for k, v := range info {
			labelMp[k] = v.LabelName
			// 查询所有一级标签下的所有子标签
			child, err := m.FindSubLabelsByPid(v)
			if err != nil || len(child) == 0 {
				if v.Pid == 1 {
					info = append(info[:k], info[k + 1:]...)
				}

				continue
			}

			// 将子类别的数据循环赋值
			for k2, _ := range child {
				info[k].Child = append(info[k].Child, child[k2])
			}

			// 将刚刚查询出来的子类别进行递归 查询出三级,四级...子类
			m.tree(child)
		}

	}

	return info
}

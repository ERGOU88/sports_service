package mlabel

import (
	"github.com/go-xorm/xorm"
	"reflect"
	"sports_service/server/global/app/log"
	"sports_service/server/models"
	"sync"
	"fmt"
)

type LabelModel struct {
	Engine         *xorm.Session
	VideoLabels    *models.VideoLabelConfig
	PostLabels     *models.PostLabelConfig
}

// 标签列表信息
type VideoLabel struct {
	CreateAt  int           `json:"create_at" `
	Icon      string        `json:"icon"`
	LabelId   int           `json:"label_id"`
	LabelName string        `json:"label_name"`
	Pid       int           `json:"pid"`
	Sortorder int           `json:"sortorder"`
	Status    int           `json:"status"`
	UpdateAt  int           `json:"update_at"`
	Child     []*VideoLabel `json:"child" xorm:"-"` // 子类标签信息
}

// 添加视频标签请求参数
type AddVideoLabelParam struct {
	LabelName    string    `json:"label_name"`     // 视频标签名称
	Icon         string    `json:"icon"`           // icon
}

// 删除视频标签请求参数
type DelVideoLabelParam struct {
	LabelId      string    `json:"label_id"`       // 视频标签id
}

var videoLabels []*VideoLabel

// labelId -> VideoLabel
var labelMp map[string]*VideoLabel

var mutex sync.Mutex

func init() {
	labelMp = make(map[string]*VideoLabel)
}

// 实栗
func NewLabelModel(engine *xorm.Session) *LabelModel {
	return &LabelModel{
		Engine: engine,
		VideoLabels: new(models.VideoLabelConfig),
		PostLabels: new(models.PostLabelConfig),
	}
}

// 删除视频标签
func (m *LabelModel) DelVideoLabel(labelId string) error {
	if _, err := m.Engine.Where("label_id=?", labelId).Delete(&models.VideoLabels{}); err != nil {
		return err
	}

	return nil
}

// 添加视频标签
func (m *LabelModel) AddVideoLabel() error {
	if _, err := m.Engine.Insert(m.VideoLabels); err != nil {
		return err
	}

	return nil
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
func (m *LabelModel) IsExistsLabel(labelId string) bool {
	if _, ok := labelMp[labelId]; !ok {
		return false
	}

	return true
}

// 清空标签map
func (m *LabelModel) DelAllLabelByMem() {
	labelMp = make(map[string]*VideoLabel)
}

// 通过标签id 获取标签信息
func (m *LabelModel) GetLabelInfoByMem(labelId string) *VideoLabel {
	label, ok := labelMp[labelId]
	if !ok {
		return nil
	}

	return label
}

// 通过标签id 获取标签名称
func (m *LabelModel) GetLabelNameByMem(labelId string) string {
	label, ok := labelMp[labelId]
	if ok {
		return label.LabelName
	}

	return ""
}

// 删除标签信息（内存）
func (m *LabelModel) DelLabelInfoByMem(labelId string) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := labelMp[labelId]; ok {
		delete(labelMp, labelId)
	}
}

// 添加标签信息（内存）
func (m *LabelModel) AddLabelInfoByMem() {
	mutex.Lock()
	defer mutex.Unlock()
	info := &VideoLabel{
		CreateAt: m.VideoLabels.CreateAt,
		Icon: m.VideoLabels.Icon,
		LabelId: m.VideoLabels.LabelId,
		LabelName: m.VideoLabels.LabelName,
		Pid: m.VideoLabels.Pid,
		Sortorder: m.VideoLabels.Sortorder,
		Status: m.VideoLabels.Status,
		UpdateAt: m.VideoLabels.UpdateAt,
	}

	labelMp[fmt.Sprint(m.VideoLabels.LabelId)] = info
}

// 从内存读取视频标签 （第一次请求 内存没有 则从数据库load到内存）
func (m *LabelModel) GetVideoLabelList() []*VideoLabel {
	if len(videoLabels) == 0 {
		var err error
		err, videoLabels = m.LoadLabelsInfoByDb()
		if err != nil {
			return []*VideoLabel{}
		}
	}

	return videoLabels
}

const (
	QUERY_PARENT_LABELS = "SELECT * FROM `video_label_config` WHERE `pid` = 0 ORDER BY sortorder DESC"
)

// 从数据库获取标签信息
func (m *LabelModel) LoadLabelsInfoByDb() (error, []*VideoLabel) {
	// 定义指针切片用来存储所有标签
	var info []*VideoLabel
	// 定义指针切片返回控制器
	var res []*VideoLabel

	// 找出所有1级类别
	if err := m.Engine.Table(&models.VideoLabelConfig{}).SQL(QUERY_PARENT_LABELS).Find(&info); err != nil {
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
func (m *LabelModel) FindSubLabelsByPid(label *VideoLabel) ([]*VideoLabel, error) {
	var child []*VideoLabel
	if err := m.Engine.Table(&models.VideoLabelConfig{}).SQL(QUERY_SUB_LABELS, label.LabelId).Find(&child); err != nil {
		log.Log.Errorf("label_trace: get child labels info err:%s", err)
		return nil, err
	}

	return child, nil
}

// 树状图重构
func (m *LabelModel) tree(info []*VideoLabel) []*VideoLabel {
	if reflect.ValueOf(info).IsValid() {
		// 循环所有1级标签
		for k, v := range info {
			labelMp[fmt.Sprint(v.LabelId)] = v
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

package mconfigure

import (
	"github.com/go-xorm/xorm"
	"sports_service/global/app/log"
	"sports_service/global/consts"
	"sports_service/models"
)

// 配置相关
type ConfigModel struct {
	Engine         *xorm.Session
	VersionControl *models.AppVersionControl
	ActionScore    *models.ActionScoreConfig
}

// 实例
func NewConfigModel(engine *xorm.Session) *ConfigModel {
	return &ConfigModel{
		Engine:         engine,
		VersionControl: new(models.AppVersionControl),
		ActionScore:    new(models.ActionScoreConfig),
	}
}

// 添加新包请求参数
type AddPackageParams struct {
	VersionName string `json:"version_name"  binding:"required"` // 版本名
	Version     string `json:"version"  binding:"required"`      // 版本
	VersionCode int    `json:"version_code"  binding:"required"` // 版本code
	Size        string `json:"size"  binding:"required"`         // 包大小
	ByteSize    int    `json:"byte_size"`                        // 字节大小
	IsForce     int32  `json:"is_force" `                        // 是否强更 0 不需要强更 1 需要强更
	Status      int32  `json:"status"`                           // 1 可用 -1 不可用
	Platform    int32  `json:"platform"`                         // 0 android 1 ios
	UpgradeUrl  string `json:"upgrade_url"  binding:"required"`  // 新包地址
	Describe    string `json:"describe"`                         // 版本说明
}

// 更新包数据请求参数
type UpdatePackageParams struct {
	Id          int64  `json:"id"  binding:"required"`           // 数据id
	VersionName string `json:"version_name"  binding:"required"` // 版本名
	Version     string `json:"version"  binding:"required"`      // 版本
	VersionCode int    `json:"version_code"  binding:"required"` // 版本code
	Size        string `json:"size"  binding:"required"`         // 包大小
	ByteSize    int    `json:"byte_size"`                        // 字节大小
	IsForce     int32  `json:"is_force"`                         // 是否强更 0 不需要强更 1 需要强更
	Status      int32  `json:"status"`                           // 1 可用 -1 不可用
	Platform    int32  `json:"platform"`                         // 0 android 1 ios
	UpgradeUrl  string `json:"upgrade_url"  binding:"required"`  // 新包地址
	Describe    string `json:"describe"`                         // 版本说明
}

// 删除包请求参数
type DelPackageParam struct {
	Id int64 `json:"id" binding:"required"` // 数据id
}

type UpgradeInfo struct {
	Version     string `json:"version"`
	VersionCode int    `json:"version_code"`
	Size        string `json:"size"`
	Describe    string `json:"describe"`
	UpgradeURL  string `json:"upgrade_url"`
	HasNewPkg   bool   `json:"has_new_pkg"`
	IsForce     bool   `json:"is_force"`
}

// 添加新包
func (m *ConfigModel) AddNewPackage() (int64, error) {
	return m.Engine.InsertOne(m.VersionControl)
}

// 更新包信息
func (m *ConfigModel) UpdatePackageInfo(id int64) (int64, error) {
	return m.Engine.Where("id=?", id).Update(m.VersionControl)
}

// 删除包
func (m *ConfigModel) DelPackage(id int64) (int64, error) {
	return m.Engine.Where("id=?", id).Delete(&models.AppVersionControl{})
}

// 获取包版本信息列表
func (m *ConfigModel) GetPackageInfoList(offset, size int) []*models.AppVersionControl {
	var list []*models.AppVersionControl
	if err := m.Engine.Desc("id").Limit(size, offset).Find(&list); err != nil {
		return nil
	}

	return list
}

// 获取包详情
func (m *ConfigModel) GetPackageDetail(id string) *models.AppVersionControl {
	m.VersionControl = new(models.AppVersionControl)
	ok, err := m.Engine.Where("id=?", id).Get(m.VersionControl)
	if !ok || err != nil {
		return nil
	}

	return m.VersionControl
}

// 通过版本code及所属平台 获取包详情
func (m *ConfigModel) GetPackageDetailByVersion(plt int32, versionCode string) *models.AppVersionControl {
	m.VersionControl = new(models.AppVersionControl)
	ok, err := m.Engine.Where("platform=? AND version_code=?", plt, versionCode).Desc("id").Limit(1).Get(m.VersionControl)
	if !ok || err != nil {
		return nil
	}

	return m.VersionControl
}

// 获取最新包信息
func (m *ConfigModel) GetLatestPackageInfo(plt int32) *models.AppVersionControl {
	m.VersionControl = new(models.AppVersionControl)
	ok, err := m.Engine.Where("platform=? AND status=1", plt).Desc("version_code").Limit(1).Get(m.VersionControl)
	if !ok || err != nil {
		log.Log.Errorf("configure_trace: get latest package fail, plt:%d err:%s", plt, err)
		return nil
	}

	return m.VersionControl
}

// 通过版本号(例如：v1.0.1) 获取包信息
func (m *ConfigModel) GetPackageByVersion(version string, plt int32) *models.AppVersionControl {
	m.VersionControl = new(models.AppVersionControl)
	ok, err := m.Engine.Where("version=? and platform=?", version, plt).Desc("version_code").Limit(1).Get(m.VersionControl)
	if !ok || err != nil {
		return nil
	}

	return m.VersionControl
}

// 添加行为得分配置
func (m *ConfigModel) AddActionScoreConf(param *models.ActionScoreConfig) (int64, error) {
	return m.Engine.InsertOne(param)
}

// 更新行为得分配置
func (m *ConfigModel) UpdateActionScoreConf(param *models.ActionScoreConfig) (int64, error) {
	return m.Engine.Table(&models.ActionScoreConfig{}).Where("id=?", param.Id).Update(param)
}

// 获取行为得分配置列表
func (m *ConfigModel) GetActionScoreConfList() ([]*models.ActionScoreConfig, error) {
	var list []*models.ActionScoreConfig
	if err := m.Engine.Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取行为得分配置
func (m *ConfigModel) GetActionScoreConf() (bool, error) {
	var (
		ok  bool
		err error
	)
	if m.ActionScore.WorkType > 0 {
		ok, err = m.Engine.Where("work_type=?", m.ActionScore.WorkType).Get(m.ActionScore)
	}

	if m.ActionScore.Id > 0 {
		ok, err = m.Engine.Where("id=?", m.ActionScore.Id).Get(m.ActionScore)
	}

	return ok, err
}

// 获取行为得分
func (m *ConfigModel) GetActionScore(workType, actionType int) int {
	m.ActionScore.WorkType = workType
	// 默认1分
	score := 1
	if ok, err := m.GetActionScoreConf(); ok && err == nil {
		switch actionType {
		case consts.ACTION_TYPE_FABULOUS:
			score = m.ActionScore.FabulousScore
		case consts.ACTION_TYPE_BROWSE:
			score = m.ActionScore.BrowseScore
		case consts.ACTION_TYPE_SHARE:
			score = m.ActionScore.ShareScore
		case consts.ACTION_TYPE_BARRAGE:
			score = m.ActionScore.BarrageScore
		case consts.ACTION_TYPE_COMMENT:
			score = m.ActionScore.CommentScore
		case consts.ACTION_TYPE_COLLECT:
			score = m.ActionScore.CollectScore
		}
	}

	return score
}

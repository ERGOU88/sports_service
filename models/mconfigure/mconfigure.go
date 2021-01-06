package mconfigure

import (
  "github.com/go-xorm/xorm"
  "sports_service/server/models"
)

// 配置相关
type ConfigModel struct {
  Engine                *xorm.Session
  VersionControl        *models.AppVersionControl
}

// 实例
func NewConfigModel(engine *xorm.Session) *ConfigModel {
  return &ConfigModel{
    Engine: engine,
    VersionControl: new(models.AppVersionControl),
  }
}

// 添加新包请求参数
type AddPackageParams struct {
  VersionName     string     `json:"version_name"  binding:"required"`     // 版本名
  Version         string     `json:"version"  binding:"required"`          // 版本
  VersionCode     int        `json:"version_code"  binding:"required"`     // 版本code
  Size            string     `json:"size"  binding:"required"`             // 包大小
  IsForce         int32      `json:"is_force" `                            // 是否强更 0 不需要强更 1 需要强更
  Status          int32      `json:"status"`                               // 0 可用 1 不可用
  Platform        int32      `json:"platform"`                             // 0 android 1 ios
  UpgradeUrl      string     `json:"upgrade_url"  binding:"required"`      // 新包地址
}

// 更新包数据请求参数
type UpdatePackageParams struct {
  Id              int64      `json:"id"  binding:"required"`               // 数据id
  VersionName     string     `json:"version_name"  binding:"required"`     // 版本名
  Version         string     `json:"version"  binding:"required"`          // 版本
  VersionCode     int        `json:"version_code"  binding:"required"`     // 版本code
  Size            string     `json:"size"  binding:"required"`             // 包大小
  IsForce         int32      `json:"is_force"`                             // 是否强更 0 不需要强更 1 需要强更
  Status          int32      `json:"status"`                               // 0 可用 1 不可用
  Platform        int32      `json:"platform"`                             // 0 android 1 ios
  UpgradeUrl      string     `json:"upgrade_url"  binding:"required"`      // 新包地址
}

// 删除包请求参数
type DelPackageParam struct {
  Id              int64      `json:"id" binding:"required"`               // 数据id
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
  ok, err := m.Engine.Where("platform=?", plt).Desc("version_code").Limit(1).Get(m.VersionControl)
  if !ok || err != nil {
    return nil
  }

  return m.VersionControl
}

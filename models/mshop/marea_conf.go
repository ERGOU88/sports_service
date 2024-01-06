package mshop

import (
	"fmt"
	"reflect"
	"sports_service/dao"
	"sports_service/global/rdskey"
	"sports_service/models"
	"sports_service/util"
	"sync"
)

// 全国地址信息
type AreaInfo struct {
	Code  int64       `json:"code" xorm:"not null pk comment('区划代码') BIGINT(12)"`
	Name  string      `json:"name" xorm:"not null default '' comment('名称') index VARCHAR(128)"`
	Level int         `json:"level" xorm:"not null comment('级别1-5,省市县镇村') index TINYINT(1)"`
	Pcode int64       `json:"pcode" xorm:"comment('父级区划代码') index BIGINT(12)"`
	Child []*AreaInfo `json:"child" xorm:"-"` // 子地区信息
}

var areaList []*AreaInfo

// code -> AreaInfo
var areaMp map[string]*AreaInfo

var areaMutex sync.Mutex

func init() {
	areaMp = make(map[string]*AreaInfo)
}

func InitAreaInfo() {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	shopModel := NewShop(socket)
	//shopModel.GetArea()
	shopModel.LoadAreaMpToMem()
}

// 通过地区code获取地区信息
func (m *ShopModel) GetAreaInfoByCode(code string) *models.Area {
	info := new(models.Area)
	ok, err := m.Engine.Where("code=?", code).Get(info)
	if !ok || err != nil {
		return nil
	}

	return info
}

// 地区是否存在
func (m *ShopModel) HasExist(code string) bool {
	if _, ok := areaMp[code]; !ok {
		return false
	}

	return true
}

// 清空地区map
func (m *ShopModel) DelAllAreaByMem() {
	areaMp = make(map[string]*AreaInfo)
}

// code 获取内存地区信息
func (m *ShopModel) GetAreaInfoByMem(code string) *AreaInfo {
	area, ok := areaMp[code]
	if !ok {
		return nil
	}

	return area
}

// 通过code 获取地区名称
func (m *ShopModel) GetAreaNameByMem(code string) string {
	area, ok := areaMp[code]
	if ok {
		return area.Name
	}

	return ""
}

// 删除地区信息（内存）
func (m *ShopModel) DelAreaInfoByMem(code string) {
	areaMutex.Lock()
	defer areaMutex.Unlock()
	if _, ok := areaMp[code]; ok {
		delete(areaMp, code)
	}
}

// 修改频率低的数据 直接清理内存数据 下次从数据库load内存
func (m *ShopModel) CleanAreaInfoByMem() {
	areaList = nil
}

// 将地区map load到内存
func (m *ShopModel) LoadAreaMpToMem() {
	if areaMp != nil {
		return
	}

	mpStr, err := m.GetAreaMpByRds()
	if err != nil {
		return
	}

	if err = util.JsonFast.UnmarshalFromString(mpStr, &areaMp); err != nil {
		return
	}

	return
}

// 从内存读取地区信息 （第一次请求 内存没有 则从数据库load到内存）
func (m *ShopModel) GetArea() []*AreaInfo {
	if len(areaList) > 0 {
		return areaList
	}

	// 从redis获取
	area, err := m.GetAreaInfoByRds()
	if err == nil {
		if err = util.JsonFast.UnmarshalFromString(area, &areaList); err == nil {
			return areaList
		}
	}

	err, areaList = m.LoadAreaInfoByDb()
	if err != nil {
		return []*AreaInfo{}
	}

	str, err := util.JsonFast.MarshalToString(areaList)
	if err == nil {
		m.SetAreaInfoToRds(str)
	}

	mpStr, err := util.JsonFast.MarshalToString(areaMp)
	if err == nil {
		m.SetAreaMpToRds(mpStr)
	}

	return areaList
}

// 从redis获取地区列表信息
func (m *ShopModel) GetAreaInfoByRds() (string, error) {
	rds := dao.NewRedisDao()
	return rds.Get(rdskey.AREA_LIST_INFO)
}

// 地区列表存储到redis
func (m *ShopModel) SetAreaInfoToRds(areaInfo string) error {
	rds := dao.NewRedisDao()
	return rds.Set(rdskey.AREA_LIST_INFO, areaInfo)
}

// 地区map存储到redis
func (m *ShopModel) SetAreaMpToRds(areaMp string) error {
	rds := dao.NewRedisDao()
	return rds.Set(rdskey.AREA_MP_INFO, areaMp)
}

// 从redis获取地区map信息
func (m *ShopModel) GetAreaMpByRds() (string, error) {
	rds := dao.NewRedisDao()
	return rds.Get(rdskey.AREA_MP_INFO)
}

const (
	QUERY_PARENT_AREA = "SELECT * FROM `area` WHERE `pcode` = 0 AND is_show=0"
)

// 从数据库获取地区信息
func (m *ShopModel) LoadAreaInfoByDb() (error, []*AreaInfo) {
	// 定义指针切片用来存储所有地区
	var info []*AreaInfo
	// 定义指针切片返回控制器
	var res []*AreaInfo

	// 找出所有1级地区
	if err := m.Engine.Table(&models.Area{}).SQL(QUERY_PARENT_AREA).Find(&info); err != nil {
		return err, nil
	}

	// 判断是否存在数据 存在 则进行树状图重构
	if reflect.ValueOf(info).IsValid() {
		res = m.trees(info)
	}

	return nil, res
}

const (
	QUERY_SUB_AREA = "SELECT * FROM `area` WHERE `pcode`=? AND is_show=0" // 查询某个父类地区下的所有子地区
)

// 查询某个父类地区下的所有子地区
func (m *ShopModel) FindSubAreaByCode(info *AreaInfo) ([]*AreaInfo, error) {
	var child []*AreaInfo
	if err := m.Engine.Table(&models.Area{}).SQL(QUERY_SUB_AREA, info.Code).Find(&child); err != nil {
		return []*AreaInfo{}, err
	}

	return child, nil
}

// 树状图重构
func (m *ShopModel) trees(info []*AreaInfo) []*AreaInfo {
	if reflect.ValueOf(info).IsValid() {
		// 循环所有1级地区
		for k, v := range info {
			areaMp[fmt.Sprint(v.Code)] = v
			// 查询所有一级地区下的所有子地区
			child, err := m.FindSubAreaByCode(v)
			if err != nil || len(child) == 0 {
				if v.Level == 1 {
					info = append(info[:k], info[k+1:]...)
				}

				continue
			}

			// 将子类别的数据循环赋值
			for k2, _ := range child {
				info[k].Child = append(info[k].Child, child[k2])
			}

			// 将刚刚查询出来的子地区进行递归 查询出三级,四级...地区
			m.trees(child)
		}
	}

	return info
}

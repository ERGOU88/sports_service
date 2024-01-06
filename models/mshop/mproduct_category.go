package mshop

import (
	"fmt"
	"reflect"
	"sports_service/dao"
	"sports_service/models"
	"sync"
)

// 商品分类列表信息
type Category struct {
	CategoryId   int         `json:"category_id"`
	CategoryName string      `json:"category_name"`
	ShortName    string      `json:"short_name"`
	Pid          int         `json:"pid"`
	Level        int         `json:"level"`
	IsShow       int         `json:"is_show"`
	Sortorder    int         `json:"sortorder"`
	Image        string      `json:"image"`
	Keywords     string      `json:"keywords"`
	Description  string      `json:"description"`
	Child        []*Category `json:"child" xorm:"-"` // 子分类信息
}

var categoryList []*Category

// categoryId -> Category
var categoryMp map[string]*Category

var mutex sync.Mutex

func init() {
	categoryMp = make(map[string]*Category)
}

func InitProductCategory() {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	shopModel := NewShop(socket)
	shopModel.GetProductCategory()
}

// 通过分类id获取商品分类信息
func (m *ShopModel) GetProductCategoryInfoById(categoryId string) *models.ProductCategory {
	info := new(models.ProductCategory)
	ok, err := m.Engine.Where("category_id=?", categoryId).Get(info)
	if !ok || err != nil {
		return nil
	}

	return info
}

// 分类是否存在
func (m *ShopModel) IsExists(categoryId string) bool {
	if _, ok := categoryMp[categoryId]; !ok {
		return false
	}

	return true
}

// 清空分类map
func (m *ShopModel) DelAllCategoryByMem() {
	categoryMp = make(map[string]*Category)
}

// 通过分类id 获取分类信息
func (m *ShopModel) GetCategoryInfoByMem(categoryId string) *Category {
	category, ok := categoryMp[categoryId]
	if !ok {
		return nil
	}

	return category
}

// 通过分类id 获取分类名称
func (m *ShopModel) GetCategoryNameByMem(categoryId string) string {
	category, ok := categoryMp[categoryId]
	if ok {
		return category.CategoryName
	}

	return ""
}

// 删除分类信息（内存）
func (m *ShopModel) DelCategoryInfoByMem(categoryId string) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := categoryMp[categoryId]; ok {
		delete(categoryMp, categoryId)
	}
}

// 修改频率低的数据 直接清理内存数据 下次从数据库load内存
func (m *ShopModel) CleanCategoryInfoByMem() {
	categoryList = nil
}

// 从内存读取分类 （第一次请求 内存没有 则从数据库load到内存）
func (m *ShopModel) GetProductCategory() []*Category {
	//if len(categoryList) == 0 {
	var err error
	err, categoryList = m.LoadCategoryInfoByDb()
	if err != nil {
		return []*Category{}
	}
	//}

	return categoryList
}

func (m *ShopModel) GetProductCategoryByBackend() (error, []*Category) {
	var list []*Category
	if err := m.Engine.Table(&models.ProductCategory{}).Find(&list); err != nil {
		return err, nil
	}

	return nil, list
}

const (
	QUERY_PARENT_CATEGORY = "SELECT * FROM `product_category` WHERE `pid` = 0 AND is_show=0 ORDER BY sortorder DESC"
)

// 从数据库获取分类信息
func (m *ShopModel) LoadCategoryInfoByDb() (error, []*Category) {
	// 定义指针切片用来存储所有分类
	var info []*Category
	// 定义指针切片返回控制器
	var res []*Category

	// 找出所有1级类别
	if err := m.Engine.Table(&models.ProductCategory{}).SQL(QUERY_PARENT_CATEGORY).Find(&info); err != nil {
		return err, nil
	}

	// 判断是否存在数据 存在 则进行树状图重构
	if reflect.ValueOf(info).IsValid() {
		res = m.tree(info)
	}

	return nil, res
}

const (
	QUERY_SUB_CATEGORY = "SELECT * FROM `product_category` WHERE `pid`=? AND is_show=0" // 查询某个父类下的所有子类别
)

// 通过父类分类id查询下属的子分类
func (m *ShopModel) FindSubCategoryByPid(info *Category) ([]*Category, error) {
	var child []*Category
	if err := m.Engine.Table(&models.ProductCategory{}).SQL(QUERY_SUB_CATEGORY, info.CategoryId).Find(&child); err != nil {
		return []*Category{}, err
	}

	return child, nil
}

// 树状图重构
func (m *ShopModel) tree(info []*Category) []*Category {
	if reflect.ValueOf(info).IsValid() {
		// 循环所有1级分类
		for k, v := range info {
			categoryMp[fmt.Sprint(v.CategoryId)] = v
			// 查询所有一级分类下的所有子分类
			child, err := m.FindSubCategoryByPid(v)
			if err != nil || len(child) == 0 {
				if v.Pid == 1 {
					info = append(info[:k], info[k+1:]...)
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

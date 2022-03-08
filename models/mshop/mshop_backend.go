package mshop

import (
	"sports_service/server/models"
	"errors"
	"sports_service/server/tools/tencentCloud"
)

type AddOrEditCategorySpecReq struct {
	CategoryId     int64        `json:"category_id"`
	SpecInfo       []SpecInfo   `json:"spec_info"`
}

// 管理后台发货请求
type DeliverProductReq struct {
	OrderId           string          `json:"order_id" binding:"required"`
	DeliveryTypeName  string          `json:"delivery_type_name" binding:"required"`     // 配送方式名称
	DeliveryCode      string          `json:"delivery_code" binding:"required"`          // 运单号
	DeliveryTelephone string          `json:"delivery_telephone" binding:"required"`     // 承运人电话
}

// 添加/编辑商品 请求参数
type AddOrEditProductReq struct {
	Id             int    `json:"id" xorm:"not null pk autoincr comment('商品id') INT(11)"`
	ProductName    string `json:"product_name" xorm:"not null default '' comment('商品名称') VARCHAR(255)"`
	ProductImage   string `json:"product_image" xorm:"not null default '' comment('商品主图路径') VARCHAR(512)"`
	ProductDetail  []string `json:"product_detail" xorm:"comment('商品详情') TEXT"`
	Status         int    `json:"status" xorm:"not null default 1 comment('商品状态（0. 正常 1. 下架）') TINYINT(2)"`
	IsFreeShip     int    `json:"is_free_ship" xorm:"not null default 0 comment('是否免邮 0 免邮') TINYINT(2)"`
	IsDelete       int    `json:"is_delete" xorm:"not null default 0 comment('是否已经删除') index TINYINT(4)"`
	Introduction   []string `json:"introduction" xorm:"not null default '' comment('促销语') VARCHAR(255)"`
	Keywords       string `json:"keywords" xorm:"not null default '' comment('关键词') VARCHAR(255)"`
	Sortorder      int    `json:"sortorder" xorm:"not null default 0 comment('排序') index INT(11)"`
	VideoUrl       string `json:"video_url" xorm:"not null default '' comment('视频地址') VARCHAR(512)"`
	CurPrice       int    `json:"cur_price" xorm:"not null default 0 comment('商品价格（分）') INT(10)"`
	MarketPrice    int    `json:"market_price" xorm:"not null default 0 comment('划线价格（分）') INT(10)"`
	IsRecommend    int    `json:"is_recommend" xorm:"not null default 0 comment('推荐（0：不推荐；1：推荐）') TINYINT(1)"`
	IsTop          int    `json:"is_top" xorm:"not null default 0 comment('置顶（0：不置顶；1：置顶；）') TINYINT(1)"`
	IsCream        int    `json:"is_cream" xorm:"not null default 0 comment('是否精华（0: 不是 1: 是）') TINYINT(1)"`
	Specifications []*SpecInfo     `json:"specifications" xorm:"not null default '' comment('全部规格参数数据') VARCHAR(3000)"`
	SpecTemplate   []*SpecTemplate  `json:"spec_template" xorm:"not null default '' comment('特有规格参数及可选值信息，json格式') VARCHAR(1000)"`
	AfterService   []int64 `json:"after_service" xorm:"default '' comment('售后服务') VARCHAR(1000)"`
	CreateAt       int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt       int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
	SkuList        []*ProductSkuInfo   `json:"sku_list"`
	CategoryId     int    `json:"category_id"`
	IsReset        bool   `json:"is_reset"`     // 是否重置sku
}

// 商品sku
type ProductSkuInfo struct {
	Id            int    `json:"id" xorm:"not null pk autoincr comment('sku id') INT(11)"`
	ProductId     int64  `json:"product_id" xorm:"not null comment('商品id') index BIGINT(20)"`
	Title         string `json:"title" xorm:"not null comment('商品标题') VARCHAR(255)"`
	SkuImage      tencentCloud.BucketURI `json:"sku_image" xorm:"not null default '' comment('sku主图') VARCHAR(255)"`
	SkuNo         string `json:"sku_no" xorm:"not null default '' comment('商品sku编码') VARCHAR(255)"`
	Images        []tencentCloud.BucketURI `json:"images" xorm:"default '' comment('sku图片[多张]') VARCHAR(1000)"`
	CurPrice      int    `json:"cur_price" xorm:"not null default 0 comment('商品价格（分）') INT(10)"`
	MarketPrice   int    `json:"market_price" xorm:"not null default 0 comment('划线价格（分）') INT(10)"`
	IsFreeShip    int    `json:"is_free_ship" xorm:"not null default 0 comment('是否免邮 0 免邮') TINYINT(2)"`
	Indexes       string `json:"indexes" xorm:"comment('特有规格属性在商品属性模板中的对应下标组合 例如0_1_0 则可能表示 颜色选项 0 远峰蓝 内存选项 1 3GB  机身存储 0  16GB') index VARCHAR(100)"`
	OwnSpec       []*OwnSpec `json:"own_spec" xorm:"comment('商品实体的特有规格参数，json格式，反序列化时保证有序') VARCHAR(1000)"`
	Status        int    `json:"status" xorm:"not null default 0 comment('是否有效，0 有效, 1 无效') TINYINT(1)"`
	DiscountPrice int    `json:"discount_price" xorm:"not null default 0 comment('活动折扣价（默认等于单价）单位：分') INT(10)"`
	StartTime     int    `json:"start_time" xorm:"not null default 0 comment('活动开始时间') INT(11)"`
	EndTime       int    `json:"end_time" xorm:"not null default 0 comment('活动结束时间') INT(11)"`
	CreateAt      int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt      int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
	Sortorder     int    `json:"sortorder" xorm:"not null default 0 comment('排序权重') INT(11)"`
	Stock         int    `json:"stock"`
	MaxBuy        int    `json:"max_buy"`
 	MinBuy        int    `json:"min_buy"`
}

// 管理后台商品详情
type ProductDetail struct {
	Id             int    `json:"id" xorm:"not null pk autoincr comment('商品id') INT(11)"`
	ProductName    string `json:"product_name" xorm:"not null default '' comment('商品名称') VARCHAR(255)"`
	ProductImage   tencentCloud.BucketURI `json:"product_image" xorm:"not null default '' comment('商品主图路径') VARCHAR(512)"`
	ProductDetail  []tencentCloud.BucketURI `json:"product_detail" xorm:"comment('商品详情') TEXT"`
	Status         int    `json:"status" xorm:"not null default 1 comment('商品状态（0. 正常 1. 下架）') TINYINT(2)"`
	IsFreeShip     int    `json:"is_free_ship" xorm:"not null default 0 comment('是否免邮 0 免邮') TINYINT(2)"`
	IsDelete       int    `json:"is_delete" xorm:"not null default 0 comment('是否已经删除') index TINYINT(4)"`
	Introduction   []string `json:"introduction" xorm:"not null default '' comment('促销语') VARCHAR(255)"`
	Keywords       string `json:"keywords" xorm:"not null default '' comment('关键词') VARCHAR(255)"`
	Sortorder      int    `json:"sortorder" xorm:"not null default 0 comment('排序') index INT(11)"`
	VideoUrl       tencentCloud.BucketURI `json:"video_url" xorm:"not null default '' comment('视频地址') VARCHAR(512)"`
	CurPrice       int    `json:"cur_price" xorm:"not null default 0 comment('商品价格（分）') INT(10)"`
	MarketPrice    int    `json:"market_price" xorm:"not null default 0 comment('划线价格（分）') INT(10)"`
	IsRecommend    int    `json:"is_recommend" xorm:"not null default 0 comment('推荐（0：不推荐；1：推荐）') TINYINT(1)"`
	IsTop          int    `json:"is_top" xorm:"not null default 0 comment('置顶（0：不置顶；1：置顶；）') TINYINT(1)"`
	IsCream        int    `json:"is_cream" xorm:"not null default 0 comment('是否精华（0: 不是 1: 是）') TINYINT(1)"`
	Specifications []*SpecInfo     `json:"specifications" xorm:"not null default '' comment('全部规格参数数据') VARCHAR(3000)"`
	SpecTemplate   []*SpecTemplate  `json:"spec_template" xorm:"not null default '' comment('特有规格参数及可选值信息，json格式') VARCHAR(1000)"`
	AfterService   []*AfterService `json:"after_service" xorm:"default '' comment('售后服务') VARCHAR(1000)"`
	CreateAt       int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt       int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
	SkuList        []*ProductSkuInfo   `json:"sku_list" xorm:"-"`
	CategoryId     int    `json:"category_id" xorm:"-"`
	CategoryName   string `json:"category_name" xorm:"-"`
}

type CategorySpecInfo struct {
	Specifications    []*SpecInfo   `json:"specifications"`
	CategoryId        int64         `json:"category_id"`
	CreateAt          int64         `json:"create_at"`
	UpdateAt          int64         `json:"update_at"`
	CategoryName      string        `json:"category_name" xorm:"-"`
}

type ProductCategory struct {
	CategoryId   int    `json:"category_id" xorm:"not null pk autoincr comment('分类id') INT(11)"`
	CategoryName string `json:"category_name" xorm:"not null default '' comment('分类名称') VARCHAR(50)"`
	IsShow       int    `json:"is_show" xorm:"not null default 0 comment('是否显示（0显示  1不显示）') INT(11)"`
	Sortorder    int    `json:"sortorder" xorm:"not null default 0 comment('排序') INT(11)"`
}

func (m *ShopModel) AddProductCategory(info *models.ProductCategory) (int64, error) {
	return m.Engine.InsertOne(info)
}

func (m *ShopModel) UpdateProductCategory(id int, mp map[string]interface{}) (int64, error) {
	return m.Engine.Table(&models.ProductCategory{}).Where("category_id=?", id).Update(mp)
}

func (m *ShopModel) GetServiceList() ([]AfterService, error) {
	var list []AfterService
	if err := m.Engine.Table(models.ShopServiceConf{}).Find(&list); err != nil {
		return list, err
	}
	
	return list, nil
}

func (m *ShopModel) AddService(info *models.ShopServiceConf) (int64, error) {
	return m.Engine.InsertOne(info)
}

func (m *ShopModel) DelService(id string) (int64, error) {
	return m.Engine.Where("id=?", id).Delete(&models.ShopServiceConf{})
}

func (m *ShopModel) UpdateService(info *models.ShopServiceConf) (int64, error) {
	return m.Engine.Where("id=?", info.Id).Update(info)
}

func (m *ShopModel) GetCategorySpec(categoryId string) (*CategorySpecInfo, error) {
	info := &CategorySpecInfo{}
	ok, err := m.Engine.Table(&models.ProductSpecification{}).Where("category_id=?", categoryId).Get(info)
	if err != nil {
		return nil, err
	}
	
	if !ok {
		return nil, errors.New("sku not found")
	}
	
	return info, nil
}

func (m *ShopModel) GetCategorySpecList() ([]*CategorySpecInfo, error) {
	var list []*CategorySpecInfo
	if err := m.Engine.Table(&models.ProductSpecification{}).Find(&list); err != nil {
		return nil, err
	}
	
	return list, nil
}

func (m *ShopModel) AddCategorySpec(info *models.ProductSpecification) (int64, error) {
	return m.Engine.InsertOne(info)
}

func (m *ShopModel) UpdateCategorySpec(info *models.ProductSpecification) (int64, error) {
	return m.Engine.Where("category_id=?", info.CategoryId).Update(info)
}

func (m *ShopModel) DelCategorySpec(categoryId string) (int64, error) {
	return m.Engine.Where("category_id=?", categoryId).Delete(&models.ProductSpecification{})
}

// 添加商品spu
func (m *ShopModel) AddProductSpu(spu *models.Products) (int64, error) {
	return m.Engine.InsertOne(spu)
}

// 更新商品spu
func (m *ShopModel) UpdateProductSpu(id string, mp map[string]interface{}) (int64, error) {
	return m.Engine.Table(&models.Products{}).Where("id=?", id).Update(mp)
}

// 添加商品sku
func (m *ShopModel) AddProductSku(sku *models.ProductSku) (int64, error) {
	return m.Engine.InsertOne(sku)
}

// 更新商品sku
func (m *ShopModel) UpdateProductSku(id string, mp map[string]interface{}) (int64, error) {
	return m.Engine.Table(&models.ProductSku{}).Where("id=?", id).Update(mp)
}

// 添加商品sku库存
func (m *ShopModel) AddProductSkuStock(stock []*models.ProductSkuStock) (int64, error) {
	return m.Engine.InsertMulti(stock)
}

func (m *ShopModel) UpdateProductSkuStockInfo(stock *models.ProductSkuStock) (int64, error) {
	return m.Engine.Where("sku_id=?", stock.SkuId).Update(stock)
}

// 添加商品分类关联
func (m *ShopModel) AddProductCategoryRelated(info *models.ProductCategoryRelated) (int64, error) {
	return m.Engine.InsertOne(info)
}

func (m *ShopModel) DelProductCategoryRelated(productId string) (int64, error) {
	return m.Engine.Where("product_id=?", productId).Update(&models.ProductCategoryRelated{Status: 1})
}

func (m *ShopModel) GetProductCategoryRelated(condition string) (*models.ProductCategoryRelated, error) {
	info := &models.ProductCategoryRelated{}
	ok, err := m.Engine.Where(condition).Get(info)
	if err != nil {
		return nil, err
	}
	
	if !ok {
		return nil, errors.New("sku not found")
	}
	
	return info, nil
}

// 软删除废弃的sku
func (m *ShopModel) DelSkuByProductId(productId string) (int64, error) {
	return m.Engine.Where("product_id=?", productId).Update(&models.ProductSku{IsDelete: 1, Status:1})
}

// 管理后台获取spu列表 及总数
func (m *ShopModel) GetSpuList(sortType, keyword string, offset, size int) ([]*ProductSimpleInfo, int64, error) {
	var list []*ProductSimpleInfo
	engine := m.Engine
	engine.Where("is_delete=0")
	if keyword != "" {
		engine.Where("product_name like '%" + keyword + "%'")
	}
	
	switch sortType {
	case "0":
		engine.OrderBy("id DESC, sale_num DESC, is_top DESC, is_recommend DESC, is_cream DESC, sortorder DESC")
	case "1":
		engine.OrderBy("id DESC, sale_num ASC, is_top DESC, is_recommend DESC, is_cream DESC, sortorder DESC")
	case "2":
		engine.OrderBy("id DESC, cur_price DESC, is_top DESC, is_recommend DESC, is_cream DESC, sortorder DESC")
	case "3":
		engine.OrderBy("id DESC, cur_price ASC, is_top DESC, is_recommend DESC, is_cream DESC, sortorder DESC")
	default:
		engine.OrderBy("id DESC, sale_num DESC, is_top DESC, is_recommend DESC, is_cream DESC, sortorder DESC")
	}
	
	count, _ := engine.Count(&models.Products{})
	if err := engine.Table(&models.Products{}).Limit(size, offset).Find(&list); err != nil {
		return nil, count, err
	}
	
	return list, count, nil
}

func (m *ShopModel) GetSpuTotal() (int64, error) {
	return m.Engine.Count(&models.Products{})
}

// 添加商品服务
func (m *ShopModel) AddProductService(list []*models.ProductService) (int64, error) {
	return m.Engine.InsertMulti(list)
}

// 删除商品服务
func (m *ShopModel) DelProductService(condition string) (int64, error) {
	return m.Engine.Where(condition).Delete(&models.ProductService{})
}

// 管理后台获取商品spu
func (m *ShopModel) GetProductById(productId string) (*ProductDetail, error) {
	spu := &ProductDetail{}
	ok, err := m.Engine.Table(&models.Products{}).Where("id=?", productId).Get(spu)
	if !ok || err != nil {
		return nil, err
	}
	
	if !ok {
		return nil, errors.New("spu not found")
	}
	
	return spu, nil
	
}

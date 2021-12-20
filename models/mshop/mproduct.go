package mshop

import (
	"sports_service/server/models"
)

type ProductSimpleInfo struct {
	Id             int    `json:"id" xorm:"not null pk autoincr comment('商品id') INT(11)"`
	ProductName    string `json:"product_name" xorm:"not null default '' comment('商品名称') VARCHAR(255)"`
	ProductImage   string `json:"product_image" xorm:"not null default '' comment('商品主图路径') VARCHAR(512)"`
	Status         int    `json:"status" xorm:"not null default 0 comment('商品状态（0. 正常 1. 下架）') TINYINT(2)"`
	IsFreeShip     int    `json:"is_free_ship" xorm:"not null default 0 comment('是否免邮 0 免邮') TINYINT(2)"`
	Introduction   string `json:"introduction" xorm:"not null default '' comment('促销语') VARCHAR(255)"`
	Keywords       string `json:"keywords" xorm:"not null default '' comment('关键词') VARCHAR(255)"`
	Sortorder      int    `json:"sortorder" xorm:"not null default 0 comment('排序') index INT(11)"`
	VideoUrl       string `json:"video_url" xorm:"not null default '' comment('视频地址') VARCHAR(512)"`
	SaleNum        int    `json:"sale_num" xorm:"not null default 0 comment('销量') INT(11)"`
	CurPrice       int    `json:"cur_price" xorm:"not null default 0 comment('商品价格（分）') INT(10)"`
	MarketPrice    int    `json:"market_price" xorm:"not null default 0 comment('划线价格（分）') INT(10)"`
}

type ProductDetail struct {

}

const (
	GET_SPU_LIST_BY_CATEGORY = "SELECT p.* FROM products AS p LEFT JOIN product_category_related AS pc " +
		"ON p.id = pc.product_id WHERE pc.category_id = ? AND p.status = 0 AND p.is_delete = 0 " +
		"ORDER BY p.sale_num DESC, p.is_top DESC, p.is_recommend DESC, p.is_cream DESC, p.sortorder DESC, p.id DESC LIMIT ?, ?"
)
// 通过分类获取spu列表
// sortType 0 销量倒序 1 价格倒序 2 价格正序
func (m *ShopModel) GetSpuListByCategory(categoryId, sortType string, offset, size int) ([]ProductSimpleInfo, error) {
	var list []ProductSimpleInfo
	sql := "SELECT p.* FROM products AS p LEFT JOIN product_category_related AS pc " +
		"ON p.id = pc.product_id WHERE pc.category_id = ? AND p.status = 0 AND p.is_delete = 0 ORDER BY "

	switch sortType {
	case "0":
		sql += "p.sale_num DESC, "
	case "1":
		sql += "p.cur_price DESC, "
	case "2":
		sql += "p.cur_price ASC, "
	default:
		sql += "p.sale_num DESC, "
	}

	sql += "p.is_top DESC, p.is_recommend DESC, p.is_cream DESC, p.sortorder DESC, p.id DESC LIMIT ?, ?"
	if err := m.Engine.SQL(sql, categoryId, offset, size).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取所有商品spu
// sortType 0 销量倒序 1 价格倒序 2 价格正序
func (m *ShopModel) GetAllSpu(sortType string, offset, size int) ([]ProductSimpleInfo, error) {
	var list []ProductSimpleInfo
	engine := m.Engine.Table(&models.Products{})
	engine.Where("status=0 AND is_delete=0")

	switch sortType {
	case "0":
		engine.OrderBy("sale_num DESC, is_top DESC, is_recommend DESC, is_cream DESC, sortorder DESC")
	case "1":
		engine.OrderBy("cur_price ASC, is_top DESC, is_recommend DESC, is_cream DESC, sortorder DESC")
	case "2":
		engine.OrderBy("cur_price DESC, is_top DESC, is_recommend DESC, is_cream DESC, sortorder DESC")
	default:
		engine.OrderBy("sale_num DESC, is_top DESC, is_recommend DESC, is_cream DESC, sortorder DESC")
	}

	if err := engine.Limit(size, offset).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

const (
	GET_RECOMMEND_PRODUCTS = "SELECT DISTINCT(p.id), p.* FROM `products` AS p JOIN (SELECT RAND() * " +
		"((SELECT MAX(id) FROM `products`)-(SELECT MIN(id) FROM `products`)) + (SELECT MIN(id) FROM `products`) AS id) AS p2 " +
		"WHERE p.id >= p2.id-1 AND p.id != ? AND p.status=0 AND p.is_delete=0 ORDER BY p.id LIMIT ?"
)
func (m *ShopModel) GetRecommendProducts(productId string, limit int) ([]ProductSimpleInfo, error) {
	var list []ProductSimpleInfo
	if err := m.Engine.SQL(GET_RECOMMEND_PRODUCTS, productId, limit).Find(&list); err != nil {
		return list, err
	}

	return list, nil
}

// 获取sku列表
func (m *ShopModel) GetProductSkuList() ([]models.ProductSku, error) {
	var list []models.ProductSku
	if err := m.Engine.Where("status=0").Desc("sortorder").Find(&list); err != nil {
		return list, err
	}

	return list, nil
}

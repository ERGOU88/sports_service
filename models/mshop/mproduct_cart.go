package mshop

import (
	"sports_service/server/models"
	//tc "sports_service/server/tools/tencentCloud"
)

type ProductCartInfo struct {
	Id              int             `json:"sku_id"`
	CartId          int             `json:"cart_id"`
	ProductId       int64           `json:"product_id"`
	Title           string          `json:"title"`
	SkuImage        string          `json:"sku_image"`
	SkuNo           string          `json:"sku_no"`
	CurPrice        int             `json:"cur_price"`
	MarketPrice     int             `json:"market_price"`
	IsFreeShip      int             `json:"is_free_ship"`
	DiscountPrice   int             `json:"discount_price"`
	StartTime       int64           `json:"start_time"`
	EndTime         int64           `json:"end_time"`
	RemainDuration  int64           `json:"remain_duration"`        // 活动剩余时长
	HasActivities   int32           `json:"has_activities"`         // 1 有活动
	OwnSpec         []OwnSpec       `json:"own_spec"`               // 商品实体的特有规格参数
	Indexes         string          `json:"indexes"`                // 特有规格属性在商品属性模板中的对应下标组合
	Stock           int             `json:"stock"`                  // 库存
	MaxBuy          int             `json:"max_buy"`                // 限购 0 表示无限制
	MinBuy          int             `json:"min_buy"`                // 起购数
	Count           int             `json:"count"`                  // 当前数量
	IsCheck         int             `json:"is_check"`               // 0选择 1 未选择
}

type UpdateProductCartParam struct {
	Params []*models.ProductCart `json:"params"`
}

func (m *ShopModel) GetProductCart(condition string) (*models.ProductCart, error) {
	cart := &models.ProductCart{}
	ok, err := m.Engine.Where(condition).Get(cart)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return cart, nil
}

// 添加商品购物车
func (m *ShopModel) AddProductCart(info *models.ProductCart) (int64, error) {
	return m.Engine.InsertOne(info)
}

const (
	UPDATE_PRODUCT_CART = "UPDATE `product_cart` SET count=count+?, is_check=?, create_at=? WHERE product_id=? AND sku_id=? AND user_id=?"
)
// 更新商品购物车
func (m *ShopModel) UpdateProductCart(info *models.ProductCart) (int64, error) {
	res, err := m.Engine.Exec(UPDATE_PRODUCT_CART, info.Count, info.IsCheck, info.CreateAt, info.ProductId, info.SkuId, info.UserId)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

// 获取用户购物车数量
func (m *ShopModel) GetProductCartNum(userId string) (int64, error) {
	return m.Engine.Where("user_id=?", userId).Count(&models.ProductCart{})
}

const (
	GET_PRODUCT_CART_LIST = "SELECT sku.*, cart.id AS cart_id, cart.count, cart.user_id, cart.is_check FROM product_sku AS sku INNER JOIN product_cart AS cart " +
		"ON sku.id=cart.sku_id WHERE cart.user_id=? ORDER BY sku.create_at DESC"
)
func (m *ShopModel) GetProductCartList(userId string) ([]*ProductCartInfo, error) {
	var list []*ProductCartInfo
	if err := m.Engine.SQL(GET_PRODUCT_CART_LIST, userId).Find(&list); err != nil {
		return list, err
	}

	return list, nil
}

func (m *ShopModel) UpdateProductCartById(info *models.ProductCart) (int64, error) {
	return m.Engine.Where("id=?", info.Id).Cols("count, is_check").Update(info)
}

// 清理用户购物车
func (m *ShopModel) DelProductCartByIds(cartIds []int) (int64, error) {
	return m.Engine.In("id", cartIds).Delete(&models.ProductCart{})
}

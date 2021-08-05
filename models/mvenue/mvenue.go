package mvenue

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

type VenueModel struct {
	Venue    *models.VenueInfo
	Engine   *xorm.Session
}

// 场馆商品
type VenueProduct struct {
	Id                int    `json:"id"`                   // 商品id
	ProductName       string `json:"product_name"`
	ProductType       int    `json:"product_type"`
	EffectiveDuration int    `json:"effective_duration"`   // 有效时长 例如体验卡 15/h
	RealAmount        int    `json:"real_amount"`          // 定价
	CurAmount         int    `json:"cur_amount"`           // 售价
	DiscountRate      int    `json:"discount_rate"`        // 折扣率
	DiscountAmount    int    `json:"discount_amount"`      // 优惠金额
	HasDiscount       int32  `json:"has_discount"`         // 是否打折 0 未打折 1 打折
	VenueId           int64  `json:"venue_id"`             // 场馆id
	Sales             int64  `json:"sales"`                // 销量
	Icon              string `json:"icon"`                 // 商品icon
	Image             string `json:"image"`                // 商品图片
	Describe          string `json:"describe"`             // 商品介绍
	Title             string `json:"title"`                // 商品简介
}

func NewVenueModel(engine *xorm.Session) *VenueModel {
	return &VenueModel{
		Venue: new(models.VenueInfo),
		Engine: engine,
	}
}

// 通过场馆id 获取场馆信息
func (m *VenueModel) GetVenueInfoById() error {
	ok, err := m.Engine.Where("id=?", m.Venue.Id).Get(m.Venue)
	if !ok || err != nil {
		return err
	}

	return nil
}

// 通过场馆id 获取场馆商品列表
func (m *VenueModel) GetVenueProducts() ([]*models.VenueProductInfo, error) {
	var list []*models.VenueProductInfo
	if err := m.Engine.Where("id=?", m.Venue.Id).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}


package mvenue

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

type VenueModel struct {
	Venue     *models.VenueInfo
	Engine    *xorm.Session
	Recommend *models.VenueRecommendConf
	Product   *models.VenueProductInfo
	Vip       *models.VenueVipInfo
	Record    *models.VenueEntryOrExitRecords
}

// 购买次卡/月卡/季卡/年卡 请求参数
type PurchaseVipCardParam struct {
	ProductId    int64  `binding:"required" json:"product_id"`   // 商品ID
	Count        int    `binding:"required" json:"count"`        // 购买数量
	UserId       string `json:"user_id"`
	VenueId      int64  `binding:"required" json:"venue_id" `    // 场馆id
	ChannelId    int    `json:"channel_id"`                      // android/ios
}

// 场馆商品
type VenueProduct struct {
	Id                int64  `json:"id"`                   // 商品id
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
		Recommend: new(models.VenueRecommendConf),
		Product: new(models.VenueProductInfo),
		Vip: new(models.VenueVipInfo),
		Record: new(models.VenueEntryOrExitRecords),
		Engine: engine,
	}
}

// 通过场馆id 获取场馆信息
func (m *VenueModel) GetVenueInfoById(id string) (bool, error) {
	m.Venue = new(models.VenueInfo)
	return m.Engine.Where("id=?", id).Get(m.Venue)
}

// 获取场馆列表
func (m *VenueModel) GetVenueList() ([]*models.VenueInfo, error) {
	var list []*models.VenueInfo
	if err := m.Engine.Where("status=0").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 通过场馆id 获取线上场馆商品列表
func (m *VenueModel) GetVenueProducts() ([]*models.VenueProductInfo, error) {
	var list []*models.VenueProductInfo
	if err := m.Engine.Where("venue_id=? AND instance_type=1", m.Venue.Id).Asc("product_type").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 通过id获取商品
func (m *VenueModel) GetVenueProductById(id string) (bool, error) {
	m.Product = new(models.VenueProductInfo)
	return m.Engine.Where("id=?", id).Get(m.Product)
}

// 通过商品类型 获取 商品信息
func (m *VenueModel) GetVenueProductByType(productType int) (bool, error) {
	m.Product = new(models.VenueProductInfo)
	return m.Engine.Where("product_type=?", productType).Get(m.Product)
}

// 通过id获取推荐信息配置
func (m *VenueModel) GetRecommendInfoById(id string) (bool, error) {
	m.Recommend = new(models.VenueRecommendConf)
	return m.Engine.Where("id=? AND status=0", id).Get(m.Recommend)
}

// 获取场馆会员信息
func (m *VenueModel) GetVenueVipInfo(userId string, venueId int64) (bool, error) {
	m.Vip = new(models.VenueVipInfo)
	return m.Engine.Where("user_id=? AND venue_id=?", userId, venueId).Get(m.Vip)
}

// 添加场馆会员数据
func (m *VenueModel) AddVenueVipInfo() (int64, error) {
	return m.Engine.InsertOne(m.Vip)
}

// 更新场馆会员数据
func (m *VenueModel) UpdateVenueVipInfo(cols string) (int64, error) {
	return m.Engine.Where("id=?", m.Vip.Id).Cols(cols).Update(m.Vip)
}

// 场馆进出场记录
func (m *VenueModel) VenueEntryOrExitRecords(userId string, offset, size int) ([]*models.VenueEntryOrExitRecords, error) {
	var list []*models.VenueEntryOrExitRecords
	if err := m.Engine.Where("user_id=? AND status=0", userId).Desc("id").Limit(size, offset).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

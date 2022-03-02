package mvenue

import "sports_service/server/models"

type AddMarkParam struct {
	Conf []*models.VenueRecommendConf  `json:"conf"`
}

// 添加店长
type VenueAdminParam struct {
	Mobile    int64  `json:"mobile" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Status    int    `json:"status"`
	VenueId   int64  `json:"venue_id"`
}

type DelMarkParam struct {
	Ids  []int    `json:"ids"`
}

// 通过场馆id 获取场馆所有商品
func (m *VenueModel) GetVenueAllProduct() ([]*models.VenueProductInfo, error) {
	var list []*models.VenueProductInfo
	if err := m.Engine.Where("venue_id=?", m.Venue.Id).Asc("product_type").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

func (m *VenueModel) UpdateVenueInfo(info *models.VenueInfo) (int64, error) {
	return m.Engine.Where("id=?", info.Id).Update(info)
}

func (m *VenueModel) AddVenueInfo(info *models.VenueInfo) (int64, error) {
	return m.Engine.InsertOne(info)
}

// 添加场馆角标配置
func (m *VenueModel) AddMark(infos []*models.VenueRecommendConf) (int64, error) {
	return m.Engine.InsertMulti(infos)
}

func (m *VenueModel) DelMark(ids []int) (int64, error) {
	return m.Engine.In("id", ids).Delete(&models.VenueRecommendConf{})
}

func (m *VenueModel) MarkList(venueId string) ([]*models.VenueRecommendConf, error) {
	var list []*models.VenueRecommendConf
	if err := m.Engine.Where("venue_id=? AND status=0", venueId).Find(&list); err != nil {
		return nil, err
	}
	
	return list, nil
}

func (m *VenueModel) AddVenueManager(admin *models.VenueAdministrator) (int64, error) {
	return m.Engine.InsertOne(admin)
}

func (m *VenueModel) UpdateVenueManager(admin *models.VenueAdministrator) (int64, error) {
	return m.Engine.Table(&models.VenueAdministrator{}).Where("id=?", admin.Id).Update(admin)
}

// 店长列表
func (m *VenueModel) VenueManagerList(offset, size int) ([]*models.VenueAdministrator, error) {
	var list []*models.VenueAdministrator
	if err := m.Engine.Where("roles=?", "ROLE_ADMIN").Limit(size, offset).Find(&list); err != nil {
		return nil, err
	}
	
	return list, nil
}

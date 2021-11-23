package mvenue

import "sports_service/server/models"

type AddMarkParam struct {
	Conf []*models.VenueRecommendConf  `json:"conf"`
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
	return m.Engine.Update(info)
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
	if err := m.Engine.Where("venue_id=?", venueId).Find(&list); err != nil {
		return nil, err
	}
	
	return list, nil
}

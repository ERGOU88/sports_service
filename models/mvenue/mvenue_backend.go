package mvenue

import "sports_service/server/models"

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

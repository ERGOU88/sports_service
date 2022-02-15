package mshop

import (
	"sports_service/server/models"
)

type AddOrEditCategorySpecReq struct {
	CategoryId     int64    `json:"category_id"`
	SpecInfo       []Spec   `json:"spec_info"`
}

func (m *ShopModel) AddProductCategory(info *models.ProductCategory) (int64, error) {
	return m.Engine.InsertOne(info)
}

func (m *ShopModel) UpdateProductCategory(info *models.ProductCategory) (int64, error) {
	return m.Engine.Where("category_id=?", info.CategoryId).Update(info)
}

func (m *ShopModel) GetServiceList() ([]models.ShopServiceConf, error) {
	var list []models.ShopServiceConf
	if err := m.Engine.Find(&list); err != nil {
		return list, err
	}
	
	return list, nil
}

func (m *ShopModel) AddService(info *models.ShopServiceConf) (int64, error) {
	return m.Engine.InsertOne(info)
}

func (m *ShopModel) UpdateService(info *models.ShopServiceConf) (int64, error) {
	return m.Engine.Where("id=?", info.Id).Update(info)
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

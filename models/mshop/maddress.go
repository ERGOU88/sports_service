package mshop

import (
	"sports_service/server/models"
	"errors"
)

func (m *ShopModel) AddUserAddr(addr *models.UserAddress) (int64, error) {
	return m.Engine.InsertOne(addr)
}

func (m *ShopModel) UpdateUserAddr(addr *models.UserAddress) (int64, error) {
	return m.Engine.Where("id=?", addr.Id).Update(addr)
}

func (m *ShopModel) UpdateUserDefaultAddr(id, userId string, isDefault int) (int64, error) {
	if id == "" && userId == "" {
		return 0, errors.New("invalid param")
	}

	engine := m.Engine
	if id != "" {
		engine.Where("id=?", id)
	}

	if userId != "" {
		engine.Where("user_id=?", userId)
	}

	info := &models.UserAddress{
		IsDefault: isDefault,
	}

	return engine.Cols("is_default").Update(info)
}

func (m *ShopModel) GetUserAddr(id, userId string) (*models.UserAddress, error) {
	addr := &models.UserAddress{}
	engine := m.Engine
	if id != "" {
		engine.Where("id=?", id)
	}
	
	if userId != "" {
		engine.Where("user_id=?", userId)
	}
	
	ok, err := engine.Desc("is_default").Get(addr)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("addr not found")
	}

	return addr, nil
}

func (m *ShopModel) GetUserAddrByUserId(userId string, offset, size int) ([]models.UserAddress, error) {
	var list []models.UserAddress
	if err := m.Engine.Where("user_id=?", userId).Limit(size, offset).Find(&list); err != nil {
		return list, err
	}

	return list, nil
}

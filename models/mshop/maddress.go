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
	if id != "" {
		m.Engine.Where("id=?", id)
	}

	if userId != "" {
		m.Engine.Where("user_id=?", userId)
	}

	info := &models.UserAddress{
		IsDefault: isDefault,
	}

	return m.Engine.Update(info)
}

func (m *ShopModel) GetUserAddrById(id string) (*models.UserAddress, error) {
	addr := &models.UserAddress{}
	ok, err := m.Engine.Where("id=?", id).Get(addr)
	if err != nil {
		return addr, err
	}

	if !ok {
		return addr, errors.New("get addr fail")
	}

	return addr, nil
}

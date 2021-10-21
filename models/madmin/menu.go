package madmin

import "sports_service/server/models"

func (m *AdminModel) GetMenuList(offset, size int) ([]*models.SystemMenu, error){
	var list []*models.SystemMenu
	if err := m.Engine.Limit(size, offset).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

func (m *AdminModel) AddMenu(menu *models.SystemMenu) (int64, error) {
	return m.Engine.InsertOne(menu)
}

func (m *AdminModel) UpdateMenu(menu *models.SystemMenu) (int64, error) {
	return m.Engine.Update(menu)
}

func (m *AdminModel) DelMenu(menuId string) (int64, error) {
	return m.Engine.Where("menu_id=?", menuId).Delete(m.Menu)
}

func (m *AdminModel) GetMenu(menuId string) (bool, error) {
	return m.Engine.Where("menu_id=?", menuId).Get(m.Menu)
}

func (m *AdminModel) GetMenuByIds(menuIds []int) ([]*models.SystemMenu, error) {
	var list []*models.SystemMenu
	if err := m.Engine.In("menu_id", menuIds).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

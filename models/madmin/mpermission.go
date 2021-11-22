package madmin

import "sports_service/server/models"

// 获取角色可查看的菜单
func (m *AdminModel) GetRoleMenu(roleId string) ([]*models.SystemRoleMenu, error) {
	var list []*models.SystemRoleMenu
	if err := m.Engine.Where("role_id=?", roleId).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

type AddRoleMenuParam struct {
	Menus      []*models.SystemRoleMenu   `json:"menus"`
}

// 添加角色可查看的菜单
func (m *AdminModel) AddRoleMenuList(list []*models.SystemRoleMenu) (int64, error) {
	return m.Engine.InsertMulti(list)
}

// 更新角色可查看的菜单
func (m *AdminModel) UpdateRoleMenu(menu *models.SystemRoleMenu) (int64, error) {
	return m.Engine.Where("role_id=? AND menu_id=?", menu.RoleId, menu.MenuId).Update(menu)
}
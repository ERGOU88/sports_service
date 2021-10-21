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

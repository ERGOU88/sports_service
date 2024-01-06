package madmin

import "sports_service/models"

func (m *AdminModel) AddRole(role *models.SystemRole) (int64, error) {
	return m.Engine.InsertOne(role)
}

func (m *AdminModel) UpdateRole(role *models.SystemRole) (int64, error) {
	return m.Engine.Update(role)
}

func (m *AdminModel) DelRole(roleId string) (int64, error) {
	return m.Engine.Where("role_id=?", roleId).Delete(&models.SystemRole{})
}

func (m *AdminModel) GetRoleList(offset, size int) ([]*models.SystemRole, error) {
	var list []*models.SystemRole
	if err := m.Engine.Where("status=0").Limit(size, offset).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

func (m *AdminModel) GetRole(roleId string) (bool, error) {
	return m.Engine.Where("role_id=?", roleId).Get(m.Role)
}

func (m *AdminModel) GetRoleByName(roleName string) (bool, error) {
	return m.Engine.Where("role_name=?", roleName).Get(m.Role)
}

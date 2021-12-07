package madmin

import "sports_service/server/models"

// 后台用户注册/登陆请求参数 （todo: 注册为测试使用）
type AdminRegOrLoginParams struct {
	UserName     string       `json:"username" binding:"required"`
	Password     string       `json:"password" binding:"required"`
}

// 禁用/解禁 管理员
type ForbidAdmin struct {
	UserName     string       `json:"username" binding:"required"`
	Status       int          `json:"status" binding:"required"`
}

// 添加后台用户
func (m *AdminModel) AddAdminUser(admin *models.SystemUser) error {
	m.User = admin
	if _, err := m.Engine.InsertOne(m.User); err != nil {
		return err
	}

	return nil
}

func (m *AdminModel) UpdateAdminUser(admin *models.SystemUser) (int64, error) {
	m.User = admin
	return m.Engine.Where("user_id=?", admin.UserId).Update(m.User)
}

// 通过用户名 查询 管理员
func (m *AdminModel) FindAdminUserByName(userName string) *models.SystemUser {
	ok, err := m.Engine.Where("username=?", userName).Get(m.User)
	if !ok || err != nil {
		return nil
	}

	return m.User
}

func (m *AdminModel) GetAdminUserList(offset, size int) ([]*models.SystemUser, error) {
	var list []*models.SystemUser
	if err := m.Engine.Limit(size, offset).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}


package madmin

import (
	"sports_service/server/models"
  "github.com/go-xorm/xorm"
)

type AdminModel struct {
  User           *models.AdminUser
  Engine         *xorm.Session
}

func NewAdminModel(engine *xorm.Session) *AdminModel {
  return &AdminModel{
    User: new(models.AdminUser),
    Engine: engine,
  }
}

// 后台用户注册/登陆请求参数 （todo: 注册为测试使用）
type AdminRegOrLoginParams struct {
  UserName     string       `json:"user_name" binding:"required"`
  Password     string       `json:"password" binding:"required"`
}

// 添加后台用户
func (m *AdminModel) AddAdminUser() error {
  if _, err := m.Engine.InsertOne(m.User); err != nil {
    return err
  }

  return nil
}

// 通过用户名 查询 管理员
func (m *AdminModel) FindAdminUserByName(userName string) *models.AdminUser {
  admin := new(models.AdminUser)
  ok, err := m.Engine.Where("user_name=?", userName).Get(admin)
  if !ok || err != nil {
    return nil
  }

  return admin
}

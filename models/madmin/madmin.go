package madmin

import (
  "sports_service/server/models"
  "github.com/go-xorm/xorm"
)

type AdminModel struct {
  User           *models.SystemUser
  Role           *models.SystemRole
  Menu           *models.SystemMenu
  RoleMenu       *models.SystemRoleMenu
  Engine         *xorm.Session
}

func NewAdminModel(engine *xorm.Session) *AdminModel {
  return &AdminModel{
    User: new(models.SystemUser),
    Role: new(models.SystemRole),
    Menu: new(models.SystemMenu),
    RoleMenu: new(models.SystemRoleMenu),
    Engine: engine,
  }
}

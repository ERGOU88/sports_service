package madmin

import (
	"github.com/go-xorm/xorm"
	"sports_service/models"
)

type AdminModel struct {
	User     *models.SystemUser
	Role     *models.SystemRole
	Menu     *models.SystemMenu
	RoleMenu *models.SystemRoleMenu
	Engine   *xorm.Session
}

func NewAdminModel(engine *xorm.Session) *AdminModel {
	return &AdminModel{
		User:     new(models.SystemUser),
		Role:     new(models.SystemRole),
		Menu:     new(models.SystemMenu),
		RoleMenu: new(models.SystemRoleMenu),
		Engine:   engine,
	}
}

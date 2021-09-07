package cinformation

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/models"
	"sports_service/server/models/minformation"
	"sports_service/server/models/muser"
)

type InformationModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	information *minformation.InformationModel
}

func New(c *gin.Context) InformationModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	return InformationModule{
		context: c,
		user: muser.NewUserModel(socket),
		information: minformation.NewInformationModel(socket),
		engine: socket,
	}
}

// 获取资讯列表
func (svc *InformationModule) GetInformationList(page, size int) (int, []*models.Information) {
	offset := (page - 1) * size
	list, err := svc.information.GetInformationList(offset, size)
	if err != nil {
		return errdef.ERROR, []*models.Information{}
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*models.Information{}
	}

	return errdef.SUCCESS, list
}

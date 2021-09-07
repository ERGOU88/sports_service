package minformation

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

type InformationModel struct {
	Engine      *xorm.Session
	Information *models.Information
}

func NewInformationModel(engine *xorm.Session) *InformationModel {
	return &InformationModel{
		Information: new(models.Information),
		Engine: engine,
	}
}

// 获取资讯列表
func (m *InformationModel) GetInformationList(offset, size int) ([]*models.Information, error) {
	var list []*models.Information
	if err := m.Engine.Where("status=0").Limit(offset, size).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

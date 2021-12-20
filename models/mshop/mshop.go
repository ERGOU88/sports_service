package mshop

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

type ShopModel struct {
	Engine    *xorm.Session
	Category  *models.ProductCategory

}

func NewShop(engine *xorm.Session) *ShopModel {
	return &ShopModel{
		Engine: engine,
	}
}

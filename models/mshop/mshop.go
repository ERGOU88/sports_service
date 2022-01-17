package mshop

import (
	"github.com/go-xorm/xorm"
)

type ShopModel struct {
	Engine    *xorm.Session
}

func NewShop(engine *xorm.Session) *ShopModel {
	return &ShopModel{
		Engine: engine,
	}
}

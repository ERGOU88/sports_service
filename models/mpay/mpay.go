package mpay

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

type PayModel struct {
	Engine           *xorm.Session
	PayChannel       *models.VenuePaymentChannel
}

func NewPayModel(engine *xorm.Session) *PayModel {
	return &PayModel{
		Engine: engine,
		PayChannel: new(models.VenuePaymentChannel),
	}
}

// 获取支付渠道配置 1 支付宝 2 微信
func (m *PayModel) GetPaymentChannel(payType int) (bool, error) {
	return m.Engine.Where("status=0 AND pay_type=?", payType).Get(m.PayChannel)
}


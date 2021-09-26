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

// 通过id获取支付渠道配置
func (m *PayModel) GetPaymentChannel(payChannelId int) (bool, error) {
	return m.Engine.Where("status=0 AND id=?", payChannelId).Get(m.PayChannel)
}


package cvenue

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/models"
	"sports_service/server/models/morder"
	"sports_service/server/models/mvenue"
)

type VenueModule struct {
	context     *gin.Context
	engine      *xorm.Session
	venue       *mvenue.VenueModel
	order       *morder.OrderModel
}

func New(c *gin.Context) *VenueModule {
	socket := dao.Engine.NewSession()
	defer socket.Close()
	return &VenueModule{
		context: c,
		venue: mvenue.NewVenueModel(socket),
		order: morder.NewOrderModel(socket),
		engine:  socket,
	}
}

// 获取首页数据
func (svc *VenueModule) GetHomePageInfo(venueId int64) (int, *models.VenueInfo, []*mvenue.VenueProduct) {
	venueInfo, err := svc.GetVenueInfo(venueId)
	if err != nil {
		return errdef.ERROR, nil, nil
	}

	if venueInfo == nil {
		return errdef.ERROR, nil, nil
	}

	productInfo, err := svc.GetVenueProducts(venueId)
	if err != nil {
		return errdef.ERROR, venueInfo, productInfo
	}

	return errdef.SUCCESS, venueInfo, productInfo
}

// 获取场馆信息
func (svc *VenueModule) GetVenueInfo(venueId int64) (*models.VenueInfo, error) {
	svc.venue.Venue.Id = venueId
	if err := svc.venue.GetVenueInfoById(); err != nil {
		return nil, err
	}

	return svc.venue.Venue, nil
}

// 获取场馆商品[月卡、年卡、体验卡 等]
func (svc *VenueModule) GetVenueProducts(venueId int64) ([]*mvenue.VenueProduct, error) {
	svc.venue.Venue.Id = venueId
	list, err := svc.venue.GetVenueProducts()
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return []*mvenue.VenueProduct{}, nil
	}

	res := make([]*mvenue.VenueProduct, len(list))
	for index, val := range list {
		info := &mvenue.VenueProduct{
			Id:  val.Id,
			Icon: val.Icon,
			ProductName: val.ProductName,
			ProductType: val.ProductType,
			EffectiveDuration: val.EffectiveDuration,
			Describe: val.Describe,
			Title: val.Title,
			Image: val.Image,
			RealAmount: val.RealAmount,
			CurAmount: val.CurAmount,
			VenueId: val.VenueId,
		}

		svc.order.OrderProduct.OrderType = val.ProductType
		svc.order.OrderProduct.ProductId = val.Id
		info.Sales, err = svc.order.GetSalesByProduct()
		if err != nil {
			log.Log.Errorf("venue_trace: get sales fail, err:%s", err)
		}

		// 如果定价 >= 售价 则表示有折扣
		if val.RealAmount >= val.CurAmount {
			info.HasDiscount = 1
		    info.DiscountAmount = val.DiscountAmount
			info.DiscountRate = val.DiscountRate
		}

		res[index] = info
	}

	return res, nil
}

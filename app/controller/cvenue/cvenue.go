package cvenue

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/models"
	"sports_service/server/models/mvenue"
)

type VenueModule struct {
	context     *gin.Context
	engine      *xorm.Session
	venue       *mvenue.VenueModel
}

func NewVenue(c *gin.Context) *VenueModule {
	socket := dao.Engine.NewSession()
	defer socket.Close()
	return &VenueModule{
		context: c,
		engine:  socket,
	}
}

// 获取场馆信息
func (svc *VenueModule) GetVenueInfo(venueId int64) (int, *models.VenueInfo) {
	svc.venue.Venue.Id = venueId
	if err := svc.venue.GetVenueInfoById(); err != nil {
		return errdef.ERROR, svc.venue.Venue
	}

	return errdef.SUCCESS, svc.venue.Venue
}

// 获取场馆商品[月卡、年卡、体验卡 等]
func (svc *VenueModule) GetVenueProduct(venueId int64) (int, []*mvenue.VenueProduct) {
	svc.venue.Venue.Id = venueId
	list, err := svc.venue.GetVenueProducts()
	if err != nil {
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*mvenue.VenueProduct{}
	}

	//res := make([]*mvenue.VenueProduct, len(list))
	//for _, val := range list {
	//	info := &mvenue.VenueProduct{
	//		Id:  val.Id,
	//		Icon: val.Icon,
	//		ProductName: val.ProductName,
	//		ProductType: val.ProductType,
	//		EffectiveDuration: val.EffectiveDuration,
	//		Describe: val.Describe,
	//		Title: val.Title,
	//		Image: val.Image,
	//		RealAmount: val.RealAmount,
	//		CurAmount: val.CurAmount,
	//		DiscountAmount: val.DiscountAmount,
	//		DiscountRate: val.DiscountRate,
	//		VenueId: val.VenueId,
	//	}
	//
	//
	//
	//}

	return errdef.SUCCESS, nil
}

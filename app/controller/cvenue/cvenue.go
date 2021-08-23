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
	"sports_service/server/util"
	"strconv"
	"fmt"
)

type VenueModule struct {
	context     *gin.Context
	engine      *xorm.Session
	venue       *mvenue.VenueModel
	order       *morder.OrderModel
}

type VenueInfoRes struct {
	Id            int64    `json:"id"`
	VenueName     string   `json:"venue_name"`
	Address       string   `json:"address"`
	Describe      string   `json:"describe"`
	Telephone     string   `json:"telephone"`
	VenueImages   []string `json:"venue_images"`
	BusinessHours string   `json:"business_hours"`
	Services      string   `json:"services"`
	Longitude     float64  `json:"longitude"`
	Latitude      float64  `json:"latitude"`
	Status        int      `json:"status"`
	ImageNum      int      `json:"image_num"`
}

func New(c *gin.Context) *VenueModule {
	venueSocket := dao.VenueEngine.NewSession()
	defer venueSocket.Close()
	return &VenueModule{
		context: c,
		venue:   mvenue.NewVenueModel(venueSocket),
		order:   morder.NewOrderModel(venueSocket),
		engine:  venueSocket,
	}
}

// 获取首页数据
func (svc *VenueModule) GetHomePageInfo(venueId int64) (int, *VenueInfoRes, []*mvenue.VenueProduct) {
	venueInfo, err := svc.GetVenueInfo(fmt.Sprint(venueId))
	if err != nil {
		log.Log.Errorf("venue_trace: get venue info fail, err:%s", err)
		return errdef.ERROR, nil, nil
	}

	if venueInfo == nil {
		log.Log.Errorf("venue_trace: get venue info fail, err:%s", err)
		return errdef.ERROR, nil, nil
	}

	res := &VenueInfoRes{
		Id: venueInfo.Id,
		VenueName: venueInfo.VenueName,
		Address: venueInfo.Address,
		Describe: venueInfo.Describe,
		Telephone: venueInfo.Telephone,
		BusinessHours: venueInfo.BusinessHours,
		Services: venueInfo.Services,
		Status: venueInfo.ServiceStatus,
	}

	if err = util.JsonFast.UnmarshalFromString(venueInfo.VenueImages, &res.VenueImages); err != nil {
		log.Log.Errorf("venue_trace: image unmarshal fail, err:%s", err)
	}

	res.ImageNum = len(res.VenueImages)

	if venueInfo.Latitude != "" && venueInfo.Longitude != "" {
		res.Longitude, err = strconv.ParseFloat(venueInfo.Longitude, 64)
		if err != nil {
			log.Log.Errorf("venue_trace: parse float fail, err:%s", err)
		}

		res.Latitude, err = strconv.ParseFloat(venueInfo.Latitude, 64)
		if err != nil {
			log.Log.Errorf("venue_trace: parse float fail, err:%s", err)
		}
	}


	productInfo, err := svc.GetVenueProducts(venueId)
	if err != nil {
		log.Log.Errorf("venue_trace: get venue products fail, err:%s", err)
		return errdef.ERROR, res, productInfo
	}

	return errdef.SUCCESS, res, productInfo
}

// 获取场馆信息
func (svc *VenueModule) GetVenueInfo(id string) (*models.VenueInfo, error) {
	ok, err := svc.venue.GetVenueInfoById(id)
	if !ok || err != nil {
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

		svc.order.OrderProduct.ProductType = val.ProductType
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

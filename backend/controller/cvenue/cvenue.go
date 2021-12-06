package cvenue

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/models"
	"sports_service/server/models/morder"
	"sports_service/server/models/mvenue"
	"sports_service/server/global/backend/log"
	"fmt"
	"sports_service/server/tools/tencentCloud"
	"sports_service/server/util"
	"time"
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
	VenueImages   []tencentCloud.BucketURI `json:"venue_images"`
	BusinessHours string   `json:"business_hours"`
	Services      string   `json:"services"`
	Longitude     float64  `json:"longitude"`
	Latitude      float64  `json:"latitude"`
	ImageNum      int      `json:"image_num"`
	Instructions  string   `json:"instructions"`
	PromotionPic  tencentCloud.BucketURI   `json:"promotion_pic"`
	ProductNum    int      `json:"product_num"`    // 商品数量
	TotalSales    int64    `json:"total_sales"`    // 销售总额
	OrderNum      int64    `json:"order_num"`      // 订单数量（成功订单）
	TotalRefund   int64    `json:"total_refund"`   // 退款总额
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

// 获取场馆列表
func (svc *VenueModule) GetVenueList() (int, []*VenueInfoRes) {
	list, err := svc.venue.GetVenueList()
	if err != nil {
		return errdef.ERROR, nil
	}

	if list == nil {
		return errdef.SUCCESS, []*VenueInfoRes{}
	}

	res := make([]*VenueInfoRes, len(list))
	for index, item := range list {
		svc.venue.Venue.Id = item.Id
		info := &VenueInfoRes{
			Id: item.Id,
			VenueName: item.VenueName,
			Address: item.Address,
			Describe: item.Describe,
			Telephone: item.Telephone,
			BusinessHours: item.BusinessHours,
			Services: item.Services,
			Instructions: item.Instructions,
			PromotionPic: tencentCloud.BucketURI(item.PromotionPic),
		}
		
		if err = util.JsonFast.UnmarshalFromString(item.VenueImages, &info.VenueImages); err != nil {
			log.Log.Errorf("venue_trace: image unmarshal fail, err:%s", err)
		}
		
		products, err := svc.venue.GetVenueAllProduct()
		if err != nil {
			log.Log.Errorf("venue_trace: get venue product fail, venueId:%d, err:%s", item.Id, err)
		}

		info.ProductNum = len(products)
		totalSales, err := svc.order.GetTotalSalesByVenue(fmt.Sprint(item.Id))
		if err != nil {
			log.Log.Errorf("venue_trace: get total sales fail, venueId:%d, err:%s", item.Id, err)
		}

		info.TotalSales = totalSales

		orderNum, err := svc.order.GetOrderCountByVenue(fmt.Sprint(item.Id))
		if err != nil {
			log.Log.Errorf("venue_trace: get order count fail, venueId:%d, err:%s", item.Id, err)
		}

		info.OrderNum = orderNum

		totalRefund, err := svc.order.GetTotalRefundByVenue(fmt.Sprint(item.Id))
		if err != nil {
			log.Log.Errorf("venue_trace: get total refund fail, venueId:%d, err:%s", item.Id, err)
		}

		info.TotalRefund = totalRefund

		res[index] = info
	}

	return errdef.SUCCESS, res
}

// 获取场馆信息
func (svc *VenueModule) GetVenueInfo(id string) (*VenueInfoRes, error) {
	ok, err := svc.venue.GetVenueInfoById(id)
	if !ok || err != nil {
		return nil, err
	}

	info := &VenueInfoRes{
		Id: svc.venue.Venue.Id,
		VenueName: svc.venue.Venue.VenueName,
		Address: svc.venue.Venue.Address,
		Describe: svc.venue.Venue.Describe,
		Telephone: svc.venue.Venue.Telephone,
		BusinessHours: svc.venue.Venue.BusinessHours,
		Services: svc.venue.Venue.Services,
		Instructions: svc.venue.Venue.Instructions,
		PromotionPic: tencentCloud.BucketURI(svc.venue.Venue.PromotionPic),
	}

	if err = util.JsonFast.UnmarshalFromString(svc.venue.Venue.VenueImages, &info.VenueImages); err != nil {
		log.Log.Errorf("venue_trace: image unmarshal fail, err:%s", err)
	}

	return info, nil
}

func (svc *VenueModule) EditVenueInfo(info *models.VenueInfo) int {
	if _, err := svc.venue.UpdateVenueInfo(info); err != nil {
		return errdef.ERROR
	}

	return errdef.SUCCESS
}

func (svc *VenueModule) AddVenueInfo(info *models.VenueInfo) int {
	if _, err := svc.venue.AddVenueInfo(info); err != nil {
		return errdef.ERROR
	}

	return errdef.SUCCESS
}

// 更新退款费率
func (svc *VenueModule) UpdateRefundRate(param *models.VenueRefundRules) int {
	if param.FeeRate <= 0 {
		return errdef.INVALID_PARAMS
	}

	param.UpdateAt = int(time.Now().Unix())
	if _, err := svc.order.UpdateRefundRate(param); err != nil {
		return errdef.ERROR
	}

	return errdef.SUCCESS
}

// 获取退款规则
func (svc *VenueModule) GetRefundRules() (int, []*models.VenueRefundRules) {
	rules, err := svc.order.GetRefundRules()
	if err != nil {
		return errdef.ERROR, []*models.VenueRefundRules{}
	}

	return errdef.SUCCESS, rules
}

func (svc *VenueModule) AddMark(param *mvenue.AddMarkParam) int {
	if _, err := svc.venue.AddMark(param.Conf); err != nil {
		return errdef.ERROR
	}
	
	return errdef.SUCCESS
}

func (svc *VenueModule) DelMark(param *mvenue.DelMarkParam) int {
	if _, err := svc.venue.DelMark(param.Ids); err != nil {
		return errdef.ERROR
	}
	
	return errdef.SUCCESS
}

func (svc *VenueModule) MarkList(venueId string) (int, []*models.VenueRecommendConf) {
	list, err := svc.venue.MarkList(venueId)
	if err != nil {
		return errdef.ERROR, nil
	}
	
	return errdef.SUCCESS, list
}

func (svc *VenueModule) AddStoreManager(admin *mvenue.VenueAdminParam) int {
	now := int(time.Now().Unix())
	pwd, err := util.GenPassword(admin.Password)
	if err != nil {
		return errdef.ERROR
	}
	
	info := &models.VenueAdministrator{
		Mobile: admin.Mobile,
		Name: admin.Name,
		Username: admin.Username,
		Password: pwd,
		Status: admin.Status,
		Roles: "ROLE_ADMIN",
		VenueId: admin.VenueId,
		CreateAt: now,
		UpdateAt: now,
	}
	
	if _, err := svc.venue.AddVenueManager(info); err != nil {
		log.Log.Errorf("venue_trace: add manager fail, err:%s", err)
		return errdef.ERROR
	}
	
	return errdef.SUCCESS
}

func (svc *VenueModule) EditStoreManage(admin *models.VenueAdministrator) int {
	admin.UpdateAt = int(time.Now().Unix())
	pwd, err := util.GenPassword(admin.Password)
	if err != nil {
		return errdef.ERROR
	}
	
	admin.Password = pwd
	
	if _, err := svc.venue.UpdateVenueManager(admin); err != nil {
		return errdef.ERROR
	}
	
	return errdef.SUCCESS
}

// 店长列表
func (svc *VenueModule) StoreManageList(page, size int) (int, []*models.VenueAdministrator) {
	offset := (page - 1) * size
	list, err := svc.venue.VenueManagerList(offset, size)
	if err != nil {
		return errdef.ERROR, nil
	}
	
	if len(list) == 0 {
		return errdef.SUCCESS, []*models.VenueAdministrator{}
	}
	
	return errdef.SUCCESS, list
}

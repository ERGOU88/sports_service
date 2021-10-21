package cvenue

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"errors"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/morder"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvenue"
	"sports_service/server/util"
	"strconv"
	"fmt"
	"strings"
	"time"
)

type VenueModule struct {
	context     *gin.Context
	engine      *xorm.Session
	venue       *mvenue.VenueModel
	order       *morder.OrderModel
	user        *muser.UserModel
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
	appSocket := dao.AppEngine.NewSession()
	defer appSocket.Close()
	return &VenueModule{
		context: c,
		venue:   mvenue.NewVenueModel(venueSocket),
		order:   morder.NewOrderModel(venueSocket),
		user: muser.NewUserModel(appSocket),
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
		Status: svc.IsOpen(time.Now().Format(consts.FORMAT_DATE), venueInfo.BusinessHours),
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

// 是否营业 0 表示营业中
func (svc *VenueModule) IsOpen(date, timeNode string) int {
	nodes := strings.Split(timeNode, "-")
	if len(nodes) == 2 {
		now := time.Now().Unix()
		// 场馆营业 开始时间 及 结束时间
		start := fmt.Sprintf("%s %s", date, nodes[0])
		end := fmt.Sprintf("%s %s", date, nodes[1])
		ts := new(util.TimeS)
		startTm := ts.GetTimeStrOrStamp(start, "YmdHi")
		endTm := ts.GetTimeStrOrStamp(end, "YmdHi")
		// 当前时间 > 开始时间 且 < 结束时间 营业中
		if now > startTm.(int64) && now < endTm.(int64) {
			return 0
		}
	}

	return 1
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
			Instructions: val.Instructions,
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

// 购买会员卡
func (svc *VenueModule) PurchaseVipCard(param *mvenue.PurchaseVipCardParam) (int, *mappointment.OrderResp) {
	if param.Count <= 0 {
		log.Log.Errorf("venue_trace: invalid count, count:%d", param.Count)
		return errdef.INVALID_PARAMS, nil
	}

	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("venue_trace: session begin fail, err:%s", err )
		return errdef.ERROR, nil
	}

	ok, err := svc.venue.GetVenueInfoById(fmt.Sprint(param.VenueId))
	if !ok || err != nil {
		log.Log.Errorf("venue_trace: venue not found, venueId:%s", param.VenueId)
		return errdef.VENUE_NOT_EXISTS, nil
	}

	user := svc.user.FindUserByUserid(param.UserId)
	if user == nil {
		log.Log.Errorf("venue_trace: user not found, userId:%s", param.UserId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS, nil
	}

	ok, err = svc.venue.GetVenueProductById(fmt.Sprint(param.ProductId))
	if !ok || err != nil {
		log.Log.Errorf("venue_trace: product not found, productId:%s", param.ProductId)
		svc.engine.Rollback()
		return errdef.VENUE_PRODUCT_NOT_EXIST, nil
	}

	totalAmount := svc.venue.Product.CurAmount * param.Count
	orderId := util.NewOrderId()
	now := int(time.Now().Unix())

	// 添加购买的会员卡记录
	if err := svc.AddVipCardRecord(orderId, param.UserId, param.VenueId, now, param.Count); err != nil {
		log.Log.Errorf("venue_trace: add card record fail, orderId:%s, err:%s", orderId, err)
		svc.engine.Rollback()
		return errdef.ORDER_ADD_CARD_RECORD_FAIL, nil
	}

	// 添加订单商品流水
	if err := svc.AddOrderProduct(orderId, now, param.Count); err != nil {
		log.Log.Errorf("venue_trace: add order products fail, err:%s", err)
		svc.engine.Rollback()
		return errdef.ORDER_PRODUCT_ADD_FAIL, nil
	}

	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	extra := &mappointment.OrderResp{
		OrderId:        orderId,
		CreateAt:       time.Unix(int64(now), 0).In(cstSh).Format(consts.FORMAT_TM),
		Count:          param.Count,
		Id:             param.ProductId,
		MobileNum:      util.HideMobileNum(fmt.Sprint(svc.user.User.MobileNum)),
		Name:           svc.venue.Product.ProductName,
		PayDuration:    consts.PAYMENT_DURATION,
		TotalAmount:    svc.order.OrderProduct.Amount,
		ExpireDuration: svc.venue.Product.ExpireDuration,
		VenueId:        param.VenueId,
		OrderType:      svc.venue.Product.ProductType,
		VenueName:      svc.venue.Venue.VenueName,
		ProductImg:     svc.venue.Product.Image,
		OriginalAmount: svc.venue.Product.RealAmount * param.Count,
	}

	// 添加订单
	if err := svc.AddOrder(extra, orderId, param.UserId, svc.venue.Product.ProductName, now, svc.venue.Product.ProductType,
		totalAmount, param.ChannelId); err != nil {
		log.Log.Errorf("venue_trace: add order fail, err:%s", err)
		svc.engine.Rollback()
		return errdef.ORDER_ADD_FAIL, nil
	}

	// 记录需处理支付超时的订单
	if _, err := svc.order.RecordOrderId(orderId); err != nil {
		log.Log.Errorf("venue_trace: record orderId fail, err:%s", err)
		svc.engine.Rollback()
		return errdef.ORDER_ADD_FAIL, nil
	}

	svc.engine.Commit()
	return errdef.SUCCESS, extra
}

// 添加购买会员卡记录
func (svc *VenueModule) AddVipCardRecord(orderId, userId string, venueId int64, now, count int) error {
	svc.order.CardRecord.ProductType = svc.venue.Product.ProductType
	svc.order.CardRecord.VenueId = venueId
	svc.order.CardRecord.PayOrderId = orderId
	svc.order.CardRecord.UserId = userId
	svc.order.CardRecord.SingleDuration = svc.venue.Product.EffectiveDuration
	svc.order.CardRecord.ExpireDuration = svc.venue.Product.ExpireDuration
	svc.order.CardRecord.Duration = svc.venue.Product.EffectiveDuration * count
	svc.order.CardRecord.CreateAt = now
	svc.order.CardRecord.UpdateAt = now
	svc.order.CardRecord.PurchasedNum = count
	affected, err := svc.order.AddVipCardRecord()
	if err != nil {
		return err
	}

	if affected != 1 {
		return errors.New("add vip card record fail, affected not 1")
	}

	return nil
}

// 添加订单商品流水
func (svc *VenueModule) AddOrderProduct(orderId string, now, count int) error {
	svc.order.OrderProduct.ProductId = svc.venue.Product.Id
	svc.order.OrderProduct.ProductType = svc.venue.Product.ProductType
	svc.order.OrderProduct.Count = count
	svc.order.OrderProduct.ProductCategory = consts.PRODUCT_CATEGORY_CARD
	svc.order.OrderProduct.RealAmount = svc.venue.Product.RealAmount
	svc.order.OrderProduct.CurAmount = svc.venue.Product.CurAmount
	svc.order.OrderProduct.DiscountRate = svc.venue.Product.DiscountRate
	svc.order.OrderProduct.DiscountAmount = svc.venue.Product.DiscountAmount
	svc.order.OrderProduct.Amount = svc.venue.Product.CurAmount * count
	svc.order.OrderProduct.CreateAt = now
	svc.order.OrderProduct.UpdateAt = now
	svc.order.OrderProduct.PayOrderId = orderId
	svc.order.OrderProduct.SnapshotId = svc.order.CardRecord.Id
	svc.order.OrderProduct.VenueId = svc.venue.Venue.Id

	affected, err := svc.order.AddOrderProduct()
	if err != nil {
		return err
	}

	if affected != 1 {
		return errors.New("add order product fail, affected not 1")
	}

	return nil
}

// 添加订单
func (svc *VenueModule) AddOrder(extra *mappointment.OrderResp, orderId, userId, subject string, now, productType,
	totalAmount, channel int) error {

	bts, _ := util.JsonFast.Marshal(extra)
	svc.order.Order.Extra = string(bts)
	svc.order.Order.PayOrderId = orderId
	svc.order.Order.UserId = userId
	svc.order.Order.OrderType = 1001
	svc.order.Order.CreateAt = now
	svc.order.Order.UpdateAt = now
	svc.order.Order.Amount = totalAmount
	svc.order.Order.ChannelId = channel
	svc.order.Order.Subject = subject
	svc.order.Order.ProductType = productType
	svc.order.Order.VenueId = svc.venue.Venue.Id
	svc.order.Order.OriginalAmount = extra.OriginalAmount
	// 次卡需要核销

	affected, err := svc.order.AddOrder()
	if err != nil {
		return err
	}

	if affected != 1 {
		log.Log.Errorf("venue_trace: add order fail, err:%s", err)
		return errors.New("add order fail, affected not 1")
	}

	return nil
}

// 获取场馆用户进出场记录
func (svc *VenueModule) GetActionRecord(userId string, page, size int) (int, []mvenue.VenueEntryOrExitRecords) {
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("venue_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS, nil
	}

	offset := (page - 1) * size
	list, err := svc.venue.VenueEntryOrExitRecords(userId, offset, size)
	if err != nil {
		log.Log.Errorf("venue_trace: get action record fail, userId:%s", userId)
		return errdef.VENUE_ACTION_RECORD_FAIL, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []mvenue.VenueEntryOrExitRecords{}
	}

	res := make([]mvenue.VenueEntryOrExitRecords, 0)
	for _, item := range list {
		if item.DtEnter == 0 && item.DtExit == 0 {
			continue
		}

		info := mvenue.VenueEntryOrExitRecords{}
		info.Id = item.Id
		info.UserId = item.UserId
		info.VenueId = item.VenueId
		info.Status = 0
		ok, err := svc.venue.GetVenueInfoById(fmt.Sprint(item.VenueId))
		if ok && err == nil {
			info.VenueName = svc.venue.Venue.VenueName
		}

		if item.DtExit > 0 {
			info.ActionType = 2
			info.CreateAt = item.DtExit
			info.UpdateAt = item.DtExit
			res = append(res, info)
		}

		if item.DtEnter > 0 {
			info.ActionType = 1
			info.CreateAt = item.DtEnter
			info.UpdateAt = item.DtEnter
			res = append(res, info)
		}
	}

	return errdef.SUCCESS, res
}

package cappointment

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/morder"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvenue"
	"sports_service/server/util"
	"time"
)

type VenueAppointmentModule struct {
	context         *gin.Context
	engine          *xorm.Session
	user            *muser.UserModel
	venue           *mvenue.VenueModel
	order           *morder.OrderModel
	*base
}

func NewVenue(c *gin.Context) *VenueAppointmentModule {
	venueSocket := dao.VenueEngine.NewSession()
	defer venueSocket.Close()
	appSocket := dao.AppEngine.NewSession()
	defer appSocket.Close()
	return &VenueAppointmentModule{
		context: c,
		user:    muser.NewUserModel(appSocket),
		venue:   mvenue.NewVenueModel(venueSocket),
		order:   morder.NewOrderModel(venueSocket),
		engine:  venueSocket,
		base:    New(venueSocket),
	}
}

// 场馆选项 todo：暂时只有一个场馆
func (svc *VenueAppointmentModule) Options(relatedId int64) (int, interface{}) {
	list, err := svc.venue.GetVenueList()
	if err != nil {
		return errdef.ERROR, nil
	}

	if list == nil {
		return errdef.SUCCESS, []interface{}{}
	}

	res := make([]*mappointment.Options, len(list))
	for index, item := range list {
		info := &mappointment.Options{
			Id: item.Id,
			Name: item.VenueName,
		}

		res[index] = info
	}

	return errdef.SUCCESS, res
}

// todo:
// 预约场馆
// 1 判断库存数据是否存在 不存在写入 如并发情况 遇写入失败[mysql err 1062 唯一索引约束错误] 说明已有用户成功预约 同时已写入库存数据
// 2 进行库存更新 判断库存是否足够 库存不够 直接返回各节点最新库存量
// 3 如库存均足够 则判断是否充值会员时长
// 4 如充值会员时长 预约的时间 按价格从高至低 抵扣时长 且每个时间节点最多只可抵扣一次 [能抵扣则进行抵扣 预约数量-1] 并 计算抵扣时长后的 订单总价
// 5 如未充值会员时长 或 会员时长不足 则剩余预约 按售价 * 预约数量 计算订单总价
// 6 记录订单、订单商品流水、预约流水
func (svc *VenueAppointmentModule) Appointment(params *mappointment.AppointmentReq) (int, interface{}) {
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("venue_trace: session begin fail, err:%s", err)
		return errdef.ERROR, nil
	}

	if len(params.Infos) == 0 {
		svc.engine.Rollback()
		return errdef.ERROR, nil
	}

	user := svc.user.FindUserByUserid(params.UserId)
	if user == nil {
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS, nil
	}

	list, err := svc.GetAppointmentConfByIds(params.Ids)
	if err != nil {
		svc.engine.Rollback()
		return errdef.ERROR, nil
	}

	if len(list) != len(params.Infos) {
		svc.engine.Rollback()
		return errdef.ERROR, nil
	}

	isEnough := true
	totalAmount := 0
	orderId := util.NewOrderId()
	now := int(time.Now().Unix())
	stockInfo := make([]*mappointment.StockInfoResp, 0)
	// 订单商品流水map
	orderMp := make(map[int64]*models.VenueOrderProductInfo)
	// 预约流水map
	recordMp := make(map[int64]*models.AppointmentRecord)
	for index, item := range params.Infos {
		if item.Count <= 0 || item.Count > svc.appointment.AppointmentInfo.QuotaNum {
			svc.engine.Rollback()
			return errdef.ERROR, nil
		}

		svc.venue.Venue.Id = item.RelatedId
		ok, err := svc.venue.GetVenueInfoById()
		if !ok || err != nil {
			log.Log.Errorf("venue_trace: get venue info fail, err:%s", err)
			svc.engine.Rollback()
			return errdef.ERROR, nil
		}

		if err := svc.GetAppointmentConf(item.Id); err != nil {
			svc.engine.Rollback()
			return errdef.ERROR, nil
		}

		svc.appointment.AppointmentInfo.Id = item.Id
		ok, err = svc.appointment.GetAppointmentConfById()
		if err != nil || !ok {
			svc.engine.Rollback()
			return errdef.ERROR, nil
		}

		date := svc.GetDateById(item.DateId, consts.FORMAT_DATE)
		if date == "" {
			svc.engine.Rollback()
			return errdef.ERROR, nil
		}

		orderMp[item.Id] = svc.SetOrderProductInfo(orderId, now, item)
		recordMp[item.Id] = svc.SetAppointmentRecordInfo(user.UserId, date, orderId, now, item)

		// 数量 * 售价
		totalAmount += item.Count * svc.appointment.AppointmentInfo.CurAmount
		ok, err = svc.appointment.HasExistsStockInfo()
		if err != nil {
			svc.engine.Rollback()
			return errdef.ERROR, nil
		}

		var (
			myErr *mysql.MySQLError
		)
		if !ok {
			// 库存数据不存在 写入
			data := &models.VenueAppointmentStock{
				Date:            date,
				TimeNode:        svc.appointment.AppointmentInfo.TimeNode,
				QuotaNum:        svc.appointment.AppointmentInfo.QuotaNum,
				PurchasedNum:    item.Count,
				AppointmentType: svc.appointment.AppointmentInfo.AppointmentType,
				RelatedId:       svc.appointment.AppointmentInfo.RelatedId,
				CreateAt:        now,
				UpdateAt:        now,
			}

			_, err = svc.appointment.AddStockInfo(data)
			// 是否为mysql错误
			if err != nil && errors.As(err, &myErr) {
				// 不是唯一索引约束错误 则直接返回错误
				if myErr.Number != 1062 {
					svc.engine.Rollback()
					return errdef.ERROR, nil
				}
			}
		}

		// 数据存在 或 插入时唯一索引约束错误
		if ok || myErr.Number == 1062 {
			// 更新库存
			affected, err := svc.appointment.UpdateStockInfo(svc.appointment.AppointmentInfo.TimeNode, date, item.Count, now,
				params.AppointmentType, int(item.RelatedId))
			if err != nil {
				svc.engine.Rollback()
				return errdef.ERROR, nil
			}

			// 更新未成功 库存不够 || 当前预约库存足够 但已出现库存不足的预约时间点 则需返回最新各预约节点的剩余库存
			if affected == 0 || affected == 1 && !isEnough {
				isEnough = false
				// 查看最新的库存 并返回 tips：读快照无问题 购买时保证数据一致性即可
				ok, err := svc.QueryStockInfo(params.AppointmentType, item.RelatedId, date, svc.appointment.AppointmentInfo.TimeNode)
				if !ok || err != nil {
					log.Log.Errorf("venue_trace: get stock info fail, err:%s", err)
				} else {
					info := &mappointment.StockInfoResp{
						Id: item.Id,
						Stock: svc.appointment.Stock.QuotaNum - svc.appointment.Stock.PurchasedNum,
					}

					stockInfo = append(stockInfo, info)
				}

				svc.engine.Rollback()
				continue
			}


		}

		// 是否遍历完所有预约节点
		if index < len(params.Infos) - 1 {
			continue
		}

		// 如果出现库存不足的情况
		if !isEnough {
			svc.engine.Rollback()
			return errdef.ERROR, stockInfo
		}

	}

	// 库存都足够 则判断用户是否充值会员时长
	//// 存在会员数据 则需要先抵扣时间
	vip, err := svc.appointment.GetVenueVipInfo(params.UserId)
	if err != nil {
		// 查询失败
		svc.engine.Rollback()
		return errdef.ERROR, nil
	}

	// 总抵扣时长
	totalDeductionTm := 0
	// 如果是会员
	if vip != nil {
		// 会员时长 > 0
		if vip.Duration > 0 {
			// 开始走抵扣流程 预约的时间节点[多个] 按价格从高至低 开始抵扣 每个时间节点最多只可抵扣一次
			for _, val := range list {
				if val.Duration <= 0 {
					continue
				}

				// 是否足够抵扣
				affected, err := svc.appointment.UpdateVenueVipInfo(val.Duration * -1, user.UserId)
				if err != nil {
					log.Log.Errorf("venue_trace: update vip duration fail, err:%s", err)
					svc.engine.Rollback()
					return errdef.ERROR, nil
				}

				// 会员时长不够 查看下一个预约节点 是否可抵扣
				if affected == 0 {
					continue
				}

				// 足够抵扣 则记录抵扣的记录
				if affected == 1 {
					// 抵扣一个 则 减去一个的售价
					totalAmount = totalAmount - val.CurAmount
					orderMp[val.Id].DeductionNum = affected
					orderMp[val.Id].DeductionTm = int64(val.Duration)
					orderMp[val.Id].DeductionAmount = int64(val.CurAmount)
					recordMp[val.Id].DeductionTm = int64(val.Duration)
					totalDeductionTm += val.Duration
				}
			}
		}
	}

	svc.order.Order.Extra = ""
	svc.order.Order.PayOrderId = orderId
	svc.order.Order.UserId = user.UserId
	svc.order.Order.CreateAt = now
	svc.order.Order.UpdateAt = now
	svc.order.Order.Amount = totalAmount
	affected, err := svc.order.AddOrder()
	if err != nil {
		svc.engine.Rollback()
		return errdef.ERROR, nil
	}

	if affected != 1 {
		svc.engine.Rollback()
		return errdef.ERROR, nil
	}

	var olst []*models.VenueOrderProductInfo
	for _, val := range orderMp {
		olst = append(olst, val)
	}

	// 添加订单商品流水
	affected, err = svc.order.AddMultiOrderProduct(olst)
	if err != nil {
		svc.engine.Rollback()
		return errdef.ERROR, nil
	}

	if affected != int64(len(olst)) {
		svc.engine.Rollback()
		return errdef.ERROR, nil
	}

	var rlst []*models.AppointmentRecord
	for _, val := range recordMp {
		rlst = append(rlst, val)
	}

	// 添加预约记录流水
	affected, err = svc.appointment.AddMultiAppointmentRecord(rlst)
	if err != nil {
		svc.engine.Rollback()
		return errdef.ERROR, nil
	}

	if affected != int64(len(rlst)) {
		svc.engine.Rollback()
		return errdef.ERROR, nil
	}

	svc.engine.Commit()

	// todo: 订单15分钟过期
	return errdef.SUCCESS, nil
}




// 取消预约
func (svc *VenueAppointmentModule) AppointmentCancel() int {
	return 2000
}

// 预约场馆选项
func (svc *VenueAppointmentModule) AppointmentOptions() (int, interface{}) {
	date := svc.GetDateById(svc.DateId, consts.FORMAT_DATE)
	if date == "" {
		return errdef.ERROR, nil
	}

	list, err := svc.GetAppointmentOptions()
	if err != nil {
		log.Log.Errorf("venue_trace: get options fail, err:%s", err)
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []interface{}{}
	}

	res := make([]*mappointment.OptionsInfo, 0)
	for _, item := range list {
		info := svc.SetAppointmentOptionsRes(date, item)
		if info == nil {
			continue
		}

		svc.venue.Venue.Id = item.RelatedId
		ok, err := svc.venue.GetVenueInfoById()
		if err != nil {
			log.Log.Errorf("venue_trace: get venue info by id fail, err:%s", err)
		}

		if ok {
			info.Name = svc.venue.Venue.VenueName
		}

		svc.venue.Labels.TimeNode = item.TimeNode
		svc.venue.Labels.Date = date
		svc.venue.Labels.VenueId = item.RelatedId
		labels, err := svc.venue.GetVenueUserLabels()
		if err != nil {
			log.Log.Errorf("venue_trace: get venue user lables fail, err:%s", err)
		}

		info.Labels = make([]*mappointment.LabelInfo, len(labels))
		for key, val := range labels {
			label := &mappointment.LabelInfo{
				UserId: val.UserId,
				LabelId: val.LabelId,
				LabelName: val.LabelName,
			}

			info.Labels[key] = label
		}

		svc.appointment.Record.AppointmentType = 0
		svc.appointment.Record.TimeNode = item.TimeNode
		svc.appointment.Record.RelatedId = item.RelatedId
		svc.appointment.Record.Date = date
		records, err := svc.appointment.GetAppointmentRecord()
		if err != nil {
			log.Log.Errorf("venue_trace: get appointment record fail, err:%s", err)
		}

		info.ReservedUsers = make([]*mappointment.ReservedUsers, 0)
		if len(records) > 0 {
			for _, val := range records {
				uinfo := &mappointment.ReservedUsers{
					UserId: val.UserId,
				}

				user := svc.user.FindUserByUserid(val.UserId)
				if user != nil {
					uinfo.NickName = user.NickName
					uinfo.Avatar = user.Avatar
				}

				info.ReservedUsers = append(info.ReservedUsers, uinfo)

				if val.PurchasedNum > 1 {
					for i := 0; i < val.PurchasedNum; i++ {
						info.ReservedUsers = append(info.ReservedUsers, uinfo)
					}
				}

			}
		}

		res = append(res, info)
	}


	return errdef.SUCCESS, res
}

func (svc *VenueAppointmentModule) AppointmentDetail() (int, interface{}) {
	return 4000, nil
}

// 场馆预约日期配置
func (svc *VenueAppointmentModule) AppointmentDate() (int, interface{}) {
	return errdef.SUCCESS, svc.AppointmentDateInfo(6, 0)
}


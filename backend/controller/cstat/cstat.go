package cstat

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mcommunity"
	"sports_service/server/models/morder"
	"sports_service/server/models/mposting"
	"sports_service/server/models/mstat"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"time"
)

type StatModule struct {
	context      *gin.Context
	engine       *xorm.Session
	post         *mposting.PostingModel
	user         *muser.UserModel
	video        *mvideo.VideoModel
	community    *mcommunity.CommunityModel
	stat         *mstat.StatModel
	order        *morder.OrderModel
}

func New(c *gin.Context) StatModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	venueSocket := dao.VenueEngine.NewSession()
	defer venueSocket.Close()
	return StatModule{
		context: c,
		post: mposting.NewPostingModel(socket),
		user: muser.NewUserModel(socket),
		video: mvideo.NewVideoModel(socket),
		community: mcommunity.NewCommunityModel(socket),
		stat: mstat.NewStatModel(socket),
		order: morder.NewOrderModel(venueSocket),
		engine: socket,
	}
}

// 管理后台首页统计数据
func (svc *StatModule) GetHomePageInfo(queryMinDate, queryMaxDate string) (int, mstat.HomePageInfo) {
	today := time.Now().Format(consts. FORMAT_DATE)
	// 日活
	dau, err := svc.stat.GetDAUByDate(today)
	if err != nil {
		log.Log.Errorf("stat_trace: get dau by date fail, err:%s", err)
		return errdef.ERROR, mstat.HomePageInfo{}
	}

	// 总用户数
	totalUser, err := svc.stat.GetTotalUser()
	if err != nil {
		log.Log.Errorf("stat_trace: get total by date fail, err:%s", err)
		return errdef.ERROR, mstat.HomePageInfo{}
	}

	month := time.Now().Format(consts.FORMAT_MONTH)
	// 月活
	mau, err := svc.stat.GetMAUByMonth(month)
	if err != nil {
		log.Log.Errorf("stat_trace: get mau by month fail, err:%s", err)
		return errdef.ERROR, mstat.HomePageInfo{}
	}

	newUsers, err := svc.stat.GetNetAdditionByDate(today)
	if err != nil {
		log.Log.Errorf("stat_trace: get net addition by date fail, err:%s", err)
		return errdef.ERROR, mstat.HomePageInfo{}
	}

	totalOrder, err := svc.order.GetOrderCount()
	if err != nil {
		log.Log.Errorf("stat_trace: get order count fail, err:%s", err)
		return errdef.ERROR, mstat.HomePageInfo{}
	}

	homepageInfo := mstat.HomePageInfo{
		TopInfo: make(map[string]interface{}, 0),
	}
	homepageInfo.TopInfo["dau"] = dau.Count
	homepageInfo.TopInfo["total_user"] = totalUser
	homepageInfo.TopInfo["mau"] = mau.Count
	homepageInfo.TopInfo["new_users"] = newUsers.Count
	homepageInfo.TopInfo["total_order"] = totalOrder

	homepageInfo.DauList, err = svc.stat.GetDAUByDays(100)
	if err != nil {
		log.Log.Errorf("stat_trace: get dau by days fail, err:%s", err)
		return errdef.ERROR, homepageInfo
	}

	homepageInfo.NewUserList, err = svc.stat.GetNetAdditionByDays(100)
	if err != nil {
		log.Log.Errorf("stat_trace: get new users by days fail, err:%s", err)
		return errdef.ERROR, homepageInfo
	}

	minDate := time.Now().AddDate(0, 0, -100).Format(consts.FORMAT_DATE)
	// 次日留存率
	homepageInfo.NextDayRetentionRate, err = svc.stat.GetUserRetentionRate("", minDate, today)
	if err != nil {
		log.Log.Errorf("stat_trace: get user retentionRate fail, err:%s", err)
		return errdef.ERROR, homepageInfo
	}

	minDate = time.Now().AddDate(0, 0, -100).Format(consts.FORMAT_DATE)
	maxDate := today
	if queryMinDate != "" && queryMaxDate != "" {
		minDate = queryMinDate
		maxDate = queryMaxDate
	}

	// 留存率详情
	homepageInfo.RetentionRate, err = svc.stat.GetUserRetentionRate("1", minDate, maxDate)
	return errdef.SUCCESS, homepageInfo
}

// 视频分区统计 [发布占比]
func (svc *StatModule) GetVideoSubareaStat() (int, map[string]interface{}) {
	//subareaStat, err := svc.stat.GetVideoSubareaStat()
	//if err != nil {
	//	log.Log.Errorf("stat_trace: get video subarea stat fail, err:%s", err)
	//	return errdef.ERROR, nil
	//}

	//total, err := svc.stat.GetVideoTotal()
	//if err != nil {
	//	log.Log.Errorf("stat_trace: get video total fail, err:%s", err)
	//	return errdef.ERROR, nil
	//}

	mp := make(map[string]interface{}, 0)
	mp["title"] = "视频各板块发布数据"

	//for _, item := range subareaStat {
		//item.Rate = float64(item.Count) / float64(total)
	//}

	return errdef.SUCCESS, mp
}

// 帖子分区统计 [发布占比]
func (svc *StatModule) GetPostSectionStat() {

}

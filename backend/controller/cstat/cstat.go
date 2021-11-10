package cstat

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/global/backend/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mcommunity"
	"sports_service/server/models/morder"
	"sports_service/server/models/mposting"
	"sports_service/server/models/mstat"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"sports_service/server/util"
	"time"
	"fmt"
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
	days := 11
	minDate := time.Now().AddDate(0, 0, -days).Format(consts.FORMAT_DATE)
	maxDate := time.Now().Format(consts. FORMAT_DATE)
	if queryMinDate != "" && queryMaxDate != "" {
		minDate = queryMinDate
		maxDate = queryMaxDate
		min, err := time.Parse(consts.FORMAT_DATE, queryMinDate)
		if err != nil {
			log.Log.Errorf("stat_trace: time.Parse fail, err:%s", err)
			return errdef.ERROR, mstat.HomePageInfo{}
		}

		max, err := time.Parse(consts.FORMAT_DATE, queryMaxDate)
		if err != nil {
			log.Log.Errorf("stat_trace: time.Parse fail, err:%s", err)
			return errdef.ERROR, mstat.HomePageInfo{}
		}

		days = util.GetDiffDays(max, min)
		log.Log.Infof("##########days:%d", days)
	}

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

	dailyLoyaltyUsers, err := svc.stat.GetLoyaltyUsers(today)
	if err != nil {
		log.Log.Errorf("stat_trace: get loyalty users fail, err:%s", err)
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
	homepageInfo.TopInfo["daily_loyalty_users"] = dailyLoyaltyUsers

	dauList, err := svc.stat.GetDAUByDays(days)
	if err != nil {
		log.Log.Errorf("stat_trace: get dau by days fail, err:%s", err)
		return errdef.ERROR, homepageInfo
	}

	homepageInfo.DauList =  svc.ResultInfoByDate(dauList, days)

	newUserList, err := svc.stat.GetNetAdditionByDays(days)
	if err != nil {
		log.Log.Errorf("stat_trace: get new users by days fail, err:%s", err)
		return errdef.ERROR, homepageInfo
	}
	homepageInfo.NewUserList = svc.ResultInfoByDate(newUserList, days)

	// 次日留存率
	homepageInfo.NextDayRetentionRate, err = svc.stat.GetUserRetentionRate("", minDate, today)
	if err != nil {
		log.Log.Errorf("stat_trace: get user retentionRate fail, err:%s", err)
		return errdef.ERROR, homepageInfo
	}

	// 留存率详情
	homepageInfo.RetentionRate, err = svc.stat.GetUserRetentionRate("1", minDate, maxDate)
	return errdef.SUCCESS, homepageInfo
}

// 视频分区统计 [发布占比]
func (svc *StatModule) GetVideoSubareaStat() (int, map[string]interface{}) {
	subareaStat, err := svc.stat.GetVideoSubareaStat()
	if err != nil {
		log.Log.Errorf("stat_trace: get video subarea stat fail, err:%s", err)
		return errdef.ERROR, nil
	}

	total, err := svc.stat.GetVideoTotal()
	if err != nil {
		log.Log.Errorf("stat_trace: get video total fail, err:%s", err)
		return errdef.ERROR, nil
	}

	mp := make(map[string]interface{}, 0)
	mp["title"] = "视频各分区发布数据"

	for _, item := range subareaStat {
		item.Rate = fmt.Sprintf("%.2f%s", float64(item.Count) / float64(total) * 100, "%")
		if item.Id == 0 {
			item.Name = "其他"
			continue
		}

		subarea, err := svc.video.GetSubAreaById(fmt.Sprint(item.Id))
		if err == nil {
			item.Name = subarea.SubareaName
		}
	}

	mp["list"] = subareaStat

	return errdef.SUCCESS, mp
}

// 帖子板块统计 [发布占比]
func (svc *StatModule) GetPostSectionStat() (int, map[string]interface{}) {
	sectionStat, err := svc.stat.GetPostSectionStat()
	if err != nil {
		log.Log.Errorf("stat_trace: get post section stat fail, err:%s", err)
		return errdef.ERROR, nil
	}

	total, err := svc.stat.GetPostTotal()
	if err != nil {
		log.Log.Errorf("stat_trace: get post total fail, err:%s", err)
		return errdef.ERROR, nil
	}

	mp := make(map[string]interface{}, 0)
	mp["title"] = "帖子各板块发布数据"

	for _, item := range sectionStat {
		item.Rate = fmt.Sprintf("%.2f%s", float64(item.Count) / float64(total) * 100, "%")
		if item.Id == 0 {
			item.Name = "其他"
			continue
		}

		section, err := svc.community.GetSectionInfo(fmt.Sprint(item.Id))
		if err == nil && section != nil {
			item.Name = section.SectionName
		}
	}

	mp["list"] = sectionStat

	return errdef.SUCCESS, mp
}

// 各板块每日发布帖子统计
func (svc *StatModule) PublishPostDaily(queryMinDate, queryMaxDate string) (int, map[int64]*PublishStat) {
	days := 5
	minDate := time.Now().AddDate(0, 0, -days).Format(consts.FORMAT_DATE)
	maxDate := time.Now().Format(consts. FORMAT_DATE)
	if queryMinDate != "" && queryMaxDate != "" {
		minDate = queryMinDate
		maxDate = queryMaxDate
		min, err := time.Parse(consts.FORMAT_DATE, queryMinDate)
		if err != nil {
			log.Log.Errorf("stat_trace: time.Parse fail, err:%s", err)
			return errdef.ERROR, nil
		}

		max, err := time.Parse(consts.FORMAT_DATE, queryMaxDate)
		if err != nil {
			log.Log.Errorf("stat_trace: time.Parse fail, err:%s", err)
			return errdef.ERROR, nil
		}

		days = util.GetDiffDays(max, min)
		log.Log.Infof("##########days:%d", days)
	}

	stat, err := svc.stat.PublishDataDailyByPost(minDate, maxDate)
	if err != nil {
		return errdef.ERROR, nil
	}

	return  errdef.SUCCESS, svc.ResultInfo(days, 2, stat, maxDate)
}


// 各分区每日发布视频统计
func (svc *StatModule) PublishVideoDaily(queryMinDate, queryMaxDate string) (int, map[int64]*PublishStat) {
	days := 5
	minDate := time.Now().AddDate(0, 0, -days).Format(consts.FORMAT_DATE)
	maxDate := time.Now().Format(consts. FORMAT_DATE)
	if queryMinDate != "" && queryMaxDate != "" {
		minDate = queryMinDate
		maxDate = queryMaxDate
		min, err := time.Parse(consts.FORMAT_DATE, queryMinDate)
		if err != nil {
			log.Log.Errorf("stat_trace: time.Parse fail, err:%s", err)
			return errdef.ERROR, nil
		}

		max, err := time.Parse(consts.FORMAT_DATE, queryMaxDate)
		if err != nil {
			log.Log.Errorf("stat_trace: time.Parse fail, err:%s", err)
			return errdef.ERROR, nil
		}

		days = util.GetDiffDays(max, min)
		log.Log.Infof("##########days:%d", days)
	}

	stat, err := svc.stat.PublishDataDailyByVideo(minDate, maxDate)
	if err != nil {
		return errdef.ERROR, nil
	}

	return  errdef.SUCCESS, svc.ResultInfo(days, 1, stat, maxDate)

}

type PublishStat struct {
	Title string             `json:"title"`
	List  map[string]int64   `json:"list"`
}

// pubType 1 视频 2 帖子
func (svc *StatModule) ResultInfo(days, pubType int, data []*mstat.Stat, maxDate string) map[int64]*PublishStat {
	list := make(map[int64]*PublishStat)
	max, err := time.Parse(consts.FORMAT_DATE, maxDate)
	if err != nil {
		return nil
	}

	for i := 0; i <= days; i++ {
		date := max.AddDate(0, 0, -i).Format("2006-01-02")
		for _, v := range data {
			if v.Dt == date {
				vs, ok := list[v.Id]
				if ok {
					vs.List[date] = v.Count
				} else {
					list[v.Id] = &PublishStat{
						Title: v.Name,
						List:  make(map[string]int64),
					}

					list[v.Id].List[date] = v.Count
				}
			}

		}
	}

	list[1000] = &PublishStat{
		Title: "总数",
		List:  make(map[string]int64),
	}
	//mp := make(map[string]int64, days-1)
	for k, v := range list {
		for i := 0; i <= days; i++ {
			date := max.AddDate(0, 0, -i).Format("2006-01-02")
			if _, ok := v.List[date]; !ok {
				list[k].List[date] = 0
			}

			if pubType == 1 {
				list[1000].List[date] = svc.stat.GetDailyPublishVideoByDate(date)
			} else {
				list[1000].List[date] = svc.stat.GetDailyPublishPostByDate(date)
			}

		}
	}

	return list
}

// 每日帖子发布总数
func (svc *StatModule) DailyTotalPost(queryMinDate, queryMaxDate string) (int, map[string]interface{}) {
	minDate := time.Now().AddDate(0, 0, -11).Format(consts.FORMAT_DATE)
	maxDate := time.Now().Format(consts. FORMAT_DATE)
	if queryMinDate != "" && queryMaxDate != "" {
		minDate = queryMinDate
		maxDate = queryMaxDate
	}

	stat, err := svc.stat.GetTotalDailyPublishByPost(minDate, maxDate)
	if err != nil {
		return errdef.ERROR, nil
	}

	return errdef.SUCCESS, svc.ResultInfoByDate(stat, 10)
}

// 每日视频发布总数
func (svc *StatModule) DailyTotalVideo(queryMinDate, queryMaxDate string) (int, map[string]interface{}) {
	minDate := time.Now().AddDate(0, 0, -11).Format(consts.FORMAT_DATE)
	maxDate := time.Now().Format(consts. FORMAT_DATE)
	if queryMinDate != "" && queryMaxDate != "" {
		minDate = queryMinDate
		maxDate = queryMaxDate
	}

	stat, err := svc.stat. GetTotalDailyPublishByVideo(minDate, maxDate)
	if err != nil {
		return errdef.ERROR, nil
	}

	return errdef.SUCCESS, svc.ResultInfoByDate(stat, 10)
}

func (svc *StatModule) ResultInfoByDate(data []*mstat.Stat, days int) map[string]interface{} {
	mapList := make(map[string]interface{})
	for i := 0; i <= days; i++ {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		for _, v := range data {
			if v.Dt == date {
				mapList[date] = v.Count
			}
		}

		if _, ok := mapList[date]; !ok {
			mapList[date] = 0
		}
	}
	return mapList
}

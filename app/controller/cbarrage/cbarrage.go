package cbarrage

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/dao"
	"sports_service/global/app/errdef"
	"sports_service/global/app/log"
	"sports_service/global/consts"
	"sports_service/models"
	"sports_service/models/mbarrage"
	"sports_service/models/mconfigure"
	"sports_service/models/mcontest"
	"sports_service/models/muser"
	"sports_service/models/mvideo"
	"sports_service/tools/tencentCloud"
	"time"
)

type BarrageModule struct {
	context *gin.Context
	engine  *xorm.Session
	user    *muser.UserModel
	video   *mvideo.VideoModel
	barrage *mbarrage.BarrageModel
	contest *mcontest.ContestModel
	config  *mconfigure.ConfigModel
}

func New(c *gin.Context) BarrageModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	return BarrageModule{
		context: c,
		user:    muser.NewUserModel(socket),
		video:   mvideo.NewVideoModel(socket),
		barrage: mbarrage.NewBarrageModel(socket),
		contest: mcontest.NewContestModel(socket),
		config:  mconfigure.NewConfigModel(socket),
		engine:  socket,
	}
}

// 发送弹幕
func (svc *BarrageModule) SendBarrage(userId string, params *mbarrage.SendBarrageParams) int {
	client := tencentCloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
	// 检测弹幕内容
	isPass, content, err := client.TextModeration(params.Content)
	if err != nil {
		log.Log.Errorf("barrage_trace: validate barrage content err: %s，pass: %v", err, isPass)
		return errdef.CLOUD_FILTER_FAIL
	}

	if !isPass {
		params.Content = content
	}

	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("barrage_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS
	}

	var code int
	switch params.BarrageType {
	case 0:
		code = svc.SendVideoBarrage(userId, params)
	case 1:
		code = svc.SendLiveBarrage(userId, params)
	default:
		return errdef.INVALID_PARAMS
	}

	return code
}

// 发送直播/直播回放弹幕
func (svc *BarrageModule) SendLiveBarrage(userId string, params *mbarrage.SendBarrageParams) int {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("barrage_trace: session begin err:%s", err)
		return errdef.ERROR
	}

	// 查询直播是否存在
	ok, err := svc.contest.GetLiveInfoByCondition(fmt.Sprintf("id=%d", params.VideoId))
	if !ok || err != nil {
		log.Log.Errorf("barrage_trace: LIVE not found, liveId:%s", params.VideoId)
		svc.engine.Rollback()
		return errdef.CONTEST_GET_LIVE_FAIL
	}

	now := int(time.Now().Unix())
	svc.barrage.Barrage.VideoId = params.VideoId
	svc.barrage.Barrage.SendTime = int64(now)
	svc.barrage.Barrage.UserId = userId
	svc.barrage.Barrage.Content = params.Content
	svc.barrage.Barrage.BarrageType = 1
	svc.barrage.Barrage.VideoCurDuration = params.VideoCurDuration
	if svc.contest.VideoLive.StartTime > 0 && params.VideoCurDuration == 0 {
		svc.barrage.Barrage.VideoCurDuration = now - svc.contest.VideoLive.StartTime
	}
	// 可选参数
	svc.barrage.Barrage.Color = params.Color
	svc.barrage.Barrage.Font = params.Font
	svc.barrage.Barrage.Location = params.Location
	// 存储到mysql
	if err := svc.barrage.RecordVideoBarrage(); err != nil {
		log.Log.Errorf("barrage_trace: record live barrage fail, liveId:%d, err:%s", svc.contest.VideoLive.Id, err)
		svc.engine.Rollback()
		return errdef.BARRAGE_LIVE_SEND_FAIL
	}

	svc.engine.Commit()

	return errdef.SUCCESS
}

// 发送视频弹幕(记录到数据库 并发布到nsq) todo：nsq替换为kafka, 发布的内容 需通过敏感词过滤
func (svc *BarrageModule) SendVideoBarrage(userId string, params *mbarrage.SendBarrageParams) int {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("barrage_trace: session begin err:%s", err)
		return errdef.ERROR
	}

	// 查询视频是否存在
	video := svc.video.FindVideoById(fmt.Sprint(params.VideoId))
	if video == nil {
		log.Log.Errorf("barrage_trace: video not found, videoId:%s", params.VideoId)
		svc.engine.Rollback()
		return errdef.VIDEO_NOT_EXISTS
	}

	now := time.Now().Unix()
	svc.barrage.Barrage.VideoId = params.VideoId
	svc.barrage.Barrage.SendTime = now
	svc.barrage.Barrage.UserId = userId
	svc.barrage.Barrage.Content = params.Content
	svc.barrage.Barrage.VideoCurDuration = params.VideoCurDuration
	// 可选参数
	svc.barrage.Barrage.Color = params.Color
	svc.barrage.Barrage.Font = params.Font
	svc.barrage.Barrage.Location = params.Location
	// 存储到mysql
	if err := svc.barrage.RecordVideoBarrage(); err != nil {
		log.Log.Errorf("barrage_trace: record video barrage err:%s", err)
		svc.engine.Rollback()
		return errdef.BARRAGE_VIDEO_SEND_FAIL
	}

	score := svc.config.GetActionScore(int(consts.WORK_TYPE_VIDEO), consts.ACTION_TYPE_BARRAGE)
	// 更新视频弹幕总计 +1
	if err := svc.video.UpdateVideoBarrageNum(video.VideoId, int(now), consts.CONFIRM_OPERATE, score); err != nil {
		log.Log.Errorf("barrage_trace: update video barrage num err:%s", err)
		svc.engine.Rollback()
		return errdef.BARRAGE_VIDEO_SEND_FAIL
	}

	svc.engine.Commit()

	// todo: nsq 已停用
	//barrageMsg := &pbBarrage.BarrageMessage{
	//	Barrage: &pbBarrage.BarrageInfo{
	//		Uid: userId,
	//		Content: params.Content,
	//		VideoId: fmt.Sprint(params.VideoId),
	//		CurDuration: int64(params.VideoCurDuration),
	//		SendTime: now,
	//	},
	//}

	//bts, err := proto.Marshal(barrageMsg)
	//if err != nil {
	//	log.Log.Errorf("barrage_trace: proto marshal err:%s", err)
	//}

	// 发布到topic
	//if err := nsq.NsqProducer.Publish(consts.MESSAGE_TOPIC, bts); err != nil {
	//	log.Log.Errorf("barrage_trace: publish msg err:%s", err)
	//}

	return errdef.SUCCESS
}

// 获取视频弹幕列表
func (svc *BarrageModule) GetVideoBarrageList(videoId, barrageType, minDuration, maxDuration string) (int, []*models.VideoBarrage) {
	if videoId == "" {
		log.Log.Errorf("barrage_trace: invalid id, id:%s", videoId)
		return errdef.INVALID_PARAMS, nil
	}
	switch barrageType {
	case "0":
		// 查询视频是否存在
		if video := svc.video.FindVideoById(videoId); video == nil {
			log.Log.Errorf("barrage_trace: video not found, videoId:%s", videoId)
			return errdef.VIDEO_NOT_EXISTS, nil
		}
	case "1":
		ok, err := svc.contest.GetLiveInfoByCondition(fmt.Sprintf("id=%s", videoId))
		if !ok || err != nil {
			log.Log.Errorf("barrage_trace: live not found, liveId:%s", videoId)
			return errdef.CONTEST_GET_LIVE_FAIL, nil
		}
	default:
		return errdef.INVALID_PARAMS, nil
	}

	list := svc.barrage.GetBarrageByDuration(videoId, barrageType, minDuration, maxDuration, 0, 1000)
	if list == nil {
		list = []*models.VideoBarrage{}
	}

	return errdef.SUCCESS, list
}

// 获取用户视频弹幕总数
func (svc *BarrageModule) GetUserTotalVideoBarrage(userId string) int64 {
	return svc.barrage.GetUserTotalVideoBarrage(userId)
}

package cbarrage

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	pbBarrage "sports_service/server/proto/barrage"
	"sports_service/server/tools/nsq"
	"sports_service/server/global/app/errdef"
	"sports_service/server/models/mbarrage"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/log"
	"github.com/golang/protobuf/proto"
	"fmt"
	"time"
)

type BarrageModule struct {
	context    *gin.Context
	engine     *xorm.Session
	user       *muser.UserModel
	video      *mvideo.VideoModel
	barrage    *mbarrage.BarrageModel
}

func New(c *gin.Context) BarrageModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return BarrageModule{
		context: c,
		user: muser.NewUserModel(socket),
		video: mvideo.NewVideoModel(socket),
		barrage: mbarrage.NewBarrageModel(socket),
		engine: socket,
	}
}

// 发送视频弹幕(记录到数据库 并发布到nsq) todo：nsq替换为kafka, 发布的内容 需通过敏感词过滤
func (svc *BarrageModule) SendVideoBarrage(userId string, params *mbarrage.SendBarrageParams) int {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("barrage_trace: session begin err:%s", err)
		return errdef.ERROR
	}
	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("barrage_trace: user not found, userId:%s", userId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS
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

	// 更新视频弹幕总计 +1
	if err := svc.video.UpdateVideoBarrageNum(video.VideoId, int(now), consts.CONFIRM_OPERATE); err != nil {
		log.Log.Errorf("like_trace: update video like num err:%s", err)
		svc.engine.Rollback()
		return errdef.LIKE_VIDEO_FAIL
	}

	svc.engine.Commit()

	barrageMsg := &pbBarrage.BarrageMessage{
		Barrage: &pbBarrage.BarrageInfo{
			Uid: userId,
			Content: params.Content,
			VideoId: fmt.Sprint(params.VideoId),
			CurDuration: int64(params.VideoCurDuration),
			SendTime: now,
		},
	}

	bts, err := proto.Marshal(barrageMsg)
	if err != nil {
		log.Log.Errorf("barrage_trace: proto marshal err:%s", err)
	}

	// 发布到topic
	if err := nsq.NsqProducer.Publish(consts.MESSAGE_TOPIC, bts); err != nil {
		log.Log.Errorf("barrage_trace: publish msg err:%s", err)
	}

	return errdef.SUCCESS
}

// 获取视频弹幕列表
func (svc *BarrageModule) GetVideoBarrageList(videoId, minDuration, maxDuration string) (int, []*models.VideoBarrage) {
	// 查询视频是否存在
	if video := svc.video.FindVideoById(videoId); video == nil {
		log.Log.Errorf("barrage_trace: video not found, videoId:%s", videoId)
		return errdef.VIDEO_NOT_EXISTS, nil
	}

	return errdef.SUCCESS, svc.barrage.GetBarrageByDuration(videoId, minDuration, maxDuration, 0, 1000)
}

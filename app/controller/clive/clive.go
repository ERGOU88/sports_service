package clive

import (
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/tools/live"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/models/mcontest"
	"sports_service/server/global/consts"
	"strconv"
	"time"
)

type LiveModule struct {
	context    *gin.Context
	engine     *xorm.Session
	contest    *mcontest.ContestModel
}

func New(c *gin.Context) LiveModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	return LiveModule{
		context: c,
		contest: mcontest.NewContestModel(socket),
		engine: socket,
	}
}

// 推流/断流回调
func (svc *LiveModule) PushOrDisconnectStreamCallback(params *mcontest.StreamCallbackInfo) int {
	log.Log.Infof("live_trace: stream callbackInfo:%+v", params)
	if code := svc.ValidateParamInfo(params); code != errdef.SUCCESS {
		return code
	}

	now := int(time.Now().Unix())
	var cols string
	svc.contest.VideoLive.Status = params.Status
	svc.contest.VideoLive.UpdateAt = now
	switch params.Status {
	case 1:
		log.Log.Errorf("live_trace: push stream, roomId:%s, startTm:%d", params.StreamID, params.EventTime)
		// 已记录过开播时长
		if svc.contest.VideoLive.StartTm > 0 {
			return errdef.SUCCESS
		}

		svc.contest.VideoLive.StartTm = params.EventTime
		cols = "start_tm, status, update_at"
	case 2:
		log.Log.Errorf("live_trace: disconnect stream, roomId:%s, startTm:%d", params.StreamID, params.EventTime)
		duration, err := strconv.Atoi(params.PushDuration)
		if err != nil {
			log.Log.Errorf("live_trace: strconv.Atoi fail, duration:%s", params.PushDuration)
		}

		svc.contest.VideoLive.Duration += int64(duration)
		svc.contest.VideoLive.EndTime = params.EventTime
		cols = "end_tm, status, update_at, duration"
	}

	affected, err := svc.contest.UpdateLiveInfo(cols)
	if affected != 1 || err != nil {
		log.Log.Errorf("live_trace: update live info fail, id:%d, err:%s", svc.contest.VideoLive.Id, err)
		return errdef.ERROR
	}

	return errdef.CALLBACK_SUCCESS
}

// 录制回调
func (svc *LiveModule) TranscribeStreamCallback(param *mcontest.StreamCallbackInfo) int {
	log.Log.Infof("live_trace: transcribe callbackInfo:%+v", param)
	if code := svc.ValidateParamInfo(param); code != errdef.SUCCESS {
		return code
	}

	now := int(time.Now().Unix())
	svc.contest.VideoLiveReplay.UserId = svc.contest.VideoLive.UserId
	svc.contest.VideoLiveReplay.UpdateAt = now
	svc.contest.VideoLiveReplay.CreateAt = now
	svc.contest.VideoLiveReplay.Cover = svc.contest.VideoLive.Cover
	svc.contest.VideoLiveReplay.Duration = param.Duration
	svc.contest.VideoLiveReplay.Labeltype = 1
	svc.contest.VideoLiveReplay.Size = int64(param.FileSize)
	svc.contest.VideoLiveReplay.TaskId = param.TaskID
	svc.contest.VideoLiveReplay.HistoryAddr = param.VideoURL
	svc.contest.VideoLiveReplay.Title = svc.contest.VideoLive.Title
	svc.contest.VideoLiveReplay.Describe = svc.contest.VideoLive.Describe
	svc.contest.VideoLiveReplay.LiveId = svc.contest.VideoLive.Id
	// 添加直播回放
	affected, err := svc.contest.AddVideoLiveReply()
	if affected != 1 || err != nil {
		log.Log.Errorf("live_trace: add video live reply fail, roomId:%s, err:%s", svc.contest.VideoLive.RoomId, err)
		return errdef.ERROR
	}

	return errdef.CALLBACK_SUCCESS
}

// 生成回调签名
func (svc *LiveModule) GenCallbackSign(t int) string {
	return live.BuildCallbackSign(t)
}

// 校验回调参数信息
func (svc *LiveModule) ValidateParamInfo(params *mcontest.StreamCallbackInfo ) int {
	sign := svc.GenCallbackSign(params.T)
	if sign != params.Sign {
		log.Log.Errorf("live_trace: sign not match, sign:%s, param sign:%s", sign, params.Sign)
		return errdef.INVALID_PARAMS
	}

	ok, err := svc.contest.GetLiveInfoByRoomId(params.StreamID)
	if !ok || err != nil {
		log.Log.Errorf("live_trace: get live info by roomId fail, streamId:%s, err:%s", params.StreamID, err)
		return errdef.INVALID_PARAMS
	}

	if params.Appid != consts.TX_APP_ID {
		log.Log.Errorf("live_trace: invalid appId, appId:%d, param appId:%d", consts.TX_APP_ID, params.Appid)
		return errdef.INVALID_PARAMS
	}

	if params.App != live.PUSH_STREAM_HOST {
		log.Log.Errorf("live_trace: invalid push stream host, host:%s, param host:%s", live.PUSH_STREAM_HOST, params.App)
		return errdef.INVALID_PARAMS
	}

	return errdef.SUCCESS
}

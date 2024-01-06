package event

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sports_service/dao"
	"sports_service/global/app/log"
	"sports_service/global/consts"
	"sports_service/global/rdskey"
	"sports_service/models/mconfigure"
	"sports_service/models/minformation"
	"sports_service/models/mposting"
	"sports_service/models/muser"
	"sports_service/models/mvideo"
	"sports_service/redismq/protocol"
	"sports_service/util"
	"time"
)

// 作品是否置顶
func LoopPopTopEvent() {
	for !closing {
		conn := dao.RedisPool().Get()
		values, err := redis.Values(conn.Do("BRPOP", rdskey.MSG_TOP_EVENT_KEY, 0))
		conn.Close()
		if err != nil {
			log.Log.Errorf("redisMq_trace: loopPop event fail, err:%s", err)
			// 防止出现错误时 频繁刷日志
			time.Sleep(time.Second)
			continue
		}

		if len(values) < 2 {
			log.Log.Errorf("redisMq_trace: invalid values, len:%d, values:%+v", len(values), values)
		}

		bts, ok := values[1].([]byte)
		if !ok {
			log.Log.Errorf("redisMq_trace: value[1] unSupport type")
			continue
		}

		if err := TopEventConsumer(bts); err != nil {
			log.Log.Errorf("redisMq_trace: event consumer fail, err:%s, msg:%s", err, string(bts))
		}

	}
}

func TopEventConsumer(bts []byte) error {
	event := protocol.Event{}
	if err := util.JsonFast.Unmarshal(bts, &event); err != nil {
		log.Log.Errorf("redisMq_trace: proto unmarshal event err:%s", err)
		return err
	}

	if err := handleTopEvent(event); err != nil {
		log.Log.Errorf("handleEvent err:%s", err)
		return err
	}

	return nil
}

func handleTopEvent(event protocol.Event) error {
	info := &protocol.WorkInfo{}
	if err := util.JsonFast.Unmarshal(event.Data, info); err != nil {
		log.Log.Errorf("redisMq_trace: proto unmarshal data err:%s", err)
		return nil
	}

	session := dao.AppEngine.NewSession()
	defer session.Close()
	umodel := muser.NewUserModel(session)
	user := umodel.FindUserByUserid(event.UserId)
	if user == nil {
		log.Log.Errorf("redisMq_trace: user not found, userId:%s", event.UserId)
		return nil
	}

	cmodel := mconfigure.NewConfigModel(session)
	cmodel.ActionScore.WorkType = int(event.EventType)
	ok, err := cmodel.GetActionScoreConf()
	if !ok || err != nil {
		return nil
	}

	switch event.EventType {
	case consts.EVENT_SET_TOP_VIDEO:
		vmodel := mvideo.NewVideoModel(session)
		videos := vmodel.FindVideoById(info.Id)
		if videos == nil {
			return errors.New("video not found")
		}

		if videos.IsTop == 1 {
			return nil
		}

		statistic := vmodel.GetVideoStatistic(info.Id)
		if statistic.HeatNum >= cmodel.ActionScore.TopScore {
			videos.IsTop = 1
			if _, err := vmodel.UpdateVideoInfo(); err != nil {
				return err
			}
		}

	case consts.EVENT_SET_TOP_POST:
		pmodel := mposting.NewPostingModel(session)
		post, err := pmodel.GetPostById(info.Id)
		if err != nil {
			return err
		}

		if post.IsTop == 1 {
			return nil
		}

		statistic, err := pmodel.GetPostStatistic(info.Id)
		if err != nil {
			return err
		}

		if statistic.HeatNum >= cmodel.ActionScore.TopScore {
			post.IsTop = 1
			if _, err := pmodel.UpdatePostInfo(post.Id, "is_top"); err != nil {
				return err
			}
		}

	case consts.EVENT_SET_TOP_INFO:
		imodel := minformation.NewInformationModel(session)
		ok, err := imodel.GetInformationById(info.Id)
		if !ok || err != nil {
			return err
		}

		if imodel.Information.IsTop == 1 {
			return nil
		}

		ok, err = imodel.GetInformationStatistic(info.Id)
		if !ok || err != nil {
			return err
		}

		if imodel.Statistic.HeatNum >= cmodel.ActionScore.TopScore {
			imodel.Information.IsTop = 1
			if _, err := imodel.UpdateInfo(fmt.Sprintf("id=%s", info.Id), "is_top"); err != nil {
				return err
			}
		}
	}

	return nil
}

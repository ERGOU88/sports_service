package job

import (
	"context"
	"sports_service/server/global/consts"
	"sports_service/server/models/mvideo"
	"sports_service/server/models/muser"
	cloud "sports_service/server/tools/tencentCloud"
	"sports_service/server/util"
	"strings"
	"sports_service/server/global/app/log"
	"time"
	"sports_service/server/dao"
)

// 主动拉取事件（腾讯云）
func PullEventsJob() {
	ticker := time.NewTicker(time.Minute * 3)
	defer ticker.Stop()

	for {
		select {
		case <- ticker.C:
			log.Log.Debugf("开始拉取事件[腾讯云]")
			if err := pullEvents(); err != nil {
				log.Log.Errorf("job_trace: pull events err:%s", err)
			}
			log.Log.Debugf("事件处理完毕")
		}
	}

}

// 主动拉取事件回调
func pullEvents() error {
	vod := cloud.New(consts.SECRET_ID, consts.SECRET_KEY, consts.VOD_API_DOMAIN)
	resp, err := vod.PullEvents()
	if err != nil {
		return err
	}

	for _, event := range resp.Response.EventSet {
		switch *event.EventType {
		// 上传事件
		case consts.EVENT_TYPE_UPLOAD:
			source := new(cloud.SourceContext)
			if err := util.JsonFast.Unmarshal([]byte(*event.FileUploadEvent.MediaBasicInfo.SourceInfo.SourceContext), source); err != nil {
				log.Log.Errorf("video_trace: jsonfast unmarshal event sourceContext err:%s", err)
				continue
			}

			vmodel := mvideo.NewVideoModel(dao.Engine.Context(context.Background()))

			// 通过任务id 获取 用户id
			userId, err := vmodel.GetUploadUserIdByTaskId(source.TaskId)
			if err != nil || userId == "" {
				log.Log.Errorf("video_trace: invalid taskId, taskId:%d", source.TaskId)
				continue
			}

			umodel := muser.NewUserModel(dao.Engine.Context(context.Background()))
			// 查询用户是否存在
			if user := umodel.FindUserByUserid(userId); user == nil {
				log.Log.Errorf("video_trace: user not found, userId:%s", userId)
				continue
			}

			// 是否为同一个用户
			if strings.Compare(userId, source.UserId) != 0 {
				log.Log.Errorf("job_trace: userId not match, eventUserId:%s, redis userId:%s", source.UserId, userId)
				continue
			}

			// todo: 确认事件通知


		}
	}

	return nil
}


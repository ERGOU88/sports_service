package job

import (
	"context"
	"sports_service/server/dao"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mlabel"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	cloud "sports_service/server/tools/tencentCloud"
	"sports_service/server/util"
	"strconv"
	"strings"
	"time"
)

// 主动拉取事件（腾讯云）
func PullEventsJob() {
	ticker := time.NewTicker(time.Minute * 2)
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
	vod := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
	resp, err := vod.PullEvents()
	if err != nil {
		return err
	}

  session := dao.Engine.Context(context.Background())
  if err := session.Begin(); err != nil {
    log.Log.Errorf("job_trace: session begin err:%s", err)
    return err
  }

	for _, event := range resp.Response.EventSet {
		switch *event.EventType {
		// 上传事件
		case consts.EVENT_TYPE_UPLOAD:
			source := new(cloud.SourceContext)
			if err := util.JsonFast.Unmarshal([]byte(*event.FileUploadEvent.MediaBasicInfo.SourceInfo.SourceContext), source); err != nil {
				log.Log.Errorf("job_trace: jsonfast unmarshal event sourceContext err:%s", err)
				continue
			}

			vmodel := mvideo.NewVideoModel(session)

			// 通过任务id 获取 用户id
			userId, err := vmodel.GetUploadUserIdByTaskId(source.TaskId)
			if err != nil || userId == "" {
				log.Log.Errorf("job_trace: invalid taskId, taskId:%d", source.TaskId)
				session.Rollback()
				continue
			}

			umodel := muser.NewUserModel(session)
			// 查询用户是否存在
			if user := umodel.FindUserByUserid(userId); user == nil {
				log.Log.Errorf("job_trace: user not found, userId:%s", userId)
				session.Rollback()
				continue
			}

			// 是否为同一个用户
			if strings.Compare(userId, source.UserId) != 0 {
				log.Log.Errorf("job_trace: userId not match, eventUserId:%s, redis userId:%s", source.UserId, userId)
				session.Rollback()
				continue
			}

			info, err := vmodel.GetPublishInfo(source.UserId, source.TaskId)
			if err != nil || info == "" {
				log.Log.Errorf("job_trace: get publish info err:%s", err)
				session.Rollback()
				continue
			}

			// 获取用户发布的视频信息
			pubInfo := new(mvideo.VideoPublishParams)
			if err := util.JsonFast.Unmarshal([]byte(info), pubInfo); err != nil {
				log.Log.Errorf("job_trace: jsonFast unmarshal err: %s", err)
				session.Rollback()
				continue
			}

			if pubInfo.TaskId != source.TaskId {
				log.Log.Errorf("job_trace: task id not match, pub taskId:%d, source taskId:%d", pubInfo.TaskId, source.TaskId)
				session.Rollback()
				continue
			}

			// 数据记录到视频审核表 同时 标签记录到 视频标签表（多条记录 同一个videoId对应N个labelId 生成N条记录）
			now := time.Now().Unix()
      vmodel.Videos.UserId = userId
      vmodel.Videos.Cover = *event.FileUploadEvent.MediaBasicInfo.CoverUrl
      vmodel.Videos.Title = pubInfo.Title
      vmodel.Videos.Describe = pubInfo.Describe
      vmodel.Videos.VideoAddr = pubInfo.VideoAddr
      vmodel.Videos.VideoDuration = pubInfo.VideoDuration
      vmodel.Videos.CreateAt = int(now)
      vmodel.Videos.UpdateAt = int(now)
      vmodel.Videos.UserType = consts.PUBLISH_VIDEO_BY_USER
      vmodel.Videos.VideoWidth = pubInfo.VideoWidth
      vmodel.Videos.VideoHeight = pubInfo.VideoHeight
			fileId, _ := strconv.Atoi(*event.FileUploadEvent.FileId)
      vmodel.Videos.FileId = int64(fileId)
      vmodel.Videos.Size = pubInfo.Size
      // todo: 如果有 记录用户自定义标签

			// 视频发布
			affected, err := vmodel.VideoPublish()
			if err != nil || affected != 1 {
				log.Log.Errorf("job_trace: publish video err:%s, affected:%d", err, affected)
				session.Rollback()
				return err
			}

			lmodel := mlabel.NewLabelModel(session)
			labelIds := strings.Split(pubInfo.VideoLabels, ",")
			// 组装多条记录 写入视频标签表
			labelInfos := make([]*models.VideoLabels, len(labelIds))
			for index, labelId := range labelIds {
				if lmodel.GetLabelInfoByMem(labelId) == nil {
					log.Log.Errorf("job_trace: label not found, labelId:%s", labelId)
					continue
				}

				info := new(models.VideoLabels)
				info.VideoId = vmodel.Videos.VideoId
				info.LabelId = labelId
				info.LabelName = lmodel.GetLabelNameByMem(labelId)
				info.CreateAt = int(now)
				labelInfos[index] = info
			}

			if len(labelInfos) > 0 {
        // 添加视频标签（多条）
        affected, err = vmodel.AddVideoLabels(labelInfos)
        if err != nil || int(affected) != len(labelInfos) {
          log.Log.Errorf("job_trace: add video labels err:%s", err)
          session.Rollback()
          continue
        }
      }


			vmodel.Statistic.VideoId = vmodel.Videos.VideoId
			vmodel.Statistic.CreateAt = int(now)
			vmodel.Statistic.UpdateAt = int(now)
			// 初始化视频统计数据
			if err := vmodel.AddVideoStatistic(); err != nil {
				log.Log.Errorf("job_trace: add video statistic err:%s", err)
				session.Rollback()
				continue
			}

			// 记录事件回调信息
			vmodel.Events.FileId = int64(fileId)
			vmodel.Events.CreateAt = int(now)
			bts, _ := util.JsonFast.Marshal(event)
			vmodel.Events.Event = string(bts)
			affected, err = vmodel.RecordTencentEvent()
			if err != nil || affected != 1 {
				log.Log.Errorf("job_trace: record tencent event err:%s, affected:%d", err, affected)
				session.Rollback()
				continue
			}

			client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
			// 确认事件回调
			if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
				log.Log.Errorf("job_trace: confirm events err:%s", err)
				session.Rollback()
				continue
			}

			session.Commit()

		default:

		}
	}

	return nil
}


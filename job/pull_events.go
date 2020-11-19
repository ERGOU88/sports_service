package job

import (
  "context"
  "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
  "sports_service/server/dao"
  "sports_service/server/global/app/log"
  "sports_service/server/global/consts"
  "sports_service/server/models"
  "sports_service/server/models/mlabel"
  "sports_service/server/models/muser"
  "sports_service/server/models/mvideo"
  cloud "sports_service/server/tools/tencentCloud"
  "fmt"
  "sports_service/server/util"
  "strconv"
  "strings"
  "time"
  "errors"
)

// 主动拉取事件（腾讯云）
func PullEventsJob() {
	ticker := time.NewTicker(time.Minute * 1)
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

	for _, event := range resp.Response.EventSet {
    log.Log.Debugf("eventType:%v", *event.EventType)
		switch *event.EventType {
		// 上传事件
		case consts.EVENT_TYPE_UPLOAD:
      log.Log.Debugf("upload event:%+v", *event.FileUploadEvent)
		  if err := uploadEvent(event); err != nil {
		    log.Log.Errorf("job_trace: uploadEvent err:%s", err)
		    continue
      }
		// 任务流状态变更（包含视频转码完成）
    case consts.EVENT_PROCEDURE_STATE_CHANGED:
      log.Log.Debugf("transcode event:%+v", *event.ProcedureStateChangeEvent)
      transCodeCompleteEvent(event)

		default:

		}
	}

	return nil
}

// 视频转码事件
func transCodeCompleteEvent(event *v20180717.EventContent) error {
  session := dao.Engine.Context(context.Background())
  if err := session.Begin(); err != nil {
    log.Log.Errorf("job_trace: session begin err:%s", err)
    return err
  }

  client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
  vmodel := mvideo.NewVideoModel(session)
  video := vmodel.GetVideoByFileId(*event.ProcedureStateChangeEvent.FileId)
  if video == nil {
    log.Log.Errorf("job_trace: video not found, fileId:%s", *event.ProcedureStateChangeEvent.FileId)
    session.Rollback()
    // 确认事件回调
    if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
      log.Log.Errorf("job_trace: confirm events err:%s", err)
    }

    return errors.New("video not found")
  }

  list := make([]*mvideo.PlayInfo, 0)
  for _, info := range event.ProcedureStateChangeEvent.MediaProcessResultSet {
    log.Log.Debugf("info:%v", info)
    // todo:
    switch *info.Type {
    case "Transcode":
      if *info.TranscodeTask.ErrCode != 0 {
        log.Log.Errorf("job_trace: media process errCode:%d", *info.TranscodeTask.ErrCode)
        continue
      }

      // 流畅（FLU） 100010	MP4  100210	HLS
      playInfo := new(mvideo.PlayInfo)
      if *info.TranscodeTask.Output.Definition == 100010 || *info.TranscodeTask.Output.Definition == 100210 {
        playInfo.Type = "1"
      }

      // 标清（SD）	100020	MP4	 100220	 HLS
      if *info.TranscodeTask.Output.Definition == 100020 || *info.TranscodeTask.Output.Definition == 100220 {
        playInfo.Type = "2"
      }

      // 高清（HD）	100030	MP4	 100230	HLS
      if *info.TranscodeTask.Output.Definition == 100030 || *info.TranscodeTask.Output.Definition == 100230 {
        playInfo.Type = "3"
      }

      // 全高清（FHD）	100040	MP4 100240	HLS
      if *info.TranscodeTask.Output.Definition == 100040 || *info.TranscodeTask.Output.Definition == 100240 {
        playInfo.Type = "4"
      }

      // 全高清（FHD）	100040	MP4 100240	HLS
      if *info.TranscodeTask.Output.Definition == 100040 || *info.TranscodeTask.Output.Definition == 100240 {
        playInfo.Type = "4"
      }

      // 2K	100070	MP4	100270	HLS
      if *info.TranscodeTask.Output.Definition == 100070 || *info.TranscodeTask.Output.Definition == 100270 {
        playInfo.Type = "5"
      }

      // 4K	100080	MP4	100280	HLS
      if *info.TranscodeTask.Output.Definition == 100080 || *info.TranscodeTask.Output.Definition == 100280 {
        playInfo.Type = "6"
      }

      playInfo.Url = *info.TranscodeTask.Output.Url
      playInfo.Size = *info.TranscodeTask.Output.Size
      playInfo.Duration = int64(*info.TranscodeTask.Output.Duration * 1000)

      list = append(list, playInfo)
    }
  }

  playBts, err := util.JsonFast.Marshal(list)
  if err != nil {
    log.Log.Errorf("job_trace: jsonFast err:%s", err)
  }

  video.PlayInfo = string(playBts)
  if err := vmodel.UpdateVideoPlayInfo(fmt.Sprint(video.VideoId)); err != nil {
    session.Rollback()
    return errors.New("job_trace: update video play info fail")
  }

  now := time.Now().Unix()
  // 记录事件回调信息
  fileId, _ := strconv.Atoi(*event.ProcedureStateChangeEvent.FileId)
  vmodel.Events.FileId = int64(fileId)
  vmodel.Events.CreateAt = int(now)
  vmodel.Events.EventType = consts.EVENT_PROCEDURE_STATE_CHANGED_TYPE
  bts, _ := util.JsonFast.Marshal(event)
  vmodel.Events.Event = string(bts)
  affected, err := vmodel.RecordTencentEvent()
  if err != nil || affected != 1 {
    log.Log.Errorf("job_trace: record tencent transcode complete event err:%s, affected:%d", err, affected)
    session.Rollback()
    return errors.New("record tencent complete event fail")
  }

  // 确认事件回调
  if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
    log.Log.Errorf("job_trace: confirm events err:%s", err)
    session.Rollback()
    return errors.New("confirm event fail")
  }

  session.Commit()
  return nil
}

// 上传事件
func uploadEvent(event *v20180717.EventContent) error {
  session := dao.Engine.Context(context.Background())
  if err := session.Begin(); err != nil {
    log.Log.Errorf("job_trace: session begin err:%s", err)
    return err
  }

  client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
  source := new(cloud.SourceContext)
  if err := util.JsonFast.Unmarshal([]byte(*event.FileUploadEvent.MediaBasicInfo.SourceInfo.SourceContext), source); err != nil {
    log.Log.Errorf("job_trace: jsonfast unmarshal event sourceContext err:%s", err)
    // 确认事件回调
    if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
      log.Log.Errorf("job_trace: confirm events err:%s", err)
    }

    session.Rollback()
    return errors.New("jsonfast unmarshal event sourceContext err")
  }

  vmodel := mvideo.NewVideoModel(session)

  // 通过任务id 获取 用户id
  userId, err := vmodel.GetUploadUserIdByTaskId(source.TaskId)
  if err != nil || userId == "" {
    log.Log.Errorf("job_trace: invalid taskId, taskId:%d", source.TaskId)
    // 确认事件回调
    if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
      log.Log.Errorf("job_trace: confirm events err:%s", err)
    }

    session.Rollback()
    return errors.New("invalid taskId")
  }

  umodel := muser.NewUserModel(session)
  // 查询用户是否存在
  if user := umodel.FindUserByUserid(userId); user == nil {
    log.Log.Errorf("job_trace: user not found, userId:%s", userId)
    session.Rollback()
    return errors.New("user not found")
  }

  // 是否为同一个用户
  if strings.Compare(userId, source.UserId) != 0 {
    log.Log.Errorf("job_trace: userId not match, eventUserId:%s, redis userId:%s", source.UserId, userId)
    session.Rollback()
    return errors.New("userId not match")
  }

  info, err := vmodel.GetPublishInfo(source.UserId, source.TaskId)
  if err != nil || info == "" {
    log.Log.Errorf("job_trace: get publish info err:%s", err)
    session.Rollback()
    return errors.New("get publish info fail")
  }

  // 获取用户发布的视频信息
  pubInfo := new(mvideo.VideoPublishParams)
  if err := util.JsonFast.Unmarshal([]byte(info), pubInfo); err != nil {
    log.Log.Errorf("job_trace: jsonFast unmarshal err: %s", err)
    session.Rollback()
    return errors.New("jsonFast unmarshal err")
  }

  if pubInfo.TaskId != source.TaskId {
    log.Log.Errorf("job_trace: task id not match, pub taskId:%d, source taskId:%d", pubInfo.TaskId, source.TaskId)
    session.Rollback()
    return errors.New("taskId not match")
  }

  // 数据记录到视频审核表 同时 标签记录到 视频标签表（多条记录 同一个videoId对应N个labelId 生成N条记录）
  now := time.Now().Unix()
  vmodel.Videos.UserId = userId
  vmodel.Videos.Cover = *event.FileUploadEvent.MediaBasicInfo.CoverUrl
  vmodel.Videos.Title = pubInfo.Title
  vmodel.Videos.Describe = pubInfo.Describe
  vmodel.Videos.VideoAddr = pubInfo.VideoAddr
  // 转为毫秒
  vmodel.Videos.VideoDuration = int(*event.FileUploadEvent.MetaData.VideoDuration * 1000)
  vmodel.Videos.CreateAt = int(now)
  vmodel.Videos.UpdateAt = int(now)
  vmodel.Videos.UserType = consts.PUBLISH_VIDEO_BY_USER
  vmodel.Videos.VideoWidth = *event.FileUploadEvent.MetaData.Width
  vmodel.Videos.VideoHeight = *event.FileUploadEvent.MetaData.Height
  // 单位：字节
  vmodel.Videos.Size = *event.FileUploadEvent.MetaData.Size
  fileId, _ := strconv.Atoi(*event.FileUploadEvent.FileId)
  vmodel.Videos.FileId = int64(fileId)
  //vmodel.Videos.Size = pubInfo.Size
  // todo: 如果有 记录用户自定义标签

  // 视频发布
  affected, err := vmodel.VideoPublish()
  if err != nil || affected != 1 {
    log.Log.Errorf("job_trace: publish video err:%s, affected:%d", err, affected)
    session.Rollback()
    return errors.New("publish video fail")
  }

  lmodel := mlabel.NewLabelModel(session)
  labelIds := strings.Split(pubInfo.VideoLabels, ",")
  // 组装多条记录 写入视频标签表
  labelInfos := make([]*models.VideoLabels, 0)
  for _, labelId := range labelIds {
    if lmodel.GetLabelInfoByMem(labelId) == nil {
      log.Log.Errorf("job_trace: label not found, labelId:%s", labelId)
      continue
    }

    info := new(models.VideoLabels)
    info.VideoId = vmodel.Videos.VideoId
    info.LabelId = labelId
    info.LabelName = lmodel.GetLabelNameByMem(labelId)
    info.CreateAt = int(now)
    labelInfos = append(labelInfos, info)
  }

  if len(labelInfos) > 0 {
    // 添加视频标签（多条）
    affected, err = vmodel.AddVideoLabels(labelInfos)
    if err != nil || int(affected) != len(labelInfos) {
      log.Log.Errorf("job_trace: add video labels err:%s", err)
      session.Rollback()
      return errors.New("add video labels fail")
    }
  }

  vmodel.Statistic.VideoId = vmodel.Videos.VideoId
  vmodel.Statistic.CreateAt = int(now)
  vmodel.Statistic.UpdateAt = int(now)
  // 初始化视频统计数据
  if err := vmodel.AddVideoStatistic(); err != nil {
    log.Log.Errorf("job_trace: add video statistic err:%s", err)
    session.Rollback()
    return errors.New("add video statistic fail")
  }

  // 记录事件回调信息
  vmodel.Events.FileId = int64(fileId)
  vmodel.Events.CreateAt = int(now)
  vmodel.Events.EventType = consts.EVENT_UPLOAD_TYPE
  bts, _ := util.JsonFast.Marshal(event)
  vmodel.Events.Event = string(bts)
  affected, err = vmodel.RecordTencentEvent()
  if err != nil || affected != 1 {
    log.Log.Errorf("job_trace: record tencent event err:%s, affected:%d", err, affected)
    session.Rollback()
    return errors.New("record tencent event fail")
  }

  // 确认事件回调
  if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
    log.Log.Errorf("job_trace: confirm events err:%s", err)
    session.Rollback()
    return errors.New("confirm event fail")
  }

  session.Commit()

  return nil
}


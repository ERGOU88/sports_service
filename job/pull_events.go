package job

import (
  "errors"
  "github.com/garyburd/redigo/redis"
  "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
  "sports_service/server/dao"
  "sports_service/server/global/app/log"
  "sports_service/server/global/consts"
  "sports_service/server/models"
  "sports_service/server/models/mlabel"
  "sports_service/server/models/mposting"
  "sports_service/server/models/muser"
  "sports_service/server/models/mvideo"
  cloud "sports_service/server/tools/tencentCloud"
  "sports_service/server/util"
  "strconv"
  "strings"
  "time"
  "fmt"
)

// 主动拉取事件（腾讯云）
func PullEventsJob() {
  ticker := time.NewTicker(time.Second * 10)
  defer ticker.Stop()

  for {
    select {
    case <- ticker.C:
      log.Log.Debugf("开始拉取事件[腾讯云]")
      if err := pullEvents(); err != nil {
        log.Log.Infof("job_trace: pull events err:%s", err)
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
    log.Log.Debugf("event:%v, eventType:%v", *event, *event.EventType)
    switch *event.EventType {
    // 上传事件
    case consts.EVENT_TYPE_UPLOAD:
      log.Log.Debugf("upload event:%+v", *event.FileUploadEvent)
      if err := newUploadEvent(event); err != nil {
        log.Log.Errorf("job_trace: uploadEvent err:%s", err)
        continue
      }

    // 任务流状态变更（包含视频转码完成、视频AI审核等）
    case consts.EVENT_PROCEDURE_STATE_CHANGED:
      if err := procedureStateChangedEvent(event); err != nil {
        log.Log.Errorf("job_trace: procedureStateChangedEvent err:%s", err)
      }

    // 文件被删除
    case consts.EVENT_FILE_DELETED:
      log.Log.Debugf("fileDeleted event:%+v", *event.FileDeleteEvent)
      fileDeletedEvent(event)
    default:

    }
  }

  return nil
}

// 任务流状态变更（包含视频转码完成、视频AI审核等）
func procedureStateChangedEvent(event *v20180717.EventContent) error {
  bts, _ := util.JsonFast.Marshal(event.ProcedureStateChangeEvent)
  mp, err := util.JsonStringToMap(string(bts))
  if err != nil {
    log.Log.Errorf("job_trace: jsonStringToMap err:%s", err)
    return err
  }

  session := dao.AppEngine.NewSession()
  defer session.Close()
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
   //if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
   //  log.Log.Errorf("job_trace: confirm events err:%s", err)
   //}

   return errors.New("video not found")
  }

  if b := util.MapExist(mp, "MediaProcessResultSet"); b {
    log.Log.Debugf("transcode event:%+v", *event.ProcedureStateChangeEvent)
    if err := transCodeCompleteEvent(event, video); err != nil {
      session.Rollback()
      log.Log.Errorf("job_trace: transCode fail, err:%s", err)
      return err
    }
  }

  if b := util.MapExist(mp, "AiContentReviewResultSet"); b {
    log.Log.Debugf("ai review event:%+v", *event.ProcedureStateChangeEvent)
    if err := aiContentReviewEvent(event, vmodel); err != nil {
      log.Log.Errorf("job_trace: ai review fail, err:%s", err)
      session.Rollback()
      return err
    }
  }

  if err := vmodel.UpdateVideoPlayInfo(fmt.Sprint(video.VideoId)); err != nil {
    log.Log.Errorf("job_trace: update video info fail, err:%s", err)
    session.Rollback()
    return errors.New("job_trace: update video info fail")
  }

  // 如果是社区发布的视频 且 视频状态不是审核中 需要修改关联的帖子状态
  if video.PubType == 2 && vmodel.Videos.Status != 0 {
    pmodel := mposting.NewPostingModel(session)
    pmodel.Posting.Status = video.Status
    if err := pmodel.UpdatePostStatus(video.UserId, fmt.Sprint(video.VideoId)); err != nil {
      log.Log.Errorf("job_trace: update post status fail, err:%s", err)
      session.Rollback()
      return err
    }
  }

  now := time.Now().Unix()
  // 记录事件回调信息
  fileId, _ := strconv.Atoi(*event.ProcedureStateChangeEvent.FileId)
  vmodel.Events.FileId = int64(fileId)
  vmodel.Events.CreateAt = int(now)
  vmodel.Events.EventType = consts.EVENT_PROCEDURE_STATE_CHANGED_TYPE
  bts, _ = util.JsonFast.Marshal(event)
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

// todo：是否需要人工复审 确定规则 记录相关审核事件处理结果 [当前版本无相关需求]
// 视频AI 内容审核
// Type为Porn、Terrorism和Political三种类型的审核结果，分别代表视频画面令人反感的信息、视频画面令人不安全的信息和视频画面令人不适宜的信息
// suggestion pass：嫌疑度不高，建议直接通过 review：嫌疑度较高，建议人工复核 block：嫌疑度很高，建议直接屏蔽
// segments 有嫌疑的视频片段，帮助定位视频中具体哪一段涉嫌违规
// confidence 审核评分（0 - 100），评分越高，嫌疑越大
func aiContentReviewEvent(event *v20180717.EventContent, vmodel *mvideo.VideoModel) error {
  //session := dao.AppEngine.NewSession()
  //defer session.Close()
  //if err := session.Begin(); err != nil {
  //  log.Log.Errorf("job_trace: session begin err:%s", err)
  //  return err
  //}
  //
  //client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
  //vmodel := mvideo.NewVideoModel(session)
  //video := vmodel.GetVideoByFileId(*event.ProcedureStateChangeEvent.FileId)
  //if video == nil {
  //  log.Log.Errorf("job_trace: video not found, fileId:%s", *event.ProcedureStateChangeEvent.FileId)
  //  session.Rollback()
  //  // 确认事件回调
  //  //if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
  //  //  log.Log.Errorf("job_trace: confirm events err:%s", err)
  //  //}
  //
  //  return errors.New("video not found")
  //}

  if len(event.ProcedureStateChangeEvent.AiContentReviewResultSet) <= 0 {
    return errors.New("AiContentReviewResultSet Not Found")
  }

  // 视频状态默认为通过
  auditState := true
  // 记录ai审核状态 1 通过 2 不通过 3 建议复审 默认为通过
  aiState := consts.AI_AUDIT_PASS
  // 有一项不满足条件 则ai审核结果为不通过
  OutLoop:
  for _, item := range event.ProcedureStateChangeEvent.AiContentReviewResultSet {
    log.Log.Debugf("AIEvent:%v", *item)
    switch *item.Type {
    case "Porn":
      if *item.PornTask.Status != "SUCCESS" {
        continue
      }

      // 建议复审
      if *item.PornTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_REVIEW {
        aiState = consts.AI_AUDIT_REVIEW
      }

      // 建议屏蔽
      if *item.PornTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_BLOCK {
        auditState = false
        aiState = consts.AI_AUDIT_BLOCK
        break OutLoop
      }

    case "Terrorism":
      if *item.TerrorismTask.Status != "SUCCESS" {
        continue
      }

      // ai建议复审
      if *item.TerrorismTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_REVIEW {
        aiState = consts.AI_AUDIT_REVIEW
      }

      if *item.TerrorismTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_BLOCK {
        auditState = false
        aiState = consts.AI_AUDIT_BLOCK
        break OutLoop
      }

    case "Political":
      if *item.PoliticalTask.Status != "SUCCESS" {
        continue
      }

      // ai建议复审
      if *item.PoliticalTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_REVIEW {
        aiState = consts.AI_AUDIT_REVIEW
      }

      if *item.PoliticalTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_BLOCK {
        auditState = false
        aiState = consts.AI_AUDIT_BLOCK
        break OutLoop
      }

    case "Porn.Asr":
      if *item.PornAsrTask.Status != "SUCCESS" {
        continue
      }

      // ai建议复审
      if *item.PornAsrTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_REVIEW {
        aiState = consts.AI_AUDIT_REVIEW
      }

      if *item.PornAsrTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_BLOCK {
        auditState = false
        aiState = consts.AI_AUDIT_BLOCK
        break OutLoop
      }

    case "Porn.Ocr":
      if *item.PornOcrTask.Status != "SUCCESS" {
        continue
      }

      // ai建议复审
      if *item.PornOcrTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_REVIEW {
        aiState = consts.AI_AUDIT_REVIEW
      }

      if *item.PornOcrTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_BLOCK {
        auditState = false
        aiState = consts.AI_AUDIT_BLOCK
        break OutLoop
      }

    case "Terrorism.Ocr":
      if *item.TerrorismOcrTask.Status != "SUCCESS" {
        continue
      }

      // ai建议复审
      if *item.TerrorismOcrTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_REVIEW {
        aiState = consts.AI_AUDIT_REVIEW
      }

      if *item.TerrorismOcrTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_BLOCK {
        auditState = false
        aiState = consts.AI_AUDIT_BLOCK
        break OutLoop
      }

    case "Political.Asr":
      if *item.PoliticalAsrTask.Status != "SUCCESS" {
        continue
      }

      // ai建议复审
      if *item.PoliticalAsrTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_REVIEW {
        aiState = consts.AI_AUDIT_REVIEW
      }

      if *item.PoliticalAsrTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_BLOCK {
        auditState = false
        aiState = consts.AI_AUDIT_BLOCK
        break OutLoop
      }

    case "Political.Ocr":
      if *item.PoliticalOcrTask.Status != "SUCCESS" {
        continue
      }

      // ai建议复审
      if *item.PoliticalOcrTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_REVIEW {
        aiState = consts.AI_AUDIT_REVIEW
      }

      if *item.PoliticalOcrTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_BLOCK {
        auditState = false
        aiState = consts.AI_AUDIT_BLOCK
        break OutLoop
      }

    case "Prohibited.Asr":
      if *item.ProhibitedAsrTask.Status != "SUCCESS" {
        continue
      }

      // ai建议复审
      if *item.ProhibitedAsrTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_REVIEW {
        aiState = consts.AI_AUDIT_REVIEW
      }

      if *item.ProhibitedAsrTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_BLOCK {
        auditState = false
        aiState = consts.AI_AUDIT_BLOCK
        break OutLoop
      }

    case "Prohibited.Ocr":
      if *item.ProhibitedOcrTask.Status != "SUCCESS" {
        continue
      }

      // ai建议复审
      if *item.ProhibitedOcrTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_REVIEW {
        aiState = consts.AI_AUDIT_REVIEW
      }

      if *item.ProhibitedOcrTask.Output.Suggestion == consts.VIDEO_AI_AUDIT_BLOCK {
        auditState = false
        aiState = consts.AI_AUDIT_BLOCK
        break OutLoop
      }

    }
  }



  // 如果视频还在审核中的状态 则修改状态
  // 人工已审核 则 不处理
  if vmodel.Videos.Status == 0 {
    // 获取设置的审核模式
    mode := vmodel.GetAuditMode()
    // 视频审核状态 true 为 审核通过 false 为不通过
    // ai审核为通过/复审 需人工进行审核
    // ai审核为true 且当前审核模式为 ai+人工 则直接设置为通过
    if auditState == true && mode == consts.AUDIT_MODE_AI_AND_MANUAL {
      vmodel.Videos.Status = 1
    }
  }

  // 1为AI审核通过 2为AI审核不通过 3为AI建议复审
  vmodel.Videos.AiStatus = aiState


  //// 更新视频审核状态
  //if err := vmodel.UpdateVideoStatus(video.UserId, fmt.Sprint(video.VideoId)); err != nil {
  //  log.Log.Errorf("job_trace: update video status fail, err:%s", err)
  //  session.Rollback()
  //  return err
  //}

  //now := time.Now().Unix()
  //// 记录事件回调信息
  //fileId, _ := strconv.Atoi(*event.ProcedureStateChangeEvent.FileId)
  //vmodel.Events.FileId = int64(fileId)
  //vmodel.Events.UpdateAt = int(now)
  //vmodel.Events.EventType = consts.EVENT_PROCEDURE_STATE_CHANGED_TYPE
  //bts, _ := util.JsonFast.Marshal(event)
  //vmodel.Events.Event = string(bts)
  //affected, err := vmodel.RecordTencentEvent()
  //if err != nil || affected != 1 {
  //  log.Log.Errorf("job_trace: record tencent aiContentReview event err:%s, affected:%d", err, affected)
  //  session.Rollback()
  //  return errors.New("record tencent complete event fail")
  //}
  //
  //// 确认事件回调
  //if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
  //  log.Log.Errorf("job_trace: confirm events err:%s", err)
  //  session.Rollback()
  //  return errors.New("confirm event fail")
  //}
  //
  //session.Commit()

  return nil
}

// 文件删除事件 todo: 修改数据状态？
func fileDeletedEvent(event *v20180717.EventContent) error {
  client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
  session := dao.AppEngine.NewSession()
  defer session.Close()
  vmodel := mvideo.NewVideoModel(session)
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
    return errors.New("record tencent complete event fail")
  }

  if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
    log.Log.Errorf("job_trace: confirm events err:%s", err)
  }

  return nil
}


// 视频转码事件
func transCodeCompleteEvent(event *v20180717.EventContent, video *models.Videos) error {
  //session := dao.AppEngine.NewSession()
  //defer session.Close()
  //if err := session.Begin(); err != nil {
  //  log.Log.Errorf("job_trace: session begin err:%s", err)
  //  return err
  //}

  //client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
  //vmodel := mvideo.NewVideoModel(session)
  //video := vmodel.GetVideoByFileId(*event.ProcedureStateChangeEvent.FileId)
  //if video == nil {
  //  log.Log.Errorf("job_trace: video not found, fileId:%s", *event.ProcedureStateChangeEvent.FileId)
  //  session.Rollback()
  //  // 确认事件回调
  //  if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
  //  log.Log.Errorf("job_trace: confirm events err:%s", err)
  //  }
  //
  //  return errors.New("video not found")
  //}


  if len(event.ProcedureStateChangeEvent.MediaProcessResultSet) <= 0 {
    return errors.New("MediaProcessResultSet Not Exists")
  }

  list := make([]*mvideo.PlayInfo, 0)
  for _, info := range event.ProcedureStateChangeEvent.MediaProcessResultSet {
    log.Log.Debugf("info:%+v", *info)
    // todo:
    switch *info.Type {
    case "Transcode":
      if *info.TranscodeTask.ErrCode != 0 {
        log.Log.Errorf("job_trace: media process errCode:%d", *info.TranscodeTask.ErrCode)
        continue
      }

      log.Log.Infof("output:%q", *info.TranscodeTask.Output)

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

    case "AdaptiveDynamicStreaming":
      if *info.AdaptiveDynamicStreamingTask.ErrCode != 0 {
        log.Log.Errorf("job_trace: media process errCode:%d", *info.TranscodeTask.ErrCode)
        continue
      }

      playInfo := new(mvideo.PlayInfo)
      // 标清（SD）834710 HLS
      if *info.AdaptiveDynamicStreamingTask.Output.Definition == 834710 {
        playInfo.Type = "2"
      }

      // 高清（HD）819081 HLS
      if *info.AdaptiveDynamicStreamingTask.Output.Definition == 819081 {
        playInfo.Type = "3"
      }

      playInfo.Url = *info.AdaptiveDynamicStreamingTask.Output.Url
      playInfo.Duration = int64(*event.ProcedureStateChangeEvent.MetaData.Duration * 1000)
      // todo: 此为转码前视频大小
      playInfo.Size = *event.ProcedureStateChangeEvent.MetaData.Size

      list = append(list, playInfo)
    }
  }

  playBts, err := util.JsonFast.Marshal(list)
  if err != nil {
    log.Log.Errorf("job_trace: jsonFast err:%s", err)
    return err
  }

  video.PlayInfo = string(playBts)
  //if err := vmodel.UpdateVideoPlayInfo(fmt.Sprint(video.VideoId)); err != nil {
  //
  //  return errors.New("job_trace: update video play info fail")
  //}

  //now := time.Now().Unix()
  //// 记录事件回调信息
  //fileId, _ := strconv.Atoi(*event.ProcedureStateChangeEvent.FileId)
  //vmodel.Events.FileId = int64(fileId)
  //vmodel.Events.UpdateAt = int(now)
  //vmodel.Events.EventType = consts.EVENT_PROCEDURE_STATE_CHANGED_TYPE
  //bts, _ := util.JsonFast.Marshal(event)
  //vmodel.Events.Event = string(bts)
  //affected, err := vmodel.RecordTencentEvent()
  //if err != nil || affected != 1 {
  //  log.Log.Errorf("job_trace: record tencent transcode complete event err:%s, affected:%d", err, affected)
  //  session.Rollback()
  //  return errors.New("record tencent complete event fail")
  //}
  //
  //// 确认事件回调
  //if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
  //  log.Log.Errorf("job_trace: confirm events err:%s", err)
  //  session.Rollback()
  //  return errors.New("confirm event fail")
  //}

  //session.Commit()
  return nil
}

// 上传事件
func uploadEvent(event *v20180717.EventContent) error {
  session := dao.AppEngine.NewSession()
  defer session.Close()
  if err := session.Begin(); err != nil {
    log.Log.Errorf("job_trace: session begin err:%s", err)
    return err
  }

  bts, _ := util.JsonFast.Marshal(event.FileUploadEvent)
  mp, err := util.JsonStringToMap(string(bts))
  if err != nil {
    log.Log.Errorf("job_trace: jsonStringToMap err:%s", err)
    session.Rollback()
    return err
  }

  if b := util.MapExist(mp, "MediaBasicInfo"); !b {
    log.Log.Error("job_trace: MediaBasicInfo Not Exists")
    session.Rollback()
    return errors.New("MediaBasicInfo Not Exists")
  }

  client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
  if event.FileUploadEvent.MediaBasicInfo.SourceInfo == nil {
    log.Log.Error("job_trace: get source info fail")
    session.Rollback()
    return errors.New("get source info fail")
  }

  source := new(cloud.SourceContext)
  if err := util.JsonFast.Unmarshal([]byte(*event.FileUploadEvent.MediaBasicInfo.SourceInfo.SourceContext), source); err != nil {
    log.Log.Errorf("job_trace: jsonfast unmarshal event sourceContext err:%s", err)
    session.Rollback()
    return errors.New("jsonfast unmarshal event sourceContext err")
  }

  if source.UserId == "" || source.TaskId == 0 {
    log.Log.Errorf("job_trace: invalid source info, source:%+v", source)
    session.Rollback()
    return errors.New("invalid source info")
  }

  // 当前时间 - 任务开始时间 >= 10分钟 结束任务
  if time.Now().Unix() - source.Tm >= 10 * 60 {
    log.Log.Errorf("job_trace: end job, source:%+v", source)
    // 确认事件回调
    if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
      log.Log.Errorf("job_trace: confirm events err:%s", err)
    }

    session.Rollback()
    return errors.New("end job")
  }

  // 修改封面 没有视频时长
  if int(*event.FileUploadEvent.MetaData.VideoDuration) == 0 {
    log.Log.Errorf("job_trace: invalid video duration, duration:%v", *event.FileUploadEvent.MetaData.VideoDuration)
    // 确认事件回调
    if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
      log.Log.Errorf("job_trace: confirm events err:%s", err)
    }

    session.Rollback()
    return errors.New("invalid video duration")
  }

  vmodel := mvideo.NewVideoModel(session)

  // 通过任务id 获取 用户id
  userId, err := vmodel.GetUploadUserIdByTaskId(source.TaskId)
  if err != nil && err != redis.ErrNil {
    log.Log.Errorf("job_trace: invalid taskId, taskId:%d", source.TaskId)
    session.Rollback()
    return errors.New("invalid taskId")
  }

  // userId 为空 表示该上传任务已过期（三天过期）
  if userId == "" {
    log.Log.Error("job_trace: user id not exists")
    // 确认事件回调
    if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
      log.Log.Errorf("job_trace: confirm events err:%s", err)
    }
    session.Rollback()
    return errors.New("user id not exists")
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
    // 确认事件回调
    //if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
    // log.Log.Errorf("job_trace: confirm events err:%s", err)
    //}

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
  if pubInfo.Cover != "" {
    vmodel.Videos.Cover = pubInfo.Cover
  }

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
  bts, _ = util.JsonFast.Marshal(event)
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

  // 删除存储发布信息的key
  if _, err := vmodel.DelPublishInfo(userId, source.TaskId); err != nil {
    log.Log.Errorf("job_trace: del publish info key err:%s", err)
  }

  return nil
}

// 新版上传事件处理 todo: 处理帖子关联逻辑
func newUploadEvent(event *v20180717.EventContent) error {
  session := dao.AppEngine.NewSession()
  defer session.Close()
  if err := session.Begin(); err != nil {
    log.Log.Errorf("job_trace: session begin err:%s", err)
    return err
  }

  bts, _ := util.JsonFast.Marshal(event.FileUploadEvent)
  mp, err := util.JsonStringToMap(string(bts))
  if err != nil {
    log.Log.Errorf("job_trace: jsonStringToMap err:%s", err)
    session.Rollback()
    return err
  }

  if b := util.MapExist(mp, "MediaBasicInfo"); !b {
    log.Log.Error("job_trace: MediaBasicInfo Not Exists")
    session.Rollback()
    return errors.New("MediaBasicInfo Not Exists")
  }

  client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
  if event.FileUploadEvent.MediaBasicInfo.SourceInfo == nil {
    log.Log.Error("job_trace: get source info fail")
    session.Rollback()
    return errors.New("get source info fail")
  }

  source := new(cloud.SourceContext)
  if err := util.JsonFast.Unmarshal([]byte(*event.FileUploadEvent.MediaBasicInfo.SourceInfo.SourceContext), source); err != nil {
    log.Log.Errorf("job_trace: jsonfast unmarshal event sourceContext err:%s", err)
    session.Rollback()
    return errors.New("jsonfast unmarshal event sourceContext err")
  }

  if source.UserId == "" || source.TaskId == 0 {
    log.Log.Errorf("job_trace: invalid source info, source:%+v", source)
    session.Rollback()
    return errors.New("invalid source info")
  }

  // 当前时间 - 任务开始时间 >= 10分钟 结束任务
  if time.Now().Unix() - source.Tm >= 10 * 60 {
    log.Log.Errorf("job_trace: end job, source:%+v", source)
    // 确认事件回调
    if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
      log.Log.Errorf("job_trace: confirm events err:%s", err)
    }

    session.Rollback()
    return errors.New("end job")
  }

  // 修改封面 没有视频时长
  if int(*event.FileUploadEvent.MetaData.VideoDuration) == 0 {
    log.Log.Errorf("job_trace: invalid video duration, duration:%v", *event.FileUploadEvent.MetaData.VideoDuration)
    // 确认事件回调
    if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
      log.Log.Errorf("job_trace: confirm events err:%s", err)
    }

    session.Rollback()
    return errors.New("invalid video duration")
  }

  vmodel := mvideo.NewVideoModel(session)

  // 通过任务id 获取 用户id
  userId, err := vmodel.GetUploadUserIdByTaskId(source.TaskId)
  if err != nil && err != redis.ErrNil {
    log.Log.Errorf("job_trace: invalid taskId, taskId:%d", source.TaskId)
    session.Rollback()
    return errors.New("invalid taskId")
  }

  // userId 为空 表示该上传任务已过期（三天过期）
  if userId == "" {
    log.Log.Error("job_trace: user id not exists")
    // 确认事件回调
    if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
      log.Log.Errorf("job_trace: confirm events err:%s", err)
    }
    session.Rollback()
    return errors.New("user id not exists")
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

  str, err := vmodel.GetPublishInfo(source.UserId, source.TaskId)
  if err != nil || str == "" {
    log.Log.Errorf("job_trace: get publish info err:%s", err)
    // 确认事件回调
    //if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
    // log.Log.Errorf("job_trace: confirm events err:%s", err)
    //}

    session.Rollback()
    return errors.New("get publish info fail")
  }

  infos := strings.Split(str, "_")
  if len(infos) != 3 {
    log.Log.Errorf("job_trace: get publish info err:%s", err)
    session.Rollback()
    return errors.New("get publish info fail")
  }

  videoId := infos[0]
  vmodel.Videos = vmodel.FindVideoById(videoId)
  if vmodel.Videos == nil {
    log.Log.Errorf("job_trace: video not found, videoId:%s", videoId)
    session.Rollback()
    return errors.New("video not found")
  }

  // 获取用户发布的视频信息
  pubInfo := new(mvideo.VideoPublishParams)
  if err := util.JsonFast.Unmarshal([]byte(infos[1]), pubInfo); err != nil {
    log.Log.Errorf("job_trace: jsonFast unmarshal err: %s", err)
    session.Rollback()
    return errors.New("jsonFast unmarshal err")
  }

  if pubInfo.TaskId != source.TaskId {
    log.Log.Errorf("job_trace: task id not match, pub taskId:%d, source taskId:%d", pubInfo.TaskId, source.TaskId)
    session.Rollback()
    return errors.New("taskId not match")
  }

  // 数据更新 标签记录到 视频标签表（多条记录 同一个videoId对应N个labelId 生成N条记录）
  now := time.Now().Unix()
  vmodel.Videos.UserId = userId
  vmodel.Videos.Cover = *event.FileUploadEvent.MediaBasicInfo.CoverUrl
  if pubInfo.Cover != "" {
    vmodel.Videos.Cover = pubInfo.Cover
  }

  vmodel.Videos.Title = pubInfo.Title
  vmodel.Videos.Describe = pubInfo.Describe
  vmodel.Videos.VideoAddr = pubInfo.VideoAddr
  // 转为毫秒
  vmodel.Videos.VideoDuration = int(*event.FileUploadEvent.MetaData.VideoDuration * 1000)
  //vmodel.Videos.UpdateAt = int(now)
  vmodel.Videos.UpdateAt = int(now)
  vmodel.Videos.UserType = consts.PUBLISH_VIDEO_BY_USER
  vmodel.Videos.VideoWidth = *event.FileUploadEvent.MetaData.Width
  vmodel.Videos.VideoHeight = *event.FileUploadEvent.MetaData.Height
  // 单位：字节
  //vmodel.Videos.Size = *event.FileUploadEvent.MetaData.Size
  fileId, _ := strconv.Atoi(*event.FileUploadEvent.FileId)
  //vmodel.Videos.FileId = int64(fileId)
  vmodel.Videos.Size = pubInfo.Size
  // todo: 如果有 记录用户自定义标签

  // 更新视频信息
  affected, err := vmodel.UpdateVideoInfo()
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


  // 记录事件回调信息
  vmodel.Events.FileId = int64(fileId)
  vmodel.Events.CreateAt = int(now)
  vmodel.Events.EventType = consts.EVENT_UPLOAD_TYPE
  bts, _ = util.JsonFast.Marshal(event.FileUploadEvent)
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

  // 删除存储发布信息的key
  //if _, err := vmodel.DelPublishInfo(userId, source.TaskId); err != nil {
  //  log.Log.Errorf("job_trace: del publish info key err:%s", err)
  //}

  return nil
}


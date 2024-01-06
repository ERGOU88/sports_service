package job

import (
	"errors"
	"fmt"
	v20180717 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	"sports_service/dao"
	"sports_service/global/backend/log"
	"sports_service/global/consts"
	"sports_service/models/medu"
	cloud "sports_service/tools/tencentCloud"
	"sports_service/util"
	"strconv"
	"time"
)

// 课程视频处理
func CourseVideoEventsJob() {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Log.Debugf("开始拉取事件[腾讯云]")
			if err := pullEvent(); err != nil {
				log.Log.Errorf("job_trace: pull events err:%s", err)
			}
			log.Log.Debugf("事件处理完毕")
		}
	}

}

// 主动拉取事件回调
func pullEvent() error {
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
			log.Log.Debugf("job_trace: upload event:%+v, fileId:%s", *event.FileUploadEvent, *event.FileUploadEvent.FileId)
			if err := uploadEventByCourse(event); err != nil {
				log.Log.Errorf("job_trace: uploadEvent err:%s", err)
				continue
			}
			// 任务流状态变更（包含视频转码完成）
		case consts.EVENT_PROCEDURE_STATE_CHANGED:
			log.Log.Debugf("job_trace: transcode event:%+v", *event.ProcedureStateChangeEvent)
			transCodeCompleteEventByCourse(event)

		default:

		}
	}

	return nil
}

// 视频转码事件
func transCodeCompleteEventByCourse(event *v20180717.EventContent) error {
	session := dao.AppEngine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		log.Log.Errorf("job_trace: session begin err:%s", err)
		return err
	}

	client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
	cmodel := medu.NewEduModel(session)
	video := cmodel.GetCourseVideoByFileId(*event.ProcedureStateChangeEvent.FileId)
	if video == nil {
		log.Log.Errorf("job_trace: video not found, fileId:%s", *event.ProcedureStateChangeEvent.FileId)
		session.Rollback()
		return errors.New("video not found")
	}

	list := make([]*medu.PlayInfo, 0)
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
			playInfo := new(medu.PlayInfo)
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
	// 课程视频状态 由 待发布 改为 发布
	video.Status = 0
	if err := cmodel.UpdateCourseVideoPlayInfo(fmt.Sprint(video.Id)); err != nil {
		session.Rollback()
		return errors.New("job_trace: update video play info fail")
	}

	now := time.Now().Unix()
	// 记录事件回调信息
	id, _ := strconv.Atoi(*event.ProcedureStateChangeEvent.FileId)
	cmodel.Events.FileId = int64(id)
	// 记录事件回调信息
	cmodel.Events.CreateAt = int(now)
	cmodel.Events.EventType = consts.EVENT_PROCEDURE_STATE_CHANGED_TYPE
	bts, _ := util.JsonFast.Marshal(event)
	cmodel.Events.Event = string(bts)
	affected, err := cmodel.RecordTencentEvent()
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
func uploadEventByCourse(event *v20180717.EventContent) error {
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

	source := new(medu.SourceContext)
	if err := util.JsonFast.Unmarshal([]byte(*event.FileUploadEvent.MediaBasicInfo.SourceInfo.SourceContext), source); err != nil {
		log.Log.Errorf("job_trace: jsonfast unmarshal event sourceContext err:%s", err)
		session.Rollback()
		return errors.New("jsonfast unmarshal event sourceContext err")
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

	cmodel := medu.NewEduModel(session)
	// 通过任务id 获取 腾讯云文件id
	fileId, err := cmodel.GetUploadFileIdByTaskId(source.EduTaskId)
	if err != nil || fileId == "" {
		log.Log.Errorf("job_trace: invalid taskId, taskId:%d", source.EduTaskId)
		session.Rollback()
		return errors.New("invalid taskId")
	}

	// 腾讯云文件id 查询 课程视频
	video := cmodel.GetCourseVideoByFileId(fileId)
	if video == nil {
		log.Log.Errorf("job_trace: get course video by tx fileId err:%s, fileId:%s", err, fileId)
		// 确认事件回调
		//if err := client.ConfirmEvents([]string{*event.EventHandle}); err != nil {
		//  log.Log.Errorf("job_trace: confirm events err:%s", err)
		//}
		session.Rollback()
		return errors.New("get publish info fail")
	}

	now := time.Now().Unix()
	video.Cover = *event.FileUploadEvent.MediaBasicInfo.CoverUrl
	// 转为毫秒
	video.VideoDuration = int(*event.FileUploadEvent.MetaData.VideoDuration * 1000)
	video.UpdateAt = int(now)
	video.VideoWidth = *event.FileUploadEvent.MetaData.Width
	video.VideoHeight = *event.FileUploadEvent.MetaData.Height
	// 单位：字节
	video.Size = *event.FileUploadEvent.MetaData.Size
	// 课程视频状态 由 待发布 1 改为 发布 0
	video.Status = 0

	log.Log.Errorf("job_trace: video Info:%+v", video)

	// 更新课程视频信息
	affected, err := cmodel.UpdateCourseVideoInfo(fmt.Sprint(video.Id))
	if err != nil || affected != 1 {
		log.Log.Errorf("job_trace: publish video err:%s, affected:%d", err, affected)
		session.Rollback()
		return errors.New("publish video fail")
	}

	// 通过id获取课程
	course := cmodel.GetCourseById(fmt.Sprint(video.CourseId))
	if course != nil {
		// 课程状态为1待发布 则 修改为 0正常
		if course.Status == 1 {
			condition := fmt.Sprintf("id=%d", video.CourseId)
			cols := "status, update_at"
			cmodel.Course.Status = 0
			cmodel.Course.UpdateAt = int(now)
			// 课程视频回调已处理完毕  将课程状态从待发布设置为已发布
			affected, err = cmodel.UpdateCourseInfo(condition, cols)
			if affected != 1 || err != nil {
				log.Log.Errorf("job_trace: update course info fail, err:%s, affected:%d", err, affected)
				session.Rollback()
				return errors.New("update course info fail")
			}
		}
	}

	// 记录事件回调信息
	id, _ := strconv.Atoi(*event.FileUploadEvent.FileId)
	cmodel.Events.FileId = int64(id)
	cmodel.Events.CreateAt = int(now)
	cmodel.Events.EventType = consts.EVENT_UPLOAD_TYPE
	bts, _ = util.JsonFast.Marshal(event)
	cmodel.Events.Event = string(bts)
	affected, err = cmodel.RecordTencentEvent()
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

	// 删除记录腾讯文件id的key
	if _, err := cmodel.DelRecordFileIdKey(source.EduTaskId); err != nil {
		log.Log.Errorf("job_trace: del record file id key err:%s", err)
	}

	return nil
}

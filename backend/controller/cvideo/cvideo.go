package cvideo

import (
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/go-xorm/xorm"
  "sports_service/server/dao"
  "sports_service/server/global/backend/log"
  "sports_service/server/global/backend/errdef"
  "sports_service/server/global/consts"
  "sports_service/server/models"
  "sports_service/server/models/mattention"
  "sports_service/server/models/mlabel"
  "sports_service/server/models/muser"
  "sports_service/server/models/mvideo"
  redismq "sports_service/server/redismq/event"
  "sports_service/server/tools/tencentCloud"
  "sports_service/server/util"
  "time"
)

type VideoModule struct {
  context      *gin.Context
  engine       *xorm.Session
  video        *mvideo.VideoModel
  label        *mlabel.LabelModel
  attention    *mattention.AttentionModel
  user         *muser.UserModel
}

func New(c *gin.Context) VideoModule {
  socket := dao.AppEngine.NewSession()
  defer socket.Close()
  return VideoModule{
    context: c,
    video: mvideo.NewVideoModel(socket),
    label: mlabel.NewLabelModel(socket),
    attention: mattention.NewAttentionModel(socket),
    user: muser.NewUserModel(socket),
    engine: socket,
  }
}

// 修改视频状态
func (svc *VideoModule) EditVideoStatus(param *mvideo.EditVideoStatusParam) int {
  if param.VideoId == "" || param.Status == 0 {
    return errdef.INVALID_PARAMS
  }

  if err := svc.engine.Begin(); err != nil {
    return errdef.ERROR
  }

  video := svc.video.FindVideoById(param.VideoId)
  if video == nil {
    svc.engine.Rollback()
    return errdef.VIDEO_NOT_EXISTS
  }

  status := fmt.Sprint(video.Status)
  // 视频已删除
  if status == consts.VIDEO_DELETE_STATUS {
    svc.engine.Rollback()
    return errdef.VIDEO_ALREADY_DELETE
  }

  // 视频已通过审核 只能执行逻辑删除
  if status == consts.VIDEO_AUDIT_SUCCESS && fmt.Sprint(param.Status) != consts.VIDEO_DELETE_STATUS {
    svc.engine.Rollback()
    return errdef.VIDEO_ALREADY_PASS
  }

  // 通过 / 不通过 / 执行删除操作 且 视频状态为审核通过 则只能逻辑删除/不通过 直接更新视频状态
  if fmt.Sprint(param.Status) == consts.VIDEO_AUDIT_SUCCESS || fmt.Sprint(param.Status) == consts.VIDEO_AUDIT_FAILURE ||
    (fmt.Sprint(param.Status) == consts.VIDEO_DELETE_STATUS && status == consts.VIDEO_AUDIT_SUCCESS) {
    video.Status = int(param.Status)
    // 更新视频状态
    if err := svc.video.UpdateVideoStatus(video.UserId, param.VideoId); err != nil {
      svc.engine.Rollback()
      return errdef.VIDEO_EDIT_STATUS_FAIL
    }

    // 如果是逻辑删除 则 同时需要修改视频标签状态
    if fmt.Sprint(param.Status) == consts.VIDEO_DELETE_STATUS {
      condition := fmt.Sprintf("video_id=%s", param.VideoId)
      cols := fmt.Sprintf("status")
      svc.video.Labels.Status = 0
      if _, err := svc.video.UpdateVideoLabelInfo(condition, cols); err != nil {
        svc.engine.Rollback()
        return errdef.VIDEO_EDIT_STATUS_FAIL
      }
    }

    // 如果是审核通过
    if fmt.Sprint(param.Status) == consts.VIDEO_AUDIT_SUCCESS {
      condition := fmt.Sprintf("video_id=%s", param.VideoId)
      cols := fmt.Sprintf("status")
      // 将视频标签置为可用
      svc.video.Labels.Status = 1
      if _, err := svc.video.UpdateVideoLabelInfo(condition, cols); err != nil {
        svc.engine.Rollback()
        return errdef.VIDEO_EDIT_STATUS_FAIL
      }

      // 获取发布者用户信息
      user := svc.user.FindUserByUserid(video.UserId)
      if user != nil {
        // 获取发布者粉丝们的userId
        userIds := svc.attention.GetFansList(user.UserId)
        for _, userId := range userIds {
          // 给发布者的粉丝 发送 发布新视频推送
          //event.PushEventMsg(config.Global.AmqpDsn, userId, user.NickName, video.Cover, "", consts.FOCUS_USER_PUBLISH_VIDEO_MSG)
          redismq.PushEventMsg(redismq.NewEvent(userId, fmt.Sprint(video.VideoId), user.NickName, video.Cover, "", consts.FOCUS_USER_PUBLISH_VIDEO_MSG))
        }
      }
    }

    svc.engine.Commit()
    return errdef.SUCCESS
  }

  // 如果执行删除操作 且 视频状态未审核通过 删除相关所有数据
  if fmt.Sprint(param.Status) == consts.VIDEO_DELETE_STATUS && status != consts.VIDEO_AUDIT_SUCCESS {

    // 视频为未审核/审核失败 物理删除发布的视频、视频标签、视频总计
    if err := svc.video.DelPublishById(video.UserId, param.VideoId); err != nil {
      svc.engine.Rollback()
      return errdef.VIDEO_DELETE_PUBLISH_FAIL
    }

    // 删除视频标签
    if err := svc.video.DelVideoLabels(param.VideoId); err != nil {
      svc.engine.Rollback()
      return errdef.VIDEO_DELETE_PUBLISH_FAIL
    }

    // 删除视频统计数据
    if err := svc.video.DelVideoStatistic(param.VideoId); err != nil {
      svc.engine.Rollback()
      return errdef.VIDEO_DELETE_PUBLISH_FAIL
    }

  }

  svc.engine.Commit()

  return errdef.SUCCESS
}

// 获取视频列表
func (svc *VideoModule) GetVideoList(page, size int) []*mvideo.VideoDetailInfo {
  offset := (page - 1) * size
  list := svc.video.GetVideoList(offset, size)
  if len(list) == 0 {
    return []*mvideo.VideoDetailInfo{}
  }
  for _, video := range list {
    video.Labels = svc.video.GetVideoLabels(fmt.Sprint(video.VideoId))
    video.VideoAddr = svc.video.AntiStealingLink(video.VideoAddr)
  }

  return list
}

// 获取已审核通过的视频总数
func (svc *VideoModule) GetVideoTotalCount() int64 {
  return svc.video.GetVideoTotalCount()
}

// 获取未审核/审核失败的视频总数
func (svc *VideoModule) GetVideoReviewTotalCount() int64 {
  return svc.video.GetVideoReviewTotalCount()
}

// 修改视频置顶状态
func (svc *VideoModule) EditVideoTopStatus(param *mvideo.EditTopStatusParam) int {
  if param.VideoId == "" {
    return errdef.INVALID_PARAMS
  }

  if param.Status != consts.VIDEO_IS_TOP && param.Status != consts.VIDEO_NOT_TOP {
    return errdef.INVALID_PARAMS
  }

  video := svc.video.FindVideoById(param.VideoId)
  if video == nil {
    return errdef.VIDEO_NOT_EXISTS
  }

  // 视频未审核成功 不能设置置顶
  if fmt.Sprint(video.Status) != consts.VIDEO_AUDIT_SUCCESS {
    return errdef.VIDEO_NOT_PASS
  }

  // 更新视频置顶状态
  video.IsTop = int(param.Status)
  if err := svc.video.UpdateVideoTopStatus(param.VideoId); err != nil {
    return errdef.VIDEO_EDIT_TOP_FAIL
  }

  return errdef.SUCCESS
}

// 编辑视频推荐状态
func (svc *VideoModule) EditVideoRecommendStatus(param *mvideo.EditRecommendStatusParam) int {
  if param.VideoId == "" {
    return errdef.INVALID_PARAMS
  }

  if param.Status != consts.VIDEO_IS_RECOMMEND && param.Status != consts.VIDEO_NOT_RECOMMEND {
    return errdef.INVALID_PARAMS
  }

  video := svc.video.FindVideoById(param.VideoId)
  if video == nil {
    return errdef.VIDEO_NOT_EXISTS
  }

  // 视频未审核成功 不能设置推荐
  if fmt.Sprint(video.Status) != consts.VIDEO_AUDIT_SUCCESS {
    return errdef.VIDEO_NOT_PASS
  }

  // 更新视频推荐状态
  video.IsRecommend = int(param.Status)
  if err := svc.video.UpdateVideoRecommendStatus(param.VideoId); err != nil {
    return errdef.VIDEO_EDIT_RECOMMEND_FAIL
  }

  return errdef.SUCCESS
}

// 获取审核中/审核失败的视频列表
func (svc *VideoModule) GetVideoReviewList(page, size int) []*models.Videos {
  offset := (page - 1) * size
  list := svc.video.GetVideoReviewList(offset, size)
  if list == nil {
    return []*models.Videos{}
  }

  for _, v := range list {
    v.VideoAddr  = svc.video.AntiStealingLink(v.VideoAddr)
  }

  return list
}

// 获取视频标签列表
func (svc *VideoModule) GetVideoLabelList() []*mlabel.VideoLabel {
  return svc.label.GetVideoLabelList()
}

// 添加视频标签
func (svc *VideoModule) AddVideoLabel(param *mlabel.AddVideoLabelParam) int {
  client := tencentCloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
  // 检测热搜内容
  isPass, err := client.TextModeration(param.LabelName)
  if !isPass || err != nil {
    return errdef.VIDEO_INVALID_LABEL_NAME
  }

  lenth := util.GetStrLen([]rune(param.LabelName))
  if lenth <= 0 || lenth > 20 {
    return errdef.VIDEO_INVALID_LABEL_LEN
  }

  now := time.Now().Unix()
  svc.label.VideoLabels.UpdateAt = int(now)
  svc.label.VideoLabels.LabelName = param.LabelName
  svc.label.VideoLabels.CreateAt = int(now)
  svc.label.VideoLabels.Icon = param.Icon
  svc.label.VideoLabels.Sortorder = param.Sortorder
  svc.label.VideoLabels.Status = 1
  if err := svc.label.AddVideoLabel(); err != nil {
    return errdef.VIDEO_INVALID_LABEL_NAME
  }

  // 添加到内存中
  svc.label.AddLabelInfoByMem()
  svc.label.CleanLabelInfoByMem()

  return errdef.SUCCESS
}

// 删除视频标签
func (svc *VideoModule) DelVideoLabel(labelId string) int {
  if info := svc.label.GetLabelInfoByMem(labelId); info == nil {
    return errdef.VIDEO_LABEL_NOT_EXISTS
  }

  if err := svc.label.DelVideoLabel(labelId); err != nil {
    return errdef.VIDEO_LABEL_DELETE_FAIL
  }

  // 从内存中删除
  svc.label.DelLabelInfoByMem(labelId)
  svc.label.CleanLabelInfoByMem()
  return errdef.SUCCESS
}

// 添加视频分区
func (svc *VideoModule) AddVideoSubarea(param *mvideo.AddSubarea) int {
  svc.video.Subarea.SubareaName = param.Name
  svc.video.Subarea.Sortorder = param.SortOrder
  if _, err := svc.video.AddSubArea(); err != nil {
    log.Log.Errorf("")
    return errdef.VIDEO_ADD_SUBAREA_FAIL
  }

  return errdef.SUCCESS
}

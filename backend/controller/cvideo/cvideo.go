package cvideo

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/dao"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mlabel"
	"sports_service/server/models/mvideo"
	"github.com/go-xorm/xorm"
	"fmt"
	"sports_service/server/util"
	"time"
)

type VideoModule struct {
	context      *gin.Context
	engine       *xorm.Session
	video        *mvideo.VideoModel
	label        *mlabel.LabelModel
}

func New(c *gin.Context) VideoModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return VideoModule{
		context: c,
		video: mvideo.NewVideoModel(socket),
		label: mlabel.NewLabelModel(socket),
		engine: socket,
	}
}

// 修改视频状态
func (svc *VideoModule) EditVideoStatus(param *mvideo.EditVideoStatusParam) int {
	if param.VideoId == "" || param.Status == 0 {
		return errdef.INVALID_PARAMS
	}

	video := svc.video.FindVideoById(param.VideoId)
	if video == nil {
		return errdef.VIDEO_NOT_EXISTS
	}

	status := fmt.Sprint(video.Status)
	// 视频已删除
	if status == consts.VIDEO_DELETE_STATUS {
		return errdef.VIDEO_ALREADY_DELETE
	}

	// 视频已通过审核 只能执行删除
	if status == consts.VIDEO_AUDIT_SUCCESS && fmt.Sprint(param.Status) != consts.VIDEO_DELETE_STATUS {
		return errdef.VIDEO_ALREADY_PASS
	}

	// 通过 / 不通过 / 执行删除操作 且 视频状态为审核通过 则只能逻辑删除 直接更新视频状态
	if fmt.Sprint(param.Status) == consts.VIDEO_AUDIT_SUCCESS || fmt.Sprint(param.Status) == consts.VIDEO_AUDIT_FAILURE ||
		(fmt.Sprint(param.Status) == consts.VIDEO_DELETE_STATUS && status == consts.VIDEO_AUDIT_SUCCESS) {
		video.Status = int(param.Status)
		// 更新视频状态
		if err := svc.video.UpdateVideoStatus(video.UserId, param.VideoId); err != nil {
			return errdef.VIDEO_EDIT_STATUS_FAIL
		}
	}

	// 如果执行删除操作 且 视频状态未审核通过 删除相关所有数据
	if fmt.Sprint(param.Status) == consts.VIDEO_DELETE_STATUS && status != consts.VIDEO_AUDIT_SUCCESS {
		if err := svc.engine.Begin(); err != nil {
			return errdef.ERROR
		}

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

		svc.engine.Commit()
	}

	return errdef.SUCCESS
}

// 获取视频列表
func (svc *VideoModule) GetVideoList(page, size int) []*mvideo.VideoDetailInfo {
	offset := (page - 1) * size
	list := svc.video.GetVideoList(offset, size)
	for _, video := range list {
    video.Labels = svc.video.GetVideoLabels(fmt.Sprint(video.VideoId))
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
 	return svc.video.GetVideoReviewList(offset, size)
}

// 获取视频标签列表
func (svc *VideoModule) GetVideoLabelList() []*mlabel.VideoLabel {
	return svc.label.GetVideoLabelList()
}

// 添加视频标签 todo:脏词过滤
func (svc *VideoModule) AddVideoLabel(param *mlabel.AddVideoLabelParam) int {
	lenth := util.GetStrLen([]rune(param.LabelName))
	if lenth <= 0 || lenth > 20 {
		return errdef.VIDEO_INVALID_LABEL_NAME
	}

	now := time.Now().Unix()
	svc.label.VideoLabels.UpdateAt = int(now)
	svc.label.VideoLabels.LabelName = param.LabelName
	svc.label.VideoLabels.CreateAt = int(now)
	svc.label.VideoLabels.Icon = param.Icon
	if err := svc.label.AddVideoLabel(); err != nil {
		return errdef.VIDEO_INVALID_LABEL_NAME
	}

	// 添加到内存中
	svc.label.AddLabelInfoByMem()

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

	return errdef.SUCCESS
}

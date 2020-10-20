package mvideo

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"fmt"
)

// todo: 视频id自增 帖子可以使用分布式唯一id
type VideoModel struct {
	Videos    *models.Videos
	Engine    *xorm.Session
	Browse    *models.UserBrowseRecord
	Labels    *models.VideoLabels
	Statistic *models.VideoStatistic
	Events    *models.TencentCloudEvents
}

// 视频发布请求参数
type VideoPublishParams struct {
	TaskId         int64   `binding:"required" json:"task_id"`        // 任务id
	Cover          string  `json:"cover"`                             // 视频封面
	Title          string  `binding:"required" json:"title"`          // 视频标题
	Describe       string  `binding:"required" json:"describe"`       // 视频描述
	VideoAddr      string  `binding:"required" json:"video_addr"`     // 视频地址
	VideoDuration  int     `binding:"required" json:"video_duration"` // 视频时长
	FileId         string  `binding:"required" json:"file_id"`        // 腾讯云文件id
	VideoWidth     int64   `binding:"required" json:"video_width"`    // 视频宽
	VideoHeight    int64   `binding:"required" json:"video_height"`   // 视频高
	VideoLabels    string  `binding:"required" json:"video_labels"`   // 视频标签id（多个用逗号分隔）
	Size           int64   `json:"size"`                              // 视频字节数
	CustomLabels   string  `json:"custom_labels"`                     // 字符串（多个用逗号分隔）
}

// 视频信息
type VideosInfoResp struct {
	VideoId       int64  `json:"video_id"`       // 视频id
	Title         string `json:"title"`          // 标题
	Describe      string `json:"describe"`       // 描述
	Cover         string `json:"cover"`          // 封面
	VideoAddr     string `json:"video_addr"`     // 视频地址
	IsRecommend   int    `json:"is_recommend"`   // 是否推荐
	IsTop         int    `json:"is_top"`         // 是否置顶
	VideoDuration int    `json:"video_duration"` // 视频时长
	VideoWidth    int64  `json:"video_width"`    // 视频宽
	VideoHeight   int64  `json:"video_height"`   // 视频高
	CreateAt      int    `json:"create_at"`      // 视频创建时间
	UserId        string `json:"user_id"`        // 发布视频的用户id
	Avatar        string `json:"avatar"`         // 头像
	Nickname      string `json:"nick_name"`      // 昵称
	IsAttention   int    `json:"is_attention"`   // 是否关注 1 关注 2 未关注
	OpTime        int    `json:"op_time"`        // 用户收藏/点赞等的操作时间
}

// 视频信息
type VideosInfo struct {
	VideoId       int64  `json:"video_id"`       // 视频id
	Title         string `json:"title"`          // 标题
	Describe      string `json:"describe"`       // 描述
	Cover         string `json:"cover"`          // 封面
	VideoAddr     string `json:"video_addr"`     // 视频地址
	IsRecommend   int    `json:"is_recommend"`   // 是否推荐
	IsTop         int    `json:"is_top"`         // 是否置顶
	VideoDuration int    `json:"video_duration"` // 视频时长
	VideoWidth    int64  `json:"video_width"`    // 视频宽
	VideoHeight   int64  `json:"video_height"`   // 视频高
	Status        int32  `json:"status"`         // 审核状态
	CreateAt      int    `json:"create_at"`      // 视频创建时间
	FabulousNum   int    `json:"fabulous_num"`   // 点赞数
	CommentNum    int    `json:"comment_num"`    // 评论数
	ShareNum      int    `json:"share_num"`      // 分享数
	BrowseNum     int    `json:"browse_num"`     // 浏览数（播放数）
	BarrageNum    int    `json:"barrage_num"`    // 弹幕数
  TimeElapsed   int    `json:"time_elapsed"`   // 已播放的时长 毫秒
  StatusCn      string `json:"status_cn"`      // 审核状态（中文展示）
}

// 视频详情信息
type VideoDetailInfo struct {
	VideoId       int64                 `json:"video_id"`       // 视频id
	Title         string                `json:"title"`          // 标题
	Describe      string                `json:"describe"`       // 描述
	Cover         string                `json:"cover"`          // 封面
	VideoAddr     string                `json:"video_addr"`     // 视频地址
	IsRecommend   int                   `json:"is_recommend"`   // 是否推荐
	IsTop         int                   `json:"is_top"`         // 是否置顶
	VideoDuration int                   `json:"video_duration"` // 视频时长
	VideoWidth    int64                 `json:"video_width"`    // 视频宽
	VideoHeight   int64                 `json:"video_height"`   // 视频高
	Status        int32                 `json:"status"`         // 审核状态
	CreateAt      int                   `json:"create_at"`      // 视频创建时间
	FabulousNum   int                   `json:"fabulous_num"`   // 点赞数
	CommentNum    int                   `json:"comment_num"`    // 评论数
	BarrageNum    int                   `json:"barrage_num"`    // 弹幕数
	ShareNum      int                   `json:"share_num"`      // 分享数
	BrowseNum     int                   `json:"browse_num"`     // 浏览数（播放数）
	UserId        string                `json:"user_id"`        // 发布视频的用户id
	Avatar        string                `json:"avatar"`         // 头像
	Nickname      string                `json:"nick_name"`      // 昵称
	IsAttention   int                   `json:"is_attention"`   // 是否关注 1 关注 2 未关注
	IsCollect     int                   `json:"is_collect"`     // 是否收藏
	IsLike        int                   `json:"is_like"`        // 是否点赞
	FansNum       int64                 `json:"fans_num"`       // 粉丝数
	Labels        []*models.VideoLabels `json:"labels"`         // 视频标签
}

// 删除历史记录请求参数
type DeleteHistoryParam struct {
	ComposeIds        []string     `binding:"required" json:"composeIds"` // 作品id列表
}

// 删除发布记录请求参数(不支持批量删除)
type DeletePublishParam struct {
	ComposeIds        string       `binding:"required" json:"composeIds"` // 作品id
}

// 后台修改视频状态请求参数
type EditVideoStatusParam struct {
	VideoId       string     `json:"video_id"`     // 视频id
	Status        int32      `json:"status"`       // 状态 1：审核通过 2：审核不通过 3：逻辑删除
}

// 修改视频置顶状态
type EditTopStatusParam struct {
	VideoId      string     `json:"video_id"`      // 视频id
	Status       int32      `json:"status"`        // 状态 0 不置顶 1 置顶
}

// 修改视频推荐状态
type EditRecommendStatusParam struct {
	VideoId      string     `json:"video_id"`      // 视频id
	Status       int32      `json:"status"`        // 状态 0 不推荐 1 推荐
}

// 自定义标签请求
type CustomLabelParams struct {
  CustomLabel   string    `json:"custom_label"`   // 自定义标签
}

// 实栗
func NewVideoModel(engine *xorm.Session) *VideoModel {
	return &VideoModel{
		Browse: new(models.UserBrowseRecord),
		Videos: new(models.Videos),
		Labels: new(models.VideoLabels),
		Statistic: new(models.VideoStatistic),
		Events: new(models.TencentCloudEvents),
		Engine: engine,
	}
}

// 视频发布
func (m *VideoModel) VideoPublish() (int64, error) {
	return m.Engine.InsertOne(m.Videos)
}

// 添加视频统计数据
func (m *VideoModel) AddVideoStatistic() error {
	if _, err := m.Engine.InsertOne(m.Statistic); err != nil {
		return err
	}

	return nil
}

// 获取视频统计数据
func (m *VideoModel) GetVideoStatistic(videoId string) *models.VideoStatistic {
	m.Statistic = new(models.VideoStatistic)
	ok, err := m.Engine.Where("video_id=?", videoId).Get(m.Statistic)
	if !ok || err != nil {
		log.Log.Errorf("video_trace: get video statistic info err:%s", err)
		return nil
	}

	return m.Statistic
}

// 删除视频统计数据
func (m *VideoModel) DelVideoStatistic(videoId string) error {
	if _, err := m.Engine.Where("video_id=?", videoId).Delete(&models.VideoStatistic{}); err != nil {
		return err
	}

	return nil
}

const (
	UPDATE_VIDEO_LIKE_NUM  = "UPDATE `video_statistic` SET `fabulous_num` = `fabulous_num` + ?, `update_at`=? WHERE `video_id`=? AND `fabulous_num` + ? >= 0 LIMIT 1"
)
// 更新视频点赞数
func (m *VideoModel) UpdateVideoLikeNum(videoId int64, now, num int) error {
	if _, err := m.Engine.Exec(UPDATE_VIDEO_LIKE_NUM, num, now, videoId, num); err != nil {
		return err
	}

	return nil
}

const (
	UPDATE_VIDEO_COLLECT_NUM  = "UPDATE `video_statistic` SET `collect_num` = `collect_num` + ?, `update_at`=? WHERE `video_id`=? AND `collect_num` + ? >= 0 LIMIT 1"
)
// 更新视频收藏数
func (m *VideoModel) UpdateVideoCollectNum(videoId int64, now, num int) error {
	if _, err := m.Engine.Exec(UPDATE_VIDEO_COLLECT_NUM, num, now, videoId, num); err != nil {
		return err
	}

	return nil
}

const (
	UPDATE_VIDEO_COMMENT_NUM = "UPDATE `video_statistic` SET `comment_num` = `comment_num` + ?, `update_at`=? WHERE `video_id`=? AND `comment_num` + ? >= 0 LIMIT 1"
)
// 更新视频评论数
func (m *VideoModel) UpdateVideoCommentNum(videoId int64, now, num int) error {
	if _, err := m.Engine.Exec(UPDATE_VIDEO_COMMENT_NUM, num, now, videoId, num); err != nil {
		return err
	}

	return nil
}

const (
	UPDATE_VIDEO_BROWSE_NUM  = "UPDATE `video_statistic` SET `browse_num` = `browse_num` + ?, `update_at`=? WHERE `video_id`=? AND `browse_num` + ? >= 0 LIMIT 1"
)
// 更新视频浏览数
func (m *VideoModel) UpdateVideoBrowseNum(videoId int64, now, num int) error {
	if _, err := m.Engine.Exec(UPDATE_VIDEO_BROWSE_NUM, num, now, videoId, num); err != nil {
		return err
	}

	return nil
}

// 更新视频分享数
func (m *VideoModel) UpdateVideoShareNum() {
	return
}

// 更新视频游币数
func (m *VideoModel) UpdateVideoYcoinNum() {
	return
}

const (
	UPDATE_VIDEO_BARRAGE_NUM  = "UPDATE `video_statistic` SET `barrage_num` = `barrage_num` + ?, `update_at`=? WHERE `video_id`=? AND `barrage_num` + ? >= 0 LIMIT 1"
)
// 更新视频弹幕数
func (m *VideoModel) UpdateVideoBarrageNum(videoId int64, now, num int) error {
	if _, err := m.Engine.Exec(UPDATE_VIDEO_BARRAGE_NUM, num, now, videoId, num); err != nil {
		return err
	}

	return nil
}

// 分页获取 用户发布的视频列表[通过审核状态和条件查询]
func (m *VideoModel) GetUserPublishVideos(offset, size int, userId, status, field string) []*VideosInfo {
	var list []*VideosInfo

	sql := "SELECT v.*, s.fabulous_num, s.share_num, s.comment_num, s.browse_num, s.barrage_num FROM videos as v " +
		"LEFT JOIN video_statistic as s ON v.video_id=s.video_id WHERE v.user_id=? "
	if status != consts.VIDEO_VIEW_ALL {
		sql += "AND v.`status` = ? "
	} else {
		sql += "AND v.`status` != 3"
	}

	// 条件为默认时间倒序 则使用videos表的时间字段
	if field == consts.VIDEO_CONDITION_TIME {
		sql += fmt.Sprintf("GROUP BY v.video_id ORDER BY v.%s DESC, v.sortorder DESC LIMIT ?, ?", field, offset, size)
	} else {
		sql += fmt.Sprintf("GROUP BY v.video_id ORDER BY s.%s DESC, v.sortorder DESC LIMIT ?, ?", field, offset, size)
	}

	if err := m.Engine.SQL(sql, userId, offset, size).Find(&list); err != nil {
		log.Log.Errorf("video_trace: get user publish videos err:%s", err)
		return nil
	}

	return list
}

// 通过id查询视频
func (m *VideoModel) FindVideoById(videoId string) *models.Videos {
	m.Videos = new(models.Videos)
	ok, err := m.Engine.Where("video_id=?", videoId).Get(m.Videos)
	if !ok || err != nil {
		return nil
	}

	return m.Videos
}

// 通过视频id查询视频列表
func (m *VideoModel) FindVideoListByIds(videoIds string) []*models.Videos {
	var list []*models.Videos
	sql := fmt.Sprintf("SELECT * FROM videos WHERE video_id in(%s) AND status=1 ORDER BY is_top DESC, is_recommend DESC, sortorder DESC, video_id", videoIds)
	if err := m.Engine.SQL(sql).Find(&list); err != nil {
		log.Log.Errorf("video_trace: get video list err:%s", err)
		return nil
	}

	return list
}

type BrowseVideos struct {
	ComposeId       int64     `json:"compose_id"`    // 视频id
	UpdateAt        int       `json:"update_at"`     // 浏览时间
}
// 获取浏览过的视频id记录
func (m *VideoModel) GetBrowseVideosRecord(userId string, composeType, offset, size int) []*BrowseVideos {
	var list []*BrowseVideos
	if err := m.Engine.Table(&models.UserBrowseRecord{}).Where("user_id=? AND compose_type=?", userId, composeType).
		Cols("compose_id, update_at").
		Desc("id").
		Limit(size, offset).
		Find(&list); err != nil {
		log.Log.Errorf("video_trace: get browse video record err:%s", err)
		return nil
	}

	return list
}

// 通过id列表删除浏览的历史记录
func (m *VideoModel) DeleteHistoryByIds(userId string, ids string) error {
	sql := fmt.Sprintf("DELETE FROM `user_browse_record` WHERE user_id=? AND compose_id in(%s)", ids)
	if _, err := m.Engine.Exec(sql, userId); err != nil {
		return err
	}

	return nil
}

const (
	DELETE_PUBLISH_SQL = "DELETE FROM `videos` WHERE `user_id`=? AND video_id=?"
)
// 删除发布的记录
func (m *VideoModel) DelPublishById(userId, videoId string) error {
	if _, err := m.Engine.Exec(DELETE_PUBLISH_SQL, userId, videoId); err != nil {
		return err
	}

	return nil
}

// 更新视频状态
func (m *VideoModel) UpdateVideoStatus(userId, videoId string) error {
	if _, err := m.Engine.Where("user_id=? AND video_id=?", userId, videoId).Cols("status").Update(m.Videos); err != nil {
		return err
	}

	return nil
}

// 修改视频置顶状态 0 不置顶 1 置顶 （置顶权重 > 推荐 ）
func (m *VideoModel) UpdateVideoTopStatus(videoId string) error {
	if _, err := m.Engine.Where("video_id=?", videoId).Cols("is_top").Update(m.Videos); err != nil {
		return err
	}

	return nil
}

// 修改视频推荐状态 0 不推荐 1 推荐
func (m *VideoModel) UpdateVideoRecommendStatus(videoId string) error {
	if _, err := m.Engine.Where("video_id=?", videoId).Cols("is_recommend").Update(m.Videos); err != nil {
		return err
	}

	return nil
}

// 获取用户总发布数 (审核通过的)
func (m *VideoModel) GetTotalPublish(userId string) int64 {
	total, err := m.Engine.Where("user_id=? AND status=1", userId).Count(m.Videos)
	if err != nil {
		log.Log.Errorf("video_trace: get user total publish err:%s, uid:%s", err, userId)
		return 0
	}

	return total
}

const (
	QUERY_VIDEO_LIST = "SELECT v.*, s.fabulous_num,s.share_num,s.comment_num, s.browse_num FROM `videos` as v " +
		"LEFT JOIN video_statistic as s ON v.video_id=s.video_id WHERE v.status = 1 GROUP BY v.video_id " +
		"ORDER BY v.is_top DESC, v.is_recommend DESC, v.sortorder DESC, v.video_id DESC LIMIT ?, ?"
)
// 获取推荐的视频列表 todo:
func (m *VideoModel) GetVideoList(offset, size int) []*VideoDetailInfo {
	var list []*VideoDetailInfo
	if err := m.Engine.SQL(QUERY_VIDEO_LIST, offset, size).Find(&list); err != nil {
		log.Log.Errorf("video_trace: get recommend videos err:%s", err)
		return nil
	}

	return list
}

// 获取视频总数（已审核通过的）
func (m *VideoModel) GetVideoTotalCount() int64 {
  count, err := m.Engine.Where("status=1").Count(&models.Videos{})
  if err != nil {
    return 0
  }

  return count
}

// 获取视频总数（未审核/未通过审核）
func (m *VideoModel) GetVideoReviewTotalCount() int64 {
  count, err := m.Engine.Where("status = 0 or status = 2").Count(&models.Videos{})
  if err != nil {
    return 0
  }

  return count
}

const (
	QUERY_ATTENTION_VIDEOS = "SELECT v.*, s.fabulous_num,s.share_num,s.comment_num, s.browse_num FROM `videos` as v " +
		"LEFT JOIN video_statistic as s ON v.video_id=s.video_id WHERE v.status = 1 AND v.user_id in(%s) GROUP BY v.video_id " +
		"ORDER BY v.is_top DESC, v.is_recommend DESC, v.sortorder DESC, v.video_id DESC LIMIT ?, ?"
)
// 获取关注的用户发布的视频
func (m *VideoModel) GetAttentionVideos(userIds string, offset, size int) []*VideoDetailInfo {
	sql := fmt.Sprintf(QUERY_ATTENTION_VIDEOS, userIds)
	var list []*VideoDetailInfo
	if err := m.Engine.SQL(sql, offset, size).Find(&list); err != nil {
		log.Log.Errorf("video_trace: get attention videos err:%s", err)
		return nil
	}

	return list
}

// 搜索视频
// sortCondition: 播放量、弹幕数、分享数排序 默认播放量
// mixDuration: 视频时长筛选 最小时长
// maxDuration: 视频时长筛选 最大时长
// publishTime: 发布时间筛选
func (m *VideoModel) SearchVideos(name, sortCondition string, minDuration, maxDuration, publishTime int64, offset, size int) []*VideoDetailInfo {
	sql :=  "SELECT v.*, s.fabulous_num, s.share_num, s.comment_num, s.browse_num FROM videos as v " +
		"LEFT JOIN video_statistic as s ON v.video_id=s.video_id WHERE v.status=1 "

	if name != "" {
		sql += "AND v.title like '%" + name + "%' "
	}

	if minDuration != 0 && maxDuration != 0 {
		sql += "AND v.video_duration >= ? AND v.video_duration <= ? "
	}

	if publishTime != 0 {
		sql += "AND v.create_at >= ? "
	}

	sql += fmt.Sprintf("GROUP BY v.video_id ORDER BY s.%s DESC, v.is_top DESC, v.is_recommend DESC, v.sortorder DESC, v.video_id DESC LIMIT ?, ?", sortCondition)

	var list []*VideoDetailInfo
	if err := m.Engine.SQL(sql, minDuration, maxDuration, publishTime, offset, size).Find(&list); err != nil {
		log.Log.Errorf("video_trace: search videos err:%s", err)
		return nil
	}

	return list
}

// 获取热门搜索
func (m *VideoModel) GetHotSearch() *models.HotSearch {
	hot := new(models.HotSearch)
	ok, err := m.Engine.Desc("id").Limit(1).Get(hot)
	if !ok || err != nil {
		log.Log.Errorf("video_trace: get hot search err:%s", err)
		return nil
	}

	return hot
}

// 获取审核中/审核失败 的视频列表
func (m *VideoModel) GetVideoReviewList(offset, size int) []*models.Videos {
	var list []*models.Videos
	if err := m.Engine.Where("status=0 OR status=2").Desc("video_id").Limit(size, offset).Find(&list); err != nil {
		return []*models.Videos{}
	}

	return list

}

// 获取用户总浏览数
func (m *VideoModel) GetUserTotalBrowse(userId string) int64 {
	count, err := m.Engine.Where("user_id=?", userId).Count(&models.UserBrowseRecord{})
	if err != nil {
		return 0
	}

	return count
}

// 记录腾讯事件回调信息
func (m *VideoModel) RecordTencentEvent() (int64, error) {
	return m.Engine.InsertOne(m.Events)
}

package mvideo

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/global/rdskey"
	"sports_service/server/models"
	"sports_service/server/util"
	"strings"
	"time"
	"net/url"
)

// todo: 视频id自增 帖子可以使用分布式唯一id
type VideoModel struct {
	Videos       *models.Videos
	Engine       *xorm.Session
	Browse       *models.UserBrowseRecord
	Labels       *models.VideoLabels
	Statistic    *models.VideoStatistic
	Events       *models.TencentCloudEvents
	HotSearch    *models.HotSearch
	Report       *models.VideoReport
	PlayRecord   *models.UserPlayDurationRecord
	Subarea      *models.VideoSubarea
	Album        *models.VideoAlbum
}

// 视频发布请求参数
type VideoPublishParams struct {
	TaskId         int64   `binding:"required" json:"task_id"`        // 任务id
	Cover          string  `json:"cover"`                             // 视频封面
	Title          string  `binding:"required" json:"title"`          // 视频标题
	Describe       string  `json:"describe"`                          // 视频描述
	VideoAddr      string  `binding:"required" json:"video_addr"`     // 视频地址
	VideoDuration  int     `json:"video_duration"`                    // 视频时长
	FileId         string  `binding:"required" json:"file_id"`        // 腾讯云文件id
	VideoWidth     int64   `json:"video_width"`                       // 视频宽
	VideoHeight    int64   `json:"video_height"`                      // 视频高
	VideoLabels    string  `binding:"required" json:"video_labels"`   // 视频标签id（多个用逗号分隔）
	Size           int64   `json:"size"`                              // 视频字节数
	CustomLabels   string  `json:"custom_labels"`                     // 字符串（多个用逗号分隔）
	PubType        int     `json:"pub_type"`                          // 发布类型 1 首页发布 2 社区发布
	SubareaId      string  `json:"subarea_id"`                        // 分区id
	AlbumId        string  `json:"album_id"`                          // 专辑id
}

// 视频信息
type VideosInfoResp struct {
	VideoId       int64  `json:"video_id" example:"1000000000"`       // 视频id
	Title         string `json:"title" example:"视频标题"`             // 标题
	Describe      string `json:"describe" example:"视频描述"`          // 描述
	Cover         string `json:"cover" example:"视频封面"`             // 封面
	VideoAddr     string `json:"video_addr" example:"视频地址"`        // 视频地址
	IsRecommend   int    `json:"is_recommend" example:"1"`           // 是否推荐
	IsTop         int    `json:"is_top" example:"1"`                 // 是否置顶
	VideoDuration int    `json:"video_duration" example:"1000000"`   // 视频时长
	VideoWidth    int64  `json:"video_width" example:"1000"`         // 视频宽
	VideoHeight   int64  `json:"video_height" example:"1000"`        // 视频高
	CreateAt      int    `json:"create_at" example:"1600000000"`     // 视频创建时间
	UserId        string `json:"user_id" example:"用户id"`            // 发布视频的用户id
	Avatar        string `json:"avatar" example:"头像"`                // 头像
	Nickname      string `json:"nick_name" example:"昵称"`             // 昵称
	IsAttention   int    `json:"is_attention" example:"1"`            // 是否关注 1 关注 2 未关注
	OpTime        int    `json:"op_time" example:"1600000000"`        // 用户收藏/点赞/浏览等的操作时间
	TimeElapsed   int    `json:"time_elapsed" example:"1"`           // 已播放的时长 秒
}

// 视频信息
type VideosInfo struct {
	VideoId       int64  `json:"video_id" example:"1000000000"`       // 视频id
	Title         string `json:"title" example:"标题"`                // 标题
	Describe      string `json:"describe" example:"描述"`             // 描述
	Cover         string `json:"cover" example:"封面"`                // 封面
	VideoAddr     string `json:"video_addr" example:"视频地址"`        // 视频地址
	IsRecommend   int    `json:"is_recommend" example:"0"`            // 是否推荐
	IsTop         int    `json:"is_top" example:"0"`                  // 是否置顶
	VideoDuration int    `json:"video_duration" example:"1000000"`    // 视频时长
	VideoWidth    int64  `json:"video_width" example:"1000"`          // 视频宽
	VideoHeight   int64  `json:"video_height" example:"1000"`         // 视频高
	Status        int32  `json:"status" example:"0"`                  // 审核状态
	CreateAt      int    `json:"create_at" example:"1600000000"`      // 视频创建时间
	FabulousNum   int    `json:"fabulous_num" example:"1"`            // 点赞数
	CommentNum    int    `json:"comment_num" example:"1"`             // 评论数
	ShareNum      int    `json:"share_num" example:"1"`               // 分享数
	BrowseNum     int    `json:"browse_num" example:"1"`             // 浏览数（播放数）
	BarrageNum    int    `json:"barrage_num" example:"1"`            // 弹幕数
	TimeElapsed   int    `json:"time_elapsed" example:"1"`           // 已播放的时长 秒
	StatusCn      string `json:"status_cn" example:"审核中"`          // 审核状态（中文展示）
	UserId        string `json:"user_id"`                            // 发布者uid
}

// 首页推荐视频信息
type RecommendVideoInfo struct {
	VideoId       int64                 `json:"video_id"  example:"1000000000"`       // 视频id
	Title         string                `json:"title"  example:"标题"`                 // 标题
	Describe      string                `json:"describe"  example:"描述"`              // 描述
	Cover         string                `json:"cover"  example:"封面"`                 // 封面
	VideoAddr     string                `json:"video_addr"  example:"视频地址"`         // 视频地址
	IsRecommend   int                   `json:"is_recommend" example:"0"`             // 是否推荐
	IsTop         int                   `json:"is_top"  example:"0"`                  // 是否置顶
	VideoDuration int                   `json:"video_duration" example:"100000"`       // 视频时长
	VideoWidth    int64                 `json:"video_width"  example:"100"`            // 视频宽
	VideoHeight   int64                 `json:"video_height"  example:"100"`           // 视频高
	Status        int32                 `json:"status"  example:"1"`                   // 审核状态
	CreateAt      int                   `json:"create_at" example:"1600000000"`        // 视频创建时间
	FabulousNum   int                   `json:"fabulous_num" example:"10"`             // 点赞数
	StatisticsTab string                `json:"statistics_tab"`                        // 统计标签 1个点赞  2分 1个收藏  5分 1个弹幕  10分  1个评论  10分  四项中，哪个分数最高，显示哪个
	BrowseNum     int                   `json:"browse_num" example:"10"`              // 浏览数（播放数）
	UserId        string                `json:"user_id" example:"发布视频的用户id"`      // 发布视频的用户id
	Avatar        string                `json:"avatar" example:"头像"`                 // 头像
	Nickname      string                `json:"nick_name"  example:"昵称"`             // 昵称
	IsCollect     int                   `json:"is_collect" example:"1"`               // 是否收藏
	IsLike        int                   `json:"is_like" example:"1"`                  // 是否点赞
	Size          int                   `json:"size"`                                 // 视频总字节数
}

type VideoDetailInfo struct {
	VideoId       int64                 `json:"video_id"  example:"1000000000"`       // 视频id
	Title         string                `json:"title"  example:"标题"`                 // 标题
	Describe      string                `json:"describe"  example:"描述"`              // 描述
	Cover         string                `json:"cover"  example:"封面"`                 // 封面
	VideoAddr     string                `json:"video_addr"  example:"视频地址"`         // 视频地址
	IsRecommend   int                   `json:"is_recommend" example:"0"`             // 是否推荐
	IsTop         int                   `json:"is_top"  example:"0"`                 // 是否置顶
	VideoDuration int                   `json:"video_duration" example:"100000"`       // 视频时长
	VideoWidth    int64                 `json:"video_width"  example:"100"`            // 视频宽
	VideoHeight   int64                 `json:"video_height"  example:"100"`           // 视频高
	Status        int32                 `json:"status"  example:"1"`                   // 审核状态
	CreateAt      int                   `json:"create_at" example:"1600000000"`        // 视频创建时间
	FabulousNum   int                   `json:"fabulous_num" example:"10"`             // 点赞数
	CommentNum    int                   `json:"comment_num" example:"10"`              // 评论数
	BarrageNum    int                   `json:"barrage_num" example:"10"`              // 弹幕数
	CollectNum    int                   `json:"collect_num" example:"10"`              // 收藏数
	ShareNum      int                   `json:"share_num" example:"10"`               // 分享数
	BrowseNum     int                   `json:"browse_num" example:"10"`              // 浏览数（播放数）
	UserId        string                `json:"user_id" example:"发布视频的用户id"`      // 发布视频的用户id
	Avatar        string                `json:"avatar" example:"头像"`                 // 头像
	Nickname      string                `json:"nick_name"  example:"昵称"`             // 昵称
	IsAttention   int                   `json:"is_attention" example:"1"`             // 是否关注 1 关注 0 未关注
	IsCollect     int                   `json:"is_collect" example:"1"`               // 是否收藏
	IsLike        int                   `json:"is_like" example:"1"`                  // 是否点赞
	FansNum       int64                 `json:"fans_num" example:"100"`               // 粉丝数
	Labels        []*models.VideoLabels `json:"labels"`                               // 视频标签
	PlayInfo      []*PlayInfo           `json:"play_info"`                            // 视频转码后数据
	StatisticsTab string                `json:"statistics_tab"`                       // 统计标签
	Album         int64                 `json:"album"`                                // 专辑
	Subarea       int                   `json:"subarea"`                              // 分区
	AlbumInfo     []*InfoByVideoAlbum   `json:"album_info"`                           // 专辑下的视频数据
}

// 专辑下的视频数据
type InfoByVideoAlbum struct {
	VideoId       int64        `json:"video_id"`
	Title         string       `json:"title"`
	Describe      string       `json:"describe"`
	Cover         string       `json:"cover"`
	VideoAddr     string       `json:"video_addr"`
	VideoDuration int64        `json:"video_duration"`
	VideoWidth    int          `json:"video_width"`
	VideoHeight   int          `json:"video_height"`
	CreateAt      int64        `json:"create_at"`
}

// 记录视频播放时长 请求参数
type PlayDurationParams struct {
	VideoId         int64   `json:"video_id"`         // 视频id
	Duration        int     `json:"duration"`         // 已播时长
	UserId          string  `json:"user_id"`          // 用户id
}

// 删除历史记录请求参数
type DeleteHistoryParam struct {
	ComposeIds        []string     `binding:"required" json:"compose_ids"` // 作品id列表
}

// 删除发布记录请求参数(不支持批量删除)
type DeletePublishParam struct {
	ComposeIds        string       `binding:"required" json:"compose_ids"` // 作品id
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

// 自定义标签请求参数
type CustomLabelParams struct {
	CustomLabel   string    `json:"custom_label"`   // 自定义标签
}

// 添加热搜配置请求参数
type AddHotSearchParams struct {
	HotSearch       string     `json:"hot_search" binding:"required"`   // 热搜内容
	Sortorder       int        `json:"sortorder"`                       // 权重
}

// 删除热搜配置请求参数
type DelHotSearchParams struct {
	Id              int       `json:"id" binding:"required"`     // 数据id
}

// 热搜配置设置权重
type SetSortParams struct {
	Id          int       `json:"id" binding:"required"`          // 数据id
	Sortorder   int       `json:"sortorder"`                      // 权重值
}

// 热搜配置设置状态（展示/隐藏）
type SetStatusParams struct {
	Id          int        `json:"id" binding:"required"`         // 数据id
	Status      int        `json:"status"`                        // 状态 0 展示 1 隐藏
}

// 视频举报
type VideoReportParam struct {
	VideoId    int64      `json:"video_id" binding:"required"`    // 视频id
	UserId     string     `json:"user_id"`
	Reason     string     `json:"reason" binding:"required"`      // 举报理由
}

// 视频转码信息
type PlayInfo struct {
	Type     string   `json:"type" example:"1 流畅（FLU） 2 标清（SD）3 高清（HD）4 全高清（FHD）5 2K 6 4K"`    // 1 流畅（FLU） 2 标清（SD）3 高清（HD）4 全高清（FHD）5 2K 6 4K
	Url      string   `json:"url" example:"对应类型的视频地址"`
	Size     int64    `json:"size" example:"1000000000"`
	Duration int64    `json:"duration" example:"1000000000"`
}

// 分区下的视频信息
type VideoInfoBySubarea struct {
	VideoId       int64                 `json:"video_id"  example:"1000000000"`       // 视频id
	Title         string                `json:"title"  example:"标题"`                 // 标题
	Describe      string                `json:"describe"  example:"描述"`              // 描述
	Cover         string                `json:"cover"  example:"封面"`                 // 封面
	VideoAddr     string                `json:"video_addr"  example:"视频地址"`         // 视频地址
	VideoDuration int                   `json:"video_duration" example:"100000"`       // 视频时长
	VideoWidth    int64                 `json:"video_width"  example:"100"`            // 视频宽
	VideoHeight   int64                 `json:"video_height"  example:"100"`           // 视频高
	CreateAt      int                   `json:"create_at" example:"1600000000"`        // 视频创建时间

	BarrageNum    int                   `json:"barrage_num" example:"10"`              // 弹幕数
	BrowseNum     int                   `json:"browse_num" example:"10"`               // 浏览数（播放数）
	UserId        string                `json:"user_id" example:"发布视频的用户id"`       // 发布视频的用户id
	Avatar        string                `json:"avatar" example:"头像"`                  // 头像
	Nickname      string                `json:"nick_name"  example:"昵称"`              // 昵称
}

// 实栗
func NewVideoModel(engine *xorm.Session) *VideoModel {
	return &VideoModel{
		Browse: new(models.UserBrowseRecord),
		Videos: new(models.Videos),
		Labels: new(models.VideoLabels),
		Statistic: new(models.VideoStatistic),
		Events: new(models.TencentCloudEvents),
		HotSearch: new(models.HotSearch),
		Report: new(models.VideoReport),
		PlayRecord: new(models.UserPlayDurationRecord),
		Subarea: new(models.VideoSubarea),
		Album: new(models.VideoAlbum),
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
		"LEFT JOIN video_statistic as s ON v.`video_id`=s.`video_id` WHERE v.`user_id`=? "
	if status != consts.VIDEO_VIEW_ALL {
		sql += fmt.Sprintf("AND v.`status` = %s ", status)
	} else {
		sql += "AND v.`status` != 3 "
	}

	// 条件为默认时间倒序 则使用videos表的时间字段
	if field == consts.VIDEO_CONDITION_TIME {
		sql += fmt.Sprintf("GROUP BY v.`video_id` ORDER BY v.`%s` DESC, v.`sortorder` DESC ", field)
	} else {
		sql += fmt.Sprintf("GROUP BY v.`video_id` ORDER BY s.`%s` DESC, v.`sortorder` DESC ", field)
	}

	sql += "LIMIT ?, ?"

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

const (
	SEARCH_VIDEOS_BY_LABEL_ID = "SELECT v.* FROM videos AS v LEFT JOIN video_labels as vl ON v.video_id=vl.video_id " +
		" WHERE v.status=1 AND vl.label_id=?  ORDER BY is_top DESC, is_recommend DESC, sortorder DESC, video_id LIMIT ?, ?"
)
// 通过标签id搜索视频
func (m *VideoModel) SearchVideosByLabelId(labelId string, offset, size int) []*models.Videos {
	var list []*models.Videos
	if err := m.Engine.SQL(SEARCH_VIDEOS_BY_LABEL_ID, labelId, offset, size).Find(&list); err != nil {
		log.Log.Errorf("video_trace: search video list by labelId err:%s", err)
		return nil
	}

	return list
}

type BrowseRecord struct {
	ComposeId       int64     `json:"compose_id"`    // 视频id
	UpdateAt        int       `json:"update_at"`     // 浏览时间
}
// 获取浏览过的作品id记录
func (m *VideoModel) GetBrowseRecord(userId string, composeType, offset, size int) []*BrowseRecord {
	var list []*BrowseRecord
	if err := m.Engine.Table(&models.UserBrowseRecord{}).Where("user_id=? AND compose_type=?", userId, composeType).
		Cols("compose_id, update_at").
		Desc("id").
		Limit(size, offset).
		Find(&list); err != nil {
		log.Log.Errorf("video_trace: get browse record err:%s", err)
		return nil
	}

	return list
}

const (
	QUERY_USER_BROWSE_VIDEOS = "SELECT ubr.`update_at` as op_time, v.* FROM user_browse_record AS ubr " +
		"LEFT JOIN videos AS v ON ubr.compose_id=v.video_id AND ubr.compose_type=0 " +
		"WHERE ubr.user_id=? AND v.status=1 ORDER BY ubr.update_at DESC, ubr.id DESC LIMIT ?, ?"
)
// 获取用户浏览过的视频
func (m *VideoModel) GetUserBrowseVideos(userId string, offset, size int) []*VideosInfoResp {
	var list []*VideosInfoResp
	if err := m.Engine.SQL(QUERY_USER_BROWSE_VIDEOS, userId, offset, size).Find(&list); err != nil {
		return nil
	}

	return list
}

// 获取用户浏览过的视频
func (m *VideoModel) GetUserBrowseVideo(userId string, composeType int, composeId int64) *models.UserBrowseRecord {
	m.Browse = new(models.UserBrowseRecord)
	ok, err := m.Engine.Where("user_id=? AND compose_type=? AND compose_id=?", userId, composeType, composeId).Get(m.Browse)
	if !ok || err != nil {
		return nil
	}

	return m.Browse
}

// 记录用户浏览的视频记录
func (m *VideoModel) RecordUserBrowseVideo() error {
	if _, err := m.Engine.InsertOne(m.Browse); err != nil {
		return err
	}

	return nil
}

// 之前有浏览记录 更新浏览时间
func (m *VideoModel) UpdateUserBrowseVideo(userId string,  composeType int, composeId int64) error {
	if _, err := m.Engine.Where("user_id=? AND compose_id=? AND compose_type=?", userId, composeId, composeType).Cols("create_at, update_at").Update(m.Browse); err != nil {
		return err
	}

	return nil
}

// 通过id列表删除浏览的历史记录
func (m *VideoModel) DeleteHistoryByIds(userId string, sql string) error {
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
	total, err := m.Engine.Where("user_id=? AND status=1", userId).Count(&models.Videos{})
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
// 获取视频列表 todo:
func (m *VideoModel) GetVideoList(offset, size int) []*VideoDetailInfo {
	var list []*VideoDetailInfo
	if err := m.Engine.SQL(QUERY_VIDEO_LIST, offset, size).Find(&list); err != nil {
		log.Log.Errorf("video_trace: get video list err:%s", err)
		return nil
	}

	return list
}

const (
	//QUERY_RECOMMEND_VIDEO_LIST = "SELECT v.*, s.fabulous_num,s.browse_num FROM `videos` as v " +
	// "LEFT JOIN video_statistic as s ON v.video_id=s.video_id WHERE v.status = 1 AND v.video_id < ? GROUP BY v.video_id " +
	// "ORDER BY v.is_top DESC, v.is_recommend DESC, v.sortorder DESC, v.video_id DESC LIMIT ?, ?"

	QUERY_RECOMMEND_VIDEO_LIST = "SELECT v.*, s.fabulous_num,s.browse_num FROM `videos` as v " +
		"LEFT JOIN video_statistic as s ON v.video_id=s.video_id WHERE v.status = 1 AND v.video_id < ? GROUP BY v.video_id " +
		"ORDER BY v.video_id DESC LIMIT ?, ?"
)
// 获取推荐的视频列表
func (m *VideoModel) GetRecommendVideoList(index string, offset, size int) []*RecommendVideoInfo {
	var list []*RecommendVideoInfo
	if err := m.Engine.SQL(QUERY_RECOMMEND_VIDEO_LIST, index, offset, size).Find(&list); err != nil {
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
		"ORDER BY v.video_id DESC LIMIT ?, ?"
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

// 获取未浏览的视频数（关注的用户发布的视频）
func (m *VideoModel) GetUnBrowsedAttentionVideos(userIds, browseTm string) int64 {
	type tmp struct {
		Count    int64  `json:"count"`
	}

	unBrowsed := new(tmp)
	sql := fmt.Sprintf("SELECT count(1) as count FROM `videos` WHERE user_id in(%s) AND status=1 AND create_at > %s", userIds, browseTm)
	ok, err := m.Engine.SQL(sql).Get(unBrowsed)
	if !ok || err != nil {
		return 0
	}

	return unBrowsed.Count
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
		sql += fmt.Sprintf("AND v.video_duration >= %d AND v.video_duration <= %d ", minDuration, maxDuration)
	}

	if publishTime != 0 {
		sql += fmt.Sprintf("AND v.create_at >= %d ", publishTime)
	}

	sql += fmt.Sprintf("GROUP BY v.video_id ORDER BY s.%s DESC, v.is_top DESC, v.is_recommend DESC, v.sortorder DESC, v.video_id DESC LIMIT ?, ?", sortCondition)

	var list []*VideoDetailInfo
	if err := m.Engine.SQL(sql, offset, size).Find(&list); err != nil {
		log.Log.Errorf("video_trace: search videos err:%s", err)
		return nil
	}

	return list
}

// 获取热门搜索配置列表
func (m *VideoModel) GetHotSearch() []*models.HotSearch {
	var list []*models.HotSearch
	if err := m.Engine.Desc("sortorder", "id").Find(&list); err != nil {
		log.Log.Errorf("video_trace: get hot search err:%s", err)
		return []*models.HotSearch{}
	}

	return list
}

// 热搜词是否重复
func (m *VideoModel) IsRepeatHotSearchName(name string) bool {
	hot := new(models.HotSearch)
	ok, err := m.Engine.Where("name=?", name).Get(hot)
	if err == nil && ok {
		return true
	}

	return false
}

// 添加热搜配置
func (m *VideoModel) AddHotSearch() error {
	if _, err := m.Engine.InsertOne(m.HotSearch); err != nil {
		return err
	}

	return nil
}

// 删除热搜配置
func (m *VideoModel) DelHotSearch(id int) error {
	if _, err := m.Engine.Where("id=?", id).Delete(&models.HotSearch{}); err != nil {
		return err
	}

	return nil
}

// 修改热搜配置（权重、更新时间）
func (m *VideoModel) UpdateSortByHotSearch(id int) error {
	if _, err := m.Engine.Where("id=?", id).Cols("sortorder", "update_at").Update(m.HotSearch); err != nil {
		return err
	}

	return nil
}

// 更新热搜配置（状态、更新时间）
func (m *VideoModel) UpdateStatusByHotSearch(id int) error {
	if _, err := m.Engine.Where("id=?", id).Cols("status", "update_at").Update(m.HotSearch); err != nil {
		return err
	}

	return nil
}

// 获取审核中/审核失败 的视频列表
func (m *VideoModel) GetVideoReviewList(offset, size int) []*models.Videos {
	var list []*models.Videos
	if err := m.Engine.Where("status=0 OR status=2").Desc("video_id").Limit(size, offset).Find(&list); err != nil {
		return nil
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

// 添加视频举报
func (m *VideoModel) AddVideoReport() (int64, error) {
	return m.Engine.InsertOne(m.Report)
}

// 更新视频转码数据 & AI审核状态
func (m *VideoModel) UpdateVideoPlayInfo(videoId string) error {
	if _, err := m.Engine.Where("video_id=?", videoId).Cols("play_info, status").Update(m.Videos); err != nil {
		return err
	}

	return nil
}

// 记录用户搜索历史
func (m *VideoModel) RecordHistorySearch(userId, name string) error {
	rds := dao.NewRedisDao()
	_, err:= rds.ZINCRBY(rdskey.MakeKey(rdskey.SEARCH_HISTORY_CONTENT, userId), time.Now().Unix(), name)
	return err
}

// 清空用户搜索历史
func (m *VideoModel) CleanHistorySearch(userId string) error {
	rds := dao.NewRedisDao()
	_, err := rds.Del(rdskey.MakeKey(rdskey.SEARCH_HISTORY_CONTENT, userId))
	return err
}

// 获取历史搜索记录
func (m *VideoModel) GetHistorySearch(userId string) []string {
	rds := dao.NewRedisDao()
	list, err := rds.ZREVRANGEString(rdskey.MakeKey(rdskey.SEARCH_HISTORY_CONTENT, userId), 0, 9)
	if err != nil {
		return nil
	}

	return list
}

const (
	QUERY_RECOMMEND_VIDEOS = "SELECT v.video_id, v.cover, v.title, v.`describe`, v.video_addr, v.user_id, v.size, " +
		"v.play_info, v.video_duration, vs.`comment_num`, vs.`browse_num` FROM `videos` as v LEFT JOIN video_statistic as vs " +
		"ON v.video_id=vs.video_id WHERE v.status=1 GROUP BY v.video_id ORDER BY v.is_top desc, v.is_recommend desc, v.create_at desc, " +
		"v.video_id desc LIMIT ?, ?"
)
// 获取相关视频列表（暂时随机2个）
func (m *VideoModel) GetRecommendVideos(offset, limit int32) []*VideoDetailInfo {
	var list []*VideoDetailInfo
	if err := dao.AppEngine.Sql(QUERY_RECOMMEND_VIDEOS, offset, limit).Find(&list); err != nil {
		return nil
	}

	return list
}


const (
	LINK_KEY  = "MgbbCBVMhIwJCi9gnhhj"
	// todo: 测试改为 3分钟
	EXPIRE_TM = 60 * 60 * 3
)
// 视频防盗链(有效时长3个小时)
func (m *VideoModel) AntiStealingLink(videoUrl string) string {
	if videoUrl == "" {
		return ""
	}

	urls, err := url.Parse(videoUrl)
	if err != nil {
		log.Log.Errorf("video_trace: url.Parse fail, err:%s", err)
		return ""
	}

	path := urls.Path
	path = path[:strings.LastIndex(path, "/") + 1]
	rand := util.GenSecret(util.MIX_MODE, 10)
	tm := fmt.Sprintf("%x", time.Now().Unix() + EXPIRE_TM)
	signStr := fmt.Sprintf("%s%s%s%s", LINK_KEY, path, tm, rand)
	signStr = strings.Trim(signStr, " ")
	signStr = strings.Replace(signStr, "\n", "", -1)
	sign := util.Md5String(signStr)
	return fmt.Sprintf("%s?t=%s&us=%s&sign=%s", videoUrl, tm, rand, sign)
}

// 记录任务id -> 用户id
func (m *VideoModel) RecordUploadTaskId(userId string, taskId int64) error {
	rds := dao.NewRedisDao()
	return rds.SETEX(rdskey.MakeKey(rdskey.VIDEO_UPLOAD_TASK, taskId), rdskey.KEY_EXPIRE_DAY * 3, userId)
}

// 通过任务id 获取 用户id
func (m *VideoModel) GetUploadUserIdByTaskId(taskId int64) (string, error) {
	rds := dao.NewRedisDao()
	return rds.Get(rdskey.MakeKey(rdskey.VIDEO_UPLOAD_TASK, taskId))
}

// 记录用户发布的视频信息
func (m *VideoModel) RecordPublishInfo(userId, pubInfo string, taskId int64) error {
	key := rdskey.MakeKey(rdskey.VIDEO_UPLOAD_INFO, userId, taskId)
	rds := dao.NewRedisDao()
	return rds.SETEX(key, rdskey.KEY_EXPIRE_DAY * 3,  pubInfo)
}

// 获取用户发布的视频信息
func (m *VideoModel) GetPublishInfo(userId string, taskId int64) (string, error) {
	key := rdskey.MakeKey(rdskey.VIDEO_UPLOAD_INFO, userId, taskId)
	log.Log.Debugf("publishInfo key:%s", key)
	rds := dao.NewRedisDao()
	return rds.Get(key)
}

// 接收回调成功 删除 存储发布视频信息的key
func (m *VideoModel) DelPublishInfo(userId string, taskId int64) (int, error) {
	key := rdskey.MakeKey(rdskey.VIDEO_UPLOAD_INFO, userId, taskId)
	rds := dao.NewRedisDao()
	return rds.Del(key)
}

// 通过腾讯云返回的文件id查询视频
func (m *VideoModel) GetVideoByFileId(fileId string) *models.Videos {
	m.Videos = new(models.Videos)
	ok, err := m.Engine.Where("file_id=?", fileId).Get(m.Videos)
	if !ok || err != nil {
		return nil
	}

	return m.Videos
}

// 获取用户视频播放时长记录
func (m *VideoModel) GetUserPlayDurationRecord(userId, videoId string) *models.UserPlayDurationRecord {
	m.PlayRecord = new(models.UserPlayDurationRecord)
	ok, err := m.Engine.Where("user_id=? AND video_id=?", userId, videoId).Get(m.PlayRecord)
	if !ok || err != nil {
		return nil
	}

	return m.PlayRecord
}

// 添加用户播放时长记录
func (m *VideoModel) AddUserPlayDurationRecord() error {
	if _, err := m.Engine.InsertOne(m.PlayRecord); err != nil {
		return err
	}

	return nil
}

// 更新用户已播视频时长记录
func (m *VideoModel) UpdateUserPlayDurationRecord() error {
	if _, err := m.Engine.ID(m.PlayRecord.Id).Cols("play_duration, update_at").Update(m.PlayRecord); err != nil {
		return err
	}

	return nil
}

// 更新视频信息
func (m *VideoModel) UpdateVideoInfo() (int64, error) {
	return m.Engine.ID(m.Videos.VideoId).Update(m.Videos)
}

const (
	QUERY_VIDEO_LIST_BY_SUBAREA = "SELECT v.*, s.barrage_num,s.browse_num FROM `videos` as v LEFT JOIN video_statistic " +
		"as s ON v.video_id=s.video_id WHERE v.status = 1 AND v.subarea=? ORDER BY v.video_id DESC, v.is_top DESC, " +
		"v.is_recommend DESC LIMIT ?, ?"
)
// 获取分区下的视频列表
func (m *VideoModel) GetVideoListBySubarea(subarea string, offset, size int) ([]*VideoInfoBySubarea, error){
	var list []*VideoInfoBySubarea
	if err := m.Engine.SQL(QUERY_VIDEO_LIST_BY_SUBAREA, subarea, offset, size).Find(&list); err != nil {
		log.Log.Errorf("video_trace: get video by subarea err:%s", err)
		return nil, err
	}

	return list, nil
}

// 视频专辑信息
type VideoAlbumInfo struct {
	Id        int64  `json:"id"`          // 专辑id
	UserId    string `json:"user_id"`     // 用户id
	AlbumName string `json:"album_name"`  // 专辑名称
	VideoNum  int    `json:"video_num"`   // 专辑下的视频总数
}

const (
	QUERY_ALBUM_LIST_BY_USER = "SELECT va.id, va.user_id, va.album_name, video_num FROM video_album AS va " +
		"LEFT JOIN (SELECT count(1) as video_num, album FROM videos as v " +
		"WHERE v.status=1 AND v.user_id=? GROUP BY v.`album`) AS v ON va.id = v.album " +
		"WHERE va.user_id=? AND va.status=0 ORDER BY va.create_at DESC LIMIT ?, ?"
)
// 获取用户发布的专辑列表 及 专辑下的视频数量
func (m *VideoModel) GetVideoAlbumListByUser(userId string, offset, size int) ([]*VideoAlbumInfo, error) {
	var list []*VideoAlbumInfo
	if err := m.Engine.SQL(QUERY_ALBUM_LIST_BY_USER, userId, userId, offset, size).Find(&list); err != nil {
		log.Log.Errorf("video_trace: get video album list by user err:%s", err)
		return nil, err
	}

	return list, nil
}

// 获取某专辑下的视频列表
func (m *VideoModel) GetVideoListByAlbum(userId string, album int64) ([]*InfoByVideoAlbum, error) {
	var list []*InfoByVideoAlbum
	if err := m.Engine.Table(&models.Videos{}).Where("status=1 AND user_id=? AND album=?", userId, album).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

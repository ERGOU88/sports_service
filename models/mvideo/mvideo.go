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
}

// 视频发布请求参数
type VideoPublishParams struct {
	Cover          string  `binding:"required" json:"cover"`          // 视频封面
	Title          string  `binding:"required" json:"title"`          // 视频标题
	Describe       string  `binding:"required" json:"describe"`       // 视频描述
	VideoAddr      string  `binding:"required" json:"videoAddr"`     // 视频地址
	VideoDuration  int     `binding:"required" json:"videoDuration"` // 视频时长
	VideoLabels    string  `binding:"required" json:"videoLabels"`   // 视频标签id（多个用逗号分隔）
}

// 视频信息
type VideosInfoResp struct {
	VideoId       int64  `json:"videoId"`       // 视频id
	Title         string `json:"title"`         // 标题
	Describe      string `json:"describe"`      // 描述
	Cover         string `json:"cover"`         // 封面
	VideoAddr     string `json:"videoAddr"`     // 视频地址
	IsRecommend   int    `json:"isRecommend"`   // 是否推荐
	IsTop         int    `json:"isTop"`         // 是否置顶
	VideoDuration int    `json:"videoDuration"` // 视频时长
	VideoWidth    int64  `json:"videoWidth"`    // 视频宽
	VideoHeight   int64  `json:"videoHeight"`   // 视频高
	CreateAt      int    `json:"createAt"`      // 视频创建时间
	UserId        string `json:"userId"`        // 发布视频的用户id
	Avatar        string `json:"avatar"`        // 头像
	Nickname      string `json:"nickName"`      // 昵称
	IsAttention   int    `json:"isAttention"`   // 是否关注 1 关注 2 未关注
	OpTime        int    `json:"collectAt"`     // 用户收藏/点赞等的操作时间
}

// 用户发布的视频信息
type PublishVideosInfo struct {
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
}

// 删除历史记录请求参数
type DeleteHistoryParam struct {
	ComposeIds        []string     `binding:"required" json:"composeIds"` // 作品id列表
}

// 删除发布记录请求参数(不支持批量删除)
type DeletePublishParam struct {
	ComposeIds        string       `binding:"required" json:"composeIds"` // 作品id
}

// 实栗
func NewVideoModel(engine *xorm.Session) *VideoModel {
	return &VideoModel{
		Browse: new(models.UserBrowseRecord),
		Videos: new(models.Videos),
		Labels: new(models.VideoLabels),
		Statistic: new(models.VideoStatistic),
		Engine: engine,
	}
}

// 视频发布
func (m *VideoModel) VideoPublish() error {
	if _, err := m.Engine.InsertOne(m.Videos); err != nil {
		return err
	}

	return nil
}

// 添加视频标签(一次插入多条)
func (m *VideoModel) AddVideoLabels(labels []*models.VideoLabels) (int64, error) {
	return m.Engine.Insert(labels)
}

// 删除视频标签
func (m *VideoModel) DelVideoLabels(videoId string) error {
	if _, err := m.Engine.Where("video_id=?", videoId).Delete(m.Labels); err != nil {
		return err
	}

	return nil
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
	ok, err := m.Engine.Where("video_id=?", videoId).Get(m.Statistic)
	if !ok || err != nil {
		log.Log.Errorf("video_trace: get video statistic info err:%s", err)
		return nil
	}

	return m.Statistic
}

// 删除视频统计数据
func (m *VideoModel) DelVideoStatistic(videoId string) error {
	if _, err := m.Engine.Where("video_id=?", videoId).Delete(m.Statistic); err != nil {
		return err
	}

	return nil
}

const (
	UPDATE_VIDEO_LIKE_NUM  = "UPDATE `video_statistic` SET `like_num` = `like_num` + ?, `update_at`=? WHERE `video_id`=? AND `like_num` + ? > 0 LIMIT 1"
)
// 更新视频点赞数
func (m *VideoModel) UpdateVideoLikeNum(now, num int) error {
	if _, err := m.Engine.Exec(UPDATE_VIDEO_LIKE_NUM, num, now, num); err != nil {
		return err
	}

	return nil
}

const (
	UPDATE_VIDEO_COLLECT_NUM  = "UPDATE `video_statistic` SET `collect_num` = `collect_num` + ?, `update_at`=? WHERE `video_id`=? AND `collect_num` + ? > 0 LIMIT 1"
)
// 更新视频收藏数
func (m *VideoModel) UpdateVideoCollectNum(now, num int) error {
	if _, err := m.Engine.Exec(UPDATE_VIDEO_COLLECT_NUM, num, now, num); err != nil {
		return err
	}

	return nil
}

// 更新视频评论数
func (m *VideoModel) UpdateVideoCommentNum() {
	return
}

const (
	UPDATE_VIDEO_BROWSE_NUM  = "UPDATE `video_statistic` SET `browse_num` = `browse_num` + ?, `update_at`=? WHERE `video_id`=? AND `browse_num` + ? > 0 LIMIT 1"
)
// 更新视频浏览数
func (m *VideoModel) UpdateVideoBrowseNum(now, num int) error {
	if _, err := m.Engine.Exec(UPDATE_VIDEO_BROWSE_NUM, num, now, num); err != nil {
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

// 更新视频弹幕数
func (m *VideoModel) UpdateVideoBarrageNum() {
	return
}

// 分页获取 用户发布的视频列表[通过审核状态和条件查询]
func (m *VideoModel) GetUserPublishVideos(offset, size int, userId, status, field string) []*PublishVideosInfo {
	var list []*PublishVideosInfo

	sql := "SELECT v.*, s.fabulous_num, s.share_num, s.comment_num, s.browse_num FROM videos as v " +
		"LEFT JOIN video_statistic as s ON v.video_id=s.video_id WHERE v.user_id=? "
	if status != consts.VIDEO_VIEW_ALL {
		sql += "AND v.`status` = ? "
	} else {
		sql += "AND v.`status` != 3"
	}

	// 条件为默认时间倒序 则使用videos表的时间字段
	if field == consts.VIDEO_CONDITION_TIME {
		sql += fmt.Sprintf("GROUP BY v.video_id ORDER BY v.%s DESC, v.sortorder DESC LIMIT ?, ?", field)
	} else {
		sql += fmt.Sprintf("GROUP BY v.video_id ORDER BY s.%s DESC, v.sortorder DESC LIMIT ?, ?", field)
	}

	if err := m.Engine.SQL(sql, userId, offset, size).Find(&list); err != nil {
		log.Log.Errorf("video_trace: get user publish videos err:%s", err)
		return nil
	}

	return list
}

// 通过id查询视频
func (m *VideoModel) FindVideoById(videoId string) *models.Videos {
	ok, err := m.Engine.Where("video_id=?", videoId).Get(m.Videos)
	if !ok || err != nil {
		return nil
	}

	return m.Videos
}

// 通过视频id查询视频列表
func (m *VideoModel) FindVideoListByIds(videoIds string, offset, size int) []*models.Videos {
	var list []*models.Videos
	sql := fmt.Sprintf("SELECT * FROM videos WHERE video_id in(%s) AND status=0 ORDER BY is_top DESC, is_recommend DESC, sortorder DESC, id LIMIT ?, ?", videoIds)
	if err := m.Engine.SQL(sql, offset, size).Find(&list); err != nil {
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
	if err := m.Engine.Where("user_id=? AND compose_type=?", userId, composeType).
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
	if _, err :=m.Engine.Where("user_id=? AND video_id=?", userId, videoId).Cols("status").Update(m.Videos); err != nil {
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


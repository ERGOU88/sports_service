package mvideo

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/global/app/log"
	"sports_service/server/models"
	"fmt"
)

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
	VideoAddr      string  `binding:"required" json:"video_addr"`     // 视频地址
	VideoDuration  int     `binding:"required" json:"video_duration"` // 视频时长
	VideoLabels    string  `binding:"required" json:"video_labels"`   // 视频标签id（多个用逗号分隔）
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
	Nickname      string `json:"nickName"`       // 昵称
	IsAttention   int    `json:"is_attention"`   // 是否关注 1 关注 2 未关注
	OpTime        int    `json:"collect_at"`     // 用户收藏/点赞等的操作时间
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

// 分页获取 用户发布的视频列表
func (m *VideoModel) GetUserPublishVideos(offset, size int) {
	return
}

// 通过id查询视频
func (m *VideoModel) FindVideoById(videoId string) *models.Videos {
	ok, err := m.Engine.Where("video_id=?").Get(m.Videos)
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



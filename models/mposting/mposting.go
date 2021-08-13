package mposting

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mshare"
	"fmt"
)

// 帖子模块
type PostingModel struct {
	Engine            *xorm.Session
	Posting           *models.PostingInfo
	PostingTopic      *models.PostingTopic
	Browse            *models.UserBrowseRecord
	Statistic         *models.PostingStatistic
	ReceiveAt         *models.ReceivedAt
	ApplyCream        *models.PostingApplyCream
	Report            *models.PostingReport
}

// 删除发布的帖子记录请求参数(不支持批量删除)
type DeletePostParam struct {
	PostId        int64       `binding:"required" json:"post_id"` // 帖子id
}

// 申请精华帖
type ApplyCreamParam struct {
	PostId        int64      `binding:"required" json:"post_id"`  // 帖子id
	Reason        string     `json:"reason"`                      // 申请理由
}

// 帖子举报
type PostReportParam struct {
	PostId     int64      `json:"post_id" binding:"required"`     // 帖子id
	UserId     string     `json:"user_id"`
	Reason     string     `json:"reason" binding:"required"`      // 举报理由
}

// 发布帖子请求参数
type PostPublishParam struct {
	Title             string              `json:"title"`                              // 标题
	Describe          string              `json:"describe"`                           // 富文本内容
	//ForwardVideo      *ForwardVideoInfo   `json:"forward_video"`                     // 转发的视频内容 todo: 结构体
	//VideoId           string              `json:"video_id"`                          // 关联的视频id
	//PostingType    int        `json:"posting_type" binding:"required"`               // 帖子类型  0 纯文本 1 图文 2 视频 + 文
	ImagesAddr        []string            `json:"images_addr"`                         // 图片地址
	SectionId         int                 `json:"section_id"`                          // 主模块id
	TopicIds          []string            `binding:"required" json:"topic_ids"`        // 话题id （多个）
	AtInfo            []string            `json:"at_info"`                             // @信息 [需@的用户uid]
	//ContentType       int                 `json:"content_type"`                       // 0 发布 1 转发视频 2 转发帖子
	//ForwardPost       *ForwardPostInfo    `json:"forward_post"`                       // 转发的帖子内容
}

// 转发的视频信息
type ForwardVideoInfo struct {
	VideoId       string                `json:"video_id"`
	Title         string                `json:"title"  example:"标题"`                 // 标题
	Describe      string                `json:"describe"  example:"描述"`              // 描述
	Cover         string                `json:"cover"  example:"封面"`                 // 封面
	VideoAddr     string                `json:"video_addr"  example:"视频地址"`         // 视频地址
	VideoDuration int                   `json:"video_duration" example:"100000"`       // 视频时长
	CreateAt      int                   `json:"create_at" example:"1600000000"`        // 视频创建时间
	BrowseNum     int                   `json:"browse_num" example:"10"`               // 浏览数（播放数）
	UserId        string                `json:"user_id" example:"发布视频的用户id"`      // 发布视频的用户id
	Avatar        string                `json:"avatar" example:"头像"`                 // 头像
	Nickname      string                `json:"nick_name"  example:"昵称"`             // 昵称
	Size          int                   `json:"size"`                                 // 视频总字节数
	BarrageNum    int                   `json:"barrage_num" example:"1"`              // 弹幕数
	LabelNames    string                `json:"label_names"`                          // 标签名称 多个用逗号分隔
}

// 转发的帖子信息 todo: 产品当前逻辑 转发的帖子皆为文本
type ForwardPostInfo struct {
	PostId        string                `json:"post_id"`                // 转发的帖子id
	PostingType   int                   `json:"posting_type"`           // 帖子类型  0 纯文本 1 图文 2 视频 + 文
	//ContentType   int                   `json:"content_type"`           // 0 社区发布 1 转发视频 2 转发帖子
	TopicNames    string                `json:"topic_names"`            // 话题名词 多个用逗号分隔
	Title         string                `json:"title"`                  // 标题
	Describe      string                `json:"describe"`               // 描述
	Content       string                `json:"content"`                // 暂时使用不到
}


// 帖子详情数据
type PostDetailInfo struct {
	Id            int64                  `json:"post_id"  example:"1000000000"`          // 帖子id
	//Cover         string                 `json:"cover"  example:"封面"`                 // 封面
	//VideoAddr     string                 `json:"video_addr"  example:"视频地址"`         // 视频地址
	//VideoDuration int                    `json:"video_duration" example:"100000"`       // 视频时长
	//Size          int                    `json:"size"`                                 // 视频大小
	Title         string                 `json:"title"  example:"标题"`                 // 标题
	Describe      string                 `json:"describe"  example:"描述"`              // 描述
	Content       string                 `json:"content,omitempty"`                    // 帖子内容 图片列表/json 例如转发的视频
	IsRecommend   int                    `json:"is_recommend" example:"0"`             // 是否推荐 1是
	IsTop         int                    `json:"is_top"  example:"0"`                   // 是否置顶 1是
	IsCream       int                    `json:"is_cream"`                              // 是否精华 1是
	Status        int32                  `json:"status"  example:"1"`                   // 审核状态 （0：审核中，1：审核通过 2：审核不通过 3：逻辑删除）
	CreateAt      int                    `json:"create_at" example:"1600000000"`        // 创建时间
	FabulousNum   int                    `json:"fabulous_num" example:"10"`             // 点赞数
	CommentNum    int                    `json:"comment_num" example:"10"`              // 评论数
	ShareNum      int                    `json:"share_num" example:"10"`                // 分享数
	BrowseNum     int                    `json:"browse_num" example:"10"`               // 浏览数
	UserId        string                 `json:"user_id" example:"发布视频的用户id"`      // 发布者用户id
	Avatar        string                 `json:"avatar" example:"头像"`                 // 头像
	Nickname      string                 `json:"nick_name"  example:"昵称"`             // 昵称
	IsAttention   int                    `json:"is_attention" example:"1"`             // 是否关注 1 关注 0 未关注
	IsLike        int                    `json:"is_like" example:"1"`                  // 是否点赞
	Topics        []*models.PostingTopic `json:"topics,omitempty"`                     // 所属话题
	ForwardVideo  *mshare.ShareVideoInfo `json:"forward_video,omitempty"`              // 转发的视频内容 todo: 结构体
	ForwardPost   *mshare.SharePostInfo  `json:"forward_post,omitempty"`               // 转发的帖子内容
	ImagesAddr    []string               `json:"images_addr,omitempty"`                // 图片地址
	ContentType   int                    `json:"content_type"`                         // 0 社区发布 1 转发视频 2 转发帖子
	PostingType   int                    `json:"posting_type"`                         // 帖子类型  0 纯文本 1 图文 2 视频 + 文字
	HeatNum       int                    `json:"heat_num"`                             // 热度
	VideoId       int64                  `json:"video_id"`                             // 关联的视频id
	RelatedVideo  *RelatedVideo          `json:"related_video,omitempty"`              // 帖子关联的视频信息
	StatusCn      string                 `json:"status_cn"`                            // 中文状态
}

type RelatedVideo struct {
	VideoId       int64                 `json:"video_id"`
	Title         string                `json:"title"  example:"标题"`                 // 标题
	Describe      string                `json:"describe"  example:"描述"`              // 描述
	Cover         string                `json:"cover"  example:"封面"`                 // 封面
	VideoAddr     string                `json:"video_addr"  example:"视频地址"`         // 视频地址
	VideoDuration int                   `json:"video_duration" example:"100000"`       // 视频时长
	CreateAt      int                   `json:"create_at" example:"1600000000"`        // 视频创建时间
	UserId        string                `json:"user_id" example:"发布视频的用户id"`      // 发布视频的用户id
	Avatar        string                `json:"avatar" example:"头像"`                 // 头像
	Nickname      string                `json:"nick_name"  example:"昵称"`             // 昵称
	Size          int64                 `json:"size"`                                 // 视频总字节数
	IsLike        int                   `json:"is_like"`                              // 是否点赞

	FabulousNum   int                   `json:"fabulous_num" example:"10"`             // 点赞数
	ShareNum      int                   `json:"share_num" example:"10"`                // 分享/转发数
	CommentNum    int                   `json:"comment_num" example:"10"`              // 弹幕数数

	Subarea       *models.VideoSubarea  `json:"subarea"`                               // 视频分区
}

// 实栗
func NewPostingModel(engine *xorm.Session) *PostingModel {
	return &PostingModel{
		Engine: engine,
		Posting: new(models.PostingInfo),
		PostingTopic: new(models.PostingTopic),
		Statistic: new(models.PostingStatistic),
		ReceiveAt: new(models.ReceivedAt),
		ApplyCream: new(models.PostingApplyCream),
		Report: new(models.PostingReport),
	}
}

// 添加多个@ (一次插入多条)
func (m *PostingModel) AddReceiveAtList(at []*models.ReceivedAt) (int64, error) {
	return m.Engine.InsertMulti(at)
}

const (
	UPDATE_RECEIVE_AT_STATUS = "UPDATE `received_at` SET status=1,update_at=? WHERE compose_id=? AND topic_type=?"
)
// 帖子通过时 需修改@数据的状态 使@生效
func (m *PostingModel) UpdateReceiveAtStatus(composeId string, topicType, tm int) error {
	if _, err := m.Engine.Exec(UPDATE_RECEIVE_AT_STATUS, tm, composeId, topicType); err != nil {
		return err
	}

	return nil
}

// 通过帖子id获取
func (m *PostingModel) GetReceivedInfoByComposeId() {

}

const (
	UPDATE_POST_TOPIC_STATUS = "UPDATE `posting_topic` SET status=1, update_at=? WHERE posting_id=?"
)
// 帖子通过时 修改帖子所属话题 数据状态
func (m *PostingModel) UpdatePostTopicStatus(postId string, tm int) error {
	if _, err := m.Engine.Exec(UPDATE_POST_TOPIC_STATUS, tm, postId); err != nil {
		return err
	}

	return nil
}

// 通过id获取帖子
func (m *PostingModel) GetPostById(id string) (*models.PostingInfo, error) {
	m.Posting = new(models.PostingInfo)
	ok, err := m.Engine.Where("id=?", id).Get(m.Posting)
	if !ok || err != nil {
		return nil, err
	}

	return m.Posting, nil
}

// 获取帖子所属话题 [1对多]
func (m *PostingModel) GetPostTopic(postId string) ([]*models.PostingTopic, error) {
	var list []*models.PostingTopic
	if err := m.Engine.Where("posting_id=?", postId).Asc("posting_id").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 添加帖子
func (m *PostingModel) AddPost() (int64, error) {
	return m.Engine.InsertOne(m.Posting)
}

// 删除帖子
func (m *PostingModel) DelPost(postId int64) error {
	if _, err := m.Engine.Where("id=?", postId).Delete(&models.PostingInfo{}); err != nil {
		return err
	}

	return nil
}

// 添加帖子所属话题(一次插入多条)
func (m *PostingModel) AddPostingTopics(topics []*models.PostingTopic) (int64, error) {
	return m.Engine.InsertMulti(topics)
}

const (
	UPDATE_POST_BROWSE_NUM  = "UPDATE `posting_statistic` SET `browse_num` = `browse_num` + ?, " +
		"`heat_num` = `heat_num` + ?, `update_at`=? WHERE `posting_id`=? AND `browse_num` + ? >= 0 LIMIT 1"
)
// 更新帖子浏览数 及 帖子热度
func (m *PostingModel) UpdatePostBrowseNum(postId int64, now, num int) error {
	if _, err := m.Engine.Exec(UPDATE_POST_BROWSE_NUM, num, num, now, postId, num); err != nil {
		return err
	}

	return nil
}

const (
	UPDATE_POST_LIKE_NUM  = "UPDATE `posting_statistic` SET `fabulous_num` = `fabulous_num` + ?, " +
		"`heat_num` = `heat_num` + ?, `update_at`=? WHERE `posting_id`=? AND `fabulous_num` + ? >= 0 LIMIT 1"
)
// 更新帖子点赞数 及 帖子热度
func (m *PostingModel) UpdatePostLikeNum(videoId int64, now, num int) error {
	if _, err := m.Engine.Exec(UPDATE_POST_LIKE_NUM, num, num, now, videoId, num); err != nil {
		return err
	}

	return nil
}

const (
	UPDATE_POST_COMMENT_NUM = "UPDATE `posting_statistic` SET `comment_num` = `comment_num` + ?, " +
		"`heat_num` = `heat_num` + ?, `update_at`=? WHERE `posting_id`=? AND `comment_num` + ? >= 0 LIMIT 1"
)
// 更新帖子评论数 及 帖子热度
func (m *PostingModel) UpdatePostCommentNum(postId int64, now, num int) error {
	if _, err := m.Engine.Exec(UPDATE_POST_COMMENT_NUM, num, num, now, postId, num); err != nil {
		return err
	}

	return nil
}


// 获取用户浏览过的帖子
func (m *PostingModel) GetUserBrowsePost(userId string, composeType int, composeId int64) *models.UserBrowseRecord {
	m.Browse = new(models.UserBrowseRecord)
	ok, err := m.Engine.Where("user_id=? AND compose_type=? AND compose_id=?", userId, composeType, composeId).Get(m.Browse)
	if !ok || err != nil {
		return nil
	}

	return m.Browse
}

// 记录用户浏览的帖子记录
func (m *PostingModel) RecordUserBrowsePost() error {
	if _, err := m.Engine.InsertOne(m.Browse); err != nil {
		return err
	}

	return nil
}

// 之前有浏览记录 更新浏览时间
func (m *PostingModel) UpdateUserBrowsePost(userId string,  composeType int, composeId int64) error {
	if _, err := m.Engine.Where("user_id=? AND compose_id=? AND compose_type=?", userId, composeId, composeType).Cols("create_at, update_at").Update(m.Browse); err != nil {
		return err
	}

	return nil
}

// 添加帖子统计数据
func (m *PostingModel) AddPostStatistic() error {
	if _, err := m.Engine.InsertOne(m.Statistic); err != nil {
		return err
	}

	return nil
}

// 获取帖子统计数据
func (m *PostingModel) GetPostStatistic(postId string) (*models.PostingStatistic, error) {
	m.Statistic = new(models.PostingStatistic)
	ok, err := m.Engine.Where("posting_id=?", postId).Get(m.Statistic)
	if !ok || err != nil {
		return nil, err
	}

	return m.Statistic, nil
}

// 获取某板块下的帖子总数
func (m *PostingModel) GetPostNumBySection(sectionId string) (int64, error) {
	return m.Engine.Where("status=1 AND section_id=?", sectionId).Count(&models.PostingInfo{})
}

type TopPost struct {
	Id       int64    `json:"id"`             // 帖子id
	Title    string   `json:"title"`          // 标题
	Describe string   `json:"describe"`       // 描述
	IsTop    int      `json:"is_top"`         // 是否置顶 1 置顶
	CreateAt int64    `json:"create_at"`      // 发布时间
}

// 获取板块下的置顶的帖子
func (m *PostingModel) GetTopPostBySectionId(offset, size int, sectionId string) ([]*TopPost, error) {
	var list []*TopPost
	if err := m.Engine.Table(&models.PostingInfo{}).Where("section_id=? AND is_top=1 AND status=1", sectionId).Limit(size, offset).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 获取某话题下的帖子总数
func (m *PostingModel) GetPostNumByTopic(topicId string) (int64, error) {
	return m.Engine.Where("status=1 AND topic_id=?", topicId).Count(&models.PostingTopic{})
}

const (
	GET_POST_LIST_BY_SECTION = "SELECT p.*, ps.fabulous_num, ps.browse_num, ps.share_num, ps.comment_num, ps.heat_num FROM " +
		"`posting_info` AS p LEFT JOIN `posting_statistic` as ps ON p.id=ps.posting_id WHERE p.status=1 AND p.section_id=? AND p.is_top=0 " +
		" ORDER BY p.is_cream DESC, p.id DESC LIMIT ?, ?"
)
// 通过板块id 获取帖子列表
func (m *PostingModel) GetPostListBySectionId(sectionId string, offset, size int) ([]*PostDetailInfo, error) {
	var list []*PostDetailInfo
	if err := m.Engine.SQL(GET_POST_LIST_BY_SECTION, sectionId, offset, size).Find(&list); err != nil {
		return nil, err
	}


	return list, nil
}

// 通过话题id 获取同话题的帖子列表 [sortHot为1 按热度排序 默认按发布时间排序]
func (m *PostingModel) GetPostListByTopicId(topicId, sortHot string, offset, size int) ([]*PostDetailInfo, error) {
	var list []*PostDetailInfo
	sql := "SELECT p.*,ps.* FROM `posting_info` AS p LEFT JOIN `posting_topic` as pt ON p.id=pt.posting_id " +
		"LEFT JOIN `posting_statistic` as ps ON p.id=ps.posting_id WHERE p.status=1 AND pt.topic_id=? AND p.is_top=0 ORDER BY "

	if sortHot == consts.POST_SORT_HOT {
		sql += "ps.`heat_num` DESC, "
	}

	sql += "sortorder DESC, id DESC LIMIT ?, ?"

	if err := m.Engine.SQL(sql, topicId, offset, size).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 搜索帖子 标题 / 内容
func (m *PostingModel) SearchPost(name string, offset, size int) ([]*PostDetailInfo, error) {
	sql := "SELECT * FROM posting_info WHERE status=1 AND video_id=0 AND title like '%" + name + "%' OR status=1 AND video_id=0 AND content LIKE '%" + name + "%' ORDER BY `id` DESC LIMIT ?, ?"
	var list []*PostDetailInfo
	if err := m.Engine.SQL(sql, offset, size).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 搜索帖子 标题 / 内容 按热度排序
func (m *PostingModel) SearchPostOrderByHeat(name string, offset, size int) ([]*PostDetailInfo, error) {
	sql := "SELECT p.*, ps.fabulous_num, ps.browse_num, ps.share_num, ps.comment_num, ps.heat_num FROM " +
	"`posting_info` AS p LEFT JOIN `posting_statistic` as ps ON p.id=ps.posting_id WHERE p.status=1 AND p.video_id=0 "

	if name != "" {
		sql += "AND p.title LIKE '%" + name + "%' OR p.status=1 AND p.video_id=0 AND p.content LIKE '%" + name + "%' "
	}


	sql += fmt.Sprintf("ORDER BY ps.`heat_num` DESC, p.is_top DESC, p.is_cream DESC, p.id DESC LIMIT ?, ?")

	var list []*PostDetailInfo
	if err := m.Engine.SQL(sql, offset, size).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

const (
	GET_PUBLISH_POST_BY_USER = "SELECT p.*, ps.* FROM `posting_info` AS p LEFT JOIN `posting_statistic` as ps " +
		"ON p.id=ps.posting_id WHERE p.user_id=? AND p.video_id=0 ORDER BY p.id DESC LIMIT ?, ?"
)
// 获取用户发布的帖子列表 [不包含视频]
func (m *PostingModel) GetPublishPostByUser(userId, status string, offset, size int) ([]*PostDetailInfo, error) {
	sql := "SELECT p.*, ps.* FROM `posting_info` AS p LEFT JOIN `posting_statistic` as ps " +
	"ON p.id=ps.posting_id WHERE p.user_id=? AND p.video_id=0 "

	if status == consts.POST_AUDIT_SUCCESS {
		sql += "AND status=1 "
	}

	sql += "ORDER BY ps.`heat_num` DESC, p.is_top DESC, p.is_cream DESC, p.id DESC LIMIT ?, ?"

	var list []*PostDetailInfo
	if err := m.Engine.SQL(GET_PUBLISH_POST_BY_USER, userId, offset, size).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

const (
	DELETE_PUBLISH_POST = "DELETE FROM `posting_info` WHERE `id`=?"
)
// 删除发布的帖子记录
func (m *PostingModel) DelPublishPostById(postId string) error {
	if _, err := m.Engine.Exec(DELETE_PUBLISH_POST, postId); err != nil {
		return err
	}

	return nil
}


// 删除帖子所属话题
func (m *PostingModel) DelPostTopics(postId string) error {
	if _, err := m.Engine.Where("posting_id=?", postId).Delete(&models.PostingTopic{}); err != nil {
		return err
	}

	return nil
}

// 删除帖子统计数据
func (m *PostingModel) DelPostStatistic(postId string) error {
	if _, err := m.Engine.Where("posting_id=?", postId).Delete(&models.PostingStatistic{}); err != nil {
		return err
	}

	return nil
}

// 更新帖子状态[关联的视频]
func (m *PostingModel) UpdatePostStatus(userId, videoId string) error {
	if _, err := m.Engine.Where("user_id=? AND video_id=? AND posting_type=2 AND content_type=0", userId, videoId).
		Cols("status").Update(m.Posting); err != nil {
		return err
	}

	return nil
}


const (
	QUERY_ATTENTION_POSTS = "SELECT p.*, ps.* FROM `posting_info` as p " +
		"LEFT JOIN posting_statistic as ps ON p.id=ps.posting_id WHERE p.status = 1 AND p.video_id=0 AND p.user_id in(%s) " +
		"ORDER BY p.id DESC, p.is_top DESC, p.is_cream DESC LIMIT ?, ?"
)
// 获取关注的用户发布的帖子
func (m *PostingModel) GetPostListByAttention(userIds string, offset, size int) ([]*PostDetailInfo, error) {
	sql := fmt.Sprintf(QUERY_ATTENTION_POSTS, userIds)
	var list []*PostDetailInfo
	if err := m.Engine.SQL(sql, offset, size).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 添加帖子举报
func (m *PostingModel) AddPostReport() (int64, error) {
	return m.Engine.InsertOne(m.Report)
}

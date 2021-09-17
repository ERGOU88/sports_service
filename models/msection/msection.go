package msection

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
)

// 推荐的模块
type SectionModel struct {
	Engine            *xorm.Session
	Section           *models.RecommendInfoSection
}

type SectionRecommendInfo struct {
	Id            int64                 `json:"id"`                            // 视频/富文本id
	Title         string                `json:"title"`                         // 标题
	Describe      string                `json:"describe"`                      // 描述
	Cover         string                `json:"cover"`                         // 封面
	VideoAddr     string                `json:"video_addr"`                    // 视频地址
	Size          int64                 `json:"size"`                          // 视频大小 字节数
	VideoDuration int                   `json:"video_duration"`                // 视频时长
	CreateAt      int                   `json:"create_at"`                     // 创建时间
	FabulousNum   int                   `json:"fabulous_num"`                  // 点赞数
	CommentNum    int                   `json:"comment_num"`                   // 评论数
	ShareNum      int                   `json:"share_num"`                     // 分享数
	UserId        string                `json:"user_id"`                       // 用户id
	Avatar        string                `json:"avatar"`                        // 头像
	Nickname      string                `json:"nick_name"`                     // 昵称
	IsAttention   int                   `json:"is_attention"`                  // 是否关注 1 关注 0 未关注
	IsLike        int                   `json:"is_like"`                       // 是否点赞
	Labels        []*models.VideoLabels `json:"labels"`                        // 视频标签
	StatisticsTab string                `json:"statistics_tab"`                // 统计标签
	ContentType   int                   `json:"content_type"`                  // 内容类型 1 视频 2 资讯
}

// 实栗
func NewSectionModel(engine *xorm.Session) *SectionModel {
	return &SectionModel{
		Engine: engine,
		Section: new(models.RecommendInfoSection),
	}
}

// sectionType  0 首页推荐板块
func (m *SectionModel) GetRecommendSectionByType(sectionType string) ([]*models.RecommendInfoSection, error) {
	var list []*models.RecommendInfoSection
	if err := m.Engine.Where("status=0 AND section_type=?", sectionType).Desc("sortorder").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 通过id获取板块
func (m *SectionModel) GetSectionById(id string) (bool, error) {
	m.Section = new(models.RecommendInfoSection)
	return m.Engine.Where("id=?", id).Get(m.Section)
}

const (
	GET_RECOMMEND_INFO_BY_SECTION_ID = "SELECT v.video_id as id, v.size, v.describe, v.cover, v.video_addr as video_addr, " +
		"v.create_at, 1 as content_type, v.`title`, v.`video_duration`,v.user_id FROM videos as v " +
		"WHERE status=1 AND pub_type=3 AND section_id=? UNION ALL SELECT i.id as id, 0, '', i.cover, '', " +
		"i.create_at, 2 as content_type, i.title, 0, i.user_id FROM information as i " +
		"WHERE status=0 AND pub_type=2 AND related_id=? ORDER BY create_at DESC LIMIT ?, ?"
)
// 通过板块id 获取推荐的内容
func (m *SectionModel) GetRecommendInfoBySectionId(sectionId string, offset, size int) ([]*SectionRecommendInfo, error) {
	var list []*SectionRecommendInfo
	if err := m.Engine.SQL(GET_RECOMMEND_INFO_BY_SECTION_ID, sectionId, sectionId, offset, size).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

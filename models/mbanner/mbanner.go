package mbanner

import (
	"github.com/go-xorm/xorm"
	"sports_service/global/app/log"
	"sports_service/models"
	"sports_service/tools/tencentCloud"
)

type BannerModel struct {
	Banners *models.Banner
	Engine  *xorm.Session
}

// 后台添加banner请求参数
type AddBannerParams struct {
	Cover     string `binding:"required" json:"cover"`
	EndTime   int    `binding:"required" json:"end_time"`
	Explain   string `json:"explain"`
	JumpUrl   string `json:"jump_url"`
	ShareUrl  string `json:"share_url"`
	Sortorder int    `json:"sortorder"`
	StartTime int    `binding:"required" json:"start_time"`
	Title     string `json:"title"`
	Type      int    `json:"type"`
	JumpType  int    `json:"jump_type"`
	VideoAddr string `json:"video_addr"`
}

// 后台更新banner请求参数
type UpdateBannerParams struct {
	Id        int64  `json:"id" binding:"required"`
	Cover     string `json:"cover" binding:"required"`
	EndTime   int    `json:"end_time" binding:"required"`
	Explain   string `json:"explain"`
	JumpUrl   string `json:"jump_url"`
	ShareUrl  string `json:"share_url"`
	Sortorder int    `json:"sortorder"`
	StartTime int    `json:"start_time" binding:"required"`
	Title     string `json:"title"`
	Type      int    `json:"type"`
	JumpType  int    `json:"jump_type"`
	Status    int    `json:"status"`
}

// 后台删除banner请求参数
type DelBannerParam struct {
	Id string `json:"id"`
}

type Banner struct {
	Id         int                    `json:"id" xorm:"not null pk autoincr comment('主键') INT(10)"`
	Title      string                 `json:"title" xorm:"not null default '' comment('标题') VARCHAR(255)"`
	Cover      tencentCloud.BucketURI `json:"cover" xorm:"not null default '' comment('banner封面') VARCHAR(512)"`
	Explain    string                 `json:"explain" xorm:"not null default '' comment('说明') VARCHAR(255)"`
	JumpUrl    tencentCloud.BucketURI `json:"jump_url" xorm:"not null default '' comment('跳转地址') VARCHAR(512)"`
	ShareUrl   tencentCloud.BucketURI `json:"share_url" xorm:"not null default '' comment('分享地址') VARCHAR(512)"`
	Type       int                    `json:"type" xorm:"not null default 1 comment('1 首页 2 赛事 3 官网banner') INT(1)"`
	StartTime  int                    `json:"start_time" xorm:"not null default 0 comment('上架时间') INT(11)"`
	EndTime    int                    `json:"end_time" xorm:"not null default 0 comment('下架时间') INT(11)"`
	Sortorder  int                    `json:"sortorder" xorm:"not null default 0 comment('排序权重') INT(11)"`
	Status     int                    `json:"status" xorm:"not null default 0 comment('0待上架 1上架 2 已过期') TINYINT(1)"`
	CreateAt   int                    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt   int                    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	JumpType   int                    `json:"jump_type" xorm:"not null comment('跳转类型 0 站内跳转 1 站外跳转') TINYINT(1)"`
	VideoAddr  tencentCloud.BucketURI `json:"video_addr" xorm:"not null default '' comment('视频地址') VARCHAR(512)"`
	BannerType int                    `json:"banner_type" xorm:"not null default 0 comment('0图片 1视频') TINYINT(1)"`
}

// 实栗
func NewBannerMolde(engine *xorm.Session) *BannerModel {
	return &BannerModel{
		Banners: new(models.Banner),
		Engine:  engine,
	}
}

const (
	QUERY_BANNER_LIST = "SELECT * FROM `banner` WHERE `type`=? AND `start_time` < ? AND `end_time` > ? ORDER BY sortorder DESC, id DESC LIMIT ?, ?"
)

// 获取首页推荐banner 符合上架时间的 (types: 1 首页 2 直播页 3 官网banner)
func (m *BannerModel) GetRecommendBanners(bannerType int32, tm int64, offset, size int) []*Banner {
	var info []*Banner
	if err := m.Engine.SQL(QUERY_BANNER_LIST, bannerType, tm, tm, offset, size).Find(&info); err != nil {
		log.Log.Errorf("banner_trace: get recommend banners err:%s", err)
		return []*Banner{}
	}

	return info
}

// 添加banner
func (m *BannerModel) AddBanner() error {
	if _, err := m.Engine.InsertOne(m.Banners); err != nil {
		log.Log.Errorf("banner_trace: add banner err:%s", err)
		return err
	}

	return nil
}

// 删除banner
func (m *BannerModel) DelBanner(bannerId string) error {
	if _, err := m.Engine.Where("id=?", bannerId).Delete(&models.Banner{}); err != nil {
		log.Log.Errorf("banner_trace: del banner err:%s", err)
		return err
	}

	return nil
}

// 更新banner
// 状态 0.待上架 1.上架 2.已过期
func (m *BannerModel) UpdateBanner(id int64, cols string) error {
	if _, err := m.Engine.Where("id = ?", id).Cols(cols).Update(m.Banners); err != nil {
		log.Log.Errorf("banner_trace: update banner err:%s", err)
		return err
	}

	return nil
}

const (
	GET_BANNER_LIST = "SELECT * FROM `banner` ORDER BY id DESC LIMIT ?, ?"
)

// 后台获取banner列表
func (m *BannerModel) GetBannerList(offset, size int) []*Banner {
	var info []*Banner
	if err := m.Engine.SQL(GET_BANNER_LIST, offset, size).Find(&info); err != nil {
		log.Log.Errorf("banner_trace: get banner list err:%s", err)
		return []*Banner{}
	}

	return info
}

// 后台获取banner总数目
func (m *BannerModel) GetBannerTotal() int64 {
	total, err := m.Engine.Count(&models.Banner{})
	if err != nil {
		return 0
	}

	return total
}

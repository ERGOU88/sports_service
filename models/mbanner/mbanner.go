package mbanner

import (
  "github.com/go-xorm/xorm"
  "sports_service/server/global/app/log"
  "sports_service/server/models"
)

type BannerModel struct {
	Banners  *models.Banner
	Engine   *xorm.Session
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
	Id         string       `json:"id"`
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
func (m *BannerModel) GetRecommendBanners(bannerType int32, tm int64, offset, size int) []*models.Banner {
	var info []*models.Banner
	if err := m.Engine.SQL(QUERY_BANNER_LIST, bannerType, tm, tm, offset, size).Find(&info); err != nil {
		log.Log.Errorf("banner_trace: get recommend banners err:%s", err)
		return []*models.Banner{}
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
func (m *BannerModel) GetBannerList(offset, size int) []*models.Banner {
	var info []*models.Banner
	if err := m.Engine.SQL(GET_BANNER_LIST, offset, size).Find(&info); err != nil {
		log.Log.Errorf("banner_trace: get banner list err:%s", err)
		return []*models.Banner{}
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


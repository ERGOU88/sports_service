package mbanner

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
	"sports_service/server/global/app/log"
)

type BannerModel struct {
	Banners  *models.Banner
	Engine   *xorm.Session
}

// 实栗
func NewBannerMolde(engine *xorm.Session) *BannerModel {
	return &BannerModel{
		Banners: new(models.Banner),
		Engine:  engine,
	}
}

const (
	QUERY_BANNER_LIST = "SELECT * FROM `banner` WHERE status=1 AND `type`=1 ORDER BY id,sortorder DESC LIMIT 10"
)

// 获取首页推荐banner (types: 1 首页 2 直播页 3 官网banner)
func (m *BannerModel) GetRecommendBanners(bannerType int32) []*models.Banner {
	var info []*models.Banner
	if err := m.Engine.SQL(QUERY_BANNER_LIST, bannerType).Find(&info); err != nil {
		log.Log.Errorf("banner_trace: get recommend banners err:%v", err)
		return nil
	}

	return info
}


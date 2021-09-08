package contest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mbanner"
	"sports_service/server/models/muser"
	"time"
)

type ContestModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	banner      *mbanner.BannerModel
}

func New(c *gin.Context) ContestModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	return ContestModule{
		context: c,
		user: muser.NewUserModel(socket),
		banner: mbanner.NewBannerMolde(socket),
		engine: socket,
	}
}

// 获取赛事首页banner
func (svc *ContestModule) GetBanner() []*models.Banner {
	banners := svc.banner.GetRecommendBanners(int32(consts.CONTEST_BANNERS), time.Now().Unix(), 0, 10)
	if len(banners) == 0 {
		return []*models.Banner{}
	}

	return banners
}

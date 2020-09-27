package configure

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mbanner"
	"sports_service/server/models/muser"
	"time"
)

type ConfigModule struct {
	context     *gin.Context
	engine      *xorm.Session
	banner      *mbanner.BannerModel
	user        *muser.UserModel
}

func New(c *gin.Context) ConfigModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return ConfigModule{
		context: c,
		banner: mbanner.NewBannerMolde(socket),
		user: muser.NewUserModel(socket),
		engine: socket,
	}
}

// 后台添加banner
func (svc *ConfigModule) AddBanner(params *mbanner.AddBannerParams) int {
	now := int(time.Now().Unix())
	svc.banner.Banners.UpdateAt = now
	svc.banner.Banners.CreateAt = now
	svc.banner.Banners.Cover = params.Cover
	svc.banner.Banners.Status = consts.WAIT_LAUNCHE
	svc.banner.Banners.Title = params.Title
	svc.banner.Banners.Explain = params.Explain
	svc.banner.Banners.JumpUrl = params.JumpUrl
	svc.banner.Banners.ShareUrl = params.ShareUrl
	svc.banner.Banners.StartTime = params.StartTime
	svc.banner.Banners.EndTime = params.EndTime
	svc.banner.Banners.Sortorder = params.Sortorder
	svc.banner.Banners.Type = params.Type
	if err := svc.banner.AddBanner(); err != nil {
		return errdef.CONFIG_ADD_BANNER_FAIL
	}

	return errdef.SUCCESS
}

// 后台删除banner
func (svc *ConfigModule) DelBanner(param *mbanner.DelBannerParam) int {
	if err := svc.banner.DelBanner(param.Id); err != nil {
		return errdef.CONFIG_DEL_BANNER_FAIL
	}

	return errdef.SUCCESS
}

// 后台获取banner列表 同时判断时间 更新状态
func (svc *ConfigModule) GetBannerList(page, size int) []*models.Banner {
	offset := (page - 1) * size
	return svc.banner.GetBannerList(offset, size)
}

// 后台添加系统头像
func (svc *ConfigModule) AddSystemAvatar(params *muser.AddSystemAvatarParams) int {
	if err := svc.user.AddSystemAvatar(params); err != nil {
		return errdef.CONFIG_DEL_AVATAR_FAIL
	}

	return errdef.SUCCESS
}

// 后台删除系统头像
func (svc *ConfigModule) DelSystemAvatar(param *muser.DelSystemAvatarParam) int {
	if err := svc.user.DelSystemAvatar(param.Id); err != nil {
		return errdef.CONFIG_DEL_AVATAR_FAIL
	}

	return errdef.SUCCESS
}

// 后台获取系统头像（一次性获取）
func (svc *ConfigModule) GetSystemAvatars() []*models.DefaultAvatar {
	return svc.user.GetSystemAvatarList()
}
package configure

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mbanner"
	"sports_service/server/models/mconfigure"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"sports_service/server/tools/tencentCloud"
	"time"
)

type ConfigModule struct {
	context     *gin.Context
	engine      *xorm.Session
	banner      *mbanner.BannerModel
	user        *muser.UserModel
	video       *mvideo.VideoModel
	configure   *mconfigure.ConfigModel
}

func New(c *gin.Context) ConfigModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	return ConfigModule{
		context: c,
		banner: mbanner.NewBannerMolde(socket),
		user: muser.NewUserModel(socket),
		video: mvideo.NewVideoModel(socket),
		configure: mconfigure.NewConfigModel(socket),
		engine: socket,
	}
}

// 后台添加banner
func (svc *ConfigModule) AddBanner(params *mbanner.AddBannerParams) int {
	now := int(time.Now().Unix())
	if params.EndTime - params.StartTime < 1800 ||
		params.EndTime <= params.StartTime ||
		params.EndTime <= now {
		return errdef.CONFIG_INVALID_END_TIME
	}

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
	svc.banner.Banners.VideoAddr = params.VideoAddr
	// todo: 后台可选跳转形式
	svc.banner.Banners.JumpType = params.JumpType
	if err := svc.banner.AddBanner(); err != nil {
		return errdef.CONFIG_ADD_BANNER_FAIL
	}

	return errdef.SUCCESS
}

// 后台更新banner
func (svc *ConfigModule) UpdateBanner(params *mbanner.UpdateBannerParams) int {
	now := int(time.Now().Unix())
	if params.EndTime - params.StartTime < 1800 ||
		params.EndTime <= params.StartTime ||
		params.EndTime <= now {
		return errdef.CONFIG_INVALID_END_TIME
	}

	svc.banner.Banners.UpdateAt = now
	svc.banner.Banners.Cover = params.Cover
	svc.banner.Banners.Status = params.Status
	svc.banner.Banners.Title = params.Title
	svc.banner.Banners.Explain = params.Explain
	svc.banner.Banners.JumpUrl = params.JumpUrl
	svc.banner.Banners.ShareUrl = params.ShareUrl
	svc.banner.Banners.StartTime = params.StartTime
	svc.banner.Banners.EndTime = params.EndTime
	svc.banner.Banners.Sortorder = params.Sortorder
	svc.banner.Banners.Type = params.Type
	svc.banner.Banners.JumpType = params.JumpType
	cols := "update_at,cover,status,title,explain,jump_url,share_url,start_time,end_time,sortorder,type,jump_type"
	if err := svc.banner.UpdateBanner(params.Id, cols); err != nil {
		return errdef.CONFIG_UPDATE_BANNER_FAIL
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

// 后台获取banner总数
func (svc *ConfigModule) GetBannerTotal() int64 {
	return svc.banner.GetBannerTotal()
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

// 获取热搜配置
func (svc *ConfigModule) GetHotSearch() []*models.HotSearch {
	return svc.video.GetHotSearch()
}

// 添加热搜
func (svc *ConfigModule) AddHotSearch(params *mvideo.AddHotSearchParams) int {
	client := tencentCloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
	// 检测热搜内容
	isPass, err := client.TextModeration(params.HotSearch)
	if !isPass || err != nil {
		return errdef.CONFIG_INVALID_HOT_SEARCH
	}

	if b := svc.video.IsRepeatHotSearchName(params.HotSearch); b {
		return errdef.CONFIG_HOT_NAME_EXISTS
	}

	now := time.Now().Unix()
	svc.video.HotSearch.HotSearchContent = params.HotSearch
	svc.video.HotSearch.Sortorder = params.Sortorder
	svc.video.HotSearch.CreateAt = int(now)
	svc.video.HotSearch.UpdateAt = int(now)
	if err := svc.video.AddHotSearch(); err != nil {
		return errdef.CONFIG_ADD_HOT_SEARCH_FAIL
	}

	return errdef.SUCCESS
}

// 删除热搜
func (svc *ConfigModule) DelHotSearch(params *mvideo.DelHotSearchParams) int {
	if err := svc.video.DelHotSearch(params.Id); err != nil {
		return errdef.CONFIG_DEL_HOT_SEARCH_FAIL
	}

	return errdef.SUCCESS
}

// 热搜内容设置权重
func (svc *ConfigModule) SetSortByHotSearch(params *mvideo.SetSortParams) int {
	now := time.Now().Unix()
	svc.video.HotSearch.Sortorder = params.Sortorder
	svc.video.HotSearch.UpdateAt = int(now)
	if err := svc.video.UpdateSortByHotSearch(params.Id); err != nil {
		return errdef.CONFIG_SET_SORT_HOT_FAIL
	}

	return errdef.SUCCESS
}

// 热搜内容设置状态
func (svc *ConfigModule) SetStatusByHotSearch(params *mvideo.SetStatusParams) int {
	if params.Status != 0 && params.Status != 1 {
		return errdef.INVALID_PARAMS
	}

	now := time.Now().Unix()
	svc.video.HotSearch.Status = params.Status
	svc.video.HotSearch.UpdateAt = int(now)
	if err := svc.video.UpdateStatusByHotSearch(params.Id); err != nil {
		return errdef.CONFIG_SET_STATUS_HOT_FAIL
	}

	return errdef.SUCCESS
}

// 添加新包
func (svc *ConfigModule) AddNewPackage(param *mconfigure.AddPackageParams) int {
	now := int(time.Now().Unix())
	svc.configure.VersionControl.CreateAt = now
	svc.configure.VersionControl.UpdateAt = now
	svc.configure.VersionControl.Platform = int(param.Platform)
	svc.configure.VersionControl.Status = int(param.Status)
	svc.configure.VersionControl.Size = param.Size
	svc.configure.VersionControl.IsForce = int(param.IsForce)
	svc.configure.VersionControl.UpgradeUrl = param.UpgradeUrl
	svc.configure.VersionControl.Version = param.Version
	svc.configure.VersionControl.VersionCode = param.VersionCode
	svc.configure.VersionControl.VersionName = param.VersionName
	svc.configure.VersionControl.Describe = param.Describe
	// 添加新包
	affected, err := svc.configure.AddNewPackage()
	if affected != 1 || err != nil {
		return errdef.CONFIG_ADD_PACKAGE_FAIL
	}

	return errdef.SUCCESS
}

// 更新包信息
func (svc *ConfigModule) UpdatePackageInfo(param *mconfigure.UpdatePackageParams) int {
	now := int(time.Now().Unix())
	svc.configure.VersionControl.UpdateAt = now
	svc.configure.VersionControl.Platform = int(param.Platform)
	svc.configure.VersionControl.Status = int(param.Status)
	svc.configure.VersionControl.Size = param.Size
	svc.configure.VersionControl.IsForce = int(param.IsForce)
	svc.configure.VersionControl.UpgradeUrl = param.UpgradeUrl
	svc.configure.VersionControl.Version = param.Version
	svc.configure.VersionControl.VersionCode = param.VersionCode
	svc.configure.VersionControl.VersionName = param.VersionName
	svc.configure.VersionControl.Describe = param.Describe

	affected, err := svc.configure.UpdatePackageInfo(param.Id)
	if affected != 1 || err != nil {
		return errdef.CONFIG_UPDATE_PACKAGE_FAIL
	}

	return errdef.SUCCESS
}

// 删除包
func (svc *ConfigModule) DelPackage(id int64) int {
	affected, err := svc.configure.DelPackage(id)
	if affected != 1 || err != nil {
		return errdef.CONFIG_DEL_PACKAGE_FAIL
	}

	return errdef.SUCCESS
}

// 获取包列表
func (svc *ConfigModule) GetPackageList(page, size int) []*models.AppVersionControl {
	offset := (page - 1) * size
	list := svc.configure.GetPackageInfoList(offset, size)
	if list == nil {
		return []*models.AppVersionControl{}
	}

	return list
}

// 获取包详情
func (svc *ConfigModule) GetPackageDetail(id string) *models.AppVersionControl {
	return svc.configure.GetPackageDetail(id)
}

package consts

type BannerType int32

// app首页banner type为->1 直播banner->2 web官网banner->3
const (
	HOMEPAGE_BANNERS  BannerType = 1
	LIVE_BANNERS      BannerType = 2
	WEB_BANNERS       BannerType = 3
)

// banner状态 0 待上架  1 上架 2 已过期
const (
	WAIT_LAUNCHE  = 0
	HAS_LAUNCHED  = 1
	NO_LAUNCHED   = 2
)

package consts

type BannerType int32

// app首页banner type为->1 直播banner->2 web官网banner->3
const (
	HOMEPAGE_BANNERS  BannerType = 1
	LIVE_BANNERS      BannerType = 2
	WEB_BANNERS       BannerType = 3
)

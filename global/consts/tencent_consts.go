package consts

const (
	VOD_API_DOMAIN      = "vod.tencentcloudapi.com"
	TX_CLOUD_SECRET_ID  = "AKIDSfbHjxhcmiV3ECAVqvzoNgme8NTIr9C0"
	TX_CLOUD_SECRET_KEY = "g6FDvRkmZ5KJ1SijWpse9OM4XmmJZZke"
)

// NewFileUpload：视频上传完成；
// ProcedureStateChanged：任务流状态变更；
// FileDeleted：视频删除完成；
// PullComplete：视频转拉完成；
// EditMediaComplete：视频编辑完成；
// WechatPublishComplete：微信发布完成；
// ComposeMediaComplete：制作媒体文件完成；
// WechatMiniProgramPublishComplete：微信小程序发布完成。
// 兼容 2017 版的事件类型：
// TranscodeComplete：视频转码完成；
// ConcatComplete：视频拼接完成；
// ClipComplete：视频剪辑完成；
// CreateImageSpriteComplete：视频截取雪碧图完成；
// CreateSnapshotByTimeOffsetComplete：视频按时间点截图完成。
const (
	// 视频上传事件
	EVENT_TYPE_UPLOAD = "NewFileUpload"
)

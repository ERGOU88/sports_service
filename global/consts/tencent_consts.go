package consts

const (
	// 点播
	VOD_API_DOMAIN      = "vod.tencentcloudapi.com"
	// 文本内容检测
	TMS_API_DOMAIN      = "tms.tencentcloudapi.com"
	TX_CLOUD_SECRET_ID  = "AKIDFI5bssiLhodBSTDgtsZz8zbx2qOffOq1"
	TX_CLOUD_SECRET_KEY = "vZw2rPgIejX5MM5WhaDZdQwA8qHRJuEH"

	TX_CLOUD_COS_SECRET_ID  = "AKIDjU25ybRyZ4EHVyzemjOoIcZGrIH6NEYk"
	TX_CLOUD_COS_SECRET_KEY = "6v96wZaOjnmKbbWS9qGZfeqofKzSfz9h"

	TX_SMS_SECRET_ID    = "AKIDFI5bssiLhodBSTDgtsZz8zbx2qOffOq1"
	TX_SMS_SECRET_KEY   = "vZw2rPgIejX5MM5WhaDZdQwA8qHRJuEH"

	// 任务流模版名称
	VOD_PROCEDURE_NAME  = "fpv-demo"
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
	EVENT_TYPE_UPLOAD             = "NewFileUpload"
	// 任务流状态变更 包含视频转码完成
	EVENT_PROCEDURE_STATE_CHANGED = "ProcedureStateChanged"
    // 文件被删除
	EVENT_FILE_DELETED = "FileDeleted"
)

// 0 视频上传事件 1 任务流状态变更（包含视频转码完成）
const (
	EVENT_UPLOAD_TYPE                   int = iota
	EVENT_PROCEDURE_STATE_CHANGED_TYPE
)

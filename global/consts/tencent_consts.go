package consts

const (
	TX_APP_ID           = 1253904687
	// 点播
	VOD_API_DOMAIN      = "vod.tencentcloudapi.com"
	// 文本内容检测
	TMS_API_DOMAIN      = "tms.tencentcloudapi.com"
	TX_CLOUD_SECRET_ID  = "AKIDFI5bssiLhodBSTDgtsZz8zbx2qOffOq1"
	TX_CLOUD_SECRET_KEY = "vZw2rPgIejX5MM5WhaDZdQwA8qHRJuEH"

	TX_CLOUD_COS_SECRET_ID  = "AKIDFI5bssiLhodBSTDgtsZz8zbx2qOffOq1"
	TX_CLOUD_COS_SECRET_KEY = "vZw2rPgIejX5MM5WhaDZdQwA8qHRJuEH"

	TX_SMS_SECRET_ID    = "AKIDFI5bssiLhodBSTDgtsZz8zbx2qOffOq1"
	TX_SMS_SECRET_KEY   = "vZw2rPgIejX5MM5WhaDZdQwA8qHRJuEH"

	// 任务流模版名称
	// [自适应转码]
	VOD_PROCEDURE_NAME  = "fpv-demo"
	// [转 标清]
	VOD_PROCEDURE_TRANSCODE_1 = "transcode_1"
	// [转 标清 + 高清]
	VOD_PROCEDURE_TRANSCODE_2 = "transcode_2"
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

// 0 视频上传事件 1 任务流状态变更（包含视频转码完成） 2 图片审核
const (
	EVENT_UPLOAD_TYPE                   int = iota
	EVENT_PROCEDURE_STATE_CHANGED_TYPE
	EVENT_VERIFY_IMAGE_TYPE
)

// 视频AI审核结果
const (
	VIDEO_AI_AUDIT_PASS   = "pass"   // 通过
	VIDEO_AI_AUDIT_REVIEW = "review" // 建议复审
	VIDEO_AI_AUDIT_BLOCK  = "block"  // 屏蔽
)

// AI审核状态
// 1 通过
// 2 不通过
// 3 复审
const (
	AI_AUDIT_PASS   = 1
	AI_AUDIT_BLOCK  = 2
	AI_AUDIT_REVIEW = 3
)

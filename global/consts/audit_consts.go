package consts

// [视频、帖子等内容审核]
// 1 AI + 人工审核 [AI审核如果不通过 则内容为不通过 其他状态都需要人工复审]
// 2 人工审核
const (
	AUDIT_MODE_AI_AND_MANUAL = "1"
	AUDIT_MODE_MANUAL        = "2"
)

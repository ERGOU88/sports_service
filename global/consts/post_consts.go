package consts

// 状态 -1 查看所有 0 审核中 1 审核成功 2 审核失败 3 逻辑删除
const (
	POST_VIEW_ALL      = "-1"
	POST_UNDER_REVIEW  = "0"
	POST_AUDIT_SUCCESS = "1"
	POST_AUDIT_FAILURE = "2"
	POST_DELETE_STATUS = "3"
)

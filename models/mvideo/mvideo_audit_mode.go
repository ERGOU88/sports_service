package mvideo

import (
	"sports_service/dao"
	"sports_service/global/consts"
	"sports_service/global/rdskey"
)

// 设置[视频、帖子等] 审核模式 1 人工 + AI 2 人工
func (m *VideoModel) SetAuditMode(mode int) error {
	rds := dao.NewRedisDao()
	return rds.Set(rdskey.AUDIT_MODE, mode)
}

// 获取内容审核模式 如redis获取不到 则默认为 AI + 人工
func (m *VideoModel) GetAuditMode() string {
	rds := dao.NewRedisDao()
	mode, err := rds.Get(rdskey.AUDIT_MODE)
	if mode == "" || err != nil {
		return consts.AUDIT_MODE_AI_AND_MANUAL
	}

	return consts.AUDIT_MODE_MANUAL
}

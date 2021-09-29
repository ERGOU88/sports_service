package mposting

type AudiPostParam struct {
	Id       string     `json:"id"`
	Status   int        `json:"status"`
}

type SettingParam struct {
	SettingType   int   `json:"setting_type"` // 1 精华 2 置顶
	ActionType    int   `json:"action_type"`  // 1 设置 0 取消
	Id            int64 `json:"id"`           // 帖子id
}

// todo: 后台查询帖子审核列表时 需过滤掉发布的视频 以及 帖子审核通过时 需给up主的粉丝们发推送通知
// 更新帖子审核状态 不包含关联视频的帖子
func (m *PostingModel) UpdateStatusByPost() error {
	if _, err := m.Engine.Where("id=? AND video_id=0", m.Posting.Id).
		Cols("status").Update(m.Posting); err != nil {
		return err
	}

	return nil
}

const (
	GET_POST_LIST = "SELECT p.*, ps.fabulous_num, ps.browse_num, ps.share_num, ps.comment_num, ps.heat_num FROM " +
		"`posting_info` AS p LEFT JOIN `posting_statistic` as ps ON p.id=ps.posting_id WHERE p.is_top=0 " +
		" ORDER BY p.is_cream DESC, p.is_top DESC, p.id DESC LIMIT ?, ?"
)
// 获取帖子列表 [管理后台]
func (m *PostingModel) GetPostList(offset, size int) ([]*PostDetailInfo, error) {
	var list []*PostDetailInfo
	if err := m.Engine.SQL(GET_POST_LIST, offset, size).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

func (m *PostingModel) UpdatePostInfo(id int64, cols string) (int64, error) {
	return m.Engine.Where("id=?", id).Cols(cols).Update(m.Posting)
}

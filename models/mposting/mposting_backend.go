package mposting

type AudiPostParam struct {
	Id       string     `json:"id"`
	Status   int        `json:"status"`
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

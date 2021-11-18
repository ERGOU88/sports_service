package mposting

import (
	"sports_service/server/models"
	"fmt"
)

type AudiPostParam struct {
	Id       string     `json:"id"`
	Status   int        `json:"status"`
}

type SettingParam struct {
	SettingType   int   `json:"setting_type"` // 1 精华 2 置顶
	ActionType    int   `json:"action_type"`  // 1 设置 0 取消
	Id            int64 `json:"id"`           // 帖子id
}

// 批量修改帖子数据
type BatchEditParam struct {
	EditType   int32       `json:"edit_type" binding:"required"` // 1 编辑帖子标题 2 编辑帖子板块 3 编辑帖子话题
	Ids        []int64     `json:"ids" binding:"required"`       // 帖子id
	Title      string      `json:"title"`
	SectionId  int         `json:"section_id"`
	TopicIds   []int64     `json:"topic_ids"`
}

type TopicInfo struct {
	Id        int    `json:"id"`
	TopicName string `json:"topic_name"`
	Sortorder int    `json:"sortorder"`
	Status    int    `json:"status"`
	IsHot     int    `json:"is_hot"`
	CreateAt  int    `json:"create_at"`
	UpdateAt  int    `json:"update_at"`
	Cover     string `json:"cover"`
	Describe  string `json:"describe"`
	SectionId int    `json:"section_id"`
	PostNum   int64  `json:"post_num"`
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
		"`posting_info` AS p LEFT JOIN `posting_statistic` as ps ON p.id=ps.posting_id " +
		" ORDER BY p.is_cream DESC, p.is_top DESC, p.id DESC LIMIT ?, ?"
)
// 获取帖子列表 [管理后台]
func (m *PostingModel) GetPostList(offset, size int, status, title string) ([]*PostDetailInfo, error) {
	sql := "SELECT p.*, ps.fabulous_num, ps.browse_num, ps.share_num, ps.comment_num, ps.heat_num FROM " +
		"`posting_info` AS p LEFT JOIN `posting_statistic` as ps ON p.id=ps.posting_id "

	// 取审核中/审核失败
	if status == "0" {
		sql += " WHERE p.status in(0, 2) "
	} else {
		sql += fmt.Sprintf(" WHERE p.status=%s ", status)
	}

	if title != "" {
		sql += " AND p.title Like '%" + title + "%' "
	}

	sql +=  fmt.Sprintf(" ORDER BY p.is_cream DESC, p.is_top DESC, p.id DESC LIMIT %d, %d", offset, size)
	var list []*PostDetailInfo
	if err := m.Engine.SQL(sql).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

func (m *PostingModel) UpdatePostInfo(id int64, cols string) (int64, error) {
	return m.Engine.Where("id=?", id).Cols(cols).Update(m.Posting)
}

// 获取帖子总数
func (m *PostingModel) GetTotalCountByPost(condition []int, title string) int64 {
	m.Engine.In("status", condition)
	if title != "" {
		m.Engine.Where("title like '%" + title + "%'")
	}

	count, err := m.Engine.Count(&models.PostingInfo{})
	if err != nil {
		return 0
	}

	return count
}

const (
	GET_APPLY_CREAM_LIST = " SELECT p.* FROM posting_apply_cream AS pac LEFT JOIN posting_info AS p ON pac.post_id=p.id " +
		"WHERE pac.status=0 ORDER BY pac.id DESC LIMIT ?,?"
)
func (m *PostingModel) GetApplyCreamList(offset, size int) ([]*PostDetailInfo, error) {
	var list []*PostDetailInfo
	if err := m.Engine.SQL(GET_APPLY_CREAM_LIST, offset, size).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}


func (m *PostingModel) GetApplyCreamCount() (int64, error) {
	sql := "SELECT p.* FROM posting_apply_cream AS pac LEFT JOIN posting_info AS p ON pac.post_id=p.id WHERE pac.status=0"
	return m.Engine.SQL(sql).Count(&models.PostingInfo{})
}

// 更新申精状态
func (m *PostingModel) UpdateApplyCreamStatus(id int64) (int64, error) {
	return m.Engine.Where("id=?", id).Cols("status").Update(m.ApplyCream)
}

// 批量编辑
func (m *PostingModel) BatchEditPost(postIds []int64) (int64, error) {
	return m.Engine.In("id", postIds).Update(m.Posting)
}

// 批量删除帖子话题
func (m *PostingModel) BatchDelPostTopic(postIds []int64) (int64, error) {
	return m.Engine.In("posting_id", postIds).Delete(m.PostingTopic)
}

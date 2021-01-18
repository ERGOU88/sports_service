package mnotify

import (
  "fmt"
  "sports_service/server/models"
)

// 管理后台获取系统通知
func (m *NotifyModel) GetSystemNotifyList(offset, size int) []*models.SystemMessage {
  sql := "SELECT * FROM system_message WHERE send_type=0 ORDER BY system_id DESC LIMIT ?, ?"
  var list []*models.SystemMessage
  if err := m.Engine.SQL(sql, offset, size).Find(&list); err != nil {
    return nil
  }

  return list
}

// 添加系统通知
func (m *NotifyModel) AddSystemNotify(msg *models.SystemMessage) (int64, error) {
  return m.Engine.InsertOne(msg)
}

// 添加系统通知（多个）
func (m *NotifyModel) AddMultiSystemNotify(msg []*models.SystemMessage) (int64, error)  {
  return m.Engine.InsertMulti(msg)
}

// 更新系统通知消息发送状态（已发送 0 /已撤回 2）
func (m *NotifyModel) UpdateSendStatus(status, id int64) error {
  sql := fmt.Sprintf("UPDATE `system_message` SET `send_status`=%d WHERE system_id = %d", status, id)
  if _, err := m.Engine.Exec(sql); err != nil {
    return err
  }

  return nil
}

// 获取所有系统通知
func (m *NotifyModel) GetAllSystemNotify() []*models.SystemMessage {
  sql := "SELECT * FROM system_message WHERE send_type=0 ORDER BY system_id DESC"
  var list []*models.SystemMessage
  if err := m.Engine.SQL(sql).Find(&list); err != nil {
    return nil
  }

  return list

}

// 删除系统通知
func (m *NotifyModel) DelSystemNotify(id int64) (int64, error) {
  return m.Engine.ID(id).Delete(&models.SystemMessage{})
}



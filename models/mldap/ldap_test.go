package mldap

import (
	"testing"
)

// 保存用户通知设置
func TestSaveUserNotifySetting(t *testing.T) {
	model := NewAdModel()
	err := model.CheckLogin("jiangkp", "bluetrans2021")
	t.Logf("ad login err:%s", err)
}

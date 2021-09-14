package im

import (
	"testing"
)


//func TestAddUser(t *testing.T) {
//	im := NewImRealize()
//	sig, err := im.AddUser("123456", "robot1", "https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/123.jpeg")
//	t.Logf("sig:%s, err:%s", sig, err)
//}

func TestGenSig(t *testing.T) {
	sig, err := GenSig("123", EXPIRE_TM_DAY * 365)
	t.Logf("sig:%s, err:%s", sig, err)
}

func TestCreateGroup(t *testing.T) {
	im := NewImRealize()
	groupId, err := im.CreateGroup("AVChatRoom", "", "test", "test",
		"test", "test")
	t.Logf("groupId:%s, err:%s", groupId, err)
}

package main

import (
	pbBarrage "sports_service/server/proto/barrage"
	"sports_service/server/util"
	"sync"
	"fmt"
)

var userMap map[string]*User
var userMutex *sync.RWMutex

func init() {
	userMap = make(map[string]*User, 10240)
	userMutex = new(sync.RWMutex)
}

// 增加用户
func PutUser(xid string, user *User) {
	userMutex.Lock()
	userMap[xid] = user
	userMutex.Unlock()
}

// 删除用户
func DelUser(xid string) {
	userMutex.Lock()
	if _, ok := userMap[xid]; ok {
		delete(userMap, xid)
	}
	userMutex.Unlock()
}

// 广播消息 消息id 统一 -1
func Broadcast(body []byte) {
	msg := &pbBarrage.Message{
		MsgId: "-1",
		MsgType: pbBarrage.MessageType_TYPE_BROADCAST,
		Body: body,
	}
	var user *User
	userMutex.RLock()
	for _, user = range userMap {
		if err := user.Push(msg); err != nil {
			fmt.Println("push msg err:", err)
		}
	}

	userMutex.RUnlock()
}

// 推送消息
func PushMessage(xid string, body []byte, msgType pbBarrage.MessageType) {
	msg := &pbBarrage.Message{
		MsgId: fmt.Sprint(util.GetSnowId()),
		MsgType: msgType,
		Body: body,
	}

	userMutex.RLock()
	user := userMap[xid]
	userMutex.RUnlock()

	if user != nil {
		if err := user.Push(msg); err != nil {
			fmt.Println("push msg err:", err)
		}
	}
}

// 剔除用户
func KickUser(xid string) {
	userMutex.RLock()
	if user, ok := userMap[xid]; user != nil && ok {
		user.Close()
	}
	userMutex.RUnlock()
}

// 剔除所有用户
func KickAllUser() {
	userMutex.Lock()
	for _, user := range userMap {
		user.Close()
	}
	userMutex.Unlock()
}



package main

import (
	"fmt"
	pbBarrage "sports_service/proto/barrage"
)

type User struct {
	Xid         string                  // 链接唯一标示（非用户真实id）
	WriteChanel chan *pbBarrage.Message // 用于序列化所有写入的数据包
	Finish      chan int
}

func NewUser(xid string) *User {
	return &User{
		Xid:         xid,
		WriteChanel: make(chan *pbBarrage.Message, MAX_MSG_PER_CONN),
		Finish:      make(chan int, 2),
	}
}

// 推送消息
func (u *User) Push(msg *pbBarrage.Message) error {
	if len(u.WriteChanel) >= MAX_MSG_PER_CONN {
		return fmt.Errorf("write channel full")
	}

	u.WriteChanel <- msg
	return nil
}

func (u *User) Close() {
	select {
	case u.Finish <- 1:
	default:
	}
}

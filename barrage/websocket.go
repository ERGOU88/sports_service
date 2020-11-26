package main

import (
	"crypto/md5"
	"github.com/gorilla/websocket"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"os"
	"sports_service/server/global/consts"
	pbBarrage "sports_service/server/proto/barrage"
	"time"
	"sports_service/server/util"
	"github.com/golang/protobuf/proto"
)

const (
	READ_BUF_SIZE  = 2048
	WRITE_BUF_SIZE = 2045
	AUTH_DEAD_LINE = 5
	READ_DEAD_LINE = 120
)

const (
	MAX_MSG_PER_CONN = 1024
)

var running = false

var upgrader = websocket.Upgrader{
	ReadBufferSize: READ_BUF_SIZE,
	WriteBufferSize: WRITE_BUF_SIZE,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 开启websocket服务
func StartWebsocket(bindAddr string) {
	r := mux.NewRouter()
	r.HandleFunc("/ws", WebSocketHandler)
	http.Handle("/", r)

	running = true

	go func() {
		if err := http.ListenAndServe(bindAddr, nil); err != nil {
			fmt.Println("listen fail", err)
			os.Exit(1)
		}
	}()

}

// 停止websocket服务
func StopWebsocket() {
	running = false
	KickAllUser()
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	if !running {
		http.Error(w, "system rebooting...", 406)
		return
	}

	// 升级协议
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil || ws == nil {
		fmt.Printf("upgrader.Upgrade() failed: %v raw header: %v", err, r.Header)
		return
	}

	fmt.Println("new connection from :", ws.RemoteAddr().String())
	// todo: 处理socket连接
	handleSocketConn(ws)
}

func handleSocketConn(webConn *websocket.Conn) {
	defer func(){
		webConn.Close()
	}()

	webConn.SetReadDeadline(time.Now().Add(time.Second * AUTH_DEAD_LINE))
	xid := auth(webConn)
	if xid == "" {
		fmt.Printf("\nxid empty, xid:%s", xid)
		return
	}
	//xid := "1"

	user := NewUser(xid)
	// 用户存储到全局map
	PutUser(xid, user)

	webConn.SetReadDeadline(time.Now().Add(time.Second * READ_DEAD_LINE))
	// 将所有写操作串行到一个 goroutine 中
	go writeRoutine(webConn, user)

	EndReadLoop: // 结束读routine标示
	// Read Routine
	for {
		// 心跳重置过期时间
		webConn.SetReadDeadline(time.Now().Add(time.Second * 120))

		msgType, readMsg, err := webConn.ReadMessage()
		if err != nil || msgType != websocket.BinaryMessage {
			fmt.Printf("\nread msg err:%v, msgType:%d", err, msgType)
			break
		}

		msg := &pbBarrage.Message{}
		if err := proto.Unmarshal(readMsg, msg); err != nil {
			fmt.Printf("\nproto unmarshal msg err:%s", err)
			break
		}

		switch msg.MsgType {
		// 链接请求（认证）
		case pbBarrage.MessageType_TYPE_CONN:
		// 心跳
		case pbBarrage.MessageType_TYPE_HEART_BEAT:
		default:
			fmt.Printf("\nunsupported message type %d", msg.MsgType)
			// 对不支持的消息，结束连接
			break EndReadLoop
		}
	}

	user.Close()
	DelUser(user.Xid)
}


func writeRoutine(webConn *websocket.Conn, user *User) {
	var msg *pbBarrage.Message
	for {
		select {
		case msg = <- user.WriteChanel:
			if data, err := proto.Marshal(msg); err == nil {
				if err := webConn.WriteMessage(websocket.BinaryMessage, data); err != nil {
					fmt.Println("\nwrite message err:", err)
					webConn.Close()
					return
				}
			}
		case <- user.Finish:
			// 读操作已检测到异常
			webConn.Close()
			return

		}
	}
}


var WsKey = map[string]string{
	string(consts.WEB_APP_ID):       "PlvZrGmBKGuQPXVb",
	string(consts.IOS_APP_ID):       "RfhHecN9zsNcy19Y",
	string(consts.AND_APP_ID):	     "InaukEwVLLpcewX6",
}

// 认证
func auth(webConn *websocket.Conn) string {
	var res *pbBarrage.ResConnMessage = nil
	defer func() {
		connMsgResponse(webConn, res)
	}()

	msgType, readMsg, err := webConn.ReadMessage()
	if err != nil || msgType != websocket.BinaryMessage {
		fmt.Printf("\nread msg err:%v, msgType:%d", err, msgType)
		res = &pbBarrage.ResConnMessage{
			Code: pbBarrage.RetCode_CODE_FAIL,
			Content: fmt.Sprintf("Read Msg Error %v", err),
		}

		return ""
	}

	msg := pbBarrage.Message{}
	if err := proto.Unmarshal(readMsg, &msg); err != nil {
		fmt.Printf("proto.Unmarshal msg err:%v", err)
		res = &pbBarrage.ResConnMessage{
			Code: pbBarrage.RetCode_CODE_FAIL,
			Content: fmt.Sprintf("Unmarshal Msg Error %v", err),
		}

		return ""
	}

	if msg.MsgType != pbBarrage.MessageType_TYPE_CONN {
		fmt.Printf("invalid msg type: %v", msg.MsgType)
		res = &pbBarrage.ResConnMessage{
			Code: pbBarrage.RetCode_CODE_FAIL,
			Content: fmt.Sprintf("Invalid Type %v ", msg.MsgType),
		}

		return ""
	}

	connReq := pbBarrage.ReqConnMessage{}
	if err := proto.Unmarshal(msg.Body, &connReq); err != nil {
		res = &pbBarrage.ResConnMessage{
			Code: pbBarrage.RetCode_CODE_FAIL,
			Content: "Unmarshal Body Error",
		}

		return ""
	}

	if connReq.Sign == "" || connReq.Timestamp == "" || connReq.AppId == "" || connReq.Version == "" ||
		connReq.Secret == "" {
		fmt.Println("invalid param")
		res = &pbBarrage.ResConnMessage{
			Code: pbBarrage.RetCode_CODE_FAIL,
			Content: "Invalid Param",
		}

		return ""
	}

	str := fmt.Sprintf("AppId=%s&Timestamp=%s&Version=%s&Secret=%s", connReq.AppId, connReq.Timestamp, connReq.Version, connReq.Secret)
	if b := verifySign(str, connReq.AppId, connReq.Sign); !b {
		res = &pbBarrage.ResConnMessage{
			Code: pbBarrage.RetCode_CODE_FAIL,
			Content: "Sign Not Match",
		}

		return ""
	}


	res = &pbBarrage.ResConnMessage{
		Code: pbBarrage.RetCode_CODE_SUCCESS,
		Content: "Success",
	}

	xid := util.GetXID()
	return fmt.Sprint(xid)
}

func connMsgResponse(webConn *websocket.Conn, res *pbBarrage.ResConnMessage) {
	bts, err := proto.Marshal(res)
	if err != nil {
		bts = []byte(fmt.Sprintf("%v", err))
	}

	msg := pbBarrage.Message{
		MsgId: fmt.Sprint(util.GetSnowId()),
		MsgType: pbBarrage.MessageType_TYPE_CONN_RES,
		Body: bts,
	}

	if data, err := proto.Marshal(&msg); err == nil {
		if err := webConn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			fmt.Println("write message err:", err)
			webConn.Close()
			return
		}
	}
}

// 校验签名是否一致
func verifySign(str, appId, sign string) bool {
	appKey := getWsKey(appId)
	str = fmt.Sprintf("&%s", appKey)
	data := []byte(str)
	has := md5.Sum(data)
	md5Str := fmt.Sprintf("%x", has)
	fmt.Printf("client sign:%v, md5Str:%v", sign, md5Str)
	if md5Str == sign {
		return true
	}

	return false
}

// 获取websocketKey
func getWsKey(appId string) string {
	appKey, ok := WsKey[appId]
	if ok {
		return appKey
	}

	return ""
}






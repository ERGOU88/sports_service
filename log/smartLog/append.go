package smartLog

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

var logTag = []string{"DEBUG", "INFO", "WARNING", "ERR"}

//AddExStr 记录文本
func AddExStr(str string) {
	s := getSLE()
	if s != nil {
		s.Ex += str
		s.Ex += "\t"
	}
}

// AddURL 1. 在中间件中添加url到日志中
func AddURL(path string) {
	s := getSLE()
	if s != nil {
		s.Path = path
	}
	s.Online = GlobleSLE.InitURL(path)
}

// AddArgs 2. 添加参数
func AddArgs(key, v string) {
	s := getSLE()
	if s != nil {
		s.Args[key] = v
	}
}

// AddLog 2. 插入任意条日志
func AddLog(level int, data string, deep int64) {
	s := getSLE()
	if s != nil {
		if level >= len(logTag) {
			level = 3
		}
		if level < 0 {
			level = 0
		}
		// 如果有大于1的则是错误模式
		if level > 1 {
			s.Status = -1
		}
		_, file, line, ok := runtime.Caller(2)
		if ok {
			n := strings.LastIndex(file, "eating_chick")
			if n > 0 {
				file = file[n:]
			}
		}

		// 获取文件名
		s.Logs = append(s.Logs, &LogInfo{Caller: fmt.Sprint(file, ":", line), Level: logTag[level], Info: data, Deep: deep})
	}
}

// AddExID 3. 任意时刻可插入扩展索引id
func AddExID(id int64) {
	s := getSLE()
	if s != nil {
		s.UID = id
	}
}

// AddResp 4. 添加返回值
func AddResp(key string, value interface{}) {
	s := getSLE()
	if s != nil {
		s.Rets[key] = value
	}
}

// End 5. 调用end来结束日志-->不调用会超时退出5s
func End(n int) {
	s := getSLE()
	if s != nil {
		if n != 200 {
			s.Status = -1
		}
		go DumpDev(s)
	}

}

func getTID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func getSLE() *SLE {
	tid := int64(getTID())
	var s *SLE
	v, ok := GlobleSLE.Sles.Load(tid)
	if !ok {
		// 如果没有,创建一个
		// log.Println("tid==>>", tid)
		s = createSLE(tid)
		GlobleSLE.AddToID(tid)
		GlobleSLE.Sles.Store(tid, s)
	} else {
		// log.Println("tid++++", tid)
		e, ok := v.(*SLE)
		if !ok {
			//?????
		}
		s = e
	}
	return s
}

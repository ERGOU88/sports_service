package smartLog

import (
	"sync"
	"time"
	"go.uber.org/zap"
)

// GlobleSLE 智能日志管理
var GlobleSLE *SLEManage
var zaplog *zap.SugaredLogger

func init() {
	GlobleSLE = createSLEManage()
	// logger, _ := zap.NewProduction()

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.DisableCaller = true
	loggerConfig.DisableStacktrace = true
	loggerConfig.EncoderConfig.CallerKey = ""
	loggerConfig.EncoderConfig.TimeKey = ""
	loggerConfig.EncoderConfig.LevelKey = ""
	loggerConfig.Encoding = "console"
	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}

	defer logger.Sync() // flushes buffer, if any
	zaplog = logger.Sugar()
}

// 聪明日志系统
// 功能:
// 1. 根据线程id来分类日志.
// 2. 用户uid作为检索条件
// 3. 时间为检索要素, 分为近,0-3天, 中,3-30天,远 三种状态.

// SLE SmartLogElement
type SLE struct {
	UID   int64  // 用户id
	TID   int64  // 线程id --> 有效时间5分钟
	Path  string // URL路径
	Start time.Time
	Args  map[string]string      // post/get 参数
	Rets  map[string]interface{} // 返回值记录
	Logs  []*LogInfo             // 日志文本记录
	Ex    string                 // 扩展文本

	Status int   // 0 正常模式, -1 出错模式
	Online int64 // 上线时间
}

func createSLE(tid int64) *SLE {
	r := &SLE{}
	r.TID = tid
	r.Start = time.Now()
	r.Args = make(map[string]string)
	r.Rets = make(map[string]interface{})
	r.Logs = make([]*LogInfo, 0, 100)
	return r
}

// LogInfo 日志信息
type LogInfo struct {
	Caller string // 函数行和名字
	Level  string // INFO WARING DEBUG ERROR FATAL
	Info   string // 日志内容文字 注意换行优化
	Deep   int64  // 调用深度 1: app接口 2 后端逻辑
}

// SLEManage ...
type SLEManage struct {
	// 缓存当前线程id(int64) hash
	Sles sync.Map
	// 制作超时列表
	toListChan chan int64
	// url 列表,保存url的创建时间
	URLs sync.Map
}

// InitURL 注册上线时间--暂时按重启算
func (s *SLEManage) InitURL(URL string) int64 {
	v, ok := s.URLs.Load(URL)
	if !ok {
		n := time.Now().Unix()
		s.URLs.Store(URL, n)
		return n
	}
	return v.(int64)
}

// AddToID 添加id
func (s *SLEManage) AddToID(id int64) {
	s.toListChan <- id
}

func (s *SLEManage) deamonTimeoutSLE() {
	for {
		n := <-GlobleSLE.toListChan
		if n != 0 {
			v, ok := s.Sles.Load(n)
			if !ok {
				continue
			}
			sle := v.(*SLE)
			dst := sle.Start.Add(time.Second * 5).Unix()
			cur := time.Now().Unix()
			if cur < dst {
				<-time.After(time.Second * time.Duration(dst-cur))
			}
			// 重新确认未打印
			v, ok = s.Sles.Load(n)
			if !ok {
				continue
			}
			sle = v.(*SLE)
			sle.Status = -2

			DumpDev(sle)
		}
	}
}

func createSLEManage() *SLEManage {
	sle := &SLEManage{}
	sle.toListChan = make(chan int64, 2000)
	go sle.deamonTimeoutSLE()
	return sle
}

package smartLog

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// json
// dev
var lk sync.Mutex

// DumpDev ...
func DumpDev(e *SLE) {
	lk.Lock()
	defer lk.Unlock()
	// dumplog := zaplog.With(zap.Namespace("EX"), zap.Int64("TID", e.TID), zap.Int64("UID", e.UID), zap.String("T", e.Start.Format("06-01-02 15:04:05")))
	tag := fmt.Sprint("[", e.TID, "\t", e.UID, "\t", e.Start.Format("06-01-02 15:04:05"), "]\t")

	// 输出URL
	// zaplog.Infow("TID", e.TID, "UID", e.UID, "time", e.Start, "URL", e.Path)
	zaplog.Info(tag, e.Path)
	// dumplog.Info(e.Path)
	// 输出参数
	for i, v := range e.Args {
		zaplog.Info(tag, i, "=", v)
		// dumplog.Info(i, "=", v)
	}

	if (e.Status != 0) || ((e.Status == 0) && (checkStage(e.Online) < 3)) {
		// 输出日志
		for i := 0; i < len(e.Logs); i++ {
			nstr := strings.Split(e.Logs[i].Info, "\r\n")
			for j := 0; j < len(nstr); j++ {
				zaplog.Info(tag, nstr[j], "\t", e.Logs[i].Caller)
				// dumplog.Info(e.Logs[i].Level, nstr[j])
			}
		}

		// 输出返回值
		for i, v := range e.Rets {
			// fmt.Print(tag, "\t", i, "=")
			zaplog.Infof("%s%s=%#v", tag, i, v)
			// dumplog.Infof("%s=%#v\n", i, v)
		}
	}

	if len(e.Ex) > 0 {
		zaplog.Info(tag, e.Ex)
	}

	if e.Status == -2 {
		zaplog.Info(tag, "timeout...")
	}
	// log.Println("rm tid,--->", e.TID)
	GlobleSLE.Sles.Delete(e.TID)
	zaplog.Info("")
}

func checkStage(t int64) int {
	d := (time.Now().Unix() - t) / 3600 / 24
	if d <= 3 {
		return 1
	} else if d < 30 {
		return 2
	}
	return 3
}

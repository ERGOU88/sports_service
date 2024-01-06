package log

import (
	"sports_service/log"
	"sports_service/log/default"
)

func init() {
	Log = _default.InitDefaultLog()
}

var (
	// Log 第三方日志接口
	Log log.ILogger
)

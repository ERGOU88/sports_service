package log

import (
	"github.com/go-xorm/core"
)

// logger interface
type ILogger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugw(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Errorw(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Infow(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Warnw(format string, v ...interface{})
	Level() core.LogLevel
	SetLevel(l core.LogLevel)
	ShowSQL(show ...bool)
	IsShowSQL() bool
	ADDSugar(key string, val interface{})
}

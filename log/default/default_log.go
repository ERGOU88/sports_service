package _default

import (
	"log"

	"github.com/go-xorm/core"
)

type defaultLog struct {

}

func (d *defaultLog) Debug(v ...interface{}) {
	var out []interface{}
	out = append(out, "[D]")
	out = append(out, v...)

	log.Print(out...)
}

func (d *defaultLog) Debugf(format string, v ...interface{}) {
	log.Printf("[D]"+format, v...)
}

func (d *defaultLog) Debugw(format string, v ...interface{}) {
	log.Printf("[DW]"+format, v...)
}

func (d *defaultLog) Error(v ...interface{}) {
	var out []interface{}
	out = append(out, "[E]")
	out = append(out, v...)
	log.Print(out...)
}

func (d *defaultLog) Errorf(format string, v ...interface{}) {
	log.Printf("[E]"+format, v...)
}

func (d *defaultLog) Errorw(format string, v ...interface{}) {
	log.Printf("[EW]"+format, v...)
}

func (d *defaultLog) Info(v ...interface{}) {
	var out []interface{}
	out = append(out, "[I]")
	out = append(out, v...)
	log.Print(out...)
}

func (d *defaultLog) Infof(format string, v ...interface{}) {
	log.Printf("[I]"+format, v...)
}

func (d *defaultLog) Infow(format string, v ...interface{}) {
	log.Printf("[IW]"+format, v...)
}

func (d *defaultLog) Warn(v ...interface{}) {
	var out []interface{}
	out = append(out, "[W]")
	out = append(out, v...)
	log.Print(out...)
}

func (d *defaultLog) Warnf(format string, v ...interface{}) {
	log.Printf("[W]"+format, v...)
}

func (d *defaultLog) Warnw(format string, v ...interface{}) {
	log.Printf("[WW]"+format, v...)
}

func InitDefaultLog() *defaultLog {
	return &defaultLog{}
}

// Level implement core.ILogger
func (d *defaultLog) Level() core.LogLevel {
	return 0
}

// SetLevel implement core.ILogger
func (d *defaultLog) SetLevel(l core.LogLevel) {
	return
}

// ShowSQL implement core.ILogger
func (d *defaultLog) ShowSQL(show ...bool) {
}

// IsShowSQL implement core.ILogger
func (d *defaultLog) IsShowSQL() bool {
	return false
}

func (d *defaultLog) ADDSugar(key string, val interface{}) {
}

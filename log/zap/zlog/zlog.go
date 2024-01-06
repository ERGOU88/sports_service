// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zlog

import (
	"fmt"
	"github.com/go-xorm/core"
	"go.uber.org/zap"
	"sports_service/log/smartLog"
)

type SimpleZapLogger struct {
	showSQL bool
	Log     *zap.SugaredLogger
	level   core.LogLevel
}

func NewSimpleZapLogger(zapLog *zap.SugaredLogger) *SimpleZapLogger {
	return &SimpleZapLogger{Log: zapLog}
}

// Error implement core.ILogger
func (s *SimpleZapLogger) Error(v ...interface{}) {
	smartLog.AddLog(3, fmt.Sprint(v...), 1)
	s.Log.Error(v...)
	return
}

// Errorf implement core.ILogger
func (s *SimpleZapLogger) Errorf(format string, v ...interface{}) {
	smartLog.AddLog(3, fmt.Sprintf(format, v...), 1)
	s.Log.Errorf(format, v...)
	return
}

// Errorf implement core.ILogger
func (s *SimpleZapLogger) Errorw(format string, v ...interface{}) {
	smartLog.AddLog(3, fmt.Sprintf(format, v...), 1)
	s.Log.Errorw(format, v...)
	return
}

// Debug implement core.ILogger
func (s *SimpleZapLogger) Debug(v ...interface{}) {
	smartLog.AddLog(4, fmt.Sprint(v...), 1)
	s.Log.Debug(v...)
	return
}

// Debugf implement core.ILogger
func (s *SimpleZapLogger) Debugf(format string, v ...interface{}) {
	smartLog.AddLog(4, fmt.Sprintf(format, v...), 1)
	s.Log.Debugf(format, v...)
	return
}

// Debugf implement core.ILogger
func (s *SimpleZapLogger) Debugw(format string, v ...interface{}) {
	smartLog.AddLog(4, fmt.Sprintf(format, v...), 1)
	s.Log.Debugw(format, v...)
	return
}

// Info implement core.ILogger
func (s *SimpleZapLogger) Info(v ...interface{}) {
	smartLog.AddLog(1, fmt.Sprint(v...), 1)
	s.Log.Info(v...)
	return
}

// Infof implement core.ILogger
func (s *SimpleZapLogger) Infof(format string, v ...interface{}) {
	smartLog.AddLog(1, fmt.Sprintf(format, v...), 1)
	s.Log.Infof(format, v...)
	return
}

// Infof implement core.ILogger
func (s *SimpleZapLogger) Infow(format string, v ...interface{}) {
	smartLog.AddLog(1, fmt.Sprintf(format, v...), 1)
	s.Log.Infow(format, v...)
	return
}

// Warn implement core.ILogger
func (s *SimpleZapLogger) Warn(v ...interface{}) {
	smartLog.AddLog(2, fmt.Sprint(v...), 1)
	s.Log.Warn(v...)
	return
}

// Warnf implement core.ILogger
func (s *SimpleZapLogger) Warnf(format string, v ...interface{}) {
	smartLog.AddLog(2, fmt.Sprintf(format, v...), 1)
	s.Log.Warnf(format, v...)
	return
}

// Warnf implement core.ILogger
func (s *SimpleZapLogger) Warnw(format string, v ...interface{}) {
	smartLog.AddLog(2, fmt.Sprintf(format, v...), 1)
	s.Log.Warnf(format, v...)
	return
}

// Level implement core.ILogger
func (s *SimpleZapLogger) Level() core.LogLevel {
	return core.LOG_DEBUG
}

// SetLevel implement core.ILogger
func (s *SimpleZapLogger) SetLevel(l core.LogLevel) {
	s.level = l
	return
}

// ShowSQL implement core.ILogger
func (s *SimpleZapLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		s.showSQL = true
		return
	}
	s.showSQL = show[0]
}

// IsShowSQL implement core.ILogger
func (s *SimpleZapLogger) IsShowSQL() bool {
	return s.showSQL
}

// IsShowSQL implement core.ILogger
func (s *SimpleZapLogger) ADDSugar(key string, val interface{}) {
}

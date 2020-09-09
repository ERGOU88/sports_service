package zap

import (
	"bytes"
	"sports_service/server/log/zap/zlog"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// InitZapLog 初始化日志
func InitZapLog(logPath string, showColor bool, level int) *zlog.SimpleZapLogger {
	var encoding string
	var encodetime zapcore.TimeEncoder
	if showColor {
		encoding = "console"
		encodetime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]-" + strconv.FormatUint(getGID(), 10))
		}
	} else {
		encoding = "json"
		encodetime = zapcore.ISO8601TimeEncoder
	}
	encoderCfg := zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "app_id",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     encodetime,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	var lev zapcore.Level
	switch level {
	case 0:
		lev = zap.DebugLevel
	case 1:
		lev = zap.InfoLevel
	case 2:
		lev = zap.WarnLevel
	case 3:
		lev = zap.ErrorLevel
	}
	currLevel := zap.NewAtomicLevelAt(lev)

	customCfg := zap.Config{
		Level:             currLevel,
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          encoding,
		EncoderConfig:     encoderCfg,
		OutputPaths:       []string{logPath},
		ErrorOutputPaths:  []string{"stderr"},
	}

	logger, err := customCfg.Build(zap.AddCallerSkip(1), zap.AddStacktrace(zap.WarnLevel))
	if err != nil {
		panic(fmt.Errorf("init log error %v", err))
	}
	l := logger.Named("job")
	defer l.Sync()
	return zlog.NewSimpleZapLogger(l.Sugar())
}

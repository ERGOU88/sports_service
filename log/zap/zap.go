package zap

import (
	"bytes"
	"fmt"
	"github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"runtime"
	"sports_service/log/zap/zlog"
	"strconv"
	"time"
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
		CallerKey:      "file",
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
	//currLevel := zap.NewAtomicLevelAt(lev)
	//
	//customCfg := zap.Config{
	//  Level:             currLevel,
	//  Development:       true,
	//  DisableCaller:     false,
	//  DisableStacktrace: false,
	//  Encoding:          encoding,
	//  EncoderConfig:     encoderCfg,
	//  OutputPaths:       []string{fmt.Sprintf("%s%s", logPath, ".%Y%m%d%H%M%S")},
	//  ErrorOutputPaths:  []string{"stderr"},
	//}

	var core zapcore.Core
	ioWriter := getWriter(logPath)
	if encoding == "json" {
		core = zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.AddSync(ioWriter), lev)
	} else {
		core = zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), zapcore.AddSync(ioWriter), lev)
	}

	logger := zap.New(core, zap.AddCallerSkip(1), zap.AddStacktrace(zap.WarnLevel))
	//logger, err := customCfg.Build(zap.AddCallerSkip(1), zap.AddStacktrace(zap.WarnLevel))
	//if err != nil {
	//  panic(fmt.Errorf("init log error %v", err))
	//}

	l := logger.Named("fpv")
	defer l.Sync()
	return zlog.NewSimpleZapLogger(l.Sugar())
}

// 日志文件切割
func getWriter(filename string) io.Writer {
	// 保存30天内的日志，每24小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		fmt.Sprintf(filename, "%Y%m%d"),
		//rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*90),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}

	return hook
}

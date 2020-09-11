package engineLog

import (
	"fmt"
	"io/ioutil"
	"sports_service/server/log"
	"time"
	"github.com/gin-gonic/gin"
	"sports_service/server/log/smartLog"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// 日志中间件
func EngineLog(log log.ILogger, showColor bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.RequestURI
		log.Debugf("url:%s,user_agent:%s,ip:%s\n", url, c.Request.UserAgent(), c.ClientIP())

		start := time.Now()
		path := c.Request.URL.Path
		smartLog.AddURL(path)

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		version := c.Request.Header.Get("Version")
		userAgent := c.Request.Header.Get("User-Agent")
		statusCode := c.Writer.Status()
		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)
		var args string

		if c.Request.Header.Get("Content-Type") == "application/json" {
			b, _ := ioutil.ReadAll(c.Request.Body)
			args = string(b)
		} else {
			c.Request.ParseForm()
			args = c.Request.Form.Encode()
		}

		smartLog.AddExStr(fmt.Sprintf("| version:%s | User-Agent:%s | IP %s %s", version, userAgent, clientIP, latency))
		smartLog.AddArgs("ARGS", args)
		smartLog.AddExStr(fmt.Sprintf("%s %s %s", methodColor, method, reset))
		smartLog.AddExStr(fmt.Sprintf("%s %03d %s", statusColor, statusCode, reset))
		smartLog.End(statusCode)

		if showColor {
			log.Debugf("[PlayMate] %v | version:%s | User-Agent:%s |%s %3d %s|耗时： %13v | %s |%s  %s %-7s %s %s|\n",
				end.Format("2006/01/02 - 15:04:05"),
				version, userAgent,
				statusColor, statusCode, reset,
				latency,
				clientIP,
				methodColor, reset, method,
				path,
				args,
			)

		} else {
			log.Debugw(fmt.Sprintf("[PlayMate] %v | User-Agent:%s |", end.Format("2006/01/02 - 15:04:05"), userAgent),
				"status", fmt.Sprintf("%3d", statusCode),
				"time_consume", fmt.Sprintf("%13v", latency),
				"client", clientIP,
				"version", version,
				"method", fmt.Sprintf("%s", method),
				"url", fmt.Sprintf("%s %s", path, args),
			)
		}
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}

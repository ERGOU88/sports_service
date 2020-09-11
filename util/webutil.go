package util

import (
	"strings"
)

const (
	IE9     = "MSIE 9.0"
	IE8     = "MSIE 8.0"
	IE7     = "MSIE 7.0"
	IE6     = "MSIE 6.0"
	MAXTHON = "Maxthon"
	QQ      = "QQBrowser"
	GREEN   = "GreenBrowser"
	SE360   = "360SE"
	FIREFOX = "Firefox"
	OPERA   = "Opera"
	CHROME  = "Chrome"
	SAFARI  = "Safari"
	OTHER   = "其它"
	IPHONE  = "iPhone"
	ANDROID = "Android"
	IPad    = "iPad"
	Ios     = "iOS"
)

func GetClient(userAgent string) string {
	if strings.Contains(userAgent, IPHONE) || strings.Contains(userAgent, IPad) || strings.Contains(userAgent, Ios) {
		return IPHONE
	}

	if strings.Contains(userAgent, ANDROID) {
		return ANDROID
	}

	if strings.Contains(userAgent, OPERA) {
		return OPERA
	}

	if strings.Contains(userAgent, CHROME) {
		return CHROME
	}

	if strings.Contains(userAgent, FIREFOX) {
		return FIREFOX
	}

	if strings.Contains(userAgent, SAFARI) {
		return SAFARI
	}

	if strings.Contains(userAgent, SE360) {
		return SE360
	}

	if strings.Contains(userAgent, GREEN) {
		return GREEN
	}

	if strings.Contains(userAgent, QQ) {
		return QQ
	}

	if strings.Contains(userAgent, MAXTHON) {
		return MAXTHON
	}

	if strings.Contains(userAgent, IE9) {
		return IE9
	}

	if strings.Contains(userAgent, IE8) {
		return IE8
	}

	if strings.Contains(userAgent, IE7) {
		return IE7
	}

	if strings.Contains(userAgent, IE6) {
		return IE6
	}

	return OTHER
}

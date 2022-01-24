package sign

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"strings"
	"sports_service/server/app/config"
)

var AppInfo = map[string]string{
	string(consts.APPLET_APP_ID): "PlvZrGmBKGuQPXVb",
	string(consts.IOS_APP_ID):    "RfhHecN9zsNcy19Y",
	string(consts.AND_APP_ID):    "InaukEwVLLpcewX6",
}

// 检查签名
func CheckSign() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(consts.CHANNEL, 0)
		// debug模式下 不校验签名
		if config.Global.Debug == true {
			c.Next()
			return
		}

		reply := errdef.New(c)
		appId := c.GetHeader("AppId")
		sign := c.GetHeader("Sign")
		secret := c.GetHeader("Secret")
		timestamp := c.GetHeader("Timestamp")
		version := c.GetHeader("Version")
		path := c.Request.URL.Path
		str := fmt.Sprintf("%s&AppId=%s&Timestamp=%s&Version=%s", path, appId, timestamp, version)

		if !strings.Contains(path, "/api/v1/client/init") && secret == "" {
			log.Log.Errorf("sign_trace: secret not exists, secret:%s", secret)
			reply.Response(http.StatusUnauthorized, errdef.UNAUTHORIZED)
			c.Abort()
			return
		}

		log.Log.Errorf("sign_trace: path:%s, match: %d", path, strings.Compare(path, "/api/v1/client/init"))
		if !strings.Contains(path, "/api/v1/client/init")  {
			log.Log.Infof("sign_trace: add secret, secret:%s", secret)
			str = fmt.Sprintf("%s&Secret=%s", str, secret)
		}

		if appId == "" || sign == "" || timestamp == "" || version == "" {
			log.Log.Errorf("sign_trace: param error,appId:%s, sign:%s, timestamp:%s, version:%s", appId, sign, timestamp, version)
			reply.Response(http.StatusUnauthorized, errdef.UNAUTHORIZED)
			c.Abort()
			return
		}

		if strings.Compare(appId, string(consts.IOS_APP_ID)) != 0 &&
			strings.Compare(appId, string(consts.AND_APP_ID)) != 0 &&
			strings.Compare(appId, string(consts.APPLET_APP_ID)) != 0 {
			log.Log.Errorf("sign_trace: appId not match, appId:%s", appId)
			reply.Response(http.StatusUnauthorized, errdef.UNAUTHORIZED)
			c.Abort()
			return
		}

		// 校验签名是否一致
		if !verifySign(str, appId, sign) {
			reply.Response(http.StatusUnauthorized, errdef.UNAUTHORIZED)
			c.Abort()
			return
		}

		plt := getChannel(appId)
		c.Set(consts.CHANNEL, plt)

		c.Next()
	}
}

// 校验签名是否一致
func verifySign(str, appId, sign string) bool {
	appKey := getAppKey(appId)
	str += fmt.Sprintf("&%s", appKey)

	log.Log.Debugf("sign_trace: str:%s", str)
	data := []byte(str)
	has := md5.Sum(data)
	md5Str := fmt.Sprintf("%x", has)
	log.Log.Debugf("client sign:%v, md5Str:%v", sign, md5Str)
	if md5Str == sign {
		return true
	}

	log.Log.Errorf("sign_trace: sign not match, client sign:%s, real sign:%s", sign, md5Str)
	return false
}

// 获取appKey
func getAppKey(appId string) string {
	appKey, ok := AppInfo[appId]
	if ok {
		return appKey
	}

	return ""
}


func getChannel(appId string) int {
	switch appId {
	case string(consts.IOS_APP_ID):
		return consts.PLT_TYPE_IOS
	case string(consts.AND_APP_ID):
		return consts.PLT_TYPE_ANDROID
	case string(consts.APPLET_APP_ID):
		return consts.PLT_TYPE_APPLET
	}

	return 0
}



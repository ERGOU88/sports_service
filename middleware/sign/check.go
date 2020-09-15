package sign

import (
	"crypto/md5"
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/consts"
	"strings"
	"sports_service/server/global/app/log"
)

var AppInfo = map[string]string{
	string(consts.WEB_APP_ID):       "PlvZrGmBKGuQPXVb",
	string(consts.IOS_APP_ID):       "RfhHecN9zsNcy19Y",
	string(consts.AND_APP_ID):	     "InaukEwVLLpcewX6",
}

// 检查签名
func CheckSign() gin.HandlerFunc {
	return func(c *gin.Context) {
		reply := errdef.New(c)
		appId := c.GetHeader("AppId")
		sign := c.GetHeader("Sign")
		secret := c.GetHeader("Secret")
		timestamp := c.GetHeader("Timestamp")
		uri := c.Request.RequestURI
		str := fmt.Sprintf("%s&AppId=%s&Timestamp=%s", uri, appId, timestamp)

		if !strings.Contains(uri, "/api/v1/client/init") && secret == "" {
			reply.Response(http.StatusUnauthorized, errdef.UNAUTHORIZED)
			c.Abort()
			return
		}

		if strings.Compare(uri, "/api/v1/client/init") == - 1 {
			str = fmt.Sprintf("%s&Secret=%s", str, secret)
		}

		if appId == "" || sign == "" || timestamp == "" {
			reply.Response(http.StatusUnauthorized, errdef.UNAUTHORIZED)
			c.Abort()
			return
		}

		if strings.Compare(appId, string(consts.IOS_APP_ID)) == -1 &&
			strings.Compare(appId, string(consts.AND_APP_ID)) == -1 &&
			strings.Compare(appId, string(consts.WEB_APP_ID)) == -1 {
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

		c.Next()
	}
}

// 校验签名是否一致
func verifySign(str, appId, sign string) bool {
	appKey := getAppKey(appId)
	str = fmt.Sprintf("&%s", appKey)
	data := []byte(str)
	has := md5.Sum(data)
	md5Str := fmt.Sprintf("%x", has)
	log.Log.Debugf("client sign:%v, md5Str:%v", sign, md5Str)
	if md5Str == sign {
		return true
	}

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



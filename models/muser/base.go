package muser

import (
	"sports_service/server/dao"
	"sports_service/server/global/login/log"
	"sports_service/server/global/consts"
	"sports_service/server/util"
	"github.com/gin-gonic/gin"
	"sports_service/server/global/rdskey"
	"fmt"
)

type base struct {

}

// 获取版本号
func (m *base) getVersion(c *gin.Context) string {
	return c.GetHeader("Version")
}

// 获取设备类型
func (m *base) getDeviceType(c *gin.Context) int {
	switch util.GetClient(c.Request.UserAgent()) {
	case util.IPHONE, util.IPad, util.Ios:
		return int(consts.IOS_PLATFORM)
	case util.ANDROID:
		return int(consts.ANDROID_PLATFORM)
	default:
		return int(consts.ANDROID_PLATFORM)
	}
}

// 默认头像
func (m *base) defaultAvatar() string {
	return ""
}

// 默认昵称
func (m *base) getNickName(nickName string) string {
	if nickName == "" {
		nickName = "FPV用户"
	}

	rds := dao.NewRedisDao()
	nickNameNum, err := rds.INCR(rdskey.USER_NICKNAME_INCR)
	if err != nil {
		log.Log.Errorf("social_trace: user nickname incr err:%s", err)
	}

	if nickNameNum < 10 {
		return fmt.Sprintf("%s0%d", nickName, nickNameNum)
	}

	return fmt.Sprintf("%s%d", nickName, nickNameNum)
}

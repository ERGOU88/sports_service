package token

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/dao"
	"sports_service/global/app/errdef"
	"sports_service/global/app/log"
	"sports_service/global/consts"
	"sports_service/models/muser"
	"strings"
	"time"
)

// token校验
func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		reply := errdef.New(c)
		var userid, hashcode, auth string
		c.Set(consts.USER_ID, userid)
		val, err := c.Request.Cookie(consts.COOKIE_NAME)
		if err != nil {
			auth = c.Request.Header.Get("auth")
			if auth == "" {
				log.Log.Errorf("c.Request.Cookie() err is %s", err.Error())
				reply.Response(http.StatusUnauthorized, errdef.INVALID_TOKEN)
				c.Abort()
				return
			}
		} else {
			auth = val.Value
		}

		log.Log.Debugf("auth:%v", auth)
		//v := c.Request.Header.Get("auth")
		ks := strings.Split(auth, "_")
		if len(ks) != 2 {
			log.Log.Errorf("len(ks) != 2")
			ks = strings.Split(auth, "%09")
		}

		if len(ks) == 2 {
			userid = ks[0]
			hashcode = ks[1]
		}

		if len(hashcode) <= 0 {
			log.Log.Errorf("len(hashcode) <= 0")
			reply.Response(http.StatusUnauthorized, errdef.INVALID_TOKEN)
			c.Abort()
			return
		}

		log.Log.Debugf("token_trace: token userid:%s", userid)
		model := new(muser.UserModel)
		token, err := model.GetUserToken(userid)
		if err != nil && err == redis.ErrNil {
			log.Log.Errorf("token_trace: get user token by redis err:%s", err)
			reply.Response(http.StatusUnauthorized, errdef.INVALID_TOKEN)
			c.Abort()
			return
		}

		// 客户端token是否和redis存储的一致
		if res := strings.Compare(auth, token); res != 0 {
			log.Log.Errorf("token_trace: token not match, server token:%s, client token:%s", token, auth)
			reply.Response(http.StatusUnauthorized, errdef.INVALID_TOKEN)
			c.Abort()
			return
		}

		log.Log.Debugf("client token:%s, server token:%s", auth, token)

		if userid != "" {
			// 给token续约
			if err := model.SaveUserToken(userid, auth); err != nil {
				log.Log.Errorf("token_trace: save user token err:%s", err)
			}

			// todo:
			go RecordInfo(userid, c.Request.URL.Path)
		}

		c.Set(consts.USER_ID, userid)
		c.Next()
	}
}

func RecordInfo(userid, path string) {
	activityType := GetUserActivityType(path)
	session := dao.AppEngine.NewSession()
	defer session.Close()
	umodel := muser.NewUserModel(session)
	condition := fmt.Sprintf("user_id='%s'", userid)
	cols := "last_login_time"
	now := int(time.Now().Unix())
	umodel.User.LastLoginTime = now
	if _, err := umodel.UpdateUserInfos(condition, cols); err != nil {
		log.Log.Errorf("token_trace: update login time fail, userId:%s, err:%s", userid, err)
	}

	if _, err := umodel.AddActivityRecord(userid, now, activityType); err != nil {
		log.Log.Errorf("token_trace: record activity user fail, userId:%s, err:%s", userid, err)
	}

	return
}

// 获取用户活跃类型
func GetUserActivityType(path string) int {
	activityType := 0
	switch {
	case strings.Contains(path, "/api/v1/like/video"):
		activityType = consts.ACTIVITY_TYPE_LIKE_VIDEO
	case strings.Contains(path, "/api/v1/like/comment"):
		activityType = consts.ACTIVITY_TYPE_LIKE_COMMENT
	case strings.Contains(path, "/api/v1/like/information"):
		activityType = consts.ACTIVITY_TYPE_LIKE_INFORMATION
	case strings.Contains(path, "/api/v1/like/post"):
		activityType = consts.ACTIVITY_TYPE_LIKE_POST
	case strings.Contains(path, "/api/v1/collect/video"):
		activityType = consts.ACTIVITY_TYPE_LIKE_COMMENT
	case strings.Contains(path, "/api/v1/comment/publish/v2"):
		activityType = consts.ACTIVITY_TYPE_COMMENT
	case strings.Contains(path, "/api/v1/collect/video"):
		activityType = consts.ACTIVITY_TYPE_COLLECT_VIDEO
	case strings.Contains(path, "/api/v1/comment/reply"):
		activityType = consts.ACTIVITY_TYPE_REPLY
	case strings.Contains(path, "/api/v1/barrage/send"):
		activityType = consts.ACTIVITY_TYPE_BARRAGE
	case strings.Contains(path, "/api/v1/video/publish"):
		activityType = consts.ACTIVITY_TYPE_PUB_VIDEO
	case strings.Contains(path, "/api/v1/post/publish"):
		activityType = consts.ACTIVITY_TYPE_PUB_POST
	case strings.Contains(path, "/api/v1/share/social"):
		activityType = consts.ACTIVITY_TYPE_SHARE_SOCIAL
	case strings.Contains(path, "/api/v1/share/community"):
		activityType = consts.ACTIVITY_TYPE_SHARE_COMMUNITY
	}

	return activityType
}

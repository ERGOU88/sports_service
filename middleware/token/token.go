package token

import (
	"github.com/gin-gonic/gin"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/models/muser"
	"strings"
	"sports_service/server/global/app/log"
)

// token校验
func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		reply := errdef.New(c)
		var userid string
		var hashcode string
		c.Set(consts.USER_ID, userid)
		val, err := c.Request.Cookie(consts.COOKIE_NAME)
		if err != nil {
			log.Log.Errorf("c.Request.Cookie() err is %s", err.Error())
			reply.Response(http.StatusUnauthorized, errdef.INVALID_TOKEN)
			c.Abort()
			return
		}

		log.Log.Debugf("val:%v", val)
		v := val.Value
		//v := c.Request.Header.Get("auth")
		ks := strings.Split(v, "_")
		if len(ks) != 2 {
			log.Log.Errorf("len(ks) != 2")
			ks = strings.Split(v, "%09")
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
		if res := strings.Compare(v, token); res != 0 {
			log.Log.Errorf("token_trace: token not match, server token:%s, client token:%s", token, v)
			reply.Response(http.StatusUnauthorized, errdef.INVALID_TOKEN)
			c.Abort()
			return
		}

		log.Log.Debugf("client token:%s, server token:%s", v, token)

		if userid != "" {
			// 给token续约
			model.SaveUserToken(userid, v)
		}

		c.Set(consts.USER_ID, userid)
		c.Next()
	}
}

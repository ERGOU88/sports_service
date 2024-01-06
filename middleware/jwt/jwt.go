package jwt

import (
	"fmt"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/global/backend/errdef"
	"sports_service/global/backend/log"
	"strings"
)

const JWT_SECRET = "!@#WE$%!SEgfcgHT#456"

// 验证token的中间件
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.Contains(path, "/backend/v1/admin/ad/login") || strings.Contains(path, "/backend/v1/admin/login") {
			c.Next()
			return
		}

		reply := errdef.New(c)
		data, err := request.ParseFromRequest(c.Request, request.OAuth2Extractor, func(token *jwt_lib.Token) (interface{}, error) {
			b := []byte(JWT_SECRET)
			return b, nil
		})

		if data != nil {
			log.Log.Debugf("jwt_data:%+v", data)
			c.Set("jwt_data", data)
		}

		// 如果jwt校验失败
		if err != nil {
			log.Log.Errorf("jwt_trace: auth fail, err:%s", err)
			reply.Response(http.StatusOK, errdef.UNAUTHORIZED)
			c.Abort()
			return
		}

		c.Next()
	}
}

// 获取jwt过期时间
func GetExpireTime(c *gin.Context) float64 {
	tokenMap, err := Decode(c)
	if err != nil {
		return 0
	}

	v, ok := tokenMap["exp"]
	if !ok {
		return 0
	}

	expire := v.(float64)
	return expire
}

// 获取管理员信息
func GetUserInfo(c *gin.Context, key string) string {
	tokenMap, err := Decode(c)
	if err != nil {
		return ""
	}

	v, ok := tokenMap[key]
	if !ok {
		return ""
	}

	return v.(string)
}

// 解析token的内容
func Decode(c *gin.Context) (map[string]interface{}, error) {
	di, exist := c.Get("jwt_data")
	if !exist {
		return nil, fmt.Errorf("no jwt_data")
	}

	data := di.(*jwt_lib.Token)
	mapData := data.Claims.(jwt_lib.MapClaims)

	return mapData, nil
}

type JwtInfo struct {
	Key string      `json:"key"`
	Val interface{} `json:"val"`
}

// 生成jwt
// username: 用户名称
func GenerateJwt(c *gin.Context, info []JwtInfo) (string, error) {
	token, err := Put(c, JWT_SECRET, info...)
	if err != nil {
		return "", err
	}

	return token, nil
}

func Put(c *gin.Context, secret string, kvs ...JwtInfo) (string, error) {
	di, exist := c.Get("jwt_data")
	var mapClaims jwt_lib.MapClaims
	var token *jwt_lib.Token
	if !exist {
		token = jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
		var ok bool
		mapClaims, ok = token.Claims.(jwt_lib.MapClaims)
		if !ok {
			mapClaims = jwt_lib.MapClaims{}
		}
	} else {
		token = di.(*jwt_lib.Token)
		mapClaims = token.Claims.(jwt_lib.MapClaims)
	}

	for _, kv := range kvs {
		mapClaims[kv.Key] = kv.Val
	}

	token.Claims = mapClaims
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign token err:%s", err)
	}

	//c.Header("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	c.Header("Authorization", tokenString)
	return tokenString, nil
}

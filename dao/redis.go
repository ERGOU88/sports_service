package dao

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	"time"
)

var redisPool *redis.Pool

func RedisPool() *redis.Pool {
	return redisPool
}

func InitRedis(host, pwd string) {
	redisPool = GetRedisPool(host, pwd, 500)
}

func GetRedisPool(server, password string, maxConn int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     100,
		MaxActive:   maxConn,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}

			if password != "" {
        if _, err := c.Do("AUTH", password); err != nil {
          c.Close()
          return nil, err
        }
      }

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

type RedisDao struct {
}

func NewRedisDao() *RedisDao {
	return &RedisDao{}
}

// RedisConn returns redis conn.
func RedisConn() redis.Conn {
	return redisPool.Get()
}

//基础redis操作
//集合操作
//SADD 可以添加多个 返回成功数量
func (this *RedisDao) SADD(key string, value interface{}) (num int, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int(conn.Do("SADD", key, value))
	if err != nil {
		return
	}
	return
}

// Set 总是成功的
func (this *RedisDao) Set(key string, value interface{}) (err error) {
	conn := redisPool.Get()
	defer conn.Close()
	_, err = conn.Do("SET", key, value)
	if err != nil {
		return
	}
	return

}

// Del 可以删除多个key 返回删除key的num和错误
func (this *RedisDao) Del(key ...interface{}) (num int, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int(conn.Do("DEL", key...))
	if err != nil {
		return
	}
	return
}

//Get
func (this *RedisDao) Get(key string) (s string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	s, err = redis.String(conn.Do("GET", key))
	if err != nil {
		return
	}
	return
}

//Get
func (this *RedisDao) GetInt(key string) (n int, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	n, err = redis.Int(conn.Do("GET", key))
	if err != nil {
		return
	}
	return
}

func (this *RedisDao) GetInt64(key string) (n int64, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	n, err = redis.Int64(conn.Do("GET", key))
	if err != nil {
		return
	}
	return
}

//EXIST
func (this *RedisDao) EXISTS(key string) (ok bool, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	ok, err = redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return
	}
	return
}

//EXIST
func (this *RedisDao) EXISTSGet(key string) (ok bool, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	ok, err = redis.Bool(conn.Do("EXISTS GET", key))
	if err != nil {
		return
	}
	return
}

//SCARD
func (this *RedisDao) SCARD(key string) (num int, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int(conn.Do("SCARD", key))
	if err != nil {
		return
	}
	return
}

//SPOP 弹出被移除的元素, 当key不存在的时候返回 nil
func (this *RedisDao) SPOP(key string) (out string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	out, err = redis.String(conn.Do("SPOP", key))
	if err != nil {
		return
	}
	return
}

//SREM
func (this *RedisDao) SREM(key string, value interface{}) (num int, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int(conn.Do("SREM", key, value))
	if err != nil {
		return
	}
	return
}

//SISMEMBER
func (this *RedisDao) SISMEMBER(key string, value interface{}) (ok bool, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	ok, err = redis.Bool(conn.Do("SISMEMBER", key, value))
	if err != nil {
		return
	}
	return
}

//SMEMBERS
func (this *RedisDao) SMEMBERS(key string) (reply []string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	reply, err = redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		return
	}
	return
}

//List操作

//LLEN
func (this *RedisDao) LLEN(key string) (out int, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	out, err = redis.Int(conn.Do("LLEN", key))
	if err != nil {
		return
	}
	return
}

//LPOP
func (this *RedisDao) LPOP(key string) (out string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	out, err = redis.String(conn.Do("LPOP", key))
	if err != nil {
		return
	}
	return
}

//LPUSH 整型回复: 在 push 操作后的 list 长度。
func (this *RedisDao) LPUSH(key string, value interface{}) (num int, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int(conn.Do("LPUSH", key, value))
	if err != nil {
		return
	}
	return
}

//RPUSH 整型回复: 在 push 操作后的 list 长度。
func (this *RedisDao) RPUSH(key string, value interface{}) (num int, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int(conn.Do("RPUSH", key, value))
	if err != nil {
		return
	}
	return
}

//BRPOP
func (this *RedisDao) BRPOP(key string) (out []string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	out, err = redis.Strings(conn.Do("BRPOP", key, 0))
	if err != nil {
		return
	}
	return
}

//LINDEX 当 key 位置的值不是一个列表的时候，会返回一个error
func (this *RedisDao) LINDEX(key string, index int) (out string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	out, err = redis.String(conn.Do("LINDEX", key, index))
	if err != nil {
		return
	}
	return
}

//LTRIM
func (this *RedisDao) LTRM(key string, start, end int) (out string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	out, err = redis.String(conn.Do("LTRIM", key, start, end))
	if err != nil {
		return
	}
	return
}

// LREM
func (this *RedisDao) LREM(key string, value interface{}) (num int, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int(conn.Do("LREM", key, 0, value))
	if err != nil {
		return
	}
	return
}

//LRANGE

//哈希操作
//HDEL
//HEXISTS
func (this *RedisDao) HEXISTS(key, field string) (ok bool, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	ok, err = redis.Bool(conn.Do("HEXISTS", key, field))
	if err != nil {
		return
	}
	return
}

//HGET 该字段所关联的值。当字段不存在或者 key 不存在时返回nil。
func (this *RedisDao) HGET(key, field string) (out string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	out, err = redis.String(conn.Do("HGET", key, field))
	if err != nil {
		return
	}
	return
}

//HGET 该字段所关联的值。当字段不存在或者 key 不存在时返回nil。
func (this *RedisDao) HGetAll(key string) (out map[string]string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	out, err = redis.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		return
	}
	return
}

//HMGETMAP
func (this *RedisDao) HMGETMAP(key string) (m map[string]string, err error) {
	m = make(map[string]string)
	conn := redisPool.Get()
	defer conn.Close()
	m, err = redis.StringMap(conn.Do("HGETALL", key))
	return
}

//HINCRBY 增值操作执行后的该字段的值。
func (this *RedisDao) HINCRBY(key, field string, in int) (num int, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int(conn.Do("HINCRBY", key, field, in))
	if err != nil {
		return
	}
	return
}

//HMGETSTRUCT
func (this *RedisDao) HMGETSTRUCT(key, value interface{}) (err error) {
	conn := redisPool.Get()
	defer conn.Close()
	v, err := redis.Values(conn.Do("HGETALL", key))
	if err != nil {
		return err
	}

	if len(v) == 0 {
		return redis.ErrNil
	}
	err = redis.ScanStruct(v, value)
	return err
}

//HKEYS
func (this *RedisDao) HKEYS(key string) (v []string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	v, err = redis.Strings(conn.Do("HKEYS", key))
	if err != nil {
		return
	}
	return
}

//HMSET
func (this *RedisDao) HMSET(key, value interface{}) (ok string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	ok, err = redis.String(conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(value)...))
	if err != nil {
		return
	}
	return
}

//HMSET
func (this *RedisDao) HMSETArray(key interface{}, value []interface{}) (ok string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	var args []interface{}
	args = append(args, key)
	args = append(args, value...)
	ok, err = redis.String(conn.Do("HMSET", args...))
	// model.Log.Debugf("hmset ")
	if err != nil {
		return
	}
	return
}

//HSCAN
//HSET 1如果field是一个新的字段  0如果field原来在map里面已经存在
func (this *RedisDao) HDEL(key, field string) (err error) {
	conn := redisPool.Get()
	defer conn.Close()
	_, err = redis.Int(conn.Do("HDEL", key, field))
	if err != nil {
		return
	}
	return
}

//HDEL 1如果field是一个新的字段  0如果field原来在map里面已经存在
func (this *RedisDao) HSET(key, field string, value interface{}) (ok bool, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	ok, err = redis.Bool(conn.Do("HSET", key, field, value))
	if err != nil {
		return
	}
	return
}

//HDEL 1如果field是一个新的字段  0如果field原来在map里面已经存在
func (this *RedisDao) HSETCover(key, field string, value interface{}) (result string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	result, err = redis.String(conn.Do("HMSET", key, field, value))
	if err != nil {
		return
	}
	return
}

//HLEN 哈希集中字段的数量，当 key 指定的哈希集不存在时返回 0
func (this *RedisDao) HLEN(key string) (num int64, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int64(conn.Do("HLEN", key))
	if err != nil {
		return
	}
	return
}

//ZREMRANGEBYRANK myzset 0 1  0 -200(保留200名)
func (this *RedisDao) ZREMRANGEBYRANK(key string, stop int) (num int, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int(conn.Do("ZREMRANGEBYRANK", key, 0, stop))
	if err != nil {
		return
	}
	return
}

//ZSCORE
func (this *RedisDao) ZSCORE(key string, member string) (num int, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int(conn.Do("ZSCORE", key, member))
	if err != nil {
		return
	}
	return
}

//ZIsMember 判断是否是有序集合的成员
func (this *RedisDao) ZISMEMBER(key interface{}, item interface{}) (isMember bool, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	_, err = redis.Float64(conn.Do("ZSCORE", key, item))
	switch err {
	case nil:
		return true, nil
	case redis.ErrNil:
		return false, nil
	default:
		return false, err
	}
}

//ZADD
func (this *RedisDao) ZADD(key string, sorce int, member string) (num int, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int(conn.Do("ZADD", key, sorce, member))
	if err != nil {
		return
	}
	return
}

//ZADD64
func (this *RedisDao) ZADD64(key string, sorce int64, member string) (num int64, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int64(conn.Do("ZADD", key, sorce, member))
	if err != nil {
		return
	}
	return
}

//ZREVRANGEBYSCORE 逆序份数  获取的 前N个数据
func (this *RedisDao) ZREVRANGEBYSCORE(key string, limit int) (list map[string]string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	list, err = redis.StringMap(conn.Do("ZREVRANGEBYSCORE", key, "+inf", "-inf", "WITHSCORES", "limit", 0, limit))
	if err != nil {
		return
	}
	return
}

//ZREVRANGEBYSCORE 逆序份数  获取的 前N个数据 从大到小
func (this *RedisDao) ZREVRANGEBYSCOREInterval(key string, max, min interface{}, limit int) (list map[string]string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	list, err = redis.StringMap(conn.Do("ZREVRANGEBYSCORE", key, min, max, "WITHSCORES", "limit", 0, limit))
	if err != nil {
		return
	}
	return
}

//ZREVRANGEBYSCORE 逆序份数  获取的 前N个数据 从小到大
func (this *RedisDao) ZRANGEBYSCOREInterval(key string, min, max int64, limit int) (list []string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	list, err = StringMap(conn.Do("ZRANGEBYSCORE", key, min, max, "WITHSCORES", "limit", 0, limit))
	if err != nil {
		return
	}
	return
}

// ZRANGE
func (this *RedisDao) ZRANGE(key string, begin, end int64) (list map[string]string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	list, err = redis.StringMap(conn.Do("ZRANGE", key, begin, end, "WITHSCORES"))
	if err != nil {
		return
	}
	return
}

// ZRANGEString ...
func (this *RedisDao) ZRANGEString(key string, begin, end int64) (list []string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	list, err = redis.Strings(conn.Do("ZRANGE", key, begin, end))
	if err != nil {
		return
	}
	return
}

// ZRANGEString ...
func (this *RedisDao) ZREVRANGEString(key string, begin, end int64) (list []string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	list, err = redis.Strings(conn.Do("ZREVRANGE", key, begin, end))
	if err != nil {
		return
	}
	return
}

// StringMap is a helper that converts an array of strings (alternating key, value)
// into a map[string]string. The HGETALL and CONFIG GET commands return replies in this format.
// Requires an even number of values in result.
func StringMap(result interface{}, err error) ([]string, error) {
	values, err := redis.Values(result, err)
	m := make([]string, 0)
	if err != nil {
		return m, err
	}
	if len(values)%2 != 0 {
		return m, errors.New("redigo: StringMap expects even number of values result")
	}
	for i := 0; i < len(values); i += 2 {
		key, okKey := values[i].([]byte)
		//value, okValue := values[i+1].([]byte)
		if !okKey {
			return m, errors.New("redigo: ScanMap key not a bulk string value")
		}
		m = append(m, string(key))
	}
	return m, nil
}

//ZREVRANGEBYSCORE 逆序份数  获取的 前N个数据 从小到大
func (this *RedisDao) ZRANGEBYSCOREInterval1(key string, min, max int64, limit int, list interface{}) (err error) {
	conn := redisPool.Get()
	defer conn.Close()
	var bList []byte
	bList, err = redis.Bytes(conn.Do("ZRANGEBYSCORE", key, min, max, "WITHSCORES", "limit", 0, limit))
	if err != nil {
		return
	}
	err = json.Unmarshal(bList, list)
	return
}

//ZREVRANGE

//ZREVRANGEBYSCORE 逆序份数  获取的 前N个数据 不要scores
func (this *RedisDao) GetSearchKeys(key string, limit int) (list []string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	list, err = redis.Strings(conn.Do("ZREVRANGEBYSCORE", key, "+inf", "-inf", "limit", 0, limit))
	if err != nil {
		return
	}
	return
}

//ZINCRBY +increment  如果没有key 插入
func (this *RedisDao) ZINCRBY(key string, increment int64, member string) (num int64, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int64(conn.Do("ZINCRBY", key, increment, member))
	if err != nil {
		return
	}
	return
}

//ZRANK 判断一个member 在key中的索引 如果不在 返回nil ,在 返回索引
func (this *RedisDao) ZRANK(key string, member string) (num int, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int(conn.Do("ZRANK", key, member))
	if err != nil {
		return
	}
	return
}

// ZCARD 用于计算集合中元素的数量
func (this *RedisDao) ZCARD(key string) (num int64, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int64(conn.Do("ZCARD", key))
	if err != nil {
		return
	}
	return
}

//EXPIRE 设置一个key 的过期时间 返回值int 1 如果设置了过期时间 0 如果没有设置过期时间，或者不能设置过期时间
func (this *RedisDao) EXPIRE(key string, expireTime int) (num int, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int(conn.Do("EXPIRE", key, expireTime))
	if err != nil {
		return
	}

	return
}

//EXPIRE 设置一个key 的过期时间 返回值int 1 如果设置了过期时间 0 如果没有设置过期时间，或者不能设置过期时间
func (this *RedisDao) EXPIRE64(key string, expireTime int64) (num int, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int(conn.Do("EXPIRE", key, expireTime))
	if err != nil {
		return
	}
	return
}

//SETEX key seconds value
func (this *RedisDao) SETEX(key string, seconds int64, value interface{}) (err error) {
	conn := redisPool.Get()
	defer conn.Close()
	_, err = conn.Do("SETEX", key, seconds, value)
	if err != nil {
		return
	}
	return
}

//SETEX key seconds value
func (this *RedisDao) SETNX(key string) (res bool, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	var num int
	num, err = redis.Int(conn.Do("SETNX", key, 1))
	if err != nil {
		return
	}
	//设置成功
	if num == 1 {
		return true, nil
	}
	return
}

func (this *RedisDao) INCR(key string) (num int64, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int64(conn.Do("INCR", key))
	if err != nil {
		return
	}
	return
}

func (this *RedisDao) INCRBY(key string, n int64) (num int64, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int64(conn.Do("INCRBY", key, n))
	if err != nil {
		return
	}
	return
}

func (this *RedisDao) DECR(key string) (num int64, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int64(conn.Do("DECR", key))
	if err != nil {
		return
	}
	return
}

//GEOADD
func (this *RedisDao) GEOADD(key, longitude, latitude string, member string) (err error) {
	conn := redisPool.Get()
	defer conn.Close()
	_, err = conn.Do("GEOADD", key, longitude, latitude, member)
	if err != nil {
		return
	}
	return
}

//GEORADIUS
func (this *RedisDao) GEORADIUS(key, longitude, latitude string, dist int, units string) (l []string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	redisDebug("GEORADIUS", key, longitude, latitude, dist, units, "asc")
	l, err = redis.Strings(conn.Do("GEORADIUS", key, longitude, latitude, dist, units, "asc"))
	if err == nil {
		return
	}
	return
}

func redisDebug(args ...interface{}) {
	dArgs := make([]interface{}, 0)
	for _, arg := range args {
		dArgs = append(dArgs, arg, " ")
	}
}

//GEORADIUS
func (this *RedisDao) GEORADIUSWITHCOORD(key, longitude, latitude string, dist float64, units string) (l []interface{}, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	l, err = redis.Values(conn.Do("GEORADIUS", key, longitude, latitude, dist, units, "WITHCOORD"))
	if err == nil {
		return
	}
	return
}

//GEORADIUS
func (this *RedisDao) GEORADIUSDIST(key, longitude, latitude string, dist int, units string) (l []interface{}, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	redisDebug("GEORADIUS", key, longitude, latitude, dist, units, "WITHDIST", "asc")
	l, err = redis.Values(conn.Do("GEORADIUS", key, longitude, latitude, dist, units, "WITHDIST", "asc"))
	if err == nil {
		return
	}
	return
}

// GEORADIUSBYMEMBER
// units m|km|ft|mi
func (this *RedisDao) GEORADIUSBYMEMBER(key, member string, dist int, units string) (l []string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	redisDebug("GEORADIUSBYMEMBER", key, member, dist, units, "WITHDIST", "asc")
	l, err = redis.Strings(conn.Do("GEORADIUSBYMEMBER", key, member, dist, units, "ASC"))
	if err == nil {
		return
	}
	return
}

// Set 总是成功的
func (this *RedisDao) ZREM(key, member string) (num int64, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err = redis.Int64(conn.Do("ZREM", key, member))
	if err != nil {
		return
	}
	return
}

func (this *RedisDao) SetKeyAfter10min(key string, time int64, value interface{}, callback func() error) error {
	ok, err := this.EXISTS(key)
	if err != nil {
		return err
	}

	if !ok {
		if err := this.SETEX(key, time, value); err != nil {
			return err
		}

		return callback()
	}

	return nil
}

// 从list中取得数据
func (this *RedisDao) LRange(key string, start, end int64) (datas []string, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	datas, err = redis.Strings(conn.Do("LRANGE", key, start, end))
	if err != nil {
		return
	}
	return
}

func (this *RedisDao) PUBLISH(channel, data string) (int64, error) {
	conn := redisPool.Get()
	defer conn.Close()
	num, err := redis.Int64(conn.Do("PUBLISH", channel, data))
	return num, err
}

// RPOP
func (this *RedisDao) RPOP(key string) (out int64, err error) {
	conn := redisPool.Get()
	defer conn.Close()
	out, err = redis.Int64(conn.Do("RPOP", key))
	if err != nil {
		return
	}
	return
}

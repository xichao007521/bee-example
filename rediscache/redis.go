package rediscache

import (
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
	"reflect"
	"strings"
	"sync"
	"time"
)

var RedisClient *redis.ClusterClient

func init() {
	BuildRedisClient()
}

var mu sync.Mutex
func BuildRedisClient()  {
	mu.Lock()
	defer mu.Unlock()
	if RedisClient != nil {
		return
	}
	conn := beego.AppConfig.String("redis.conn")

	RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: strings.Split(conn, ","),
		Password: "",
	})

	_, err := RedisClient.Ping().Result()

	if err != nil {
		beego.Error("Redis Connection error. %s", err)
		RedisClient = nil
	}
}

const (
	// 因为redis不能缓存空值，但是我们又会经常需要缓存空值防止频繁击穿cache
	// 因此用一个标识存储空的值
	EmptyFlag = "###+--**-+###"
)

type Options struct {
	// 非空
	Key          string
	// 非空
	Rt           reflect.Type
	Expires      time.Duration
	EmptyExpires time.Duration
}



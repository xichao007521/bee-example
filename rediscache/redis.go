package rediscache

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

// NOTICE!!!
var RedisClient *redis.ClusterClient

func init() {
	BuildRedisClient()
}

var mu sync.Mutex

func BuildRedisClient() {
	mu.Lock()
	defer mu.Unlock()
	if RedisClient != nil {
		return
	}
	conn := beego.AppConfig.String("redis.conn")

	RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    strings.Split(conn, ","),
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
	Key string
	// 非空, 如果是集合类型，表示集合里面元素的类型
	Rt           reflect.Type
	Expires      time.Duration
	EmptyExpires time.Duration
}

// 返回的三个参数，依次是: cache值，是否空，错误信息
func GetCacheValueItem(v interface{}) (string, bool, error) {
	jsonB, err := json.Marshal(v)
	if err != nil {
		return "", true, err
	}
	cacheV := string(jsonB)
	cLength := 0
	switch v.(type) {
	// string 类型单独处理
	case string:
		cLength = len(v.(string))
	default:
		cLength = len(cacheV)
	}
	return cacheV, cLength == 0, nil
}

func testSetup() {
	ap, _ := os.Getwd()
	beego.TestBeegoInit(ap + "/..")
	BuildRedisClient()
}

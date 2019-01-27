package cache

import (
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
	"strings"
)

var RedisClient *redis.ClusterClient

func init() {
	conn := beego.AppConfig.String("redis.conn")

	RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    strings.Split(conn, ","),
		Password: "",
	})

	_, err := RedisClient.Ping().Result()

	if err != nil {
		beego.Error("Redis Connection error. %s", err)
	}

}

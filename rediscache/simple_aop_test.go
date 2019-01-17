package rediscache

import (
	"context"
	"github.com/astaxie/beego"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func setup() {
	ap, _ := os.Getwd()
	beego.TestBeegoInit(ap + "/..")
	BuildRedisClient()
}

func TestSimpleAop(t *testing.T) {

	setup()
	ctx := context.TODO()

	cacheKey := "testtest" + strconv.Itoa(rand.Intn(100000000))

	RedisClient.Del(cacheKey)

	options := &SimpleOptions{}
	options.Key = cacheKey
	options.Rt = reflect.TypeOf("")
	options.Expires = 30 * time.Second

	// 第一次保证cache里面取不到值
	val, fromCache, err := SimpleAop(&ctx, options, func(i *context.Context) (interface{}, error) {
		return "123", nil
	})

	if fromCache || err != nil || val != "123" {
		t.Fatal("make sure not from cache error!", err)
	}

	// 第二次，肯定从cache里面取值
	val, fromCache, err = SimpleAop(&ctx, options, func(i *context.Context) (interface{}, error) {
		return "123", nil
	})

	if !fromCache || err != nil || val != "123" {
		t.Fatal("make sure not from cache error!", err)
	}

	// 第三次，缓存空值
	RedisClient.Del(cacheKey)
	options.EmptyExpires = 10 * time.Second
	val, fromCache, err = SimpleAop(&ctx, options, func(i *context.Context) (interface{}, error) {
		return "", nil
	})

	if fromCache || err != nil || val != "" {
		t.Fatal("make sure not from cache error!", err)
	}

	cacheV, _ := RedisClient.Get(cacheKey).Result()
	if cacheV != EmptyFlag {
		t.Error("empty cache fail")
	}

	// 第四次，空值从cache里面取出
	val, fromCache, err = SimpleAop(&ctx, options, func(i *context.Context) (interface{}, error) {
		return "123", nil
	})

	if !fromCache || err != nil || val != nil {
		t.Fatal("make sure from empty cache error!", err)
	}

}


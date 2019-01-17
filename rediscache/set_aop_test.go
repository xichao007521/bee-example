package rediscache

import (
	"context"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestSetAop(t *testing.T) {
	testSetup()
	ctx := context.TODO()

	cacheKey := "testtest_set_" + strconv.Itoa(rand.Intn(100000000))

	RedisClient.Del(cacheKey)

	options := &SetOptions{}
	options.Key = cacheKey
	options.Rt = reflect.TypeOf("")
	options.Expires = 30 * time.Second

	RedisClient.Del(cacheKey)

	// 第一次保证不从cache里面取值
	val, fromCache, err := SetAop(&ctx, options, func(i *context.Context) ([]interface{}, error) {
		var result []interface{}
		result = append(result, "11", "22")
		return result, nil
	})

	if err != nil || fromCache || len(val) != 2 || val[0] != "11" || val[1] != "22" {
		t.Fatal("1. must not be from cache FAIL")
	}

	// 第二次保证从cache里面取值
	val, fromCache, err = SetAop(&ctx, options, func(i *context.Context) ([]interface{}, error) {
		var result []interface{}
		result = append(result, "11", "22")
		return result, nil
	})

	if err != nil || !fromCache || len(val) != 2 || !(val[0] != "11" || val[1] != "11") || !(val[0] != "22" || val[1] != "22") {
		t.Fatal("2. must not be from cache FAIL")
	}

	// 第三次保证保证存进去EmptyFlag
	RedisClient.Del(cacheKey)
	options.EmptyExpires = 30 * time.Second
	val, fromCache, err = SetAop(&ctx, options, func(i *context.Context) ([]interface{}, error) {
		return nil, nil
	})

	if err != nil || fromCache || val != nil {
		t.Fatal("3. must be empty FAIL")
	}
	vals, _ := RedisClient.SMembers(cacheKey).Result()
	if len(vals) != 1 && vals[0] != EmptyFlag {
		t.Fatal("3-1. must be empty FAIL")
	}

}

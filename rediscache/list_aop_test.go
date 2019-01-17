package rediscache

import (
	"context"
	"do-global.com/bee-example/models"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestListAop(t *testing.T) {
	testSetup()
	ctx := context.TODO()

	cacheKey := "testtest_list_" + strconv.Itoa(rand.Intn(100000000))

	RedisClient.Del(cacheKey)

	options := &ListOptions{}
	options.Key = cacheKey
	options.Rt = reflect.TypeOf(models.User{})
	options.Expires = 30 * time.Second
	options.Start = 0
	options.Stop = -1

	RedisClient.Del(cacheKey)

	// 第一次保证不从cache里面取值
	val, fromCache, err := ListAop(&ctx, options, func(i *context.Context) ([]interface{}, error) {
		var result []interface{}
		result = append(result, models.User{Id: 1}, models.User{Id: 2})
		return result, nil
	})

	if err != nil || fromCache || len(val) != 2 {
		t.Fatal("1. must not be from cache FAIL")
	}
	for index, vv := range val {
		u := vv.(models.User)
		if index == 0 && u.Id != 1 {
			t.Fatal("1. must not be from cache FAIL, get 0 val")
		}
		if index == 1 && u.Id != 2 {
			t.Fatal("1. must not be from cache FAIL, get 1 val")
		}
	}

	// 第二次保证从cache里面取值
	val, fromCache, err = ListAop(&ctx, options, func(i *context.Context) ([]interface{}, error) {
		var result []interface{}
		result = append(result, models.User{Id: 1}, models.User{Id: 2})
		return result, nil
	})

	if err != nil || !fromCache || len(val) != 2 {
		t.Fatal("2. must not be from cache FAIL")
	}
	for index, vv := range val {
		u := vv.(models.User)
		if index == 0 && u.Id != 1 {
			t.Fatal("2. must not be from cache FAIL, get 0 val")
		}
		if index == 1 && u.Id != 2 {
			t.Fatal("2. must not be from cache FAIL, get 1 val")
		}
	}

	// 第三次保证保证存进去EmptyFlag
	RedisClient.Del(cacheKey)
	options.EmptyExpires = 30 * time.Second
	val, fromCache, err = ListAop(&ctx, options, func(i *context.Context) ([]interface{}, error) {
		return nil, nil
	})

	if err != nil || fromCache || val != nil {
		t.Fatal("3. must be empty FAIL")
	}
	vals, _ := RedisClient.LRange(cacheKey, 0, 0).Result()
	if len(vals) != 1 && vals[0] != EmptyFlag {
		t.Fatal("3-1. must be empty FAIL")
	}

	// 第四保证取空值，并且从cache里取
	val, fromCache, err = ListAop(&ctx, options, func(i *context.Context) ([]interface{}, error) {
		return nil, nil
	})

	if err != nil || !fromCache || len(val) != 0 {
		t.Fatal("4. must be empty from cache FAIL")
	}

}

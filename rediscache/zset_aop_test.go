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

func TestZSetAop(t *testing.T) {
	testSetup()

	ctx := context.TODO()

	cacheKey := "testtest_zset_" + strconv.Itoa(rand.Intn(100000000))

	RedisClient.Del(cacheKey)

	options := &ZSetOptions{}
	options.Key = cacheKey
	options.Rt = reflect.TypeOf("")
	options.Start = 0
	options.Stop = -1
	options.Expires = 30 * time.Second

	// 第一次, map 类型 cache里没有值，从fallback取到并回填
	options.IsMap = true
	val, fromCache, err := ZSetAop(&ctx, options, func(i *context.Context) (interface{}, error) {
		r := make(map[interface{}]float64)
		r["1"] = 1
		r["2"] = 2
		return r, nil
	})
	valm := val.(map[interface{}]float64)
	if err != nil || fromCache || valm["1"] != 1 || valm["2"] != 2 {
		t.Fatal("1. must not be from cache FAIL")
	}
	// 第二次, map 类型 cache里有值
	val, fromCache, err = ZSetAop(&ctx, options, func(i *context.Context) (interface{}, error) {
		r := make(map[interface{}]float64)
		r["1"] = 1
		r["2"] = 2
		return r, nil
	})
	valm = val.(map[interface{}]float64)
	if err != nil || !fromCache || valm["1"] != 1 || valm["2"] != 2 {
		t.Fatal("2. must be from cache FAIL")
	}

	RedisClient.Del(cacheKey)
	// 第三次, 非map 类型 cache没值
	options.IsMap = false
	options.ScoreField = "Id"
	options.Rt = reflect.TypeOf(models.User{})
	val, fromCache, err = ZSetAop(&ctx, options, func(i *context.Context) (interface{}, error) {
		var r []interface{}
		r = append(r, models.User{Id: 1, Name: "name1"})
		r = append(r, models.User{Id: 2, Name: "name2"})
		return r, nil
	})
	valarr := val.([]interface{})
	u1 := valarr[0].(models.User)
	u2 := valarr[1].(models.User)
	if err != nil || fromCache || len(valarr) != 2 || u1.Name != "name1" || u1.Id != 1 || u2.Name != "name2" || u2.Id != 2 {
		t.Fatal("3. must not be from cache FAIL")
	}

	// 第四次, 非map 类型 cache里有值
	val, fromCache, err = ZSetAop(&ctx, options, func(i *context.Context) (interface{}, error) {
		var r []interface{}
		r = append(r, models.User{Id: 1, Name: "name1"})
		r = append(r, models.User{Id: 2, Name: "name2"})
		return r, nil
	})
	valarr = val.([]interface{})
	u1 = valarr[0].(models.User)
	u2 = valarr[1].(models.User)
	if err != nil || !fromCache || len(valarr) != 2 || u1.Name != "name1" || u1.Id != 1 || u2.Name != "name2" || u2.Id != 2 {
		t.Fatal("4. must not be from cache FAIL")
	}
}

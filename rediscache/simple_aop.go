package rediscache

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"reflect"
)

type SimpleOptions struct {
	Options
}

func SimpleAop(ctx *context.Context, options *SimpleOptions, fallback func(*context.Context) (interface{}, error)) (interface{}, bool, error) {
	cacheV, err := RedisClient.Get(options.Key).Result()
	if cacheV != "" {
		rtv := reflect.New(options.Rt)
		rv := rtv.Interface()
		if cacheV == EmptyFlag {
			return nil, true, nil
		}
		err := json.Unmarshal([]byte(cacheV), rv)
		return reflect.ValueOf(rv).Elem().Interface(), true, err
	}
	beego.Warn("[REDIS] cant get value from redis cache, maybe load from db!")
	var result interface{} = nil
	result, err = fallback(ctx)
	if err != nil {
		return nil, false, err
	}
	// 是否回填cache成功
	rewriteSuccess := false
	if result != nil {
		jsonB, err := json.Marshal(result)
		if err != nil {
			return nil, false, err
		}
		cacheV = string(jsonB)
		cLength := 0
		switch result.(type) {
		// string 类型单独处理
		case string:
			cLength = len(result.(string))
		default:
			cLength = len(cacheV)
		}
		if cLength > 0 {
			RedisClient.Set(options.Key, cacheV, options.Expires)
			rewriteSuccess = true
		}
	}
	// 是否需要存储空值
	if !rewriteSuccess && options.EmptyExpires > 0 {
		RedisClient.Set(options.Key, EmptyFlag, options.EmptyExpires)
		beego.Warn("[REDIS] cache empty value, key ", options.Key)
	}
	return result, false, nil
}

package rediscache

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"reflect"
)

type SetOptions struct {
	Options
}

func SetAop(ctx *context.Context, options *SetOptions, fallback func(*context.Context) ([]interface{}, error)) ([]interface{}, bool, error) {
	cacheVs, err := RedisClient.SMembers(options.Key).Result()

	var result []interface{}
	// 从cache里取到值
	if len(cacheVs) > 0 {
		if len(cacheVs) == 1 && cacheVs[0] == EmptyFlag {
			return result, true, nil
		}
		for _, cacheV := range cacheVs {
			if cacheV == EmptyFlag {
				continue
			}
			rtv := reflect.New(options.Rt)
			rv := rtv.Interface()
			err := json.Unmarshal([]byte(cacheV), rv)
			if err != nil {
				return nil, false, err
			}
			result = append(result, reflect.ValueOf(rv).Elem().Interface())
		}
		return result, true, err
	}
	beego.Warn("[REDIS][SET] cant get value from redis cache, maybe load from db!")
	result, err = fallback(ctx)
	if err != nil {
		return nil, false, err
	}
	// 回填
	rewriteCount := 0
	if result != nil && len(result) > 0 {
		for _, item := range result {
			cacheV, isEmpty, err := GetCacheValueItem(item)
			if err != nil {
				return nil, false, err
			}
			if !isEmpty {
				RedisClient.SAdd(options.Key, cacheV)
				RedisClient.Expire(options.Key, options.Expires)
				rewriteCount++
			}
		}
	}

	// 空值回填
	if rewriteCount == 0 && options.EmptyExpires > 0 {
		RedisClient.SAdd(options.Key, EmptyFlag)
		RedisClient.Expire(options.Key, options.EmptyExpires)
		beego.Warn("[REDIS][SET] cache empty value, key:", options.Key)
	}

	return result, false, nil

}

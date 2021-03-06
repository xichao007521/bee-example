package rediscache

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"reflect"
)

type ListOptions struct {
	Options
	// 此处一定要注意，除非要取第一条数据，否则一定要设置Start和Stop
	Start int64
	Stop  int64
}

func ListAop(options *ListOptions, fallback func() ([]interface{}, error)) ([]interface{}, bool, error) {
	if options.Key == "" {
		return nil, false, errors.New("Key must not be empty!")
	}
	cacheVs, err := RedisClient.LRange(options.Key, options.Start, options.Stop).Result()
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
	beego.Warn("[REDIS][LIST] cant get value from redis cache, maybe load from db!")
	result, err = fallback()
	if err != nil {
		return nil, false, err
	}
	// 回填
	rewriteCount := 0
	if result != nil && len(result) > 0 {
		for _, item := range result {
			cacheV, isEmpty, err := GetCacheValueItem(item)
			if err != nil {
				beego.Warn("[REDIS][LIST] GetCacheValueItem error!", err)
				continue
			}
			if !isEmpty {
				RedisClient.RPush(options.Key, cacheV)
				rewriteCount++
			}
		}
		if rewriteCount > 0 {
			RedisClient.Expire(options.Key, options.Expires)
		}
	}

	// 空值回填
	if rewriteCount == 0 && options.EmptyExpires > 0 {
		RedisClient.RPush(options.Key, EmptyFlag)
		RedisClient.Expire(options.Key, options.EmptyExpires)
		beego.Warn("[REDIS][LIST] cache empty value, key:", options.Key)
	}

	return result, false, nil

}

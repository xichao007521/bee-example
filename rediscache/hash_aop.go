package rediscache

import (
	"do-global.com/bee-example/globals"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/pkg/errors"
	"reflect"
)

type HashOptions struct {
	Options

	Fields []string

	// field attribute in model
	FieldAttr string
}

func HashAop(options *HashOptions, fallback func() ([]interface{}, error)) ([]interface{}, bool, error) {
	if options.Key == "" {
		return nil, false, errors.New("Key must not be empty!")
	}
	if len(options.Fields) == 0 {
		return nil, false, errors.New("Fields must not be empty!")
	}
	if options.FieldAttr == "" {
		return nil, false, errors.New("FieldAttr must not be empty!")
	}
	cacheVs, err := RedisClient.HMGet(options.Key, options.Fields...).Result()
	var result []interface{}
	shouldCallback := len(cacheVs) == 0
	if !shouldCallback {
		for _, cacheV := range cacheVs {
			if cacheV == nil {
				shouldCallback = true
				beego.Warn("[REDIS][HASH] key ", options.Key, " has nil value, values", cacheVs)
				break
			}
			rtv := reflect.New(options.Rt)
			rv := rtv.Interface()
			err := json.Unmarshal([]byte(cacheV.(string)), rv)
			if err != nil {
				return nil, false, err
			}
			result = append(result, reflect.ValueOf(rv).Elem().Interface())
		}
		if !shouldCallback {
			return result, true, nil
		}
	}

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
				beego.Warn("[REDIS][HASH] GetCacheValueItem error!", err)
				continue
			}
			// 看一下struct里面的作为field的Field是否有正确的值
			fieldV := ""
			if options.FieldAttr != "" {
				iv := reflect.ValueOf(&item)
				if reflect.TypeOf(item).Kind() == reflect.Struct {
					ivf := iv.Elem().Elem().FieldByName(options.FieldAttr)
					if ivf.IsValid() {
						vv, success := globals.Primary2String(ivf.Interface(), ivf.Kind())
						if success {
							fieldV = vv
						}
					}
				}
			}
			if fieldV == "" {
				beego.Warn("[REDIS][HASH] key ", options.Key, " value ", item, " has not valid fieldValue!!!")
			}
			if !isEmpty && fieldV != ""{
				RedisClient.HSet(options.Key, fieldV, cacheV)
				rewriteCount++
			}
		}
		if rewriteCount > 0 {
			RedisClient.Expire(options.Key, options.Expires)
		}
	}

	return result, false, nil
}

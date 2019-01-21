package rediscache

import (
	"context"
	"do-global.com/bee-example/globals"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
	"reflect"
	"strconv"
)

// ZsetAop 中，EmptyExpires 使用无效
type ZSetOptions struct {
	Options

	// 如果为true，表示返回map，
	// key是member，value是score
	// NOTICE!!! 并且返回结果的map的value一定是float64类型
	IsMap      bool

	Desc       bool
	ByScore    bool

	// 如果是struct类型，需要指定一下，使用哪个字段作为score，否则score默认是0
	ScoreField string

	// 如果是ByScore，必须需要指定Min和Max
	Min    int
	Max    int
	Offset int64
	Count  int64

	// 如果不是ByScore, 需要指定Start, Stop
	Start int64
	Stop  int64
}

// NOTICE!!! 如果fallback返回结果的map的value一定是float64类型
func ZSetAop(ctx *context.Context, options *ZSetOptions, fallback func(*context.Context) (interface{}, error)) (interface{}, bool, error) {
	zrangeBy := redis.ZRangeBy{
		Min:    strconv.Itoa(options.Min),
		Max:    strconv.Itoa(options.Max),
		Offset: options.Offset,
		Count:  options.Count,
	}
	if options.Stop == 0 {
		options.Stop = -1
	}
	var cacheVs []redis.Z
	var err error = nil
	if options.Desc && options.ByScore {
		cacheVs, err = RedisClient.ZRevRangeByScoreWithScores(options.Key, zrangeBy).Result()
	} else if options.Desc {
		cacheVs, err = RedisClient.ZRevRangeWithScores(options.Key, options.Start, options.Stop).Result()
	} else if options.ByScore {
		cacheVs, err = RedisClient.ZRangeByScoreWithScores(options.Key, zrangeBy).Result()
	} else {
		cacheVs, err = RedisClient.ZRangeWithScores(options.Key, options.Start, options.Stop).Result()
	}
	if err != nil {
		return nil, false, err
	}
	var result []interface{}
	var mapResult = make(map[interface{}]float64)
	if len(cacheVs) > 0 {
		for _, cacheV := range cacheVs {
			rtv := reflect.New(options.Rt)
			rv := rtv.Interface()
			err := json.Unmarshal([]byte(cacheV.Member.(string)), rv)
			if err != nil {
				return nil, false, err
			}
			vv := reflect.ValueOf(rv).Elem().Interface()
			if options.IsMap {
				mapResult[vv] = cacheV.Score
			} else {
				result = append(result, vv)
			}
		}
		if options.IsMap {
			return mapResult, true, nil
		} else {
			return result, true, nil
		}
	}

	beego.Warn("[REDIS][ZSET] cant get value from redis cache, maybe load from db!")

	fResult, err := fallback(ctx)
	if err != nil {
		return nil, false, err
	}
	rewriteCount := 0
	if options.IsMap {
		for k, score := range fResult.(map[interface{}]float64) {
			cacheV, isEmpty, err := GetCacheValueItem(k)
			if err != nil {
				beego.Warn("[REDIS][ZSET] GetCacheValueItem error!", err)
				continue
			}
			if !isEmpty {
				RedisClient.ZAdd(options.Key, redis.Z{Member: cacheV, Score: score})
				rewriteCount++
			}
		}
	} else {
		for _, resultItem := range fResult.([]interface{}) {
			var score float64 = 0
			// 看一下struct里面的作为score的field是否有正确的值
			if options.ScoreField != "" {
				iv := reflect.ValueOf(&resultItem)
				if reflect.TypeOf(resultItem).Kind() == reflect.Struct {
					ivf := iv.Elem().Elem().FieldByName(options.ScoreField)
					if ivf.IsValid() {
						vv, success := globals.Number2Float64(ivf.Interface(), ivf.Kind())
						if success {
							score = vv
						}
					}
				}
			}
			cacheV, isEmpty, err := GetCacheValueItem(resultItem)
			if err != nil {
				beego.Warn("[REDIS][ZSET] GetCacheValueItem error!", err)
				continue
			}
			if !isEmpty {
				RedisClient.ZAdd(options.Key, redis.Z{Member: cacheV, Score: score})
				rewriteCount++
			}
		}
	}
	if rewriteCount > 0 {
		RedisClient.Expire(options.Key, options.Expires)
	}
	return fResult, false, nil
}


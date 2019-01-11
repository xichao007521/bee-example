package services

import (
	"context"
	"do-global.com/bee-example/cache"
	"do-global.com/bee-example/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"reflect"
	"time"
)

type UserService struct {
	Service
}

func (*UserService) Login(username string, password string) *models.User {
	// TODO do something
	return &models.User{
		Id:   1,
		Name: "u1",
	}
}

func (*UserService) GetUser(reqCtx *context.Context, uid string) (*models.User, error)  {
	rt := reflect.TypeOf(models.User{})
	rv, err := GetFromCache(reqCtx, rt, uid, func(i *context.Context) (interface{}, error) {
		// get from db
		return &models.User{
			Id:   1,
			Name: "u1",
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return rv.(*models.User), nil
}

func GetFromCache(ctx *context.Context, rt reflect.Type, key string, fallback func(*context.Context) (interface{}, error)) (interface{}, error) {
	cacheV, err := cache.RedisClient.Get(key).Result()
	if cacheV != "" {
		rtv := reflect.New(rt)
		rv := rtv.Interface()
		err := json.Unmarshal([]byte(cacheV), rv)
		return rv, err
	}
	beego.Warn("[REDIS] cant get value from redis cache, maybe load from db!")
	var result interface{} = nil
	result, err = fallback(ctx)
	if err != nil {
		return nil, err
	}
	if result != nil {
		jsonB, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}
		cacheV = string(jsonB)
		cache.RedisClient.Set(key, cacheV, 24 * 30 * time.Hour)
	}
	return result, nil
}

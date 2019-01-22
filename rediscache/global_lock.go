package rediscache

import (
	"errors"
	"github.com/astaxie/beego"
	"math/rand"
	"time"
)

type GlobalLockOptions struct {
	Key           string
	Timeout time.Duration
	Expire  time.Duration
}

func GlobalLock(options *GlobalLockOptions, fallback func()) (bool, error) {
	if options.Key == "" {
		return false, errors.New("Key must not be empty!")
	}
	defer RedisClient.Del(options.Key)

	if options.Timeout == 0 {
		options.Timeout = 60 * time.Second
	}
	if options.Expire == 0 {
		options.Expire = 30 * time.Second
	}
	startTime := time.Now()
	for ; ;  {
		success, err := RedisClient.SetNX(options.Key, "1", options.Expire).Result()
		if err == nil && success {
			beego.Debug("add lock success: ", options.Key)
			fallback()
			return true, nil
		}
		beego.Debug("waiting lock: ", options.Key)
		// 短暂休眠，避免可能的活锁
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)

		spendTime := time.Now().Sub(startTime)

		if spendTime > options.Timeout {
			beego.Warn("acquireLock timeout: ", options.Key)
			return false, errors.New("acquireLock timeout " + options.Key)
		}
	}
}

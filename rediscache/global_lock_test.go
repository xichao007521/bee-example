package rediscache

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestGlobalLock(t *testing.T) {
	testSetup()

	cacheKey := "testtest_global_lock_" + strconv.Itoa(rand.Intn(100000000))

	options := &GlobalLockOptions{
		Key: cacheKey,
		Expire: 5 * time.Second,
		Timeout: 10 * time.Second,
	}

	RedisClient.Del(cacheKey)

	go func() {
		success, err := GlobalLock(options, func() {
			t.Log("hahaha1")
			// 休眠3秒
			time.Sleep(3 * time.Second)
		})

		if !success || err != nil {
			t.Fatal("1. acquire global lock FAIL")
		}
	}()

	// 确保上面的go先执行
	time.Sleep(300 * time.Millisecond)
	go func() {
		startTime := time.Now()
		success, err := GlobalLock(options, func() {
			t.Log("hahaha2")
		})
		if !success || err != nil {
			t.Fatal("2. acquire global lock FAIL")
		}
		if time.Now().Sub(startTime) < 2 * time.Second {
			t.Fatal("2.1. should not be happened, last operation hold lock 3s!")
		}
	}()

	time.Sleep(6 * time.Second)
}

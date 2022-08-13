package repository

import (
	"fmt"
	"go.uber.org/zap"
	"gohub/init/cusRedis"
	"gohub/init/cusZap"
	"gorm.io/gorm"
	"time"
)

func dbOperations(function func() error, funcName string) (err error) {
	funcNameMap[funcName] = struct{}{}
	start := time.Now()
	err = function()
	tmConsuming := time.Since(start)
	sqlPerformance(funcName, true, false, false)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果报错为“记录没有找到”，说明不存在
			return gorm.ErrRecordNotFound
		}
		sqlPerformance(funcName, false, false, true)
		cusZap.Error(fmt.Sprintf("%v failed...", funcName), zap.String("err", err.Error()))
		return err
	}
	if tmConsuming > 200*time.Millisecond {
		sqlPerformance(funcName, true, true, false)
	}
	return nil
}

func sqlPerformance(funcName string, count, slowSQL, err bool) {
	redisHSet := func(filed string) {
		exist, _ := cusRedis.Rdb.HExists(funcName, filed).Result()
		if exist {
			num, _ := cusRedis.Rdb.HGet(funcName, filed).Int()
			num++
			cusRedis.Rdb.HSet(funcName, filed, num)
		} else {
			cusRedis.Rdb.HSet(funcName, filed, 1)
		}
	}
	if count {
		redisHSet("count")
	}
	if slowSQL {
		redisHSet("slowSQL")
	}
	if err {
		redisHSet("err")
	}
}

package cusRedis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gohub/init/cusZap"
)

var Rdb *redis.Client

func RedisInit() bool {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", viper.GetString("cusRedis.host"), viper.GetInt("cusRedis.port")),
		Password: viper.GetString("cusRedis.password"),
		DB:       viper.GetInt("cusRedis.DB"),
	})
	if _, err := Rdb.Ping().Result(); err != nil {
		cusZap.Error("cusRedis init failed...", zap.String("err", err.Error()))
		return false
	}
	cusZap.Info("cusRedis init success...")
	return true
}

package main

import (
	"gohub/init/configs"
	"gohub/init/cusRedis"
	initZap "gohub/init/cusZap"
	"gohub/init/mysql"
	"gohub/init/router"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}
	dir = strings.Replace(dir, "\\", "/", -1)
	if configsInit := configs.ConfigsInit("configs", "yaml", dir+"/../configs/"); !configsInit {
		return
	}
	if zapInit := initZap.ZapInit(dir); !zapInit {
		return
	}
	if mySQLInit := mysql.MySQLInit(); !mySQLInit {
		return
	}
	if redisInit := cusRedis.RedisInit(); !redisInit {
		return
	}
	if routerInit := router.RouterInit(); !routerInit {
		return
	}
}

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
	curPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return
	}
	curPath = strings.Replace(curPath, "\\", "/", -1)
	if configsInit := configs.ConfigsInit("configs", "yaml", curPath+"/../configs/"); !configsInit {
		return
	}
	if zapInit := initZap.ZapInit(curPath); !zapInit {
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

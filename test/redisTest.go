package main

import (
	"fmt"
	"gohub/init/configs"
	"gohub/init/cusRedis"
	initZap "gohub/init/cusZap"
	"gohub/init/mysql"
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
	//err = cusRedis.Rdb.Set("15090386881", "123456", 5*time.Minute).Err()
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	time.Sleep(2 * time.Second)
	//	result, err := cusRedis.Rdb.Get("15090386881").Result()
	//	if err != nil {
	//		fmt.Println("获取失败")
	//		fmt.Println(err)
	//		return
	//	} else {
	//		fmt.Println(result)
	//	}
	//}
	result, err := cusRedis.Rdb.Get("18104077689").Result()
	if err != nil {
		return
	}
	fmt.Println(result)

}

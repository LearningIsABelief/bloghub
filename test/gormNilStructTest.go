package main

import (
	"fmt"
	"gohub/init/configs"
	"gohub/init/cusRedis"
	initZap "gohub/init/cusZap"
	"gohub/init/mysql"
	"gohub/internal/model"
	"gohub/internal/repository"
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
	expr := 1
	switch expr {
	case 1:
		{
			user := &model.User{}
			fmt.Println(user == nil)
			err = repository.CreateAUser(user)
			if err != nil {
				fmt.Printf("err:%v\n", err)
				return
			}
		}
	case 2:
		{
			var user *model.User
			fmt.Println(user == nil)
			err = repository.CreateAUser(user)
			if err != nil {
				fmt.Printf("err:%v\n", err)
				return
			}
		}
	}

}

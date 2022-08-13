package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gohub/init/cusZap"
	"gohub/internal/router"
)

func RouterInit() bool {
	e := gin.Default()
	// e.Use(cusZap.GinLogger(), cusZap.GinRecovery(true))
	router.RegisterRouters(e)
	err := e.Run(fmt.Sprintf(":%v", viper.GetInt("addr")))
	if err != nil {
		fmt.Println("router初始化错误")
		return false
	}
	cusZap.Info("router init success...")
	return true
}

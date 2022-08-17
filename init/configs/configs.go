package configs

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var RunMode string

func ConfigsInit(configName, configType, configPath string) bool {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("没有找到指定的配置文件")
		} else {
			fmt.Println(err)
		}
		return false
	}
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file changed:", in.Name)
	})
	viper.WatchConfig()
	// cusZap.Info("configRead init success...", zap.String("status", "success"))
	RunMode = viper.GetString("run_mode")
	viper.Set("cusPath", configPath)
	return true
}

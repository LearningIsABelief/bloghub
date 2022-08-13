package mysql

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gohub/init/cusZap"
	model "gohub/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func MySQLInit() bool {
	user := viper.GetString("mysql.user")
	password := viper.GetString("mysql.password")
	host := viper.GetString("mysql.host")
	port := viper.GetInt("mysql.port")
	dbname := viper.GetString("mysql.dbname")
	charset := viper.GetString("mysql.charset")
	if len(user) == 0 || len(password) == 0 || len(host) == 0 || len(dbname) == 0 || len(charset) == 0 {
		cusZap.Logger.Error("mysql configs wrong...")
	}
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&loc=Local", user, password, host, port, dbname, charset)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		cusZap.Logger.Error("mysql connect failed...", zap.String("err", err.Error()))
		return false
	}
	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		cusZap.Logger.Error("mysql autoMigrate failed...", zap.String("err", err.Error()))
		return false
	}
	cusZap.Logger.Info("mysql init success...")
	return true
}

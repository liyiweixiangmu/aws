package config

import (
	"github.com/spf13/viper"
	"log"
)

var Configuration *viper.Viper

func Init() error {
	//Env := os.Getenv("ENVIRON")
	Configuration = viper.New()
	Configuration.SetConfigName("config")
	Configuration.SetConfigType("toml")
	Configuration.AddConfigPath("./")
	//Configuration.Set("Environment", "config")
	err := Configuration.ReadInConfig()
	if err != nil {
		log.Panic("配置文件读取错误", err)
		return err
	}
	return nil
}

package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Service  ServiceConfig  `mapstructure:"service"`
	Dingtalk DingtalkConfig `mapstructure:"dingtalk"`
}

type ServiceConfig struct {
	Port int
}

type DingtalkConfig struct {
	Webhook string `mapstructure:"webhook"`
	Secret  string `mapstructure:"secret"`
}

var AppConfig Config

func init() {

	viper.SetConfigName("config")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error read config file: %s", err))
	}
	parseConfig()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		parseConfig()
	})
}

func parseConfig() {
	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error parse config file: %s", err))
	}
}

package config

import (
	"log"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func Load(configName string) {
	log.Println("InitConfig  ------" + configName)
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	// 3. 环境变量配置文件查找的路径，相对于 main.go
	viper.AddConfigPath(".")
	viper.WatchConfig()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

type Config struct {
	Server struct {
		Prefix string `yaml:"prefix"`
	} `yaml:"server"`

	DB struct {
		DSN string `yaml:"dsn"`
	} `yaml:"db"`

	Redis struct {
		Addr string `yaml:"addr"`
	} `yaml:"redis"`
}

func Get(key string) interface{} {
	return viper.Get(key)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func All() map[string]interface{} {
	return viper.AllSettings()
}

func GetInt64(path string, defaultValue ...interface{}) int64 {
	return cast.ToInt64(internalGet(path, defaultValue...))
}

func GetFloat64(path string, defaultValue ...interface{}) float64 {
	return cast.ToFloat64(internalGet(path, defaultValue...))
}

func internalGet(path string, defaultValue ...interface{}) interface{} {
	// config 或者环境变量不存在的情况
	if !viper.IsSet(path) || viper.Get(path) == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return viper.Get(path)
}

package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// pflag 解析入参，通过viper绑定后读取
func init() {
	pflag.StringP("conf", "c", "conf.yml", "app yaml config path")
	pflag.BoolP("debug", "D", false, "print log debug info")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatalf("app config, viper bind pflags fail, %s", err)
	}
	viper.SetConfigFile(viper.GetString("conf"))
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("app config parse fail, %s", err)
	}

	// 是否debug
	if IsDebug() {
		for k, v := range viper.AllSettings() {
			fmt.Printf("%s\t=> %v\n", k, v)
		}
		// os.Exit(0)
	}
}

// IsDebug 是否开启debug
func IsDebug() bool  {
	return viper.GetBool("debug")
}

// GetString 获取字符串性配置
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt 获取整型配置，app.xx
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetInt64 返回int64整型
func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

// GetDuration 获取持续时间
func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

// GetFloat 获取浮点型配置
func GetFloat(key string) float64 {
	return viper.GetFloat64(key)
}

// GetStringSlice 获取字符串切片
func GetStringSlice(key string)[]string {
	return viper.GetStringSlice(key)
}


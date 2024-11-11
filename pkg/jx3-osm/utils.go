package jx3osm

import (
	"log"
	"os"
	"time"
	_ "time/tzdata"
)

// 获取需要监控的服务器列表
func GetServerListToMonitor() []string {
	return GLO_CONF.Servers
}

// 初始化程序配置
func init_jx3_osm_configs() {

	// 配置文件相关设置
	_GLO_VIPER.SetConfigName("config")
	_GLO_VIPER.SetConfigType("toml")
	_GLO_VIPER.AddConfigPath(".")

	// 设置默认值(仅当任意来源的键值对配置不存在时有效)
	_GLO_VIPER.SetDefault("servers", []string{"斗转星移"}) // 设置默认需要监控的服务器列表
	_GLO_VIPER.SetDefault("tz", "Asia/Shanghai")       // 设置默认时区变量
	_GLO_VIPER.SetDefault("useracct", "your_xiaomi_account")
	_GLO_VIPER.SetDefault("password", "your_account_password")
	_GLO_VIPER.SetDefault("xiaoaisn", "your_xiaoai_serial_number")

	// 注意：环境变量的优先级高于 config 文件
	// 在环境变量中，由空格分隔的字符串可被识别为 slice 切片列表
	_GLO_VIPER.SetEnvPrefix("JX3_OSM")
	_GLO_VIPER.AutomaticEnv() // 尝试从环境变量读取默认值

	// 读取配置文件，当失败时尝试回写默认值
	if err := _GLO_VIPER.ReadInConfig(); err != nil {
		log.Println("Read config failed, err:", err)
		_GLO_VIPER.SafeWriteConfig()
	}

	if err := _GLO_VIPER.Unmarshal(GLO_CONF); err != nil {
		log.Println("Unmarshal config failed, err:", err)
	}

	// 监听配置文件变化，当配置文件发生变更时，重新读取并应用新值
	_GLO_VIPER.WatchConfig()                       // 仅当未配置环境变量时生效
	os.Setenv("TZ", GLO_CONF.TZ)                   // 设置时区变量
	time.Local, _ = time.LoadLocation(GLO_CONF.TZ) // 应用时区设置
}

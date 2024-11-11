package jx3osm

import (
	"log"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ServersState map[string]bool // 服务器状态记录

// 定义备用的全局变量
var _GLO_VIPER = viper.New()
var GLO_CONF = &Config{}
var GLO_SERVER_MAPS = &ServerMaps{}
var GLO_SERVERS_STATE = &ServersState{}
var GLO_MAIN_SRV_STATES = make(MainServersState) // 主服务器状态

// 初始化 viper 配置
func init() {
	init_jx3_osm_configs() // 初始化程序配置

	_GLO_VIPER.OnConfigChange(func(event fsnotify.Event) {
		// 清空已有值，防止列表缩减时已有值无法正常删除
		GLO_CONF.Servers = []string{}
		if err := _GLO_VIPER.Unmarshal(GLO_CONF); err != nil {
			log.Println("Unmarshal config failed, err:", err)
		}

		// 重新配置全局时区
		os.Setenv("TZ", GLO_CONF.TZ)
		time.Local, _ = time.LoadLocation(GLO_CONF.TZ)

		// 重新创建 ServerMaps 映射表
		GLO_SERVER_MAPS = &ServerMaps{}
		GLO_SERVER_MAPS.ParseServerList()
		GLO_SERVERS_STATE.InitServersStates()
	})
}

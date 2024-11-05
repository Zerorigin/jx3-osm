package jx3osm

// 定义 Server 结构体
type Server struct {
	ZoneName    string `csv:"ZoneName"` // 所属大区
	ServName    string `csv:"ServName"` // 服务器名(唯一主键)
	Reserved_03 string `csv:"Reserved_03"`
	ServIP      string `csv:"ServIP"`   // 服务器 IP (同一主服的 IP/Port 二元组一致)
	ServPort    int    `csv:"ServPort"` // 服务器 Port (同一主服的 IP/Port 二元组一致)
	ZoneType    string `csv:"ZoneType"` // 大区类型(点月卡区...)
	ServType    string `csv:"ServType"` // 服务器类型(点月卡服...)
	Reserved_08 string `csv:"Reserved_08"`
	Reserved_09 string `csv:"Reserved_09"`
	ZoneCode    string `csv:"ZoneCode"`  // 所属大区代码(ID)
	ServZone    string `csv:"ServZone"`  // 服务器所属主服区域(同一主服数据互通)
	ZoneAlias   string `csv:"ZoneAlias"` // 所属大区别名(ISP, ZoneName, etc.)
	Reserved_13 string `csv:"Reserved_13"`
	Reserved_14 string `csv:"Reserved_14"`
	Reserved_15 string `csv:"Reserved_15"`
	Reserved_16 string `csv:"Reserved_16"`
}

// 定义 Config 结构体，方便使用 viper 自动反系列化
// 键名要对应，且不能有特殊字符
type Config struct {
	Servers  []string `toml:"servers"`  // 需要监控的服务器列表
	TZ       string   `toml:"tz"`       // 时区
	UserAcct string   `toml:"useracct"` // 小米小爱账号
	Password string   `toml:"password"` // 小米账号密码
	XiaoAiSN string   `toml:"xiaoaisn"` // 小米小爱 SerialNumber
}

type ServerMaps map[string]Server     // 定义 ServerMaps 映射表
type MainServersState map[string]bool // 定义主服务器状态映射表

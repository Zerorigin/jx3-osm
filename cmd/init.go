package cmd

import (
	jx3osm "jx3-osm/pkg/jx3-osm"
	"math/rand"
	"time"
)

var _GLO_RAND *rand.Rand // 全局随机数生成器

func init() {
	_GLO_RAND = rand.New(rand.NewSource(time.Now().UnixNano())) // 随机化初始种子
	jx3osm.GLO_SERVER_MAPS.ParseServerList()                    // 初始化线上服务器列表
	jx3osm.GLO_SERVERS_STATE.InitServersStates()                // 初始化服务状态
}

package jx3osm

import (
	"strings"

	"github.com/tiendc/go-deepcopy"
)

// 对外提供 Server 的 TestTcpConnection 接口
func (s *Server) TestTcpConnection() bool {
	return testTcpConnection(s)
}

// 对外提供 Server 的 TestUdpConnection 接口
func (s *Server) TestUdpConnection() bool {
	return testUdpConnection(s)
}

// 对外提供 ServerMaps 的 ParseServerList 接口
func (smaps *ServerMaps) ParseServerList(svrs_text ...string) {
	// 若预处理后的 data 为空白字符串，则从网络上下载文本数据
	data := strings.Join(svrs_text, "\n")
	if len(svrs_text) == 0 || strings.TrimSpace(data) == "" {
		data = smaps.downloadServerList()
	}
	smaps.parseServerList(data)
}

func (s ServersState) InitServersStates() {

	default_state := false           // 新服务器的默认初始状态
	svrs := GetServerListToMonitor() // 获取需要监控的服务器列表

	if len(s) == 0 {
		for _, srv := range svrs {
			s[srv] = default_state
		}
		return
	}

	var old ServersState
	_ = deepcopy.Copy(&old, &s)
	s = make(ServersState)
	smaps := *GLO_SERVER_MAPS
	for _, srv := range svrs {
		if _bool, ok := old[srv]; ok {
			s[srv] = _bool // 迁移现有的状态
		} else {
			zsrv := smaps[srv].ServZone
			if state, ok := GLO_MAIN_SRV_STATES[zsrv]; ok {
				s[srv] = state // 遵从主服务器的状态
			} else {
				s[srv] = default_state
			}
		}
	}
	GLO_SERVERS_STATE = &s
}

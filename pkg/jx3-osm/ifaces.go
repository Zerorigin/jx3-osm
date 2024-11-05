package jx3osm

// 约束 Server 对象必须具备的能力接口
type ServerCaps interface {
	testNetConnection(string) bool
}

package main

import (
	jx3osm "jx3-osm/pkg/jx3-osm"
	"time"
	_ "time/tzdata"
)

func init() {
	time.LoadLocation(jx3osm.GLO_CONF.TZ) // 设置全局时区
}

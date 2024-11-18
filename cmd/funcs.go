package cmd

import (
	"errors"
	"fmt"
	xiaoai_tts "jx3-osm/internal/xiaoai-tts"
	jx3osm "jx3-osm/pkg/jx3-osm"
	"strings"
	"time"
)

// 端口关闭状态下的健康检查
func healthCheckWithinPortClosed() {

	// 获取所有子服务器的端口状态
	var bools = make(map[bool]bool)
	servers_state := *jx3osm.GLO_SERVERS_STATE
	for _, cstate := range servers_state {
		if _, ok := bools[cstate]; !ok {
			bools[cstate] = cstate
		}
		if len(bools) == 2 {
			break
		}
	}

	if _, ok := bools[false]; !ok {
		return // 都不是关服状态，不必检查了
	}

	var xiaoai *xiaoai_tts.Xiaoai = nil

	// 对子服务器进行检查
	smaps := *jx3osm.GLO_SERVER_MAPS
	for _srv_, cstate := range servers_state {
		if cstate {
			continue // 都开服了，还检测个得噢
		}

		// 备用变量及信息
		srv := smaps[_srv_] // 获取子服务器信息
		msg := fmt.Sprintf("剑网叁“%s”服务器开服啦~!", srv.ServName)

		// 子服务器所属主服务器在列表内，且主服务器是开服状态
		if lstate, ok := jx3osm.GLO_MAIN_SRV_STATES[srv.ServZone]; ok && lstate {
			// 准备小爱音箱终端设备
			if xiaoai == nil {
				var e error
				xiaoai, e = getXiaoai()
				if e != nil {
					continue // 获取失败，跳过标记，继续下一个
				}
			}

			// 子服和主服互通数据，直接尝试通知和标记
			if err := xiaoai.Say(msg); err != nil {
				jx3osm.GLO_MAIN_SRV_STATES[srv.ServZone] = true
				continue // 发送通知失败，跳过标记子服务器，继续下一个
			}
			// 标记为开服状态
			servers_state[_srv_] = true
			jx3osm.GLO_MAIN_SRV_STATES[srv.ServZone] = true
			time.Sleep(time.Second * 3)
			continue
		}

		// 子服务器所属主服务器不在列表内，则需要进行检查和通知
		if srv.TestTcpConnection() {
			// 准备小爱音箱终端设备
			if xiaoai == nil {
				var e error
				xiaoai, e = getXiaoai()
				if e != nil {
					continue // 获取失败，跳过标记，继续下一个
				}
			}

			// 通知，然后更新状态
			if err := xiaoai.Say(msg); err != nil {
				jx3osm.GLO_MAIN_SRV_STATES[srv.ServZone] = true
				continue // 发送通知失败，跳过标记子服务器，继续下一个
			}
			// 标记为开服状态
			servers_state[_srv_] = true
			jx3osm.GLO_MAIN_SRV_STATES[srv.ServZone] = true
			time.Sleep(time.Second * 3)
		}
	}
}

// 端口打开状态下的健康检查
func healthCheckWithinPortOpened() {

	// 获取所有子服务器的端口状态
	var bools = make(map[bool]bool)
	servers_state := *jx3osm.GLO_SERVERS_STATE
	for _, cstate := range servers_state {
		if _, ok := bools[cstate]; !ok {
			bools[cstate] = cstate
		}
		if len(bools) == 2 {
			break
		}
	}

	if _, ok := bools[true]; !ok {
		return // 服务器都在维护，不必检查了
	}

	var xiaoai *xiaoai_tts.Xiaoai = nil

	// 对子服务器进行检查
	smaps := *jx3osm.GLO_SERVER_MAPS
	for _srv_, cstate := range servers_state {
		if !cstate {
			continue // 都关服了，还检测个得噢
		}

		// 备用变量及信息
		srv := smaps[_srv_] // 获取子服务器信息
		msg := fmt.Sprintf("剑网叁“%s”服务器维护啦~!", srv.ServName)

		// 子服务器所属主服务器在列表内，且主服务器是关服状态
		if lstate, ok := jx3osm.GLO_MAIN_SRV_STATES[srv.ServZone]; ok && !lstate {
			// 准备小爱音箱终端设备
			if xiaoai == nil {
				var e error
				xiaoai, e = getXiaoai()
				if e != nil {
					continue // 获取失败，跳过标记，继续下一个
				}
			}

			// 子服和主服互通数据，直接尝试通知和标记
			if err := xiaoai.Say(msg); err != nil {
				jx3osm.GLO_MAIN_SRV_STATES[srv.ServZone] = false
				continue // 发送通知失败，跳过标记子服务器，继续下一个
			}
			// 标记为维护状态
			servers_state[_srv_] = false
			jx3osm.GLO_MAIN_SRV_STATES[srv.ServZone] = false
			time.Sleep(time.Second * 3)
			continue
		}

		// 子服务器所属主服务器不在列表内，则需要进行检查和通知
		if !srv.TestTcpConnection() {
			// 准备小爱音箱终端设备
			if xiaoai == nil {
				var e error
				xiaoai, e = getXiaoai()
				if e != nil {
					continue // 获取失败，跳过标记，继续下一个
				}
			}

			// 通知，然后更新状态
			if err := xiaoai.Say(msg); err != nil {
				jx3osm.GLO_MAIN_SRV_STATES[srv.ServZone] = false
				continue // 发送通知失败，跳过标记子服务器，继续下一个
			}
			// 标记为维护状态
			servers_state[_srv_] = false
			jx3osm.GLO_MAIN_SRV_STATES[srv.ServZone] = false
			time.Sleep(time.Second * 3)
		}
	}
}

func getXiaoai() (*xiaoai_tts.Xiaoai, error) {
	var xiaoai *xiaoai_tts.Xiaoai = nil

	// 5次重试获取小爱音箱
	for i := 1; i <= 5; i++ {
		var err error = nil

		// 准备小爱音箱终端设备
		if xiaoai, err = xiaoai_tts.New(
			jx3osm.GLO_CONF.UserAcct,
			jx3osm.GLO_CONF.Password,
		); err != nil {
			time.Sleep(time.Microsecond * 365)
			continue // 获取小爱音箱失败，重试
		}

		if msgs_, err := xiaoai.GetDevices(); err != nil {
			time.Sleep(time.Microsecond * 365)
			continue // 获取小爱音箱列表失败，重试
		} else {
			for idx, dev := range msgs_.Data {
				if strings.EqualFold(dev.SerialNumber, jx3osm.GLO_CONF.XiaoAiSN) {
					xiaoai.SwitchDevice(int64(idx))
					return xiaoai, nil // 找到需要的小爱音箱了
				}
			}
		}

	}
	return nil, errors.New("小爱音箱都连接不上，还通知个得儿！")
}

// 定时解析并更新线上服务器列表
func parseServerList() {
	jx3osm.GLO_SERVER_MAPS.ParseServerList()
}

package jx3osm

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gocarina/gocsv"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// 实现 Server 的 testNetConnection 接口能力
func (s *Server) testNetConnection(typ string) bool {
	// 尝试建立网络连接，以检查端口是否打开
	addr := fmt.Sprintf("%s:%d", s.ServIP, s.ServPort)
	conn, err := net.DialTimeout(typ, addr, time.Second*1)
	if err != nil {
		return false
	}

	defer conn.Close()
	return true
}

// 通过约束实现 Server 的 TestTcpConnection 接口
func testTcpConnection[T ServerCaps](s T) bool {
	return s.testNetConnection("tcp")
}

// 通过约束实现 Server 的 TestUdpConnection 接口
func testUdpConnection[T ServerCaps](s T) bool {
	return s.testNetConnection("udp")
}

// 实现 ServerMaps 的 downloadServerList 接口能力
// 下载服务器列表文本(通常情况下为 GB2312 编码格式)，并返回 UTF-8 编码的文本
func (smaps *ServerMaps) downloadServerList() string {
	// 定义 http client
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://jx3comm.xoyocdn.com/jx3hd/zhcn_hd/serverlist/serverlist.ini", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-GB;q=0.8,en-US;q=0.7,en;q=0.6")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36 Edg/130.0.0.0")
	req.Header.Set("sec-ch-ua-mobile", "?0")

	// 发送 http request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 解析 body 文本
	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyText := string(bodyData)

	// 尝试将编码格式从 GB2312 转为 UTF-8
	if !utf8.Valid(bodyData) {
		bodyText, err = simplifiedchinese.GBK.NewDecoder().String(bodyText)
		if err != nil {
			log.Fatal(err)
		}
	}

	return bodyText
}

// 解析 UTF-8 编码格式的服务器列表文本数据为 map[string]Server 结构体映射表
func (smaps ServerMaps) parseServerList(svrs_text string) {

	// 定义 CSV 解析器并指定分隔符
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '\t'
		return r
	})

	// 解码 CSV 并转换为 Server 结构体切片
	var servers []Server
	if err := gocsv.UnmarshalWithoutHeaders(strings.NewReader(svrs_text), &servers); err != nil {
		log.Fatal(err)
	}

	// 整理为映射表
	for _, svr := range servers {
		smaps[svr.ServName] = svr
	}
}

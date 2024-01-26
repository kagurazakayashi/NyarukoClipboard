package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
)

var (
	confServer string
	confClient string
	noSend     bool
	noReceive  bool
	refresh    int

	clipboardNow string
	isServer     bool = false
	protocol     string
	address      string
	conn         net.Conn
)

func init() {
	flag.StringVar(&confServer, "s", "tcp://:5888", "服务器模式，作为服务器连接的地址")
	flag.StringVar(&confClient, "c", "", "客户端模式，作为客户端连接的地址")
	flag.BoolVar(&noSend, "ns", false, "禁止发送")
	flag.BoolVar(&noReceive, "nr", false, "禁止接收")
	flag.IntVar(&refresh, "r", 1000, "剪贴板检查间隔（毫秒）")
}

func main() {
	fmt.Println("NyarukoClipboard v1.0.0")
	flag.Parse()
	if len(confClient) > 0 {
		fmt.Println("客户端模式: " + confClient)
		protocolAndAddress(confClient)
		client()
	} else if len(confServer) > 0 {
		isServer = true
		fmt.Println("服务器模式: " + confServer)
		protocolAndAddress(confServer)
		server()
	} else {
		fmt.Println("请指定服务器模式还是客户端模式")
		return
	}
	fmt.Println("剪贴板监控开始")
	go clipboardMonitoring()

	var signalch chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(signalch, os.Interrupt)
	var signal os.Signal = <-signalch
	conn.Close()
	fmt.Println("程序退出。", signal)
}

func protocolAndAddress(uri string) {
	var uriArr []string = strings.Split(uri, "://")
	protocol = uriArr[0]
	address = uriArr[1]
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

var (
	confServer string
	confClient string
	certFile   string
	keyFile    string
	noSend     bool = false
	noReceive  bool = false
	refresh    int

	clipboardNow []byte
	isServer     bool = false
	protocol     string
	address      string
	conn         net.Conn
	verbose      bool
	cdTime       int64
)

const bufSize = 1048576

func init() {
	flag.StringVar(&confServer, "s", "tcp://:7976", "服务器模式，作为服务器连接的地址")
	flag.StringVar(&confClient, "c", "", "客户端模式，作为客户端连接的地址")
	flag.StringVar(&certFile, "e", "", "证书文件")
	flag.StringVar(&keyFile, "k", "", "密钥文件")

	flag.BoolVar(&noSend, "ns", false, "禁止发送")
	flag.BoolVar(&noReceive, "nr", false, "禁止接收")
	flag.IntVar(&refresh, "r", 1000, "剪贴板检查间隔（毫秒）")
	flag.Int64Var(&cdTime, "cd", 2000, "剪贴板读取冷却时间（毫秒）")
	flag.BoolVar(&verbose, "v", false, "显示调试信息")
}

func main() {
	log.Println("NyarukoClipboard v1.0.0")
	flag.Parse()
	if len(confClient) > 0 {
		log.Println("客户端模式: " + confClient)
		protocolAndAddress(confClient)
		client()
	} else if len(confServer) > 0 {
		isServer = true
		log.Println("服务器模式: " + confServer)
		protocolAndAddress(confServer)
		server()
	} else {
		fmt.Println("请指定服务器模式还是客户端模式")
		return
	}

	var signalch chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(signalch, os.Interrupt)
	var signal os.Signal = <-signalch
	conn.Close()
	log.Println("程序退出。", signal)
}

func protocolAndAddress(uri string) {
	var uriArr []string = strings.Split(uri, "://")
	if len(uriArr) < 2 {
		protocol = "tcp"
		address = uriArr[0]
	} else {
		protocol = uriArr[0]
		address = uriArr[1]
	}
}

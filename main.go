package main

import (
	"flag"
)

var (
	confServer string
	confClient string
	noSend     bool
	noReceive  bool
)

func init() {
	flag.StringVar(&confServer, "s", "tcp://127.0.0.1:5888", "服务器模式，作为服务器连接的地址")
	flag.StringVar(&confClient, "c", "tcp://127.0.0.1:5888", "客户端模式，作为客户端连接的地址")
	flag.BoolVar(&noSend, "ns", false, "禁止发送")
	flag.BoolVar(&noReceive, "nr", false, "禁止接收")
}

func main() {
	flag.Parse()
}

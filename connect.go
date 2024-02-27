package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func dataProcess() {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		var buf [bufSize]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			running = false
			if isServer {
				fmt.Printf("客户端断开连接: %v\n", err)
			} else {
				fmt.Printf("服务端断开连接: %v\n", err)
				time.Sleep(3 * time.Second)
				fmt.Println("重试连接到服务端")
				client()
			}
			break
		}
		recv := string(buf[:n])
		clipboardPaste(recv)
	}
}

func serverSend(bytes []byte) {
	if conn == nil {
		return
	}
	i, err := conn.Write(truncateBytes(bytes))
	if err != nil {
		log.Println("发送失败: ", err)
		return
	}
	log.Println("发送成功: ", i)
}

func truncateBytes(b []byte) []byte {
	if len(b) <= bufSize {
		return b
	}
	return b[:bufSize]
}

func server() {
	log.Println("正在等待连接")
	listen, err := net.Listen(protocol, address)
	if err != nil {
		log.Printf("打开端口失败: %v\n", err)
		return
	}
	for {
		conn, err = listen.Accept()
		if err != nil {
			log.Printf("客户端断开连接: %v\n", err)
			continue
		}
		log.Println("客户端接入成功: " + conn.RemoteAddr().String())
		running = true
		go clipboardMonitoring()
		go dataProcess()
	}
}

func client() {
	log.Println("正在连接服务端")
	listen, err := net.Dial(protocol, address)
	if err != nil {
		log.Printf("未能连接到服务端: %v\n", err)
		time.Sleep(3 * time.Second)
		log.Println("重试连接到服务端")
		client()
		return
	}
	conn = listen
	log.Println("连接到服务端成功: " + conn.RemoteAddr().String())
	running = true
	go clipboardMonitoring()
	go dataProcess()
}

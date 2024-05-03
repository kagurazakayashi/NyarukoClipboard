package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func dataProcess() {
	// defer conn.Close()
	var nbuf []byte = []byte{}
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
		if verbose {
			log.Println("接收到数据: ", string(buf[0]), len(buf[:n]), buf[len(buf)-1])
		}
		if n == 0 {
			continue
		}
		if nbuf == nil {
			nbuf = buf[:n]
			var allLen = len(nbuf)
			log.Println("接收到全部数据: ", string(nbuf[0]), len(nbuf[:allLen]))
			clipboardPaste(nbuf[:allLen])
			nbuf = []byte{}
		} else {
			nbuf = append(nbuf, buf[:n]...)
		}
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
	if verbose {
		log.Println("发送成功: ", string(bytes[0]), i)
	}
}

func truncateBytes(b []byte) []byte {
	if len(b) <= bufSize {
		return b
	}
	return b[:bufSize]
}

func server() {
	log.Println("正在等待连接")
	if len(certFile) > 0 && len(keyFile) > 0 {
		serverS()
		return
	}
	listen, err := net.Listen(protocol, address)
	if err != nil {
		log.Printf("错误: 监听端口失败: %v\n", err)
		return
	}
	for {
		conn, err = listen.Accept()
		if err != nil {
			log.Printf("客户端断开连接: %v\n", err)
			continue
		}
		log.Println("客户端接入成功: " + conn.RemoteAddr().String())
		// defer conn.Close()
		running = true
		go clipboardMonitoring()
		go dataProcess()
	}
}
func serverS() {
	cer, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Printf("错误: 未能加载加密证书: %v\n", err)
		return
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	listen, err := tls.Listen(protocol, address, config)
	if err != nil {
		log.Printf("错误: 监听端口失败: %v\n", err)
		return
	}
	// defer listen.Close()
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
	if len(certFile) > 0 {
		clientS()
		return
	}
	listen, err := net.Dial(protocol, address)
	if err != nil {
		log.Printf("未能连接到服务端: %v\n", err)
		time.Sleep(3 * time.Second)
		log.Println("重试连接到服务端")
		client()
		return
	}
	conn = listen
	// defer conn.Close()
	log.Println("连接到服务端成功: " + conn.RemoteAddr().String())
	running = true
	go clipboardMonitoring()
	go dataProcess()
}
func clientS() {
	cert, err := os.ReadFile(certFile)
	if err != nil {
		log.Printf("错误: 未能加载加密证书: %v\n", err)
		return
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)

	conf := &tls.Config{
		RootCAs: certPool,
	}
	listen, err := tls.Dial(protocol, address, conf)
	conn = listen
	if err != nil {
		log.Printf("未能连接到服务端: %v\n", err)
		time.Sleep(3 * time.Second)
		log.Println("重试连接到服务端")
		client()
		return
	}
	// defer conn.Close()
	log.Println("连接到服务端成功: " + conn.RemoteAddr().String())
	running = true
	go clipboardMonitoring()
	go dataProcess()
}

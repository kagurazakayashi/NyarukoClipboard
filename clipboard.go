package main

import (
	"log"
	"time"

	"github.com/atotto/clipboard"
)

var running bool = false

func clipboardMonitoring() {
	if running {
		log.Println("剪贴板监控开始")
	} else {
		log.Println("剪贴板监控结束")
	}
	for {
		if running {
			time.Sleep(time.Duration(refresh) * time.Millisecond)
			clipboardCopy()
		} else {
			log.Println("剪贴板监控结束")
			break
		}
	}
}

func clipboardCopy() bool {
	clipboardContent, err := clipboard.ReadAll()
	if verbose {
		log.Printf("剪贴板内容: %s\n", clipboardContent)
		log.Printf("上次的内容: %s\n", clipboardNow)
	}
	if err != nil {
		if len(clipboardNow) > 0 {
			clipboardNow = ""
		}
		if verbose {
			log.Println("剪贴板读取失败", err)
		}
		return false
	}
	if clipboardNow == clipboardContent {
		log.Println("与上次内容相同。")
		return false
	}
	log.Printf("#<- %s\n", clipboardContent)
	clipboardNow = clipboardContent
	if !noSend {
		serverSend(clipboardContent)
	}
	return true
}

func clipboardPaste(text string) {
	if len(text) == 0 || clipboardNow == text {
		return
	}
	log.Printf("#-> %s\n", text)
	if noReceive {
		return
	}
	var err error = clipboard.WriteAll(text)
	if err != nil {
		log.Println("剪贴板写入失败")
		return
	}
}

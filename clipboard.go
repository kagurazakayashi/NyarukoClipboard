package main

import (
	"bytes"
	"io"
	"log"
	"strconv"
	"time"

	clipboardText "github.com/atotto/clipboard"
	clipboardImage "github.com/skanehira/clipboard-image/v2"
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

func typeArr(bytes []byte) (byte, []byte) {
	return bytes[0], bytes[1:]
}

func viewData(bytes []byte) string {
	dataType, data := typeArr(bytes)
	if len(data) == 0 {
		return "无数据"
	}
	dataType, data = typeArr(data)
	switch dataType {
	case 'T':
		return "文本: " + string(data)
	case 'I':
		return "图片: " + strconv.Itoa(len(data))
	}
	return "无法识别的数据"
}

func clipboardCopy() bool {
	clipboardTextContent, err := clipboardText.ReadAll()
	var clipboardContent []byte = []byte(clipboardTextContent)
	// if verbose {
	// 	log.Printf("剪贴板内容: %s\n", clipboardContent)
	// 	log.Printf("上次的内容: %s\n", clipboardNow)
	// }
	if err != nil {
		if verbose {
			log.Println("使用文本格式读取剪贴板失败", err)
		}
		// 嘗試使用圖片格式讀取
		reader, err := clipboardImage.Read()
		if err != nil {
			log.Println("使用图片格式读取剪贴板失败", err)
			return false
		}
		byteData, err := io.ReadAll(reader)
		if err != nil || len(byteData) == 0 {
			log.Println("使用图片格式读取剪贴板失败", err)
			return false
		}
		if len(clipboardTextContent) > 0 {
			clipboardContent = append([]byte{byte('I')}, byteData...)
		}
	} else if len(clipboardTextContent) > 0 {
		clipboardContent = append([]byte{byte('T')}, clipboardContent...)
	}
	if bytes.Equal(clipboardNow, clipboardContent) {
		log.Println("与上次内容相同。")
		return false
	}
	log.Printf("#<- %s\n", viewData(clipboardContent))
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
	var err error = clipboardText.WriteAll(text)
	if err != nil {
		log.Println("剪贴板写入失败")
		return
	}
}

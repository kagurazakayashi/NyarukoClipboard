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

type datatype byte

const (
	None  datatype = '0'
	Text  datatype = 'T'
	Image datatype = 'I'
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

func typeArr(bytes []byte) (datatype, []byte) {
	if len(bytes) < 2 {
		return None, []byte{}
	}
	return datatype(bytes[0]), bytes[1:]
}

func viewData(bytes []byte) string {
	format, data := typeArr(bytes)
	if len(data) == 0 {
		return "无数据"
	}
	switch format {
	case Text:
		return "文本: " + string(data)
	case Image:
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
		// if verbose {
		// 	log.Println("使用文本格式读取剪贴板失败", err)
		// }
		// 嘗試使用圖片格式讀取
		reader, err := clipboardImage.Read()
		if err != nil {
			// if verbose {
			// 	log.Println("使用图片格式读取剪贴板失败", err)
			// }
			return false
		}
		byteData, err := io.ReadAll(reader)
		if err != nil || len(byteData) == 0 {
			if verbose {
				log.Println("读取剪贴板失败", err)
			}
			return false
		}
		if len(clipboardTextContent) > 0 || len(byteData) > 0 {
			clipboardContent = append([]byte{byte('I')}, byteData...)
		}
	} else if len(clipboardTextContent) > 0 {
		clipboardContent = append([]byte{byte('T')}, clipboardContent...)
	}
	if bytes.Equal(clipboardNow, clipboardContent) {
		// log.Println("与上次内容相同。")
		return false
	}
	log.Printf("#<- %s\n", viewData(clipboardContent))
	clipboardNow = clipboardContent
	if !noSend {
		serverSend(clipboardContent)
	}
	return true
}

func clipboardPaste(data []byte) {
	if len(data) == 0 || bytes.Equal(clipboardNow, data) {
		return
	}
	log.Printf("#-> %s\n", viewData(data))
	if noReceive {
		return
	}
	dataType, data := typeArr(data)
	var err error
	switch dataType {
	case Text:
		err = clipboardText.WriteAll(string(data))
	case Image:
		var reader *bytes.Reader = bytes.NewReader(data)
		err = clipboardImage.Write(reader)
	}
	if err != nil {
		log.Println("剪贴板写入失败:", err)
		return
	}
}

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"
	"unicode/utf8"

	clipboardText "github.com/atotto/clipboard"
	clipboardImage "github.com/skanehira/clipboard-image/v2"
)

type datatype byte

const (
	None  datatype = '0'
	Text  datatype = 'T'
	Image datatype = 'I'
)

var (
	running   bool        = false
	skipOne   bool        = false
	skipTimer *time.Timer = nil
)

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
	// log.Println("===数据格式：", string(format), "数据长度：", len(data))
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
	if !utf8.Valid(clipboardContent) {
		err = fmt.Errorf("NO UTF-8")
		clipboardContent = []byte{}
	}
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
			if verbose && err != nil {
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
	// log.Println("===发送前比较：", string(clipboardNow), "?=", string(clipboardContent))
	if bytes.Equal(clipboardNow, clipboardContent) {
		return false
	}
	// log.Println("===发送前比较：", len(clipboardNow), "?=", len(clipboardContent))
	clipboardNow = clipboardContent
	if skipOne {
		skipOne = false
		if skipTimer != nil {
			skipTimer.Stop()
			skipTimer = nil
		}
		if verbose {
			log.Println("跳过一次")
		}
		return false
	}
	log.Printf("[发送] %s\n", viewData(clipboardContent))
	if !noSend {
		serverSend(append(clipboardContent, byte(0), byte(0)))
	}
	return true
}

func clipboardPaste(data []byte) {
	// log.Println("===接收前比较：", string(clipboardNow), "?=", string(data))
	// if len(data) == 0 || bytes.Equal(clipboardNow, data) {
	// 	return
	// }
	clipboardNow = data
	if !skipOne {
		if skipTimer != nil {
			skipTimer.Stop()
		}
		skipOne = true
		skipTimer = time.AfterFunc(time.Duration(cdTime)*time.Millisecond, func() {
			skipOne = false
			skipTimer.Stop()
			skipTimer = nil
		})
	}
	log.Printf("[接收] %s\n", viewData(data))
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

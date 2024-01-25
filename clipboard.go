package main

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
)

func clipboardMonitoring() {
	for {
		time.Sleep(time.Duration(refresh) * time.Millisecond)
		clipboardCopy()
	}
}

func clipboardCopy() bool {
	clipboardContent, err := clipboard.ReadAll()
	if err != nil {
		if len(clipboardNow) > 0 {
			clipboardNow = ""
		}
		return false
	}
	if clipboardNow == clipboardContent {
		return false
	}
	fmt.Printf("#<- %s\n", clipboardContent)
	clipboardNow = clipboardContent
	serverSend(clipboardContent)
	return true
}

func clipboardPaste(text string) {
	if len(text) == 0 || clipboardNow == text {
		return
	}
	fmt.Printf("#-> %s\n", text)
	var err error = clipboard.WriteAll(text)
	if err != nil {
		fmt.Println("剪贴板写入失败")
		return
	}
}

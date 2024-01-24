package main

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
)

func clipboardMonitoring() {
	for {
		time.Sleep(time.Duration(refresh) * time.Millisecond)
		clipboardContent, err := clipboard.ReadAll()
		if err != nil {
			if clipboardNow != "" {
				clipboardNow = ""
			}
			continue
		}
		if clipboardNow == clipboardContent {
			continue
		}
		fmt.Printf("#<- %s\n", clipboardContent)
		clipboardNow = clipboardContent
	}
}

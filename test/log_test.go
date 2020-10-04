package test

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestWriteLog(t *testing.T) {
	file, err := os.OpenFile("/Users/xcz/go/src/xcz/gdm/log/2020-10-2.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("openfile is fail for ", err)
		return
	}
	logger := log.New(file, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
	logger.Println("ceshi")
	_ = file.Close()
}

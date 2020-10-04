package test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"xcz/gdm/utils"
)

func TestGetCurrentPath(t *testing.T) {
	path, _ := os.Getwd()
	fmt.Println(path)
	absolutePath, err := utils.GetParentAbsolutePath()
	if err != nil {
		fmt.Println("err: ", err)
	}
	fmt.Println(absolutePath)
}

func TestGetParentPath(t *testing.T) {
	path := "/Users/xcz/go/src/xcz/gdm/test"
	index := strings.LastIndex(path, "/")
	var parentPath string
	//如果索引不存在
	if index == -1 {
		log.Println("index is -1")
		return
	}
	//根目录下
	prefixPath := path[:index]
	log.Println("prefixPath: ", prefixPath)
	parentPath = prefixPath + "/"
	log.Println("parentPath: ", parentPath)
	filePath := path[index+1:]
	log.Println("filePath: ", filePath)
}

func TestGetFileNameFromUrl(t *testing.T) {
	url := "http://wppkg.baidupcs.com/issue/netdisk/yunguanjia/BaiduYunGuanjia_7.0.1.1.exe"
	lastIndex := strings.LastIndex(url, "/")
	fileName := url[lastIndex+1:]
	fmt.Println(fileName)
}

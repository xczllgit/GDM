package main

import (
	"flag"
	"fmt"
	"os"
	"xcz/gdm/genesis"
	"xcz/gdm/http"
	"xcz/gdm/logs"
	"xcz/gdm/utils"
)

var targetUrl = flag.String("url", "https://iterm2.com/downloads/stable/iTerm2-3_3_12.zip", "Input Your Resource URL Here")
var threadNum = flag.Int64("threadNum", 6, "Input Thread Nums Here Which You Want To Use For Downloading")
var localAddress = flag.String("address", "", "Input Your Local Address To Storage Resource")

func init() {
	//读取配置文件
	genesis.InitConfig()
	//初始化配置变量
	genesis.InitConfigInfo()

}

func main() {
	//清除30天前的日志文件
	go utils.DeleteLogFileBefore(30)
	//解析命令行参数
	flag.Parse()
	//校验url的正确性
	isRight, err, head := utils.VerifyUrl(*targetUrl)
	defer logs.CloseLogFile()
	if err != nil {
		genesis.Logger.Fatal(utils.GetRunFuncName(), "ParseUrl encounter some errors: ", err.Error())
	}
	if !isRight {
		genesis.Logger.Fatal("ParseUrl have a illegal URL for download, url is : ", *targetUrl)
	}
	if isRight {
		genesis.Logger.Println("VerifyUrl is success, downloading will begin in a moment")
	}
	//解析线程数量的正确性
	//解析本地存储地址的正确性
	local, err := os.Stat(*localAddress)
	if err != nil {
		fmt.Println("sorry, the selected target address encounter some errors: ", err)
		genesis.Logger.Fatal("target Address encounter some error: ", err)
	}
	if !local.IsDir() {
		fmt.Println("sorry, the selected target address is not a directory")
		genesis.Logger.Fatal("sorry, the selected target address is not a directory")
	}
	http.DownloadFromUrl(*targetUrl, *threadNum, *localAddress, head)
}

package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"xcz/gdm/genesis"
	"xcz/gdm/utils"
)

var (
	group sync.WaitGroup
)

//下载url资源，指定线程数量与目标存储路径
func DownloadFromUrl(URL string, threadNum int64, localAddress string, resourceHead *utils.ResourceHead) {
	if !resourceHead.SupportMultiThread {
		threadNum = 1
	}
	subfiles := SplitResourceByThread(resourceHead.ContentLength, threadNum)
	fmt.Println("测试埋点")
	if subfiles == nil || len(subfiles) == 0 {
		genesis.Logger.Fatal("SplitResourceByThread failed")
	}
	//资源不支持多线程
	if !resourceHead.SupportMultiThread || threadNum == 1 {
		group.Add(1)
		_, err := downLoadBySingleThread(URL, localAddress, subfiles[0], 0, resourceHead)
		group.Wait()
		if err != nil {
			genesis.Logger.Fatal("[downLoadBySingleThread] fail, encounter some errors: ", err)
		}
		//下载完成，将对应子文件改名为目标文件
		targetFileName := parseFileNameFromUrl(URL)
		err = os.Rename(utils.GetFileAbsolutePath(localAddress, subfiles[0].tempFileName), utils.GetFileAbsolutePath(localAddress, targetFileName))
		if err != nil {
			genesis.Logger.Println("download completely, but rename subfile to targetfile failed, errors: ", err)
		}
		fmt.Println("Download completely")
		return
	}
	//TODO 按照多线程模式下载
	downLoadByMulThread(URL, threadNum, localAddress, subfiles, 0, resourceHead)
	return
}

//单线程下载
func downLoadBySingleThread(URL string, localAddress string, subFile *SubFile, retryCount int64, resourceHead *utils.ResourceHead) (*os.File, error) {
	//创建临时文件，避免进度条出现panic
	file, err := utils.CreateFile(localAddress, subFile.tempFileName)
	if err != nil {
		genesis.Logger.Fatal("localAddress is an invalid directory path, error info is: ", err)
	}
	//1、生成http request请求
	request, err := http.NewRequest("GET", URL, nil)
	if err != nil || request == nil {
		if err != nil {
			genesis.Logger.Println("[downLoadBySingleThread] initiating a http request fail, errro info: ", err, " and retryCount is ", retryCount)
		} else {
			genesis.Logger.Println("[downLoadBySingleThread] initiating a http request fail, and retryCount is ", retryCount)
		}
		if retryCount >= genesis.RetryCount {
			genesis.Logger.Fatal("sorry, [downLoadBySingleThread] retry count exceeds our retry value, then our program will exit")
		}
		return downLoadBySingleThread(URL, localAddress, subFile, retryCount+1, resourceHead)
	}

	//指定本次http请求的数据范围,用于断点续传
	ranges := fmt.Sprintf("bytes=%d-%d", subFile.startIndex, subFile.endIndex)
	request.Header.Set("Range", ranges)
	//传输完成后，断开http连接
	request.Header.Set("Connection", "close")

	//2、执行request请求，获取响应数据
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		genesis.Logger.Println("[downLoadBySingleThread] http connection fail, errro info: ", err, " and retryCount is ", retryCount)
		if retryCount >= genesis.RetryCount {
			genesis.Logger.Fatal("sorry, [downLoadBySingleThread] retry count exceeds our retry value, then our program will exit")
		}
		return downLoadBySingleThread(URL, localAddress, subFile, retryCount+1, resourceHead)
	}
	//http返回的response必须关闭，否则会造成内存泄漏
	defer response.Body.Close()

	//3、将数据写入到临时文件中
	length, err := io.Copy(file, response.Body)
	if err != nil || length < (subFile.endIndex-subFile.startIndex+1) {
		if err != nil {
			genesis.Logger.Println("[downLoadBySingleThread] io copy fail , errro info: ", err, " and retryCount is ", retryCount)
		} else {
			genesis.Logger.Println("[downLoadBySingleThread] data content is not enough, and retryCount is ", retryCount)
		}
		if retryCount >= genesis.RetryCount {
			genesis.Logger.Fatal("sorry, [downLoadBySingleThread] retry count exceeds our retry value, then our program will exit")
		}
		return downLoadBySingleThread(URL, localAddress, subFile, retryCount+1, resourceHead)
	}
	group.Done()
	return file, nil
}

func downLoadByMulThread(URL string, ThreadNum int64, localAddress string, subFiles []*SubFile, retryCount int64, resourceHead *utils.ResourceHead) {
	//根据子文件并发下载资源的各个部分
	for _, subFile := range subFiles {
		group.Add(1)
		go downLoadBySingleThread(URL, localAddress, subFile, retryCount, resourceHead)
	}
	time.Sleep(2 * time.Second)
	//显示进度条
	go showDownloadProgressBar(localAddress, subFiles, resourceHead)
	group.Wait()
	//等待子文件下载完成，合并各个子文件到目标文件
	fmt.Println("Download completely, wait for merging subFile")
	targetFileName := parseFileNameFromUrl(URL)
	targetFile, err := utils.CreateFile(localAddress, targetFileName)
	if err != nil {
		//清除子文件
		removeSubFile(localAddress, subFiles)
		fmt.Println("sorry, create ", targetFileName, " fail, errors info: ", err)
		genesis.Logger.Fatal("sorry, create ", targetFileName, " fail, errors info: ", err)
	}
	defer targetFile.Close()
	for _, subFile := range subFiles {
		//寻找目标文件末尾索引值（即下一次需要写入的起始索引）
		endIndex, _ := targetFile.Seek(0, os.SEEK_END)
		//将子文件内容转换为byte数组，方便写入目标文件中
		subContent := getFileContent(localAddress, subFile.tempFileName)
		_, err := targetFile.WriteAt(subContent, endIndex)
		if err != nil {
			//清除子文件
			removeSubFile(localAddress, subFiles)
			fmt.Println("sorry, merge ", targetFileName, " fail, errors info: ", err)
			genesis.Logger.Fatal("sorry, merge ", targetFileName, " fail, errors info: ", err)
		}
	}
	//清除子文件
	removeSubFile(localAddress, subFiles)
	fmt.Println("success! download completely")
	genesis.Logger.Println("success! download completely")
}

//展示下载进度条
func showDownloadProgressBar(localAddress string, subFiles []*SubFile, resourceHead *utils.ResourceHead) {
	for {
		var currentSize float64 = 0
		for _, subFile := range subFiles {
			filePath := utils.GetFileAbsolutePath(localAddress, subFile.tempFileName)
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				genesis.Logger.Println("[showDownloadProgressBar].Stat encounter some errors: ", err)
			}
			if fileInfo == nil {
				genesis.Logger.Println("[showDownloadProgressBar].Stat encounter some errors: ", err)
			}
			currentSize += float64(fileInfo.Size())
		}
		pert := currentSize / float64(resourceHead.ContentLength)
		pertStr := strconv.FormatFloat(pert*100, 'f', 2, 64)
		bars := showBar(int(pert * 100))
		showStr := "Downloading: [" + bars + "]  " + pertStr + "%"
		//加上"\r%s"，就能让进度条保持一行，为什么？
		fmt.Printf("\r%s", showStr)
		if pert >= 1 {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func showBar(num int) string {
	result := ""
	for i := 1; i <= num; i++ {
		result = result + "="
	}
	return result
}

//展示网速
func showNetSpeed() {
	//1、获取当前子文件大小、一秒前子文件大小，相减即为网速
}

func removeSubFile(localAddress string, subFiles []*SubFile) {
	for _, subFile := range subFiles {
		filePath := utils.GetFileAbsolutePath(localAddress, subFile.tempFileName)
		err := os.Remove(filePath)
		if err != nil {
			genesis.Logger.Println("remove subFile encounter some errors: ", err)
		}
	}
}

func getFileContent(localAddress string, fileName string) []byte {
	filePath := utils.GetFileAbsolutePath(localAddress, fileName)
	subfileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		genesis.Logger.Fatal("Read subfile fail, filePath: ", filePath, " ,errors info: ", err)
	}
	return subfileContent
}

func parseFileNameFromUrl(url string) string {
	lastIndex := strings.LastIndex(url, "/")
	fileName := url[lastIndex+1:]
	return fileName
}

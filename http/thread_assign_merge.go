package http

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"os"
	"time"
	"xcz/gdm/genesis"
)

/**
此文件主要包含函数： 根据多线程数量切分下载资源到各个子文件、合并各个子文件为最终文件
*/

type SubFile struct {
	tempFileName       string   //子文件名称
	startIndex         int64    //下载资源的开始范围
	endIndex           int64    //下载资源的结束范围
	previousSubFile    *SubFile //前一个子文件地址
	subsequenceSubFile *SubFile //后一个子文件地址
}

//为多线程切分对应数据范围的子文件
func SplitResourceByThread(contentLength int64, threadNum int64) []*SubFile {
	subFiles := []*SubFile{}
	fileNames := generateTempFileNames(threadNum, 0)
	if threadNum == 1 {
		subfile := &SubFile{
			tempFileName: fileNames[0],
			startIndex:   0,
			endIndex:     contentLength - 1,
		}
		subFiles = append(subFiles, subfile)
		return subFiles
	}
	//每个线程平分的数据量
	avgCap := contentLength / threadNum
	var index int64 = 1
	subFile1 := &SubFile{
		tempFileName: fileNames[0],
		startIndex:   0,
		endIndex:     avgCap - 1,
	}
	subFiles = append(subFiles, subFile1)
	index++
	for ; index <= threadNum; index++ {
		subfile := &SubFile{
			tempFileName:    fileNames[index-1],
			startIndex:      (index - 1) * avgCap,
			previousSubFile: subFiles[len(subFiles)-1],
		}
		if index == threadNum {
			subfile.endIndex = contentLength - 1
		} else {
			subfile.endIndex = (index)*avgCap - 1
		}
		subFiles[len(subFiles)-1].subsequenceSubFile = subfile
		subFiles = append(subFiles, subfile)
	}
	return subFiles
}

//为子文件生成随机名称
func generateTempFileNames(threadNum int64, retryCount int64) []string {
	prefix := genesis.SubFilePrefix
	var result []string
	tempMap := make(map[string]int64)
	var i int64
	for i = 1; i <= threadNum; i++ {
		randString := prefix + generateRandString()
		tempMap[randString] = i
		result = append(result, randString)
	}
	if len(tempMap) == len(result) && len(result) == int(threadNum) {
		return result
	} else {
		if retryCount >= genesis.RetryCount {
			fmt.Println("generateTempFileNames fail")
			genesis.Logger.Fatal("generateTempFileNames fail")
		}
		return generateTempFileNames(threadNum, retryCount+1)
	}
}

//从26个字母中随机取出5个，使用md5加密
func generateRandString() string {
	str := "abcdefghijklmnopqrstuvwxyz"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s1 := r.Int63n(26)
	s2 := r.Int63n(26)
	s3 := r.Int63n(26)
	s4 := r.Int63n(26)
	s5 := r.Int63n(26)
	tempString := str[s1:s1+1] + str[s2:s2+1] + str[s3:s3+1] + str[s4:s4+1] + str[s5:s5+1]
	data := []byte(tempString)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

//合并子文件到targetFile中
func MergeMulSubFile(subFiles []*SubFile, targetFile *os.File, localAddress string) {
	for _, subFile := range subFiles {
		//寻找目标文件末尾索引值（即下一次需要写入的起始索引）
		endIndex, _ := targetFile.Seek(0, os.SEEK_END)
		//将子文件内容转换为byte数组，方便写入目标文件中
		subContent := getFileContent(localAddress, subFile.tempFileName)
		_, err := targetFile.WriteAt(subContent, endIndex)
		if err != nil {
			//清除子文件
			removeSubFile(localAddress, subFiles)
			fmt.Println("sorry, merge ", targetFile.Name(), " fail, errors info: ", err)
			genesis.Logger.Fatal("sorry, merge ", targetFile.Name(), " fail, errors info: ", err)
		}
	}
}

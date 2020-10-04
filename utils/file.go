package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//获取当前路径
func GetAbsolutePath() string {
	path, _ := os.Getwd()
	return path
}

//目前是针对linux和macos获取上一级绝对路径
//比较奇怪，main方法调用此方法或者main直接调用os.Getwd()都只能打印出项目根路径，估计与golang的编译后运行文件位置有关
func GetParentAbsolutePath() (string, error) {
	path, _ := os.Getwd()
	s := path[len(path)-1:]
	var index int
	if s == "/" {
		path = path[:len(path)]
		index = strings.LastIndex(path, "/")
	} else {
		index = strings.LastIndex(path, "/")
	}
	var parentPath string
	//如果索引不存在
	if index == -1 {
		err := errors.New("[GetParentAbsolutePath] fail, cannot get right index for '/' ")
		return parentPath, err
	}
	//获取 path[0,index)
	prefixPath := path[:index]
	parentPath = prefixPath + "/"
	return parentPath, nil
}

func GetProjectRootPath() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		log.Println("[OpenLogFile].[GetProjectRootPath] encounter some errors: ", err)
		return "", err
	}
	s := path[len(path)-1:]
	if s == "/" {
		return path, nil
	}
	return path + "/", nil
}

//获取当前天数，生成日志文件名,文件位于当前方法上级目录/log下
func OpenLogFile() (*os.File, error) {
	year, month, day := time.Now().Date()
	var fileName string
	fileName = strconv.FormatInt(int64(year), 10) + "-" + strconv.FormatInt(int64(month), 10) +
		"-" + strconv.FormatInt(int64(day), 10) + ".log"

	absolutePath := getLogFileDirectoryPath()
	if !VerifyFileOrDirectoryIsExist(absolutePath) {
		err := os.Mkdir(absolutePath, os.ModePerm)
		if err != nil {
			log.Println("[OpenLogFile].[Mkdir] encounter some errors: ", err)
			os.Exit(-1)
			return nil, err
		}
	}
	absoluteFilePath := absolutePath + "/" + fileName
	fmt.Println("logFile path: ", absoluteFilePath)
	//0666代表读写权限，0777代表读写+执行
	file, err := os.OpenFile(absoluteFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("[OpenLogFile].[OpenFile] encounter some errors: ", err)
		os.Exit(-1)
		return nil, err
	}
	return file, nil
}

func VerifyFileOrDirectoryIsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		log.Println("VerifyFileOrDirectoryIsExist encounter some errors: ", err)
		return false
	}
	return true
}

//列出指定路径下的所有文件名（不包括子目录）
func ListFileInFolder(folder string) []string {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		fmt.Println("[ListAllFileInFolder] encounter some errors: ", err)
		return nil
	}
	var filesPath []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		absPath, err := filepath.Abs(folder + "/" + file.Name())
		if err != nil {
			fmt.Println(file.Name(), " get absPath fail for ", err)
			continue
		}
		filesPath = append(filesPath, absPath)
	}
	return filesPath
}

//列出指定路径下的所有文件名（包括子目录）
func ListAllFileInFolder(folder string) []string {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		fmt.Println("[ListAllFileInFolder] encounter some errors: ", err)
		return nil
	}
	var filesPath []string
	for _, file := range files {
		if file.IsDir() {
			inFolder := ListAllFileInFolder(folder + "/" + file.Name())
			filesPath = append(filesPath, inFolder...)
			continue
		}
		absPath, err := filepath.Abs(folder + "/" + file.Name())
		if err != nil {
			fmt.Println(file.Name(), " get absPath fail for ", err)
			continue
		}
		filesPath = append(filesPath, absPath)
	}
	return filesPath
}

//删除最近修改时间在指定日期之前的日志文件
func DeleteLogFileBefore(day int64) {
	absolutePath := getLogFileDirectoryPath()
	if !VerifyFileOrDirectoryIsExist(absolutePath) {
		return
	}
	logFiles := ListFileInFolder(absolutePath)
	now := time.Now().Unix()
	for _, file := range logFiles {
		modifyTime := GetFileLastModifyTime(file)
		if now-modifyTime >= day*3600*24 {
			err := os.Remove(file)
			if err != nil {
				fmt.Println("[DeleteLogFileBefore] Remove file ", file, "encounter some errors: ", err)
			}
		}
	}
}

//获取文件最近修改时间，如果出错，返回系统当前时间
func GetFileLastModifyTime(absoluteFilePath string) int64 {
	fileStat, err := os.Stat(absoluteFilePath)
	if err != nil {
		fmt.Println("[GetFileCreateTime] encounter some error ", err)
		return time.Now().Unix()
	}
	modTime := fileStat.ModTime().Unix()
	return modTime
}

func getLogFileDirectoryPath() string {
	projectPath, err := GetProjectRootPath()
	if err != nil {
		log.Println("[getLogFileDirectoryPath] encounter some errors: ", err)
		os.Exit(-1)
		return ""
	}
	absolutePath := projectPath + "log"
	return absolutePath
}

func CreateFile(filePath string, fileName string) (*os.File, error) {
	pathStat, err := os.Stat(filePath)
	if err != nil {
		fmt.Println(GetRunFuncName(), "target filePath encounter some errors: ", err)
		errInfo := fmt.Sprintf("target filePath encounter some errors: %+v", err)
		return nil, errors.New(errInfo)
	}
	if !pathStat.IsDir() {
		fmt.Println(GetRunFuncName(), "target filePath is not a directory")
		return nil, errors.New("target filePath is not a directory")
	}
	//这里要注意传入的文件路径是否末尾有/
	var fullFileName string
	if filePath[len(filePath)-1] == '/' {
		fullFileName = filePath + fileName
	} else {
		fullFileName = filePath + "/" + fileName
	}
	fmt.Println("target file full name is : ", fullFileName)
	_, err = os.Stat(fullFileName)
	if !os.IsNotExist(err) && err != nil {
		fmt.Println(GetRunFuncName(), "createFile fail for errors: ", err)
		os.Exit(-1)
		//return nil,errors.New("target file is already exists")
	}
	file, err := os.Create(fullFileName)
	if err != nil {
		fmt.Println(GetRunFuncName(), "create file failed, errors: ", err)
		os.Exit(-1)
	}
	return file, nil
}

//根据文件路径与文件名获取文件完整路径
func GetFileAbsolutePath(dirPath string, fileName string) string {
	//这里要注意传入的文件路径是否末尾有/
	var fullFileName string
	if dirPath[len(dirPath)-1] == '/' {
		fullFileName = dirPath + fileName
	} else {
		fullFileName = dirPath + "/" + fileName
	}
	return fullFileName
}

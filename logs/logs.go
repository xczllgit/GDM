package logs

import (
	"log"
	"os"
	"xcz/gdm/utils"
)

var logFile *os.File

func GetLogger() *log.Logger {
	var err error
	logFile, err = utils.OpenLogFile()
	if err != nil {
		log.Println("[LogInfo].[OpenLogFile] encounter some errors: ", err)
		os.Exit(-1)
	}
	logger := log.New(logFile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
	return logger
}

//打开日志文件后，不要忘记close日志文件资源
func CloseLogFile() {
	err := logFile.Close()
	if err != nil {
		logger := log.New(logFile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
		logger.Fatal("close File fail for ", err)
	}
}

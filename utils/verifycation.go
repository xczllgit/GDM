package utils

import (
	"errors"
	"log"
	"net/http"
	"runtime"
	"strings"
)

type ResourceHead struct {
	ContentLength      int64 //资源大小，Byte
	SupportMultiThread bool  //是否支持断点续传/多线程
}

//检测URL是否可达，利用curl工具
func VerifyUrl(url string) (bool, error, *ResourceHead) {
	index := strings.LastIndex(url, "/")
	if index == -1 || index == 0 || index == len(url)-1 {
		log.Println(GetRunFuncName(), " have a illegal URL for downloading ")
		return false, nil, nil
	}
	//head方法不会返回消息实体部分，用于测试url的链接是否正常比较合适
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		log.Println(GetRunFuncName(), " create http request head fail")
		return false, err, nil
	}
	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		log.Println(GetRunFuncName(), "HEAD response failed")
		return false, err, nil
	}
	if resp.StatusCode >= 400 {
		log.Println(GetRunFuncName(), "the Url status is out of order")
		err := errors.New("the Url status is out of order")
		return false, err, nil
	}
	//TODO 如何处理1xx和3xx
	resource := &ResourceHead{
		ContentLength: resp.ContentLength,
	}
	if resp.Header.Get("Accept-Ranges") != "" {
		resource.SupportMultiThread = true
	} else {
		resource.SupportMultiThread = false
	}
	return true, nil, resource
}

//获取当前函数名称
func GetRunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

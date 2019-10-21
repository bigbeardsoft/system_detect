package main

import (
	"fmt"
	"system_detect/detect/process"
	"system_detect/service"
)

/**
说明:
	第一步启动一个线程,定时采集,时间间隔从配置文件中获取.
**/
func main() {

	msg2 := service.CreateRegisterMsg("A10103")
	println(msg2)
	service.DispatherMsg(msg2)
}

func readConfig() {

}

func startDetect(f func(json string)) {

	var p process.Process
	result, err := p.GetAllProcess()
	if nil != err {
		fmt.Printf("获取进程信息发生错误:%v\n", err)
	} else {
		for _, r := range result {
			fmt.Printf("%v\n", r)
		}
	}
}

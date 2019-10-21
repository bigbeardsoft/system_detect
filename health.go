package main

import (
	"fmt"
	"system_detect/activemq"
	"system_detect/detect/process"
	system "system_detect/detect/system/linux"
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
	var p process.Process
	proc, erro := p.GetAllProcess()
	if nil != erro {
		println(erro)
	}
	println("=1==========")
	for _, r := range proc {
		fmt.Printf("%v\n", r)
	}
	println("=2==========")
	x := new(system.SysUsedInfo)
	pd, er := x.GetSystemUsedInfo()
	if nil != er {
		fmt.Printf("%v\n", er)
	}
	fmt.Printf("%v\n", pd)
	activemq.CallActiveMq()
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

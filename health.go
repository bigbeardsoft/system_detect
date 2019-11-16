package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"system_detect/service"
	"system_detect/tools"
)

type appConfig struct {
	mqIP                string
	mqPort              int
	mqUser              string
	mqPwd               string
	mqQueueName         string
	mqSendQueueName     string
	collectTimeInterval int
}

/**
说明:
	第一步启动一个线程,定时采集,时间间隔从配置文件中获取.
**/
func main() {

	ShowCmdInfo := "please in put command:\n1\tstart \n2\tstop\n"
	println(ShowCmdInfo)
	inputReader := bufio.NewReader(os.Stdin)
	b, _, _ := inputReader.ReadLine()
	cmd := string(b)
	s := new(service.CollectService)
	c := readConfig()
	mq := new(service.MQService)
	service.Init()

	for cmd != "quit" {
		if cmd == "start" {
			mq.Init(c.mqIP, strconv.Itoa(c.mqPort), c.mqUser, c.mqPwd, c.mqQueueName)
			s.Notify = func(json string) {
				mq.SendMsg(c.mqSendQueueName, json)
			}
			s.StartDetect()
		} else if cmd == "stop" {
			s.StopDetect()
		} else {
			fmt.Printf("未知命令:%s\n", cmd)
		}
		println(ShowCmdInfo)
		b, _, _ = inputReader.ReadLine()
		cmd = string(b)
	}
}

func readConfig() *appConfig {
	configInfo, err := tools.ReadConfigFile("./config.yml")
	c := new(appConfig)
	if nil != err {
		return nil
	}
	c.mqIP = configInfo["mq.ip"].(string)
	c.mqPort = configInfo["mq.port"].(int)
	c.mqUser = configInfo["mq.user"].(string)
	c.mqPwd = configInfo["mq.pwd"].(string)
	c.mqQueueName = configInfo["mq.queues"].(string)
	c.mqSendQueueName = configInfo["mq.sendqueue"].(string)
	c.collectTimeInterval = configInfo["collect.interval"].(int)
	return c
}

func testA() {

	// var m map[string]interface{}
	// m = make(map[string]interface{}, 16)
	// m["a"] = 1
	// m["b"] = "asdfas"
	// fmt.Printf("%v\n", m)
	// return
	// service.StartDetect()
	// service.StopDetect()
	// return
	// println(time.Now().String())
	// d := time.Now()
	// sd := fmt.Sprintf("%d-%02d-%02d %0d:%02d:%02d ", d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second())
	// println(sd)
	// println(fmt.Sprintf("%03d", 02))

	// return

	// msg2 := service.CreateRegisterMsg("A10103")
	// println(msg2)
	// service.DispatherMsg(msg2)
	// var p process.Process
	// proc, erro := p.GetAllProcess()
	// if nil != erro {
	// 	println(erro)
	// }
	// println("=1==========")
	// for _, r := range proc {
	// 	fmt.Printf("%v\n", r)
	// }
	// println("=2==========")
	// x := new(system.SysUsedInfo)
	// pd, er := x.GetSystemUsedInfo()
	// if nil != er {
	// 	fmt.Printf("%v\n", er)
	// }
	// fmt.Printf("%v\n", pd)
	// activemq.CallActiveMq()
}

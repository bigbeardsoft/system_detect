package main

import (
	"bufio"
	"fmt"
	"log"
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

	ShowCmdInfo := "please in put command:\nstart:开启服务 \n stop:停止服务\n"
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
			err := mq.Init(c.mqIP, strconv.Itoa(c.mqPort), c.mqUser, c.mqPwd, c.mqQueueName)
			if nil != err {
				log.Fatal(fmt.Printf("初始化mq发生异常,异常信息:\n%v", err))
			} else {
				log.Printf("连接到mq服务器[%s:%d]成功\n", c.mqIP, c.mqPort)
				s.Notify = func(json string) {
					mq.SendMsg(c.mqSendQueueName, json)
				}
				s.StartDetect()
			}
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
		fmt.Printf("读配置文件错误,错误信息:%v", err)
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

package main

import (
	"bufio"
	"os"
	"strconv"
	"system_detect/service"
	"system_detect/tools"
	"time"
)

type appConfig struct {
	mqIP                string
	mqPort              int
	mqUser              string
	mqPwd               string
	statusQueue         string
	registerQueue       string
	acceptQueues        string
	collectTimeInterval int
	clientKey           string
}

var logger tools.Logger

/**
说明:
	第一步启动一个线程,定时采集,时间间隔从配置文件中获取.
**/
func main() {
	ShowCmdInfo := "please in put command:\nstart:开启服务 \nstop:停止服务\n"
	println(ShowCmdInfo)
	inputReader := bufio.NewReader(os.Stdin)
	b, _, _ := inputReader.ReadLine()
	cmd := string(b)
	s := new(service.CollectService)
	c := readConfig()
	mq := new(service.MQService)
	service.Init()
	mq.Callback = func(msg, queue string) {
		logger.Debugf("收到:%v\n", msg)
		service.DispatherMsg(msg)
	}
	registerResult := false
	service.SetRegisterCallback(func(result bool) {
		registerResult = result
		if result {
			logger.Debugf("签到成功,tooken:%s", service.GetToken())
			s.StartDetect()
		} else {
			logger.Debugf("向服务器签到失败\n")
		}
	})
	for cmd != "quit" {
		if cmd == "start" {
			err := mq.Init(c.mqIP, strconv.Itoa(c.mqPort), c.mqUser, c.mqPwd, c.acceptQueues)
			if nil != err {
				logger.Errorf("初始化mq发生异常,异常信息:%v", err)
			} else {
				logger.Debugf("连接到mq服务器[%s:%d]成功", c.mqIP, c.mqPort)
				s.Notify = func(json string) {
					mq.SendMsg(c.statusQueue, json)
				}
				regjosn := service.CreateRegisterMsg(c.clientKey)
				mq.SendMsg(c.registerQueue, regjosn)
				go func() {
					for {
						<-time.After(time.Duration(1) * time.Minute)
						if registerResult == false {
							logger.Debugf("重发签到:%s\n", regjosn)
							mq.SendMsg(c.registerQueue, regjosn)
						} else {
							break
						}
					}
				}()
				logger.Debugf("向服务器发送签到信息:%s\n", regjosn)
			}
		} else if cmd == "stop" {
			s.StopDetect()
		} else {
			logger.Warringf("未知命令:%s\n", cmd)
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
		logger.Errorf("读配置文件错误,错误信息:%v", err)
		return nil
	}
	c.mqIP = configInfo["mq.ip"].(string)
	c.mqPort = configInfo["mq.port"].(int)
	c.mqUser = configInfo["mq.user"].(string)
	c.mqPwd = configInfo["mq.pwd"].(string)
	c.statusQueue = configInfo["mq.status_queue"].(string)
	c.registerQueue = configInfo["mq.register_queue"].(string)
	c.collectTimeInterval = configInfo["collect.interval"].(int)
	c.clientKey = configInfo["client.client_key"].(string)
	c.acceptQueues = configInfo["mq.accept_queues"].(string)
	return c
}

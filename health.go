package main

import (
	"bufio"
	"flag"
	"os"
	"strconv"
	"strings"
	"system_detect/service"
	"system_detect/tools"
	"system_detect/web"
	"time"
)

type appConfig struct {
	mqIP                string
	mqPort              int
	mqUser              string
	mqPwd               string
	statusQueue         string
	registerQueue       string
	acceptTopics        string
	collectTimeInterval int
	serviceCode         string
	showlog             bool
	webPort             int
}

var logger tools.Logger
var start string = ""
var registerResult bool = false
var stop string = ""
var quit string = ""

func init() {
	flag.StringVar(&start, "start", "start", "启动服务")
	flag.StringVar(&stop, "stop", "", "停止服务")
	flag.StringVar(&quit, "quit", "", "退出服务")
}

/**
说明:
	第一步启动一个线程,定时采集,时间间隔从配置文件中获取.
**/
func main() {

	flag.Parse()
	ShowCmdInfo := "please in put command:\n\tstart:开启服务 \n\tstop:停止服务\n\tquit:退出系统"
	inputReader := bufio.NewReader(os.Stdin)
	cmd := start
	var b []byte
	if cmd == "" {
		println(ShowCmdInfo)
		b, _, _ = inputReader.ReadLine()
		cmd = string(b)
	}
	s := new(service.CollectService)
	c := readConfig()
	mq := new(service.MQService)

	service.Init()
	mq.Callback = func(msg, queue string) {
		logger.Debugf("收到:%v\n", msg)
		service.DispatherMsg(msg)
	}
	service.SetRegisterCallback(func(result bool) {
		registerResult = result
		if result {
			logger.Debugf("签到成功,tooken:%s", service.GetToken())
			s.StartDetect()
		} else {
			logger.Debugf("向服务器签到失败\n")
		}
	})
	logger.Info("启动web服务")
	go web.StartWebServer(c.webPort)
	logger.Info("启动本地采集")
	for cmd != "quit" {
		if cmd == "start" {
			go func() {
				for cmd != "quit" {
					err := mq.Open(c.mqIP, strconv.Itoa(c.mqPort), c.mqUser, c.mqPwd, "")
					if nil != err {
						logger.Errorf("连接到mq发生异常,异常信息:%v", err)
					} else {
						logger.Debugf("连接到mq服务器[%s:%d]成功", c.mqIP, c.mqPort)
						topicName := strings.Split(c.acceptTopics, ",")
						mq.SubscribeTopic(topicName)
						s.Notify = func(json string) {
							mq.SendMsg(c.statusQueue, json)
						}
						sendReg(mq, c.serviceCode, c.registerQueue)
						break
					}
					<-time.After(time.Duration(1) * time.Minute)
				}
			}()
		} else if cmd == "stop" {
			s.StopDetect()
			s.Notify = nil
			mq.Close()
		} else if cmd == "help" {
			println(ShowCmdInfo)
		} else {
			if cmd != "" {
				logger.Warringf("未知命令:%s\n", cmd)
			}
		}
		b, _, _ = inputReader.ReadLine()
		cmd = string(b)
	}
}

func sendReg(mq *service.MQService, clientKey, registerQueue string) {
	regjson := service.CreateRegisterMsg(clientKey)
	go func() {
		for {
			err := mq.SendMsg(registerQueue, regjson)
			if err != nil {
				logger.Debugf("发送注册信息发生错误,错误信息:%v", err)
				break
			}
			logger.Debugf("向服务器发送签到信息:%s\n", regjson)
			<-time.After(time.Duration(10) * time.Second)
			if registerResult == false {
				logger.Debugf("未收到签到回复或者签到失败,重发签到:%s\n", regjson)
			} else {
				break
			}
		}
	}()
}

func readConfig() *appConfig {
	configInfo, err := tools.ReadConfigFile("./health_config.yml")
	c := new(appConfig)
	if nil != err {
		logger.Errorf("读配置文件错误,错误信息:%v", err)
		return nil
	}
	c.mqIP = configInfo["mq.ip"].(string)
	c.mqPort = configInfo["mq.port"].(int)
	c.mqUser = configInfo["mq.user"].(string)
	c.mqPwd = configInfo["mq.pwd"].(string)
	c.statusQueue = configInfo["mq.queues.status_queue"].(string)
	c.registerQueue = configInfo["mq.queues.register_queue"].(string)
	c.collectTimeInterval = configInfo["collect.interval"].(int)
	c.serviceCode = configInfo["client.client_key"].(string)
	c.acceptTopics = configInfo["mq.topic.accept_topic"].(string)
	c.webPort = configInfo["web.port"].(int)
	tools.SetServiceCode(c.serviceCode)
	return c
}

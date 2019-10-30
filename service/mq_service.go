package service

import "system_detect/activemq"

// MQService 处理消息服务器
type MQService struct {
	msgqueue  activemq.MsgQueue
	queueName string
}

type configInfo struct {
	host      string
	port      int
	queueName string
}

func readConfig() (*configInfo, error) {
	var config configInfo

	return &config, nil
}

// Init 初始化
func (msv *MQService) Init() {

}

// SendToMQServer 链接到服务
func (msv *MQService) SendToMQServer(msg string) {
	msv.msgqueue.SendMsg(msv.queueName, msg)
}

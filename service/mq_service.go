package service

import (
	"fmt"
	"runtime/debug"
	"strings"
	"system_detect/activemq"

	"github.com/fwhezfwhez/errorx"
)

// MQService 处理消息服务器
type MQService struct {
	msgqueue activemq.MsgQueue
	Callback func(msg, queueName string)
}

// Init 初始化
func (msv *MQService) Init(host, port, user, pwd, queueName string) error {
	queueNames := strings.Split(queueName, ",")
	err := msv.msgqueue.ConnectToServer(host, port, user, pwd, queueNames, msv.Callback)
	if nil != err {
		debug.PrintStack()
		return errorx.GroupErrors(fmt.Errorf("连接到服务器[%s:%s]失败,错误信息:%v", host, port, err))
	}
	return nil
}

// SendToMQServer 链接到服务
func (msv *MQService) SendMsg(queueName, msg string) {
	msv.msgqueue.SendMsg(queueName, msg)
}

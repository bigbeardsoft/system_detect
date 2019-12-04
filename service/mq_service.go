package service

import (
	"fmt"
	"strings"
	"system_detect/activemq"

	"github.com/fwhezfwhez/errorx"
)

// MQService 处理消息服务器
type MQService struct {
	msgqueue *activemq.MsgQueue
	Callback func(msg, queueName string)
}

// Open 初始化
func (msv *MQService) Open(host, port, user, pwd, queueName string) error {
	queueNames := strings.Split(queueName, ",")
	msv.msgqueue = new(activemq.MsgQueue)
	err := msv.msgqueue.ConnectToServer(host, port, user, pwd, queueNames, msv.Callback)
	if nil != err {
		return errorx.GroupErrors(fmt.Errorf("连接到服务器[%s:%s]失败,错误信息:%v", host, port, err))
	}
	return nil
}

// SendMsg 发送消息
func (msv *MQService) SendMsg(queueName, msg string) error {
	return msv.msgqueue.SendMsg(queueName, msg)
}

// Close 关闭连接
func (msv *MQService) Close() {
	if msv.msgqueue != nil {
		msv.msgqueue.Disconnect()
	}
}

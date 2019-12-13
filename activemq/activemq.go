package activemq

import (
	"container/list"
	"fmt"
	"net"
	"system_detect/tools"

	"github.com/fwhezfwhez/errorx"
	"github.com/go-stomp/stomp"
)

//MsgQueue 消息队列的相关操作
type MsgQueue struct {
	Connection *stomp.Conn
	Subs       *list.List
	Host       string
	Port       string
	user       string
	password   string
}

var logger tools.Logger

// Connect 连接到消息服务器,使用tcp的方式连接
// host:主机地址
// port:连接端口
func (msgQueue *MsgQueue) connect(host, port, user, password string) (*stomp.Conn, error) {
	var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
		stomp.ConnOpt.HeartBeat(0, 0),
	}
	conn, err := stomp.Dial("tcp", net.JoinHostPort(host, port), options...)
	if nil == err {
		msgQueue.Host = host
		msgQueue.Port = port
		msgQueue.password = password
		msgQueue.user = user
		return conn, nil
	} else {
		logger.Errorf("连接到[%s:%s %s,%s]失败,错误信息:%v\n", host, port, user, password, err)
		return nil, errorx.New(err)
	}
}

//SubscriptQueue 订阅消息
//   queueName:队列名称
//   conn:连接
// 返回:
//   stomp.Subscription:订阅对象
//   error:错误信息
func (msgQueue *MsgQueue) SubscriptQueue(queueName string, conn *stomp.Conn) (*stomp.Subscription, error) {
	if nil == conn {
		return nil, fmt.Errorf("server is not connected")
	}
	sub, err := conn.Subscribe(queueName, stomp.AckMode(stomp.AckAuto))
	return sub, errorx.New(err)
}

//UnSubscriptQueue 取消订阅
//	sub:订阅对象
func (msgQueue *MsgQueue) UnSubscriptQueue(sub *stomp.Subscription) {
	if nil != sub {
		sub.Unsubscribe()
	}
}

//Send 发送消息到队列
// 参数:
//   queueName:队列名称
//	msg:待发送的消息
//	conn:连接对象
// 返回:
//   error:错误信息
func (msgQueue *MsgQueue) Send(queueName, msg string, conn *stomp.Conn) error {
	if nil == conn {
		return fmt.Errorf("connection is nil or not connected to server")
	}
	err := conn.Send(queueName, "", []byte(msg))
	logger.Infof("send active mq %s-->%s", queueName, msg)
	if err != nil {
		logger.Errorf("active mq message send error: " + err.Error())
	}
	return err
}

//SendTopic 发送主题信息(此方法有bug,无法发送到主题)
func (msgQueue *MsgQueue) SendTopic(queue, msg string, conn *stomp.Conn) error {
	if nil == conn {
		return fmt.Errorf("connection is nil or not connected to server")
	}
	queueName := fmt.Sprintf("/topic/%s", queue)
	err := conn.Send(queueName, "", []byte(msg))
	if err != nil {
		logger.Errorf("active mq message send error: " + err.Error())
	}
	logger.Debugf("Send data to %s : %s ", queueName, msg)
	return err
}

//Receive 接收队列消息,需要采用异步或者多线程的方式调用本函数,否则会阻塞
// 参数:
//   sub:订阅对象
//   callback:回调函数
// 返回:
//   error:错误信息
func (msgQueue *MsgQueue) Receive(sub *stomp.Subscription, callback func(msg, queueName string)) error {
	if nil == sub {
		return errorx.New(fmt.Errorf("subscription is nil"))
	}
	logger.Infof("ready to receive from [%s]", sub.Destination())

	for {
		v := <-sub.C
		if v != nil {
			msg := string(v.Body)
			if nil != callback {
				go callback(msg, sub.Destination())
			}
		} else {
			break
		}
	}
	logger.Infof("end receive data")
	return nil
}

//ConnectToServer 连接到服务器,并开始接收来自队列的消息
// 参数:
//	host:主机
//	port:端口
//	queues:需要接收消息的队列
//	callback:接收到消息的回调函数
func (msgQueue *MsgQueue) ConnectToServer(host, port, user, pwd string, queues []string, callback func(msg, queueName string)) error {

	conn, err := msgQueue.connect(host, port, user, pwd)
	if err != nil {
		return errorx.New(err)
	}
	if queues == nil {
		logger.Error("queues is nil")
	}
	msgQueue.Connection = conn
	msgQueue.Subs = list.New()
	for index := range queues {
		sub, err := msgQueue.SubscriptQueue(queues[index], msgQueue.Connection)
		if nil == err {
			msgQueue.Subs.PushBack(sub)
			go msgQueue.Receive(sub, callback)
		} else {
			msgQueue.Disconnect()
			return errorx.New(err)
		}
	}
	return nil
}

//Disconnect 断开连接
func (msgQueue *MsgQueue) Disconnect() {
	if msgQueue.Subs != nil {
		for sub := msgQueue.Subs.Front(); sub != nil; sub = sub.Next() {
			o, r := sub.Value.(*stomp.Subscription)
			if r {
				msgQueue.UnSubscriptQueue(o)
			}
		}
	}
	if msgQueue.Connection != nil {
		msgQueue.Connection.Disconnect()
	}
}

//SendMsg 发送消息到消息队列
//	参数:
//		queueName:队列名称
//		msg:待发送的消息
//	返回:
//		error:错误消息
func (msgQueue *MsgQueue) SendMsg(queueName, msg string) error {
	if msgQueue.Connection == nil {
		return fmt.Errorf("连接未初始化")
	}
	return msgQueue.SendTopic(queueName, msg, msgQueue.Connection)

}

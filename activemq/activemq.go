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

const (
	//Topic 主题
	Topic string = "/topic/"
	//Queue 队列
	Queue string = "/queue/"
)

var logger tools.Logger

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
		return nil
	}
	msgQueue.Connection = conn
	if msgQueue.Subs == nil {
		msgQueue.Subs = list.New()
	}
	msgQueue.SubscribeQueue(queues, callback)
	return nil
}

//Disconnect 断开连接
func (msgQueue *MsgQueue) Disconnect() {
	if msgQueue.Subs != nil {
		var next *list.Element
		for sub := msgQueue.Subs.Front(); sub != nil; sub = next {
			o, r := sub.Value.(*stomp.Subscription)
			if r {
				o.Unsubscribe()
			}
			next = sub.Next()
			msgQueue.Subs.Remove(sub)
		}
	}
	if msgQueue.Connection != nil {
		msgQueue.Connection.Disconnect()
	}
}

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
	}
	logger.Errorf("连接到[%s:%s %s,%s]失败,错误信息:%v\n", host, port, user, password, err)
	return nil, errorx.New(err)
}

//Send 发送消息到队列
// 参数:
//   queueName:队列名称
//	msg:待发送的消息
//	conn:连接对象
// 返回:
//   error:错误信息
func (msgQueue *MsgQueue) send(queueName, msg string, conn *stomp.Conn) error {
	if nil == conn {
		return fmt.Errorf("connection is nil or not connected to server")
	}
	err := conn.Send(queueName, "", []byte(msg))
	logger.Infof("send active mq %s-->%s", queueName, msg)
	if err != nil {
		logger.Errorf("active mq message send error: " + err.Error())
	}
	logger.Infof("Send data to %s : %s ", queueName, msg)
	return err
}

//SendToTopic 发送主题信息
func (msgQueue *MsgQueue) SendToTopic(queue, msg string, conn *stomp.Conn) error {
	queueName := fmt.Sprintf("%s%s", Topic, queue)
	return msgQueue.send(queueName, msg, conn)
}

//SendToQueue 发送消息到消息队列
//	参数:
//		queue:队列名称
//		msg:待发送的消息
//	返回:
//		error:错误消息
func (msgQueue *MsgQueue) SendToQueue(queue, msg string) error {
	queueName := fmt.Sprintf("%s%s", Queue, queue)
	return msgQueue.send(queueName, msg, msgQueue.Connection)
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
			logger.Debugf("收到[%s]==>%s", sub.Destination(), msg)
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

//SubscribeTopic 订阅主题消息
//参数:
//  queues: 队列名称
//  callback: 收到消息之后的回调
func (msgQueue *MsgQueue) SubscribeTopic(queues []string, callback func(msg, queueName string)) error {
	return msgQueue.subscribe(queues, Topic, callback)
}

//SubscribeQueue 订阅消息
//参数:
//  queues: 队列名称
//  callback: 收到消息之后的回调
func (msgQueue *MsgQueue) SubscribeQueue(queues []string, callback func(msg, queueName string)) error {
	return msgQueue.subscribe(queues, Topic, callback)
}

func (msgQueue *MsgQueue) subscribe(queues []string, queueType string, callback func(msg, queueName string)) error {
	if msgQueue.Connection == nil {
		return fmt.Errorf("连接未初始化")
	}
	if queues == nil {
		logger.Warring("queues is nil")
		return nil
	}
	if msgQueue.Subs == nil {
		msgQueue.Subs = list.New()
	}
	for index := range queues {
		q := fmt.Sprintf("%s%s", queueType, queues[index])
		sub, err := msgQueue.Connection.Subscribe(q, stomp.AckMode(stomp.AckAuto))
		if nil == err {
			msgQueue.Subs.PushBack(sub)
			go msgQueue.Receive(sub, callback)
		} else {
			return errorx.New(err)
		}
	}
	return nil
}

//UnSubscriptTopic 取消主题订阅
func (msgQueue *MsgQueue) UnSubscriptTopic(queue string) error {
	return msgQueue.unSub(queue, Topic)
}

//UnSubscriptQueue 取消订阅队列
//	sub:订阅对象
func (msgQueue *MsgQueue) UnSubscriptQueue(queue string) error {
	return msgQueue.unSub(queue, Queue)
}

func (msgQueue *MsgQueue) unSub(queue, queueType string) error {
	var next *list.Element
	q := fmt.Sprintf("%s%s", queueType, queue)
	for sub := msgQueue.Subs.Front(); sub != nil; sub = next {
		o, r := sub.Value.(*stomp.Subscription)
		if r && o.Destination() == q {
			msgQueue.Subs.Remove(sub)
			o.Unsubscribe()
			return nil
		}
		next = sub.Next()
	}
	return nil
}

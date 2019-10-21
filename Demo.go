package activemq

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/go-stomp/stomp"
)

// 使用IP和端口连接到ActiveMQ服务器
// 返回ActiveMQ连接对象
func connActiveMq(host, port string) (stompConn *stomp.Conn) { // @todo 实现断开重连
	stompConn, err := stomp.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		fmt.Printf("connect to active_mq server  [%s] service, error: %v\r\n", host, err.Error())
		os.Exit(1)
	}

	return stompConn
}

// 将消息发送到ActiveMQ中
func activeMqProducer(c string, queue string, conn *stomp.Conn) {
	for r := 0; r < 10; r++ {
		err := conn.Send(queue, "text/plain", []byte(c+strconv.Itoa(r)))
		fmt.Println("send active mq " + queue)
		if err != nil {
			fmt.Println("active mq message send error: " + err.Error())
		}
	}
}

//func receiveFromMq(r chan string, queue string, conn *stomp.Conn) {
func receiveFromMq(queue string, conn *stomp.Conn) {
	sub, _ := conn.Subscribe(queue, stomp.AckMode(stomp.AckAuto))
	for {
		m := <-sub.C
		fmt.Printf("收到ActiveMQ的数据:%v\r\n", string(m.Body))
		//r <- string(m.Body)
	}
	sub.Unsubscribe()

}

func (demo Demo) Send(msg, queueName string, conn *stomp.Conn) {
	activeMqProducer(msg, queueName, conn)
}

func (demo Demo) DoAcitveMQDemo() *stomp.Conn {

	host := "127.0.0.1"
	port := "61613"

	queue := "SPF300006"

	activeMq := connActiveMq(host, port)
	//defer activeMq.Disconnect()
	// c := make(chan string)
	//r := make(chan string)
	// go func() {
	// 	for {
	// 		v := <-r
	// 		println("===>" + v)
	// 	}
	// }()
	// 启动Go routine发送消息
	go activeMqProducer("aBCD", queue, activeMq)

	go receiveFromMq(queue, activeMq)

	// for i := 0; i < 100; i++ {
	// 	// 发送1万个消息
	// 	c <- "hello world" + strconv.Itoa(i)
	// }
	return activeMq
}

type Demo struct {
	Running bool
}

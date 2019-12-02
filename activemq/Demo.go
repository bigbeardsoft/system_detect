package activemq

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/go-stomp/stomp"
)

// func testTopic() {
// 	m := topic.Manager()
// 	m2 := topic.Topic()

// }

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

func send(msg, queueName string, conn *stomp.Conn) {
	activeMqProducer(msg, queueName, conn)
}

func doActiveMQDemo() *stomp.Conn {

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

func CallActiveMq() {
	//host, port, user, pwd := "127.0.0.1", "61613", "admin", "123456"
	println("please input host:")
	inputReader := bufio.NewReader(os.Stdin)
	b, _, _ := inputReader.ReadLine()
	host := string(b)
	println("please input port:")
	b, _, _ = inputReader.ReadLine()
	port := string(b)
	println("please input user:")
	b, _, _ = inputReader.ReadLine()
	user := string(b)
	println("please input password:")
	b, _, _ = inputReader.ReadLine()
	pwd := string(b)
	var queues = []string{"SPF300006", "QueueA"}
	connectToServerB(host, port, user, pwd, queues)
}

func connectToServerB(host, port, user, pwd string, queues []string) {
	var a *MsgQueue
	var b MsgQueue

	a = &b
	a.Host = "aaaa"
	a.Port = "12312"
	err := a.ConnectToServer(host, port, user, pwd, queues, func(msg, queueName string) {
		fmt.Printf("from [%s] data :[%s]\r\n", queueName, msg)
	})
	if err != nil {
		println(err.Error())
		return
	}
	println("press [q]  to exit...")
	for {
		println("please input command:")
		inputReader := bufio.NewReader(os.Stdin)
		b, _, _ := inputReader.ReadLine()
		str := string(b)
		if str == "q" {
			a.Disconnect()
			break
		} else if str == "send" {
			println("please input content which you need to send:")
			b, _, _ := inputReader.ReadLine()
			send := string(b)
			println("please select queue:")
			for r := range queues {
				fmt.Printf("%d %s\r\n", r, queues[r])
			}
			index, _, _ := inputReader.ReadLine()
			ivalue, err := strconv.Atoi(string(index))
			if err == nil {
				errx := a.SendMsg(queues[ivalue], send)
				if errx != nil {
					println("send error,error info :" + errx.Error())
				}
			} else {
				println("select error:" + err.Error())
			}

		} else {
			println("unknown command : " + str)
		}
	}
}

func main() {
	CallActiveMq()
}

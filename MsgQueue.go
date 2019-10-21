package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"go-mycode/ActiveMQ/activemq"
)

func main() {
	CallActiveMq()
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
	ConnectToServerB(host, port, user, pwd, queues)
}

func ConnectToServerB(host, port, user, pwd string, queues []string) {
	var a *activemq.MsgQueue // = &activemq.MsgQueue
	var b activemq.MsgQueue
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

package service

import (
	"fmt"
	"sync"
)

var token string

//ParseRegMsg 解析设备签到结果,重点解析token
func ParseRegMsg(msg Msg) error {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	if msg.H.F != FRegCode {
		return fmt.Errorf("当前方法无法处理%s的协议", msg.H.F)
	}
	if msg.H.R {
		token = msg.H.K
		return nil
	}
	return fmt.Errorf("设备签到,服务器返回错误,错误消息:%v", msg)
}

//GetToken 获取token
func GetToken() string {
	return token
}

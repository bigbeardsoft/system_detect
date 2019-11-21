package parser

import "fmt"

// DealResponse 定义处理接收消息的接口
type DealResponse interface {
	// Deal 声明处理接口
	Deal(map[string]interface{})
}

// DealReg 实现处理接收消息功能
type DealReg struct {
	NotifyToken func(token string)
}

// Deal 解析签到成功之后的token
func (d *DealReg) Deal(x map[string]interface{}) {

	
	if _, ok := x["B"]; ok == false {
		fmt.Printf("被解析的信息中不包括B部分,%v", x)
		return
	}

	bodyMap := x["B"].(map[string]interface{})
	if _, ok := bodyMap["Token"]; ok == false {
		fmt.Printf("被解析的信息中不包括Token部分,%v", bodyMap)
		return
	}
	token := bodyMap["Token"].(string)
	if d.NotifyToken != nil {
		d.NotifyToken(token)
	}
}

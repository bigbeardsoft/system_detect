package service

import (
	tools "go-mycode/tools"
	"sync"
)

var msgIndex uint64 = 0

func getMsgIndex() uint64 {
	var mutex sync.Mutex
	mutex.Lock()
	msgIndex++
	defer mutex.Unlock()
	return msgIndex
}

//CreateRegisterMsg 创建注册用的json消息
func CreateRegisterMsg(serviceCode string) string {

	msgH := createHead(FRegCode, "")
	var msg Msg
	msg.H = msgH
	var msgbody MsgRegBody
	msgbody.ServiceCode = serviceCode

	msg.B = msgbody
	json, err := tools.ConvertToJSON(msg)
	if nil != err {
		return ""
	}
	return json
}

func createHead(code, token string) MsgHead {
	var m MsgHead
	m.F = code
	m.V = ProtocolVersion
	m.R = true
	m.K = token
	m.S = getMsgIndex() //strconv.Itoa(getMsgIndex())
	m.T = 1
	m.I = 1
	m.M = ""
	return m
}

// CreateDetectMsg 创建检测消息
func CreateDetectMsg() string {

	json := ""
	return json
}

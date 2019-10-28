package service

import (
	"fmt"
	"net"
	"sync"
	system "system_detect/detect/system/linux"
	tools "system_detect/tools"
	"time"
)

var msgIndex uint64

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
func CreateDetectMsg(process []ProcessInfo, sys *system.SysUsedInfo, diskinfo []DiskSpaceInfo) string {

	msgbody := new(MsgServerInfoBody)
	d := time.Now()
	sd := fmt.Sprintf("%d-%02d-%02d %0d:%02d:%02d ", d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second())
	msgbody.CollectTime = sd
	msgbody.CPU = sys.CPUFree
	msgbody.MEM = sys.MemFree / (sys.MemAll * 1.0)
	msgbody.DiskFreeSpace = diskinfo
	msgbody.HanderCount = -1
	msgbody.NetWork = -1
	msgbody.ProcessCount = len(process)
	msgbody.ProcessInfo = process

	var msg Msg
	var msghead MsgHead

	msghead = createHead(FServerStatReport, token)
	msg.B = msgbody
	msg.H = msghead
	json, err := tools.ConvertToJSON(msg)
	if nil != err {
		fmt.Printf("%v\n", err)
		return ""
	}
	return json
}

func getLocalIPAddress() []string {

	var ret []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
				ret = append(ret, ipnet.IP.String())
			}
		}
	}
	return ret
}

func getConvertRangeToString(val []string) string {
	resultStr := ""
	for _, r := range val {
		resultStr += r + ","
	}
	return resultStr
}

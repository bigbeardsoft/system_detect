package service

import (
	"encoding/json"
	"fmt"
	tools "system_detect/tools"

	"system_detect/service/parser"
)

var methods map[string]parser.DealResponse

// Init 初始化
func Init() {
	methods = make(map[string]parser.DealResponse, 1)
	dreg := new(parser.DealReg)
	dreg.NotifyToken = f
	methods[FRegCode] = dreg
}

var token string

var f = func(t string) {
	token = t
	go registerCallback(t != "")
}

var registerCallback func(result bool)

// SetRegisterCallback 设置注册结果回调
func SetRegisterCallback(callback func(result bool)) {
	registerCallback = callback
}

// GetToken 获取从服务器上返回的token信息
func GetToken() string {
	return token
}

// DispatherMsg 解析消息并完成消息调度
func DispatherMsg(jsonstr string) error {

	var tmp map[string]interface{}
	err := json.Unmarshal([]byte(jsonstr), &tmp)
	if nil != err {
		return err
	}
	headMap := tmp["H"].(map[string]interface{})

	serviceCode := headMap["C"].(string)
	if serviceCode != tools.GetServiceCode() {
		return nil
	}

	headResult, err := checkHead(headMap)
	if headResult == false {
		return err
	}

	code := headMap["F"].(string)
	if f, ok := methods[code]; ok {
		f.Deal(tmp)
	} else {
		logger.Warringf("未实现[%s]协议的处理,json内容:%s", code, jsonstr)
	}

	return nil
}

func checkHead(head map[string]interface{}) (bool, error) {
	returnResult := head["R"]
	if returnResult == false {
		returnMsg := head["M"]
		return false, fmt.Errorf("服务器返回错误,错误描述:%s", returnMsg)
	}
	version := head["V"]
	if version != ProtocolVersion {
		return false, fmt.Errorf("当前系统无法处理版本:%s的协议,只能处理:%s的协议", version, ProtocolVersion)
	}
	return true, nil
}

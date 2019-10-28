package service

import (
	"encoding/json"
	"fmt"
)

// DispatherMsg 解析消息并完成消息调度
func DispatherMsg(jsonstr string) error {

	var tmp map[string]interface{}
	err := json.Unmarshal([]byte(jsonstr), &tmp)
	if nil != err {
		return err
	}
	headMap := tmp["H"].(map[string]interface{})
	headResult, err := checkHead(headMap)
	if headResult == false {
		return err
	}
	code := headMap["F"]
	if code == FRegCode {
		fmt.Printf("%v\n", tmp)
		//调用处理函数,存储token等信息
	} else if code == FServerStatReport {
		fmt.Printf("%v\n", tmp)
	}
	return nil
}

func checkHead(head map[string]interface{}) (bool, error) {
	returnResult := head["R"]
	if returnResult == false {
		returnMsg := head["M"]
		returnCode := head["C"]
		return false, fmt.Errorf("服务器返回错误,错误编码:%s,错误描述:%s", returnCode, returnMsg)
	}
	version := head["V"]
	if version != ProtocolVersion {
		return false, fmt.Errorf("当前系统无法处理版本:%s的协议,只能处理:%s的协议", version, ProtocolVersion)
	}
	return true, nil
}

package tools

import "fmt"

//CreateErrorResult 创建错误返回结果
//errinfo 错误信息
func CreateErrorResult(errinfo string) string {
	str := fmt.Sprintf("{\"Result\":\"false\",\"ErrorMsg\":\"%s\"}", errinfo)
	return str
}

//CreateSuccessResult 创建成功返回结果
//v 需要返回的数据
func CreateSuccessResult(v interface{}) string {
	str := fmt.Sprintf("{\"Result\":\"true\",\"Data\":\"%v\"}", v)
	return str
}

var serviceCode string

//GetServiceCode 获取服务编码
func GetServiceCode() string {
	return serviceCode
}

//SetServiceCode 设置服务编码
func SetServiceCode(code string) {
	serviceCode = code
}

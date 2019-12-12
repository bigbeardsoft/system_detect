package controll

import (
	"system_detect/detect/process"
	"system_detect/tools"
)

var logger = new(tools.Logger)

//GetProcessList 获取系统中的进程列表
func GetProcessList() string {

	var p = new(process.Process)

	r, e := p.GetAllProcess()
	if e != nil {
		logger.ErrorX(e)
		return tools.CreateErrorResult("获取进程列表发生错误")
	}
	json, e := tools.ConvertToJSON(r)
	if e != nil {
		logger.ErrorX(e)
		return tools.CreateErrorResult("数据转换发生错误")
	}
	return tools.CreateSuccessResult(json)
}

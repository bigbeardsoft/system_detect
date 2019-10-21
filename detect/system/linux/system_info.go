package system

import (
	"fmt"
	"strconv"
	"strings"
	tools "system_detect/tools"
)

//SysUsedInfo 获取系统内存,CPU和进程数量信息.
type SysUsedInfo struct {
	CPUFree      float64
	MemFree      uint64
	MemAll       uint64
	NetDown      float32
	NetUp        float32
	ProcessCount int
}

//GetSystemUsedInfo 获取系统的进程数量,CPU使用率,内存大小和空闲内存.
func (p *SysUsedInfo) GetSystemUsedInfo() (*SysUsedInfo, error) {
	cmdresult, err := tools.ExecuteCommand("top -bn 1")
	if nil != err {
		return nil, err
	}
	var lines [5]string
	resultstr := strings.Split(cmdresult, "\n")
	var index int
	for _, r := range resultstr {
		r = strings.Trim(r, " \\r\\t\\n\\v\\f")
		if "" != r && len(r) > 0 {
			lines[index] = r
			index++
		}
	}
	var r SysUsedInfo

	processTokens := strings.Split(lines[1], " ")
	count, err := strconv.Atoi(processTokens[1])
	if nil == err {
		r.ProcessCount = count
	} else {
		fmt.Printf("Convert process count  happend an  error :%v\n", err)
		fmt.Println(lines[1])
	}
	cpuTokens := strings.Split(lines[2], " ")
	var cpuInfo [5]string
	i := 0

	for _, s := range cpuTokens {
		r := tools.IsNumeric(s)
		if r {
			cpuInfo[i] = strings.Trim(s, " \\r\\t\\n\\v\\f")
			i++
		}
		if i >= len(cpuInfo) {
			break
		}
	}

	cpufree, err := strconv.ParseFloat(cpuInfo[3], 64)
	if nil == err {
		r.CPUFree = cpufree
	} else {
		fmt.Printf(" Convert CPU used precentage happend an  error :%v\n", err)
		fmt.Println(lines[2])
	}

	memTokens := strings.Split(lines[3], " ")
	i = 0
	var memInfo [8]string
	for _, s := range memTokens {
		r := tools.IsNumeric(s)
		if r {
			memInfo[i] = strings.Trim(s, " \\r\\t\\n\\v\\f")
			i++
		}
		if i >= len(memInfo) {
			break
		}
	}
	memTokens = strings.Split(lines[4], " ")
	for _, s := range memTokens {
		r := tools.IsNumeric(s)
		if r {
			memInfo[i] = strings.Trim(s, " \\r\\t\\n\\v\\f")
			i++
		}
		if i >= len(memInfo) {
			break
		}
	}

	memall, err := strconv.ParseUint(memInfo[0], 10, 64)
	if nil == err {
		r.MemAll = memall
	} else {
		fmt.Printf(" Convert Mem All   (use index 0) happend an  error :%v\n", err)
		fmt.Println(lines[3])
	}
	memfree, err := strconv.ParseUint(memInfo[1], 10, 64)
	if nil == err {
		r.MemFree = memfree
	} else {
		fmt.Printf(" Convert Mem free (use index 1) happend an  error :%v\n", err)
		fmt.Println(lines[3])
	}
	memcatch, err := strconv.ParseUint(memInfo[3], 10, 64)
	if nil == err {
		r.MemFree = r.MemFree + memcatch
	} else {
		fmt.Printf(" Convert Mem catch (use index 3) happend an  error :%v\n", err)
		fmt.Println(lines[3])
	}

	return &r, nil
}

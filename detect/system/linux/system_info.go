package system

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	tools "system_detect/tools"

	"github.com/fwhezfwhez/errorx"
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

var logger tools.Logger

//GetSystemUsedInfo 获取系统的进程数量,CPU使用率,内存大小和空闲内存.
func (p *SysUsedInfo) GetSystemUsedInfo() (*SysUsedInfo, error) {
	cmdresult, err := tools.ExecuteCommand("top -bn 1")
	if nil != err {
		return nil, errorx.Wrap(err)
	}
	var lines [5]string
	resultstr := strings.Split(cmdresult, "\n")
	var index int
	for _, r := range resultstr {
		r = strings.Trim(r, " \r\t\n\v\f")
		if "" != r && len(r) > 0 {
			lines[index] = r
			index++
		}
		if index == len(lines) {
			break
		}
	}

	var r = new(SysUsedInfo)

	processTokens := strings.Split(lines[1], " ")
	count, err := strconv.Atoi(processTokens[1])
	if nil == err {
		r.ProcessCount = count
	} else {
		logger.Errorf("Convert process count  happend an  error :%v,错误数据:%s", err, lines[1])
	}
	cpuTokens := strings.Split(lines[2], " ")
	var cpuInfo [5]string
	i := 0

	for _, s := range cpuTokens {
		r := tools.IsNumeric(s)
		if r {
			cpuInfo[i] = strings.Trim(s, " \r\t\n\v\f")
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
		logger.Errorf(" Convert CPU used precentage happend an  error :%v,错误数据:%s", err, lines[2])
	}

	memTokens := strings.Split(lines[3], " ")
	i = 0
	var memInfo [8]string
	for _, s := range memTokens {
		r := tools.IsNumeric(s)
		if r {
			memInfo[i] = strings.Trim(s, " \r\t\n\v\f")
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
			memInfo[i] = strings.Trim(s, " \r\t\n\v\f")
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
		logger.Errorf(" Convert Mem All   (use index 0) happend an  error :%v,错误数据:%s", err, lines[3])
	}
	memfree, err := strconv.ParseUint(memInfo[1], 10, 64)
	if nil == err {
		r.MemFree = memfree
	} else {
		logger.Errorf(" Convert Mem free (use index 1) happend an  error :%v,错误数据:%s", err, lines[3])
	}
	memcatch, err := strconv.ParseUint(memInfo[3], 10, 64)
	if nil == err {
		r.MemFree = r.MemFree + memcatch
	} else {
		logger.Errorf(" Convert Mem catch (use index 3) happend an  error :%v,错误数据:%s", err, lines[3])
	}

	return r, nil
}

// GetSystemName 获取系统名称
func GetSystemName() string {
	cmd := "uname -sn"
	cmdresult, err := tools.ExecuteCommand(cmd)
	if nil != err {
		logger.Errorf("获取系统名称发生错误,执行[%s]命令错误", cmd)
		return ""
	}
	return cmdresult
}

// GetLocalIP 获取系统中的IP地址
func GetLocalIP() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		logger.Errorf("net.Interfaces failed, err:", err.Error())
		return ""
	}
	var ips = ""
	c := len(netInterfaces)
	for i := 0; i < c; i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						fmt.Println(ipnet.IP.String())
						if (i + 1) < c {
							ips += (ipnet.IP.String() + ",")
						} else {
							ips += ipnet.IP.String()
						}
					}
				}
			}
		}
	}
	return ips
}

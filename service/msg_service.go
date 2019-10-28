package service

import (
	"fmt"
	"sync"
	"system_detect/detect/process"
	system "system_detect/detect/system/linux"

	"github.com/fwhezfwhez/errorx"
)

var isRun bool

var mutex sync.Mutex

// StartDetect 启动监测
func StartDetect() {
	mutex.Lock()
	if isRun {
		fmt.Println("服务正在运行....")
		return
	}
	isRun = true
	mutex.Unlock()
	var sysUsed = new(system.SysUsedInfo)
	var diskStatus system.DiskStatus
	var paths []string
	var prc = new(process.Process)
	for isRun {
		s, err := sysUsed.GetSystemUsedInfo()
		if nil != err {
			fmt.Printf("========>%v\n<=======", err.(errorx.Error).PrintStackTrace())
			return
		}
		var xpath []DiskSpaceInfo
		for _, pth := range paths {
			d, e := diskStatus.DiskUsage(pth)
			if nil == e {
				xpath = append(xpath, DiskSpaceInfo{Path: pth, TotalSpace: d.All, FreeSpace: d.Free})
			}
		}
		ap, errx := prc.GetAllProcess()
		var pcc []ProcessInfo
		if nil == errx {
			for _, p := range ap {
				pcc = append(pcc, ProcessInfo{PID: p.Pid, ProcessName: p.ProcessName,
					ProcessPath: p.ProcessPath, CPU: p.CPU, MEM: p.Memory, ThreadCount: p.ThreadCount})
			}
		}

		result := CreateDetectMsg(pcc, s, xpath)
		println(result)
	}
}

// StopDetect 停止监测
func StopDetect() {
	isRun = false
}

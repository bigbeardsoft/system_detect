package service

import (
	"fmt"
	"sync"
	"system_detect/detect/process"
	system "system_detect/detect/system/linux"

	"github.com/fwhezfwhez/errorx"
)

// CollectService 系统检测
type CollectService struct {
	isRun  bool
	mutex  sync.Mutex
	Notify func(string)
}

// StartDetect 启动监测
func (servicePoint *CollectService) StartDetect() {
	servicePoint.mutex.Lock()
	if servicePoint.isRun {
		fmt.Println("服务正在运行....")
		return
	}
	servicePoint.isRun = true
	servicePoint.mutex.Unlock()
	var sysUsed = new(system.SysUsedInfo)
	var diskStatus system.DiskStatus
	var paths []string
	var prc = new(process.Process)
	for servicePoint.isRun {
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
		//println(result)
		go servicePoint.Notify(result)
	}
}

// StopDetect 停止监测
func (servicePoint *CollectService) StopDetect() {
	servicePoint.isRun = false
}

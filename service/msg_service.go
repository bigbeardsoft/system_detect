package service

import (
	"sync"
	"system_detect/detect/process"
	system "system_detect/detect/system/linux"
	"system_detect/tools"
	"time"
)

// CollectService 系统检测
type CollectService struct {
	isRun        bool
	mutex        sync.Mutex
	Notify       func(string)
	TimeInterval int
}

var logger tools.Logger

// StartDetect 启动监测
func (servicePoint *CollectService) StartDetect() {
	servicePoint.mutex.Lock()
	if servicePoint.isRun {
		logger.Debugf("服务正在运行....")
		return
	}
	servicePoint.isRun = true

	servicePoint.mutex.Unlock()
	var sysUsed = new(system.SysUsedInfo)
	var diskStatus system.DiskStatus
	var paths []string
	var prc = new(process.Process)
	systenName := system.GetSystemName()
	ip := system.GetLocalIP()
	for servicePoint.isRun {
		s, err := sysUsed.GetSystemUsedInfo()
		if nil != err {
			logger.Errorf("读取系统使用信息发生错误:%v", err)
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

		result := CreateDetectMsg(pcc, s, xpath, systenName, ip)
		go servicePoint.Notify(result)

		<-time.After(time.Duration(servicePoint.TimeInterval) * time.Second)
	}
}

// StopDetect 停止监测
func (servicePoint *CollectService) StopDetect() {
	servicePoint.isRun = false
}

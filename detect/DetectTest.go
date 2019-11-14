package detect

import (
	"fmt"

	"system_detect/detect/process"
	linux "system_detect/detect/system/linux"
)

func test() {
	d := new(linux.SysUsedInfo)
	px, _ := d.GetSystemUsedInfo()
	fmt.Printf("process count :%d,Cpu Used:%f,MemAll:%d KiB,MemFree:%d KiB\n", px.ProcessCount, px.CPUUsed, px.MemAll, px.MemFree)
}

func testDiskInfo() {
	var x linux.DiskStatus
	p, _ := x.DiskUsage("/home")
	fmt.Printf("all:%dKB,free:%dKB,used:%dKB\n", p.All, p.Free, p.Used)
}

func testProcess() {
	var a process.Process
	x := &a
	result, err := x.GetAllProcess()
	if err != nil {
		println(err)
		return
	}
	println("==========================================")
	for p := range result {
		o := result[p]
		fmt.Printf("user:%s,\tstartTime:%s,\tcpu:%f,\tThreadCount:%d,\tmem:%f,\tpid:%d,\tcmd:%s\r", o.User, o.StartTime, o.CPU, o.ThreadCount, o.Memory, o.Pid, o.ProcessPath)
	}
}

package process

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	tools "system_detect/tools"

	"github.com/fwhezfwhez/errorx"
)

var logger tools.Logger

//Process 获取系统的进程信息
type Process struct {
	ProcessName string
	ProcessPath string
	CPU         float64
	Memory      float64
	ThreadCount int
	Pid         int
	StartTime   string
	User        string
}

//GetAllProcess 获取系统的所有进程信息
func (p *Process) GetAllProcess() ([]Process, error) {
	cmd := exec.Command("ps", "aux")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		logger.Errorf("执行命令发生错误:%v", err)
		return nil, err
	}
	processes := make([]Process, 0)
	for {
		line, err := out.ReadString('\n')
		if err != nil {
			break
		}

		tokens := strings.Split(line, " ")
		ft := make([]string, 0)
		for _, t := range tokens {
			if t != "" && t != "\t" {
				ft = append(ft, t)
			}
		}
		user := ft[0]
		pid, err := strconv.Atoi(ft[1])
		if err != nil {
			continue
		}
		cpu, err := strconv.ParseFloat(ft[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		mem, err := strconv.ParseFloat(ft[5], 64)
		startTime := ft[8]
		cmd := strings.Trim(ft[10], " ")
		if len(ft) > 10 {
			for index := 11; index < len(ft); index++ {
				cmd = cmd + " " + strings.Trim(ft[index], " ")
			}
		}
		threadCount, err := GetProcessThreadCount(pid)
		if nil != err {
			logger.Errorf("获取进程[%v]的线程信息发生错误%v\n", ft, err)
		}

		processes = append(processes, Process{User: user, Pid: pid,
			CPU: cpu, Memory: mem, StartTime: startTime, ProcessPath: cmd, ThreadCount: threadCount})
	}
	for index := 0; index < len(processes); index++ {
		for j := index; j < len(processes); j++ {
			if processes[index].Pid > processes[j].Pid {
				x := processes[index]
				processes[index] = processes[j]
				processes[j] = x
			}
		}
	}
	return processes, nil
}

// GetProcessThreadCount 获取指定进程的线程数量
func GetProcessThreadCount(pid int) (int, error) {

	if ok := checkProcessExits(pid); ok == false {
		return -1, errorx.New(fmt.Errorf("进程[%d]不存在", pid))
	}
	cmd := fmt.Sprintf("ps -T -p %s", strconv.Itoa(pid))
	result, err := tools.ExecuteCommand(cmd)
	if nil != err {
		return -1, errorx.Wrap(err)
	}

	threadStrings := strings.Split(result, "\n")
	threadCount := len(threadStrings) - 1
	return threadCount, nil
}

func checkProcessExits(pid int) bool {

	cmd := fmt.Sprintf("ps -ef|grep %s", strconv.Itoa(pid))
	result, err := tools.ExecuteCommand(cmd)
	if nil == err {
		countStr := strings.Split(result, "\n")
		return len(countStr) > 1
	}
	return false
}

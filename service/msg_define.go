package service

//FRegCode 注册的F码
const FRegCode string = "100001"

// FServerStatReport 操作系统状态F码
const FServerStatReport string = "100011"

//ProtocolVersion 协议版本
const ProtocolVersion string = "V1.0.1"

// Msg 完整的消息
type Msg struct {
	H MsgHead
	B interface{}
}

// MsgHead 发送的消息头部
type MsgHead struct {
	F string
	V string
	K string
	R bool
	M string
	S uint64
	T int
	I int
	C string
}

// MsgRegBody 注册消息体
type MsgRegBody struct {
	ServiceCode string
}

// DiskSpaceInfo 磁盘空间消息
type DiskSpaceInfo struct {
	Path       string
	TotalSpace uint64
	FreeSpace  uint64
}

//ProcessInfo 进程信息消息
type ProcessInfo struct {
	PID         int
	ProcessName string
	ProcessPath string
	CPU         float64
	MEM         float64
	ThreadCount int
	HanderCount int
	NetWork     int
	PortInfo    string
}

// MsgServerInfoBody 服务器整体信息
type MsgServerInfoBody struct {
	ServerIP      string
	ServerName    string
	CollectTime   string
	ProcessCount  int
	CPU           float64
	MEM           uint64
	ThreadCount   int
	HanderCount   int
	NetWork       int
	DiskFreeSpace []DiskSpaceInfo
	ProcessInfo   []ProcessInfo
}

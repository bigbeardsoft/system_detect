package system

import "syscall"

//DiskStatus 获取硬盘使用情况.
type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

// DiskUsage disk usage of path/disk.
func (p DiskStatus) DiskUsage(path string) (DiskStatus, error) {
	var disk DiskStatus
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return disk, err
	}

	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return disk, nil
}

package disk

import "syscall"

// Creating structure for DiskStatus
type Status struct {
	All  uint64 `json:"All"`
	Used uint64 `json:"Used"`
	Free uint64 `json:"Free"`
}

// Function to get
// disk usage of path/disk
func Usage(path string) (disk Status) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

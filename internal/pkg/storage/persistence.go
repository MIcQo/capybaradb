package storage

//
//type DiskStorage struct {
//}
//
//func NewDiskStorage() *DiskStorage {
//	return &DiskStorage{}
//}
//
//func (d *DiskStorage) Free() (uint64, error) {
//	var path = "/"
//	if runtime.GOOS == "windows" {
//		path = "C:\\"
//	}
//	var usage = disk2.Usage(path)
//	return usage.Free, nil
//}
//
//func (d *DiskStorage) CreateSchema(name string) error {
//	return nil
//}
//
//func (d *DiskStorage) HasSchema(name string) bool {
//	return false
//}

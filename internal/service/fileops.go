package service

import "os"

type OsFileOps struct{}

func (o OsFileOps) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (o OsFileOps) Create(name string) (*os.File, error) {
	return os.Create(name)
}

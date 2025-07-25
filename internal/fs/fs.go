package fs

import (
	"mime/multipart"
	"os"
)

type FileOps interface {
	MkdirAll(path string, perm os.FileMode) error
	Create(path string) (*os.File, error)
	Remove(name string) error
}

type FileOpener interface {
	Open() (multipart.File, error)
}

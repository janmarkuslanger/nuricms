package mockservices

import (
	"os"
)

type MockFileOps struct {
	MkdirErr   error
	CreateErr  error
	Created    *os.File
	CreatePath string
}

func (m *MockFileOps) MkdirAll(path string, perm os.FileMode) error {
	return m.MkdirErr
}

func (m *MockFileOps) Create(name string) (*os.File, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	m.CreatePath = name
	return m.Created, nil
}

package mockservices

import (
	"os"
)

type MockFileOps struct {
	MkdirErr   error
	CreateErr  error
	RemoveErr  error
	Created    *os.File
	CreatePath string
	MkdirPath  string
	Removed    string
}

func (m *MockFileOps) MkdirAll(path string, perm os.FileMode) error {
	m.MkdirPath = path
	return m.MkdirErr
}

func (m *MockFileOps) Create(name string) (*os.File, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	m.CreatePath = name
	return m.Created, nil
}

func (m *MockFileOps) Remove(name string) error {
	m.Removed = name
	return m.RemoveErr
}

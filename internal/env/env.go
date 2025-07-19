package env

import "os"

type Env struct {
	Secret string
}

type EnvSource interface {
	Getenv(string) string
}

type OsEnv struct{}

func (OsEnv) Getenv(k string) string {
	return os.Getenv(k)
}

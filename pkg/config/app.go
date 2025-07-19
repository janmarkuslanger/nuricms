package config

import (
	"github.com/janmarkuslanger/nuricms/pkg/plugin"
	"gorm.io/gorm"
)

type Config struct {
	Port        string
	HookPlugins []plugin.HookPlugin
	Dialector   *gorm.Dialector
}

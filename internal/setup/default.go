package setup

import (
	"github.com/janmarkuslanger/nuricms/pkg/config"
	"gorm.io/driver/sqlite"
)

func SetDefaultConfig(conf config.Config) config.Config {
	if conf.Port == "" {
		conf.Port = "8080"
	}

	if conf.Dialector == nil {
		dl := sqlite.Open("nuricms.db")
		conf.Dialector = &dl
	}

	return conf
}

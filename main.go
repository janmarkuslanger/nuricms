package nuricms

import (
	"log"
	"net/http"

	"github.com/janmarkuslanger/nuricms/internal/env"
	"github.com/janmarkuslanger/nuricms/internal/setup"
	"github.com/janmarkuslanger/nuricms/pkg/config"
)

func Run(config config.Config) {
	envs := env.OsEnv{}
	a, err := setup.SetupApp(config, envs)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(":"+a.Config.Port, a.Server))
}

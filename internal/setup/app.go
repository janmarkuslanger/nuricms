package setup

import (
	"github.com/janmarkuslanger/nuricms/internal/env"
	"github.com/janmarkuslanger/nuricms/internal/modules/api"
	"github.com/janmarkuslanger/nuricms/internal/modules/apikey"
	"github.com/janmarkuslanger/nuricms/internal/modules/asset"
	"github.com/janmarkuslanger/nuricms/internal/modules/collection"
	"github.com/janmarkuslanger/nuricms/internal/modules/content"
	"github.com/janmarkuslanger/nuricms/internal/modules/field"
	"github.com/janmarkuslanger/nuricms/internal/modules/home"
	"github.com/janmarkuslanger/nuricms/internal/modules/user"
	"github.com/janmarkuslanger/nuricms/internal/modules/webhook"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/pkg/config"
)

type App struct {
	Server   *server.Server
	Services *service.Set
	Config   *config.Config
}

func SetupApp(opts config.Config, envs env.EnvSource) (*App, error) {
	conf := SetDefaultConfig(opts)
	hooks := InitHookRegistry(conf.HookPlugins)

	env, err := LoadEnv(envs)

	db, err := InitDatabase(*conf.Dialector)
	if err != nil {
		return nil, err
	}

	MigrateModels(db)
	repos := repository.NewSet(db)

	fs := service.OsFileOps{}

	services, err := service.NewSet(repos, hooks, db, env, fs)
	if err != nil {
		return nil, err
	}
	InitAdminUser(services.User)

	s := server.NewServer()
	ctrl := []server.Controller{
		collection.NewController(services),
		field.NewController(services),
		content.NewController(services),
		asset.NewController(services),
		user.NewController(services),
		home.NewController(services),
		api.NewController(services),
		apikey.NewController(services),
		webhook.NewController(services),
	}
	InitController(ctrl, s)

	return &App{Server: s, Services: services, Config: &conf}, nil
}

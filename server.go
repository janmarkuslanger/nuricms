package nuricms

import (
	"log"
	"net/http"

	"github.com/janmarkuslanger/nuricms/internal/controller/api"
	"github.com/janmarkuslanger/nuricms/internal/controller/apikey"
	"github.com/janmarkuslanger/nuricms/internal/controller/asset"
	"github.com/janmarkuslanger/nuricms/internal/controller/collection"
	"github.com/janmarkuslanger/nuricms/internal/controller/content"
	"github.com/janmarkuslanger/nuricms/internal/controller/field"
	"github.com/janmarkuslanger/nuricms/internal/controller/home"
	"github.com/janmarkuslanger/nuricms/internal/controller/user"
	"github.com/janmarkuslanger/nuricms/internal/controller/webhook"
	"github.com/janmarkuslanger/nuricms/internal/core"
	"github.com/janmarkuslanger/nuricms/internal/db"
	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/pkg/plugin"
	"gorm.io/driver/sqlite"
)

type ServerConfig struct {
	Port        string
	HookPlugins []plugin.HookPlugin
}

func StartServer(config *ServerConfig) {
	if config.Port == "" {
		config.Port = "8080"
	}

	hr := plugin.NewHookRegistry()

	for _, p := range config.HookPlugins {
		p.Register(hr)
	}

	database, err := db.Init(sqlite.Open("nuricms.db"))
	if err != nil {
		panic("failed to init db")
	}

	repos := repository.NewSet(database)
	services, _ := service.NewSet(repos, hr)

	db.DB.AutoMigrate(
		&model.Collection{},
		&model.Field{},
		&model.Content{},
		&model.ContentValue{},
		&model.Asset{},
		&model.User{},
		&model.Apikey{},
		&model.Webhook{},
	)

	_, count, _ := services.User.List(1, 1)
	if count == 0 {
		services.User.Create(dto.UserData{
			Email:    "admin@admin.com",
			Password: "mysecret",
			Role:     string(model.RoleAdmin),
		})
	}

	modules := []core.Controller{
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

	server := server.NewServer()

	for _, module := range modules {
		module.RegisterRoutes(server)
	}

	log.Fatal(http.ListenAndServe(":"+config.Port, server.Mux))
}

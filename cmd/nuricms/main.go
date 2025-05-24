package main

import (
	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/controller/api"
	"github.com/janmarkuslanger/nuricms/controller/apikey"
	"github.com/janmarkuslanger/nuricms/controller/asset"
	"github.com/janmarkuslanger/nuricms/controller/collection"
	"github.com/janmarkuslanger/nuricms/controller/content"
	"github.com/janmarkuslanger/nuricms/controller/field"
	"github.com/janmarkuslanger/nuricms/controller/home"
	"github.com/janmarkuslanger/nuricms/controller/user"
	"github.com/janmarkuslanger/nuricms/controller/webhook"
	"github.com/janmarkuslanger/nuricms/core"
	"github.com/janmarkuslanger/nuricms/db"
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/repository"
	"github.com/janmarkuslanger/nuricms/service"
)

func main() {
	router := gin.Default()
	database := db.Init()
	repos := repository.NewSet(database)
	services := service.NewSet(repos)

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

	for _, module := range modules {
		module.RegisterRoutes(router)
	}

	router.Static("/public/assets", "./public/assets")

	router.Run("0.0.0.0:8080")
}

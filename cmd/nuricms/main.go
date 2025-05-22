package main

import (
	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/core/handler"
	"github.com/janmarkuslanger/nuricms/db"
	"github.com/janmarkuslanger/nuricms/handler/api"
	"github.com/janmarkuslanger/nuricms/handler/apikey"
	"github.com/janmarkuslanger/nuricms/handler/asset"
	"github.com/janmarkuslanger/nuricms/handler/collection"
	"github.com/janmarkuslanger/nuricms/handler/content"
	"github.com/janmarkuslanger/nuricms/handler/field"
	"github.com/janmarkuslanger/nuricms/handler/home"
	"github.com/janmarkuslanger/nuricms/handler/user"
	"github.com/janmarkuslanger/nuricms/handler/webhook"
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

	modules := []handler.Handler{
		collection.NewHandler(services),
		field.NewHandler(services),
		content.NewHandler(services),
		asset.NewHandler(services),
		user.NewHandler(services),
		home.NewHandler(services),
		api.NewHandler(services),
		apikey.NewHandler(services),
		webhook.NewHandler(services),
	}

	for _, module := range modules {
		module.RegisterRoutes(router)
	}

	router.Static("/public/assets", "./public/assets")

	router.Run("0.0.0.0:8080")
}

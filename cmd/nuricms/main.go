package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/core/handler"
	"github.com/janmarkuslanger/nuricms/db"
	"github.com/janmarkuslanger/nuricms/handler/asset"
	"github.com/janmarkuslanger/nuricms/handler/collection"
	"github.com/janmarkuslanger/nuricms/handler/content"
	"github.com/janmarkuslanger/nuricms/handler/field"
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/repository"
	"github.com/janmarkuslanger/nuricms/service"
	"github.com/janmarkuslanger/nuricms/utils"
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
	)

	modules := []handler.Handler{
		collection.NewHandler(services),
		field.NewHandler(services),
		content.NewHandler(services),
		asset.NewHandler(services),
	}

	for _, module := range modules {
		module.RegisterRoutes(router)
	}

	router.GET("/", func(c *gin.Context) {
		utils.RenderWithLayout(c, "home.tmpl", gin.H{}, http.StatusOK)
	})

	router.Run(":8080")
}

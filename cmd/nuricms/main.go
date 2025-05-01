package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/core/handler"
	"github.com/janmarkuslanger/nuricms/db"
	"github.com/janmarkuslanger/nuricms/entities/collection"
	"github.com/janmarkuslanger/nuricms/entities/content"
	"github.com/janmarkuslanger/nuricms/entities/field"
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
	)

	modules := []handler.Handler{
		collection.NewHandler(services),
		field.NewHandler(services),
		content.NewHandler(services),
	}

	for _, module := range modules {
		module.RegisterRoutes(router)
	}

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from nuricms 👋")
	})

	router.Run(":8080")
}

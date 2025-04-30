package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/core/db"
	"github.com/janmarkuslanger/nuricms/core/handler"
	"github.com/janmarkuslanger/nuricms/entities/collection"
	"github.com/janmarkuslanger/nuricms/entities/content"
	"github.com/janmarkuslanger/nuricms/entities/field"
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/repository"
)

func main() {
	router := gin.Default()
	database := db.Init()
	repos := repository.NewSet(database)

	db.DB.AutoMigrate(
		&model.Collection{},
		&model.Field{},
		&model.Content{},
		&model.ContentValue{},
	)

	modules := []handler.Handler{
		collection.NewHandler(repos),
		field.NewHandler(repos),
		content.NewHandler(repos),
	}

	for _, module := range modules {
		module.RegisterRoutes(router)
	}

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from nuricms 👋")
	})

	router.Run(":8080")
}

package handler

import "github.com/gin-gonic/gin"

type Handler interface {
	RegisterRoutes(r *gin.Engine)
}

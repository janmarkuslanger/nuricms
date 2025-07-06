package core

import (
	"github.com/janmarkuslanger/nuricms/internal/server"
)

type Controller interface {
	RegisterRoutes(s *server.Server)
}

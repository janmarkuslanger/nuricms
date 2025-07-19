package setup

import (
	"github.com/janmarkuslanger/nuricms/internal/server"
)

func InitController(ctrl []server.Controller, s *server.Server) {
	for _, c := range ctrl {
		c.RegisterRoutes(s)
	}
}

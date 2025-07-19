package server

type Controller interface {
	RegisterRoutes(server *Server)
}

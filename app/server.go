package app

//go:generate go run github.com/99designs/gqlgen
import (
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	router   *gin.Engine
	resolver *Resolver
}

func NewServer(router *gin.Engine, resolver *Resolver) *Server {
	server := &Server{
		router:   router,
		resolver: resolver,
	}

	server.routes()

	return server
}

func (s *Server) Run() error {
	err := s.router.Run(":8081")

	if err != nil {
		log.Printf("Server - there was an error calling Run on router: %v", err)
		return err
	}

	return nil
}

package core

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/rs/zerolog"
)

type Server struct {
	App    Application
	Config *AppConfig
	Logger *zerolog.Logger
	Router *gin.Engine
	Redis  *redis.Client
}

func (s *Server) Run(addr ...string) error {
	return s.Router.Run(addr...)
}

package api

import (
	"everyflavor/internal/core"
	v1 "everyflavor/internal/http/api/v1"
)

func SetupHandlers(s *core.Server) {
	api := s.Router.Group("/api")
	v1.SetupV1Handlers(api, s)
}

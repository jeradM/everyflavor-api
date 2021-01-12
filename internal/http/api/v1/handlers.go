package v1

import (
	"everyflavor/internal/core"

	"github.com/gin-gonic/gin"
)

func SetupV1Handlers(router gin.IRouter, server *core.Server) {
	v1 := router.Group("/v1")

	SetupAuthHandlers(v1, server)
	SetupBatchHandlers(v1, server.App)
	SetupFlavorHandlers(v1, server.App)
	setupRecipeHandlers(v1, server.App)
	SetupUserHandlers(v1, server.App)
	SetupVendorHandlers(v1, server.App)
}

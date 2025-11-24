package http

import (
	"net/http"

	"github.com/brunosprado/api-order-processor/domain"
	"github.com/brunosprado/api-order-processor/pkg/log"
	"github.com/gin-gonic/gin"
)

type handler struct {
	clientService domain.ClientService
	log           log.Logger
}

func NewHandler(clientService domain.ClientService, logger log.Logger) http.Handler {
	handler := &handler{
		clientService: clientService,
		log:           logger,
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.POST("/orders", handler.postOrder)

	return router
}

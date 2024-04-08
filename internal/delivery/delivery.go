package delivery

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qrcode-generation-service/internal/service"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) initAPI(router *gin.Engine) {
	handler := NewHandler(h.services)
	api := router.Group("/api")
	{
		handler.initHandlers(api)
	}
}

func (h *Handler) Init() *gin.Engine {

	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		corsMiddleware,
	)

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initHandlers(api *gin.RouterGroup) {

}

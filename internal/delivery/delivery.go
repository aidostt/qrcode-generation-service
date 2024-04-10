package delivery

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qrcode-generation-service/internal/service"
)

type Handler struct {
	services *service.Services
}

type QrCodeInput struct {
	Content string `json:"content"`
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
	api.POST("/generate", h.Generate)
}

func (h *Handler) Generate(c *gin.Context) {
	var inp QrCodeInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	qr, err := h.services.QrCode.GenerateQR(inp.Content)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//c.Header("Content-Type", "image/png")
	c.JSON(http.StatusOK, qr)
}

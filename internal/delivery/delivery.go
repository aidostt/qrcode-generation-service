package delivery

import (
	proto "github.com/aidostt/protos/gen/go/reservista/qr"
	"qrcode-generation-service/internal/service"
)

type Handler struct {
	proto.UnimplementedQRServer
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

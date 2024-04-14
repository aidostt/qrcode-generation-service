package delivery

import (
	"context"
	proto "github.com/aidostt/protos/gen/go/reservista/qr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"qrcode-generation-service/pkg/logger"
)

func (h *Handler) Generate(ctx context.Context, input *proto.GenerateRequest) (*proto.GenerateResponse, error) {
	if input.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}
	qr, err := h.services.QrCode.GenerateQR(input.GetContent())
	if err != nil {
		logger.Error(err)
		return nil, status.Error(codes.Internal, "failed to generate QR")
	}
	return &proto.GenerateResponse{QR: qr}, nil
}

func (h *Handler) Scan(ctx context.Context, input *proto.ScanRequest) (*proto.ScanResponse, error) {
	return nil, nil
}

package service

import (
	"fmt"
	"github.com/skip2/go-qrcode"
	"qrcode-generation-service/pkg/dialog"
)

const (
	size = 256
)

type Service struct {
	domain string
	dialog *dialog.Dialog
}

func NewGeneratorService(domain string, dial *dialog.Dialog) *Service {
	return &Service{
		domain: domain,
		dialog: dial,
	}
}

func (s *Service) GenerateQR(content string) ([]byte, error) {
	qrCode, err := qrcode.Encode(content, qrcode.Low, size)
	if err != nil {
		return nil, fmt.Errorf("could not generate a QR code: %v", err)
	}
	return qrCode, nil
}

func (s *Service) ScanQR(userID string, reservationID string) (user UserInfo, restaurant RestaurantInfo, err error) {

	//TODO: call method GetUserByID from another microservice
	//TODO: retrieve data from microservice (synchronously)
	//TODO: write it in userInfo structure
	return
}

package service

import (
	"fmt"
	"github.com/skip2/go-qrcode"
)

const size = 256

type GeneratorService struct {
	domain string
}

func NewGeneratorService(domain string) *GeneratorService {
	return &GeneratorService{
		domain: domain,
	}
}

func (s *GeneratorService) GenerateQR(content string) ([]byte, error) {
	qrCode, err := qrcode.Encode(content, qrcode.Low, size)
	if err != nil {
		return nil, fmt.Errorf("could not generate a QR code: %v", err)
	}
	return qrCode, nil
}

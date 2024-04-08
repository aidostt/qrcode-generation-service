package service

import (
	"fmt"
	"github.com/skip2/go-qrcode"
)

type GeneratorService struct {
	domain string
}

func NewGeneratorService(domain string) *GeneratorService {
	return &GeneratorService{
		domain: domain,
	}
}

func (s *GeneratorService) Generate() ([]byte, error) {

	qrCode, err := qrcode.Encode(code.Content, qrcode.Medium, code.Size)
	if err != nil {
		return nil, fmt.Errorf("could not generate a QR code: %v", err)
	}
	return qrCode, nil
}

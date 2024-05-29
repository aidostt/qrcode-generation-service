package service

import (
	"bytes"
	"context"
	"fmt"
	proto_reservation "github.com/aidostt/protos/gen/go/reservista/reservation"
	"github.com/aidostt/protos/gen/go/reservista/user"
	"github.com/nfnt/resize"
	"github.com/skip2/go-qrcode"
	"image"
	"image/draw"
	"image/png"
	"io"
	"qrcode-generation-service/internal/domain"
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

func (s *Service) ScanQR(ctx context.Context, userID string, reservationID string) (UserInfo, RestaurantInfo, ReservationInfo, error) {
	userConn, err := s.dialog.NewConnection(s.dialog.Addresses.Users)
	defer userConn.Close()
	if err != nil {
		return UserInfo{}, RestaurantInfo{}, ReservationInfo{}, err
	}
	userClient := proto_user.NewUserClient(userConn)
	userResponse, err := userClient.GetByID(ctx, &proto_user.GetRequest{UserId: userID})
	if err != nil {
		return UserInfo{}, RestaurantInfo{}, ReservationInfo{}, err
	}
	user := UserInfo{
		Name:    userResponse.GetName(),
		Surname: userResponse.GetSurname(),
		Phone:   userResponse.GetPhone(),
		Email:   userResponse.GetEmail(),
	}
	reservationConn, err := s.dialog.NewConnection(s.dialog.Addresses.Reservations)
	defer reservationConn.Close()
	if err != nil {
		return UserInfo{}, RestaurantInfo{}, ReservationInfo{}, err
	}
	reservationClient := proto_reservation.NewReservationClient(reservationConn)
	reservationResponse, err := reservationClient.GetReservation(ctx, &proto_reservation.IDRequest{Id: reservationID})
	if reservationResponse.GetUserID() != userID {
		return UserInfo{}, RestaurantInfo{}, ReservationInfo{}, domain.ErrUnauthorized
	}
	if err != nil {
		return UserInfo{}, RestaurantInfo{}, ReservationInfo{}, err
	}
	reservation := ReservationInfo{
		Table:           reservationResponse.Table.GetTableNumber(),
		ReservationTime: reservationResponse.GetReservationTime(),
	}
	restaurant := RestaurantInfo{
		Name:    reservationResponse.Table.Restaurant.GetName(),
		Contact: reservationResponse.Table.Restaurant.GetContact(),
		Address: reservationResponse.Table.Restaurant.GetAddress(),
	}
	return user, restaurant, reservation, nil
}

func (s *Service) GenerateQRWithWatermark(watermark []byte, content string) ([]byte, error) {
	qrCode, err := s.GenerateQR(content)
	if err != nil {
		return nil, err
	}

	qrCode, err = s.AddWatermark(qrCode, watermark)
	if err != nil {
		return nil, fmt.Errorf("could not add watermark to QR code: %v", err)
	}

	return qrCode, nil
}

func (s *Service) AddWatermark(qrCode []byte, watermarkData []byte) ([]byte, error) {
	qrCodeData, err := png.Decode(bytes.NewBuffer(qrCode))
	if err != nil {
		return nil, fmt.Errorf("could not decode QR code: %v", err)
	}

	watermarkWidth := uint(float64(qrCodeData.Bounds().Dx()) * 0.25)
	watermark, err := s.ResizeWatermark(bytes.NewBuffer(watermarkData), watermarkWidth)
	if err != nil {
		return nil, fmt.Errorf("could not resize the watermark image: %v", err)
	}

	watermarkImage, err := png.Decode(bytes.NewBuffer(watermark))
	if err != nil {
		return nil, fmt.Errorf("could not decode watermark: %v", err)
	}

	halfQrCodeWidth, halfWatermarkWidth := qrCodeData.Bounds().Dx()/2, watermarkImage.Bounds().Dx()/2
	offset := image.Pt(
		halfQrCodeWidth-halfWatermarkWidth,
		halfQrCodeWidth-halfWatermarkWidth,
	)

	watermarkImageBounds := qrCodeData.Bounds()
	m := image.NewRGBA(watermarkImageBounds)

	draw.Draw(m, watermarkImageBounds, qrCodeData, image.Point{}, draw.Src)
	draw.Draw(
		m,
		watermarkImage.Bounds().Add(offset),
		watermarkImage,
		image.Point{},
		draw.Over,
	)

	watermarkedQRCode := bytes.NewBuffer(nil)
	png.Encode(watermarkedQRCode, m)

	return watermarkedQRCode.Bytes(), nil
}

func (s *Service) ResizeWatermark(watermark io.Reader, width uint) ([]byte, error) {
	decodedImage, err := png.Decode(watermark)
	if err != nil {
		return nil, fmt.Errorf("could not decode watermark image: %v", err)
	}

	m := resize.Resize(width, 0, decodedImage, resize.Lanczos3)
	resized := bytes.NewBuffer(nil)
	png.Encode(resized, m)

	return resized.Bytes(), nil
}

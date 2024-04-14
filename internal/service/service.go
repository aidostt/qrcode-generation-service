package service

import "qrcode-generation-service/pkg/dialog"

type RestaurantInfo struct {
	Name            string
	Address         string
	Table           int32
	ReservationTime string
}

type UserInfo struct {
	Name    string
	Surname string
	Phone   string
	Email   string
}

type QRCodeInput struct {
	Restaurant RestaurantInfo
	User       UserInfo
	Size       int
}

type QRCode interface {
	GenerateQR(string) ([]byte, error)
	ScanQR(string, string) (UserInfo, RestaurantInfo, error)
	//AddWatermark(ctx context.Context) ([]byte, error)
}

type Services struct {
	QrCode QRCode
}

type Dependencies struct {
	Environment string
	Domain      string
	Dialog      *dialog.Dialog
}

func NewServices(deps Dependencies) *Services {
	generatorService := NewGeneratorService(deps.Domain, deps.Dialog)
	return &Services{
		QrCode: generatorService,
	}
}

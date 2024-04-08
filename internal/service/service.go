package service

type RestaurantInfo struct {
	Name    string
	Address string
	Table   int
}

type UserInfo struct {
	Name    string
	Surname string
	Phone   string
}

type QRCodeInput struct {
	Restaurant RestaurantInfo
	User       UserInfo
	Size       int
}

type QRCode interface {
	Generate() ([]byte, error)
	//AddWatermark(ctx context.Context) ([]byte, error)
}

type Services struct {
	QrCode QRCode
}

type Dependencies struct {
	Environment string
	Domain      string
}

func NewServices(deps Dependencies) *Services {
	generatorService := NewGeneratorService(deps.Domain)
	return &Services{
		QrCode: generatorService,
	}
}

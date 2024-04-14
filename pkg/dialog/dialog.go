package dialog

import (
	"google.golang.org/grpc"
	"qrcode-generation-service/pkg/logger"
)

type DialogService interface {
	NewConnections(string) (*grpc.ClientConn, error)
}

type Dialog struct {
	Addresses Addresses
	authority string
}
type Addresses struct {
	Users        string
	Reservations string
}

func NewDialog(authority, users, reservations string) *Dialog {
	return &Dialog{authority: authority, Addresses: Addresses{Users: users, Reservations: reservations}}
}

func (d *Dialog) NewConnections(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithAuthority(d.authority))
	if err != nil {
		logger.Errorf("Failed to connect: %v", err)
		conn.Close()
		return nil, err
	}
	return conn, nil
}

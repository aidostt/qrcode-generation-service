package app

import (
	"errors"
	"fmt"
	proto "github.com/aidostt/protos/gen/go/reservista/qr"
	"net"
	"net/http"
	"os"
	"os/signal"
	"qrcode-generation-service/internal/config"
	"qrcode-generation-service/internal/delivery"
	"qrcode-generation-service/internal/server"
	"qrcode-generation-service/internal/service"
	"qrcode-generation-service/pkg/dialog"
	"qrcode-generation-service/pkg/logger"
	"syscall"
)

func Run(configPath, envPath string) {
	cfg, err := config.Init(configPath, envPath)
	if err != nil {
		logger.Error(err)

		return
	}
	// Dependencies
	dial := dialog.NewDialog(cfg.Authority, fmt.Sprintf("%v:%v", cfg.Users.Host, cfg.Users.Port), fmt.Sprintf("%v:%v", cfg.Reservations.Host, cfg.Reservations.Port))
	services := service.NewServices(service.Dependencies{
		Environment: cfg.Environment,
		Domain:      cfg.GRPC.Host,
		Dialog:      dial,
	})
	handlers := delivery.NewHandler(services)
	// GRPC Server
	srv := server.NewServer()
	proto.RegisterQRServer(srv.GrpcServer, handlers)
	l, err := net.Listen("tcp", fmt.Sprintf("%v:%v", cfg.GRPC.Host, cfg.GRPC.Port))
	if err != nil {
		logger.Errorf("error occurred while getting listener for the server: %s\n", err.Error())
		return
	}
	go func() {
		if err := srv.Run(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running grpc server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started at: " + cfg.GRPC.Host + ":" + cfg.GRPC.Port)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	srv.Stop()
}

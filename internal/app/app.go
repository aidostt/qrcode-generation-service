package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"qrcode-generation-service/internal/config"
	"qrcode-generation-service/internal/delivery"
	"qrcode-generation-service/internal/server"
	"qrcode-generation-service/internal/service"
	"qrcode-generation-service/pkg/logger"
	"syscall"
	"time"
)

func Run(configPath, envPath string) {
	cfg, err := config.Init(configPath, envPath)
	if err != nil {
		logger.Error(err)

		return
	}

	// Dependencies

	services := service.NewServices(service.Dependencies{
		Environment: cfg.Environment,
		Domain:      cfg.HTTP.Host,
	})
	handlers := delivery.NewHandler(services)

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init())

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
}

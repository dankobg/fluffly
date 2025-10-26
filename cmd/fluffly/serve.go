//go:generate go run ../../db/generator/generator.go
package fluffly

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dankobg/fluffly/auth/keto"
	"github.com/dankobg/fluffly/auth/kratos"
	"github.com/dankobg/fluffly/config"
	"github.com/dankobg/fluffly/geocoding/nominatim"
	"github.com/dankobg/fluffly/httpserver"
	"github.com/dankobg/fluffly/logging"
	"github.com/dankobg/fluffly/mailer"
	"github.com/dankobg/fluffly/media"
	"github.com/dankobg/fluffly/media/local"
	"github.com/dankobg/fluffly/media/minio"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/server"
)

type ServeCommand struct{}

func (s *ServeCommand) Run() error {
	cfg, _, err := config.New()
	if err != nil {
		slog.Error("failed to initialize config", slog.Any("error", err))
		return err
	}

	logger := logging.New(
		logging.WithConsolePretty(cfg.ENV != "production" && cfg.Logger.Pretty),
		logging.WithLevel(slog.LevelDebug),
	)

	smtpClient := mailer.NewSmtpClient(
		mailer.WithEnabled(cfg.ENV == "production"),
		mailer.WithDevHost(cfg.Email.DevSMTPHost),
		mailer.WithDevPort(cfg.Email.DevSMTPPort),
		mailer.WithDevUsername(cfg.Email.DevSMTPUsername),
		mailer.WithDevPassword(cfg.Email.DevSMTPPassword),
		mailer.WithHost(cfg.Email.SMTPHost),
		mailer.WithPort(cfg.Email.SMTPPort),
		mailer.WithUsername(cfg.Email.SMTPUsername),
		mailer.WithPassword(cfg.Email.SMTPPassword),
		mailer.WithTLS(true),
		mailer.WithFromName(cfg.Email.FromName),
		mailer.WithFromAddress(cfg.Email.FromAddress),
		mailer.WithLog(logger),
	)

	// rdb, err := redis.New(cfg.Redis)
	// if err != nil {
	// 	return fmt.Errorf("failed to connect to redis: %w", err)
	// }

	kratosClient := kratos.NewClient(cfg.KratosPublicURL, cfg.KratosAdminURL)
	ketoClient, err := keto.NewClient()
	if err != nil {
		return err
	}

	db, err := postgres.Connect(cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	pg := postgres.New(db)

	var upl media.Uploader
	switch cfg.FileStorage {
	case local.StorageKind:
		upl, err = local.NewLocalUploader(cfg.BaseURL+"/uploads", cfg.UploadDir)
		if err != nil {
			fmt.Println("failed to init local uploader: %w", err)
		}
	case minio.StorageKind:
		upl, err = minio.NewMinioUploader(cfg.Minio)
		if err != nil {
			fmt.Println("failed to init minio uploader: %w", err)
		}
	default:
		panic("unknown file storage: " + cfg.FileStorage)
	}

	httpc := httpserver.NewHttpClient()

	geoc, err := nominatim.NewNominatimGeocoder(httpc)
	if err != nil {
		return fmt.Errorf("failed to init a geocoder: %w", err)
	}

	apiHandler := server.New(cfg, logger, kratosClient, ketoClient, smtpClient, pg, upl, geoc, httpc)

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	defer stop()

	h := apiHandler.SetupRoutes(cfg.ENV, cfg.UploadDir)

	srv := httpserver.New(
		httpserver.WithHostPort("", cfg.Port),
		httpserver.WithHandler(h),
		httpserver.WithReadTimeout(cfg.Server.ReadTimeout),
		httpserver.WithReadHeaderTimeout(cfg.Server.ReadHeaderTimeout),
		httpserver.WithWriteTimeout(cfg.Server.WriteTimeout),
		httpserver.WithIdleTimeout(cfg.Server.IdleTimeout),
		httpserver.WithErrorSlog(logger, slog.LevelDebug),
		httpserver.WithBaseContext(func(l net.Listener) context.Context { return rootCtx }),
	)

	logger.Info("fluffly info", slog.String("env", cfg.ENV), slog.String("website_url", cfg.WebsiteURL), slog.String("logger_level", cfg.Logger.Level), slog.Bool("mailer_enabled", cfg.Email.Enabled))

	srvErr := make(chan error, 1)

	go func() {
		logger.Info("server is listening", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed to start", slog.Any("error", err))
			srvErr <- err
		}
	}()

	select {
	case <-rootCtx.Done():
		logger.Info("received interruption signal")
	case err := <-srvErr:
		logger.Error("received server err", slog.Any("error", err))
		stop()
	}

	logger.Info("starting shutdown", slog.Duration("graceful_timeout", cfg.Server.GracefulTimeout))
	shutdownCtx, cancel := context.WithTimeoutCause(context.Background(), cfg.Server.GracefulTimeout, fmt.Errorf("graceful shutdown timeout"))
	defer cancel()

	var shutdownErrors []error

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", slog.Any("error", err))
		shutdownErrors = append(shutdownErrors, err)
	}

	logger.Info("server shut down")
	return errors.Join(shutdownErrors...)
}

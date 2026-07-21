//go:generate go tool -modfile=../../tools.mod bobgen-psql -c ../../bobgen.yaml

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
	"github.com/dankobg/fluffly/geocoding"
	"github.com/dankobg/fluffly/geocoding/fakegeo"
	"github.com/dankobg/fluffly/geocoding/nominatim"
	"github.com/dankobg/fluffly/httpserver"
	"github.com/dankobg/fluffly/logging"
	"github.com/dankobg/fluffly/mailer"
	"github.com/dankobg/fluffly/media"
	"github.com/dankobg/fluffly/media/local"
	"github.com/dankobg/fluffly/media/rustfs"
	"github.com/dankobg/fluffly/persistence/postgres"
	"github.com/dankobg/fluffly/server"
)

type ServeCommand struct{}

func (sc *ServeCommand) Run() error {
	cfg, _, err := config.New()
	if err != nil {
		slog.Error("failed to initialize config", slog.Any("error", err))
		return err
	}

	logger := logging.New(
		logging.WithConsolePretty(cfg.App.ENV != "production" && cfg.Logger.Pretty),
		logging.WithLevel(slog.LevelDebug),
	)

	smtpClient := mailer.NewSmtpClient(
		mailer.WithEnabled(cfg.App.ENV == "production"),
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

	kratosClient := kratos.NewClient(cfg.App.KratosPublicURL, cfg.App.KratosAdminURL)

	ketoClient, err := keto.NewClient(cfg.App.KetoReadURL, cfg.App.KetoWriteURL)
	if err != nil {
		return err
	}

	pool, err := postgres.NewPool(context.Background(), cfg.Database)
	if err != nil {
		return fmt.Errorf("postgres.NewPool: %w", err)
	}
	defer pool.Close()

	pg := postgres.New(pool)

	var upl media.Uploader

	switch cfg.App.FileStorage {
	case local.StorageKind:
		upl, err = local.NewLocalUploader(cfg.App.BaseURL+"/uploads", cfg.App.UploadDir)
		if err != nil {
			return fmt.Errorf("failed to init local uploader: %w", err)
		}
	case rustfs.StorageKind:
		upl, err = rustfs.NewRustfsUploader(cfg.Rustfs)
		if err != nil {
			return fmt.Errorf("failed to init rustfs uploader: %w", err)
		}
	default:
		panic("unknown file storage: " + cfg.App.FileStorage)
	}

	httpc := httpserver.NewHttpClient()

	var geoc geocoding.Geocoder = fakegeo.FakeGeocoder{}
	if cfg.Geocoding.Enabled {
		geoc, err = nominatim.NewNominatimGeocoder(httpc)
		if err != nil {
			return fmt.Errorf("failed to init a geocoder: %w", err)
		}
	}

	apiHandler := server.New(cfg, logger, kratosClient, ketoClient, smtpClient, pg, upl, geoc, httpc)
	if err := apiHandler.PrecompileSpeciesPropertiesJsonSchemas(context.Background()); err != nil {
		return fmt.Errorf("failed to precompile species properties json schemas: %w", err)
	}

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	defer stop()

	h := apiHandler.SetupRoutes(cfg.App.ENV, cfg.App.UploadDir)

	srv := httpserver.New(
		httpserver.WithHostPort("", cfg.App.Port),
		httpserver.WithHandler(h),
		httpserver.WithReadTimeout(cfg.Server.ReadTimeout),
		httpserver.WithReadHeaderTimeout(cfg.Server.ReadHeaderTimeout),
		httpserver.WithWriteTimeout(cfg.Server.WriteTimeout),
		httpserver.WithIdleTimeout(cfg.Server.IdleTimeout),
		httpserver.WithErrorSlog(logger, slog.LevelDebug),
		httpserver.WithBaseContext(func(l net.Listener) context.Context { return rootCtx }),
	)

	logger.Info("fluffly info", slog.String("env", cfg.App.ENV), slog.String("website_url", cfg.App.WebsiteURL), slog.String("logger_level", cfg.Logger.Level), slog.Bool("mailer_enabled", cfg.Email.Enabled), slog.Bool("geocoding_enabled", cfg.Geocoding.Enabled))

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

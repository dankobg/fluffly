package server

import (
	"log/slog"

	"github.com/dankobg/fluffly/auth/keto"
	"github.com/dankobg/fluffly/auth/kratos"
	"github.com/dankobg/fluffly/config"
	"github.com/dankobg/fluffly/mailer"
	"github.com/dankobg/fluffly/persistence"
	"github.com/redis/go-redis/v9"
)

// var _ api.StrictServerInterface = (*ApiHandler)(nil)

type ApiHandler struct {
	Cfg       *config.Config
	Log       *slog.Logger
	Kratos    *kratos.Client
	Keto      *keto.Client
	Rdb       *redis.Client
	persistor persistence.Persistor
	Mailer    mailer.Mailer
}

func New(cfg *config.Config, log *slog.Logger, kratos *kratos.Client, keto *keto.Client, mailer mailer.Mailer, p persistence.Persistor) *ApiHandler {
	return &ApiHandler{
		Cfg:       cfg,
		Log:       log,
		Kratos:    kratos,
		Keto:      keto,
		persistor: p,
		Mailer:    mailer,
	}
}

package server

import (
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/dankobg/fluffly/auth/keto"
	"github.com/dankobg/fluffly/auth/kratos"
	"github.com/dankobg/fluffly/config"
	"github.com/dankobg/fluffly/geocoding"
	"github.com/dankobg/fluffly/mailer"
	"github.com/dankobg/fluffly/media"
	"github.com/dankobg/fluffly/persistence"
	"github.com/redis/go-redis/v9"
)

// var _ api.StrictServerInterface = (*ApiHandler)(nil)

type ApiHandler struct {
	Cfg        *config.Config
	Log        *slog.Logger
	Kratos     *kratos.Client
	Keto       *keto.Client
	Rdb        *redis.Client
	persistor  persistence.Persistor
	Mailer     mailer.Mailer
	openapiTpl *template.Template
	uploader   media.Uploader
	httpc      *http.Client
	geocoder   geocoding.Geocoder
}

func New(cfg *config.Config, log *slog.Logger, kratos *kratos.Client, keto *keto.Client, mailer mailer.Mailer, p persistence.Persistor, upl media.Uploader, g geocoding.Geocoder) *ApiHandler {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.IdleConnTimeout = 60 * time.Second
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	httpc := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	return &ApiHandler{
		Cfg:       cfg,
		Log:       log,
		Kratos:    kratos,
		Keto:      keto,
		persistor: p,
		Mailer:    mailer,
		uploader:  upl,
		httpc:     httpc,
		geocoder:  g,
	}
}

func (a *ApiHandler) SetOpenapiTemplates(tpl *template.Template) {
	a.openapiTpl = tpl
}

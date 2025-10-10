package server

import (
	"html/template"
	"log/slog"
	"net/http"

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

func New(
	cfg *config.Config,
	log *slog.Logger,
	kratos *kratos.Client,
	keto *keto.Client,
	mailer mailer.Mailer,
	p persistence.Persistor,
	upl media.Uploader,
	g geocoding.Geocoder,
	c *http.Client,
) *ApiHandler {
	if c == nil {
		c = http.DefaultClient
	}
	return &ApiHandler{
		Cfg:       cfg,
		Log:       log,
		Kratos:    kratos,
		Keto:      keto,
		persistor: p,
		Mailer:    mailer,
		uploader:  upl,
		httpc:     c,
		geocoder:  g,
	}
}

func (a *ApiHandler) SetOpenapiTemplates(tpl *template.Template) {
	a.openapiTpl = tpl
}

package server

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/auth/keto"
	"github.com/dankobg/fluffly/auth/kratos"
	"github.com/dankobg/fluffly/config"
	"github.com/dankobg/fluffly/geocoding"
	"github.com/dankobg/fluffly/mailer"
	"github.com/dankobg/fluffly/media"
	"github.com/dankobg/fluffly/persistence"
	"github.com/dankobg/fluffly/ptr"
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

func newNotFoundErr(code, message string, reason ...string) api.APIError {
	e := api.APIError{
		Code:       fmt.Sprintf("ERR_%s", strings.ToUpper(code)),
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
	if len(reason) > 0 {
		e.Reason = ptr.Of(reason[0])
	}
	return e
}

func newUnauthenticatedErr(code, message string, reason ...string) api.APIError {
	e := api.APIError{
		Code:       fmt.Sprintf("ERR_%s", strings.ToUpper(code)),
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
	if len(reason) > 0 && reason[0] != "" {
		e.Reason = ptr.Of(reason[0])
	}
	return e
}

func newUnauthorizedErr(code, message string, reason ...string) api.APIError {
	e := api.APIError{
		Code:       fmt.Sprintf("ERR_%s", strings.ToUpper(code)),
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
	if len(reason) > 0 && reason[0] != "" {
		e.Reason = ptr.Of(reason[0])
	}
	return e
}

func newGenericErr(statusCode int32, code, message string, reason ...string) api.APIError {
	e := api.APIError{
		Code:       fmt.Sprintf("ERR_%s", strings.ToUpper(code)),
		Message:    message,
		StatusCode: statusCode,
	}
	if len(reason) > 0 && reason[0] != "" {
		e.Reason = ptr.Of(reason[0])
	}
	return e
}

func newNotFoundResp(code, message string, reason ...string) api.NotFoundErrorResponseJSONResponse {
	e := newNotFoundErr(code, message, reason...)
	return api.NotFoundErrorResponseJSONResponse(e)
}

func newUnauthenticatedResp(code, message string, reason ...string) api.UnauthenticatedErrorResponse {
	e := newUnauthenticatedErr(code, message, reason...)
	return api.UnauthenticatedErrorResponse(e)
}

func newUnauthorizedResp(code, message string, reason ...string) api.UnauthorizedErrorResponseJSONResponse {
	e := newUnauthorizedErr(code, message, reason...)
	return api.UnauthorizedErrorResponseJSONResponse(e)
}

func newGenericResp(statusCode int32, code, message string, reason ...string) api.GenericErrorResponseJSONResponse {
	e := newGenericErr(statusCode, code, message, reason...)
	return api.GenericErrorResponseJSONResponse(e)
}

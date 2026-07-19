package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/auth/keto"
	"github.com/dankobg/fluffly/auth/kratos"
	"github.com/dankobg/fluffly/config"
	"github.com/dankobg/fluffly/geocoding"
	"github.com/dankobg/fluffly/mailer"
	"github.com/dankobg/fluffly/media"
	"github.com/dankobg/fluffly/persistence"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/redis/go-redis/v9"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

// var _ api.StrictServerInterface = (*ApiHandler)(nil)

type ApiHandler struct {
	Cfg         *config.Config
	Log         *slog.Logger
	Kratos      *kratos.Client
	Keto        *keto.Client
	Rdb         *redis.Client
	persistor   persistence.Persistor
	Mailer      mailer.Mailer
	openapiTpl  *template.Template
	uploader    media.Uploader
	httpc       *http.Client
	geocoder    geocoding.Geocoder
	schemaCache map[string]*jsonschema.Schema
	mu          sync.RWMutex
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
		Cfg:         cfg,
		Log:         log,
		Kratos:      kratos,
		Keto:        keto,
		persistor:   p,
		Mailer:      mailer,
		uploader:    upl,
		httpc:       c,
		geocoder:    g,
		schemaCache: make(map[string]*jsonschema.Schema),
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
		e.Reason = new(reason[0])
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
		e.Reason = new(reason[0])
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
		e.Reason = new(reason[0])
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
		e.Reason = new(reason[0])
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

func (a *ApiHandler) PrecompileSpeciesPropertiesJsonSchemas(ctx context.Context) error {
	species, err := a.persistor.Animal().ListAnimalSpecies(ctx, dbtype.ListAnimalSpeciesFilters{})
	if err != nil {
		return err
	}

	for _, sp := range species.Data {
		b, err := json.Marshal(sp.PropertiesSchema)
		if err != nil {
			return fmt.Errorf("failed to marshal properties json schema")
		}

		name := fmt.Sprintf("properties-schema-%d.json", sp.ID)

		compiler := jsonschema.NewCompiler()
		if err := compiler.AddResource(name, bytes.NewReader(b)); err != nil {
			return fmt.Errorf("failed to add compiler jsonschema resource")
		}

		schema, err := compiler.Compile(name)
		if err != nil {
			return fmt.Errorf("failed to compile jsonschema")
		}

		a.schemaCache[name] = schema
	}

	return nil
}

func (a *ApiHandler) LoadSpeciesPropertiesJsonSchema(specieID int64) (*jsonschema.Schema, error) {
	name := fmt.Sprintf("properties-schema-%d.json", specieID)

	a.mu.RLock()
	s, ok := a.schemaCache[name]
	a.mu.RUnlock()

	if ok {
		return s, nil
	}

	schema, err := jsonschema.Compile(name)
	if err != nil {
		return nil, err
	}

	a.mu.Lock()
	a.schemaCache[name] = schema
	a.mu.Unlock()

	return schema, nil
}

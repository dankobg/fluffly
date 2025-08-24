package server

import (
	"expvar"
	"html/template"
	"net/http"
	"net/http/pprof"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/data"
)

func (a *ApiHandler) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// debug routes
	mux.Handle("GET /debug/vars", expvar.Handler())
	mux.HandleFunc("GET /debug/pprof/", pprof.Index)
	mux.Handle("GET /debug/pprof/allocs", pprof.Handler("allocs"))
	mux.Handle("GET /debug/pprof/block", pprof.Handler("block"))
	mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
	mux.Handle("GET /debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("GET /debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("GET /debug/pprof/mutex", pprof.Handler("mutex"))
	mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("POST /debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
	mux.Handle("GET /debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)

	// static files
	mux.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.FS(data.MustPublicFS()))))

	cors, err := NewCORS(a.Cfg.Cors)
	if err != nil {
		panic("middleware failed: " + err.Error())
	}

	middlewareChain := MiddlewareChain(
		PanicRecovery,
		RequestID,
		BodyLimit(10<<20),
		cors,
		a.AttachSessionData,
	)
	// @TODO: add rate limit mw

	mux.HandleFunc("GET /api/v1/lol", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("lol route"))
	})

	mux.HandleFunc("POST /api/v1/webhooks/kratos/registration_after_password", a.registrationAfterPassword)
	mux.HandleFunc("POST /api/v1/webhooks/kratos/registration_after_oidc", a.registrationAfterOidc)

	openapi, err := api.GetSwagger()
	if err != nil {
		panic("error loading openapi spec: " + err.Error())
	}

	openapiB, err := openapi.MarshalJSON()
	if err != nil {
		panic("failed to marshal oapi schema to json: " + err.Error())
	}
	openapiTpl, err := template.ParseFS(data.MustTemplatesFS(), "openapi/*")
	if err != nil {
		panic("failed to parse openapi templates: " + err.Error())
	}
	a.SetOpenapiTemplates(openapiTpl)

	mux.HandleFunc("GET /spec", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(openapiB)
	})
	mux.HandleFunc("GET /docs/rapidoc", a.openapiRapidocPage)
	mux.HandleFunc("GET /docs/redoc", a.openapiRedocPage)
	mux.HandleFunc("GET /docs/stoplight", a.openapiStoplightPage)
	mux.HandleFunc("GET /docs/scalar", a.openapiScalarPage)
	mux.HandleFunc("GET /docs/swagger", a.openapiSwaggerPage)

	apiSrv := api.NewStrictHandler(a, make([]api.StrictMiddlewareFunc, 0))
	oapiHandler := api.HandlerFromMuxWithBaseURL(apiSrv, mux, "/api/v1")
	return middlewareChain(oapiHandler)
}

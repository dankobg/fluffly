package server

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/dankobg/fluffly/config"
	"github.com/google/uuid"
	"github.com/jub0bs/cors"
)

type Middleware func(http.Handler) http.Handler

func MiddlewareChain(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
			r.Header.Set("X-Request-ID", reqID)
		}
		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r)
	})
}

func BodyLimit(limit int64) func(http.Handler) http.Handler {
	type fileRoute struct {
		path   string
		method string
	}
	fileRoutes := []fileRoute{
		{path: "/api/v1/organizations", method: http.MethodPost},
		{path: "/api/v1/animals", method: http.MethodPost},
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
				if limit == 0 {
					limit = 10 << 20
				}
				if slices.ContainsFunc(fileRoutes, func(fr fileRoute) bool {
					return r.Method == fr.method && r.URL.Path == fr.path
				}) {
					limit = 100 << 20
				}
			}
			r.Body = http.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}

func NewCORS(cfg config.CorsConfig) (func(http.Handler) http.Handler, error) {
	corsCfg := cors.Config{
		Origins:         cfg.AllowOrigins,
		Credentialed:    cfg.AllowCredentials,
		Methods:         cfg.AllowMethods,
		RequestHeaders:  cfg.AllowHeaders,
		MaxAgeInSeconds: cfg.MaxAge,
		ResponseHeaders: cfg.ExposeHeaders,
		ExtraConfig:     cors.ExtraConfig{},
	}
	corsMw, err := cors.NewMiddleware(corsCfg)
	if err != nil {
		return nil, err
	}
	mw := func(next http.Handler) http.Handler {
		return corsMw.Wrap(next)
	}
	return mw, nil
}

package server

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/dankobg/fluffly/config"
	"github.com/google/uuid"
	"github.com/rs/cors"
)

type Middleware func(http.Handler) http.Handler

func MiddlewareChain(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, x := range slices.Backward(xs) {
			next = x(next)
		}

		return next
	}
}

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("fluffly recover: %s\n", string(debug.Stack()))
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
	const (
		normalLimit = 10 << 20
		fileLimit   = 100 << 20
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqLimit := limit
			if reqLimit == 0 {
				reqLimit = normalLimit
			}

			if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
				reqLimit = fileLimit
			}

			r.Body = http.MaxBytesReader(w, r.Body, reqLimit)
			next.ServeHTTP(w, r)
		})
	}
}

func newCORS(cfg config.CorsConfig) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowOrigins,
		AllowedMethods:   cfg.AllowMethods,
		AllowedHeaders:   cfg.AllowHeaders,
		ExposedHeaders:   cfg.ExposeHeaders,
		MaxAge:           cfg.MaxAge,
		AllowCredentials: cfg.AllowCredentials,
		// AllowPrivateNetwork:  false,
		OptionsPassthrough: true,
		// OptionsSuccessStatus: 0,
		Debug:  cfg.Debug,
		Logger: nil,
	})
}

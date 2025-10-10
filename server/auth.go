package server

import (
	"context"
	"net/http"
	"strings"

	orykratos "github.com/ory/client-go"
)

type contextKey string

const (
	oryKratosSessionCtxKey contextKey = "ory_kratos_session_key"
)

const (
	oryKratosCsrfCookiePrefix  = "csrf_token"
	oryKratosSessionCookieName = "ory_kratos_session"

	prefixBearer = "Bearer"
)

type authHeadersResult struct {
	csrf         *http.Cookie
	session      *http.Cookie
	authHeader   string
	cookieHeader string
}

func ExtractAuthHeadersFromRequest(r *http.Request) *authHeadersResult {
	result := &authHeadersResult{
		cookieHeader: r.Header.Get("Cookie"),
		authHeader:   r.Header.Get("Authorization"),
	}

	for _, c := range r.Cookies() {
		if c != nil {
			if ok := strings.HasPrefix(c.Name, oryKratosCsrfCookiePrefix); ok {
				result.csrf = c
			}
		}
	}

	sessionCookie, _ := r.Cookie(oryKratosSessionCookieName)
	if sessionCookie != nil {
		result.session = sessionCookie
	}

	return result
}

func WithSession(ctx context.Context, sess *orykratos.Session) context.Context {
	return context.WithValue(ctx, oryKratosSessionCtxKey, sess)
}

func GetSession(ctx context.Context) *orykratos.Session {
	sess, ok := ctx.Value(oryKratosSessionCtxKey).(*orykratos.Session)
	if !ok {
		return nil
	}
	return sess
}

func (a *ApiHandler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := GetSession(r.Context())
		if sess == nil {
			http.Error(w, "session is required", http.StatusUnauthorized)
			return
		}
		if sess.Active != nil && !*sess.Active {
			http.Error(w, "session is invalid or has already expired", http.StatusUnauthorized)
		}
		next.ServeHTTP(w, r)
	})
}

func (a *ApiHandler) RequireAnonymous(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := GetSession(r.Context())
		if sess != nil && sess.Active != nil && *sess.Active {
			http.Error(w, "must have no session", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *ApiHandler) AttachSessionData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		info := ExtractAuthHeadersFromRequest(r)

		var hasAuthHeader bool
		var hasCookieHeader bool

		if info.authHeader != "" && strings.HasPrefix(info.authHeader, prefixBearer) {
			hasAuthHeader = true
		}
		if info.cookieHeader != "" {
			hasCookieHeader = true
		}

		if !hasAuthHeader && !hasCookieHeader {
			next.ServeHTTP(w, r)
			return
		}

		toSessionReq := a.Kratos.Public.FrontendAPI.ToSession(ctx).Cookie(info.session.String())
		session, sessionResp, err := toSessionReq.Execute()
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		defer sessionResp.Body.Close()

		if session != nil && session.Active != nil && !*session.Active {
			next.ServeHTTP(w, r)
			return
		}

		if session != nil {
			ctxWithSession := WithSession(ctx, session)
			req := r.WithContext(ctxWithSession)
			next.ServeHTTP(w, req)
			return
		}

		next.ServeHTTP(w, r)
	})
}

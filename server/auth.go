package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
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
		if !sessionExists(sess) {
			http.Error(w, "session is required", http.StatusUnauthorized)
			return
		}
		if !sessionValid(sess) {
			http.Error(w, "session is invalid or has already expired", http.StatusUnauthorized)
		}
		next.ServeHTTP(w, r)
	})
}

func (a *ApiHandler) RequireAnonymous(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := GetSession(r.Context())
		if sessionValid(sess) {
			http.Error(w, "must have no session", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func sessionExists(sess *orykratos.Session) bool {
	return sess != nil
}

func sessionValid(sess *orykratos.Session) bool {
	return sess != nil && sess.Active != nil && *sess.Active
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

func authFunc(ctx context.Context, in *openapi3filter.AuthenticationInput) error {
	if in.SecuritySchemeName == "cookieAuth" && in.SecurityScheme.Name == "ory_kratos_session" {
		sess := GetSession(ctx)
		if !sessionExists(sess) {
			return fmt.Errorf("session is required")
		}
		if !sessionValid(sess) {
			return fmt.Errorf("session is invalid or has already expired")
		}
	}
	return nil
}

func authzIdentityID(id string) string {
	return fmt.Sprintf("identity:%s", id)
}

func authzCourierMessageID(id string) string {
	return fmt.Sprintf("courier_message:%s", id)
}

func authzSchemaID(id string) string {
	return fmt.Sprintf("schema:%s", id)
}

func authzSessionID(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func authzCountryID(id int64) string {
	return fmt.Sprintf("country:%d", id)
}

func authzAnimalID(id int64) string {
	return fmt.Sprintf("animal:%d", id)
}

func authzOrganizationID(id int64) string {
	return fmt.Sprintf("organization:%d", id)
}

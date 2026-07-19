package server

import (
	"log/slog"
	"net/http"
)

func (a *ApiHandler) renderOpenapiDocsPage(w http.ResponseWriter, tplName string) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	tplData := map[string]string{"SpecURL": a.Cfg.OpenapiSpecURL}
	if err := a.openapiTpl.ExecuteTemplate(w, tplName, tplData); err != nil {
		a.Log.Error("failed to execute openapi docs page", slog.String("name", tplName), slog.Any("error", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

func (a *ApiHandler) openapiRapidocPage(w http.ResponseWriter, r *http.Request) {
	a.renderOpenapiDocsPage(w, "rapidoc.tpl")
}

func (a *ApiHandler) openapiRedocPage(w http.ResponseWriter, r *http.Request) {
	a.renderOpenapiDocsPage(w, "redoc.tpl")
}

func (a *ApiHandler) openapiScalarPage(w http.ResponseWriter, r *http.Request) {
	a.renderOpenapiDocsPage(w, "scalar.tpl")
}

func (a *ApiHandler) openapiStoplightPage(w http.ResponseWriter, r *http.Request) {
	a.renderOpenapiDocsPage(w, "stoplight.tpl")
}

func (a *ApiHandler) openapiSwaggerPage(w http.ResponseWriter, r *http.Request) {
	a.renderOpenapiDocsPage(w, "swagger.tpl")
}

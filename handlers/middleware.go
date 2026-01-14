package handlers

import (
	"context"
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/tbellembois/gochimitheque/locales"
	"github.com/tbellembois/gochimitheque/logger"
	"github.com/tbellembois/gochimitheque/models"
	"github.com/tbellembois/gochimitheque/request"
)

// AppMiddleware is the application handlers wrapper handling the "func() *models.AppError" functions.
func (env *Env) AppMiddleware(h models.AppHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if e := h(w, r); e != nil {

			if e.OriginalError != nil {

				logger.Log.Error(e.Message + "-" + e.OriginalError.Error())

				if e.Code == http.StatusInternalServerError {
					logger.Log.Error(e.Message + "-" + e.OriginalError.Error())
				}

				http.Error(w, e.Message+"-"+e.OriginalError.Error(), e.Code)

			}

			logger.Log.Error(e.Error())
			http.Error(w, e.Error(), e.Code)

		}
	})
}

// LogingMiddleware is the application handlers wrapper logging the requests.
func (env *Env) LogingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		h.ServeHTTP(w, req)
	})
}

// HeadersMiddleware is the application handlers wrapper managing the common response HTTP headers.
func (env *Env) HeadersMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		h.ServeHTTP(w, req)
	})
}

// ContextMiddleware initialize the request context and setup initial variables.
func (env *Env) ContextMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// localization setup
		accept := r.Header.Get("Accept-Language")
		locales.Localizer = i18n.NewLocalizer(locales.Bundle, accept)

		ctx := context.WithValue(
			r.Context(),
			request.ChimithequeContextKey("container"),
			request.Container{
				AppURL:         env.AppURL,
				PersonLanguage: accept,
				BuildID:        env.BuildID,
			},
		)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

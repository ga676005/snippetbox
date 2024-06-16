package main

import (
	"log/slog"
	"net/http"
)

func (app *application) serveError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		url    = r.URL.RequestURI()
	)

	slog_method := slog.Any("method", method)
	slog_url := slog.Any("url", url)
	// slog_trace := slog.Any("trace", string(debug.Stack()))

	app.logger.Error(err.Error(), slog_method, slog_url)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

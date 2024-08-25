package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"os"
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

func createLogger() *slog.Logger {
	// 1 一般 logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 2 File logger
	//   另外開一個 terminal 跑 tail -f ./mylog.log 看 log
	// file, err := os.OpenFile("mylog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// logger := slog.New(slog.NewTextHandler(file, nil))
	// defer file.Close()

	// 3 JSON logger
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	// 	AddSource: true, // 會寫第幾行
	// }))

	return logger
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data TemplateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serveError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serveError(w, r, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

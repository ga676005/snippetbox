package main

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/form/v4"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
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
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) TemplateData {
	return TemplateData{
		CurrentYear: time.Now().Year(),

		// global 的訊息，如果有的話
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
	}
}

func (app *application) decodePostForm(r *http.Request, destination any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// 傳的東西不是 &abcStruct 的時候會有 form.InvalidDecoderError，代表 code 寫錯了所以用 Panic
	err = app.formDecoder.Decode(destination, r.PostForm)
	if err != nil {
		var invalidDecodeError *form.InvalidDecoderError

		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}

		return err
	}

	return nil
}

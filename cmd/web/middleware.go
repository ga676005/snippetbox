package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = slog.Any("ip", r.RemoteAddr)
			proto  = slog.Any("proto", r.Proto)
			method = slog.Any("method", r.Method)
			uri    = slog.Any("uri", r.URL.RequestURI())
		)

		app.logger.Info("received request", ip, proto, method, uri)

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 用 deferred function 因為它最後會執行到
		defer func() {
			// 用內建的 recover function
			if err := recover(); err != nil {
				// 設這個 HEADER，GO 會在送出 response 後自動關閉目前連線
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// 如果 handler 有另外開 goroutine 做一些事，那上面加在 router 的 recoverPanic middleware 救不到
// 所以要另外寫在 goroutine 裡
func (app *application) myHandler(w http.ResponseWriter, r *http.Request) {
	// some code ...

	// 像這樣
	go func() {
		defer func() {
			if err := recover(); err != nil {
				app.logger.Error(fmt.Sprint(err))
			}
		}()

		// doSomeBackgroundProcess()
	}()

	w.Write([]byte("OK"))
}

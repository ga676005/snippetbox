package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

// Chapter 3.1. 用環境變數 + cli flag
// $ export SNIPPETBOX_ADDR=":9999"
// $ go run ./cmd/web -addr=$SNIPPETBOX_ADDR

type application struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	// 	AddSource: true,
	// }))

	app := &application{
		logger: logger,
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux.Handle("GET /static", http.NotFoundHandler())
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	logger.Info("starting server", slog.Any("addr", *addr))

	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

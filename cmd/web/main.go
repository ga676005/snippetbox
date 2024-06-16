package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
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

	logger.Info("starting server", slog.Any("addr", *addr))

	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

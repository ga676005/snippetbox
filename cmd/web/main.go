package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
)

// Chapter 3.1. 用環境變數 + cli flag
// $ export SNIPPETBOX_ADDR=":9999"
// $ go run ./cmd/web -addr=$SNIPPETBOX_ADDR

// MYSQL 連線 db
// ubuntu 執行 mysql -D snippetbox -u web -p
// 輸入密碼

// redirect log
// $ go run ./cmd/web >>/tmp/web.log

type application struct {
	logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	file, err := os.OpenFile("mylog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logger := slog.New(slog.NewTextHandler(file, nil))
	defer file.Close()
	// 另外開一個 terminal 看 log
	// tail -f ./mylog.log

	// 一般 logger
	// logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// JSON logger
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	// 	AddSource: true, // 會寫第幾行
	// }))

	app := &application{
		logger: logger,
	}

	logger.Info("starting server", slog.Any("addr", *addr))

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

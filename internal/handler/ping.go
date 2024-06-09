package handler

import (
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func PingHandler(log *slog.Logger, db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("PingHandler called")
		// Your ping logic here, which can use db
		w.Write([]byte("pong"))
	}
}

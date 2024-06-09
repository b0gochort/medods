package main

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"medods/internal/handler"
	"medods/internal/repository/postgres"
	"medods/model"
	"medods/pkg/logging"

	"log/slog"

	"net/http"
	"os"
)

var config model.Config

func main() {
	log := logging.InitLog()

	configData, err := os.ReadFile("config/config.json")
	if err != nil {
		log.Error("could not read config file", "err", err)
		return
	}
	if err = json.Unmarshal(configData, &config); err != nil {
		log.Error("unmarshal config", "err", err)
		return
	}

	db, err := postgres.New(config)
	if err != nil {
		log.Error("could not connect to db", "err", err)
		return
	}

	server := &http.Server{
		Addr:    config.HTTP.Port,
		Handler: initRoutes(log, db),
	}

	log.Info("starting server on port", "port", config.HTTP.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server", "err", err)
	}

}

func initRoutes(log *slog.Logger, db *sqlx.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("received request for", "path", r.URL.Path, "method", r.Method)

		switch r.URL.Path {
		case "/ping":
			handler.Ping(w, r)
		case "/helloworld":
			handler.HelloWorld(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}
	})

	return mux
}

package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"medods/internal/handler"
	"medods/internal/repository/postgres"
	"medods/model"
	"medods/pkg/logging"

	"log/slog"

	"net/http"
)

var Config model.Config

func getConf() *model.Config {
	viper.SetConfigFile("config/config.json")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("%v", err)
	}

	conf := &model.Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}

func main() {
	log := logging.InitLog()

	Config = *getConf()

	db, err := postgres.New(Config)
	if err != nil {
		log.Error("could not connect to db", "err", err)
		return
	}

	server := &http.Server{
		Addr:    Config.HTTP.Port,
		Handler: initRoutes(log, db),
	}

	log.Info("starting server on port", "port", Config.HTTP.Port)
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
			handler.PingHandler(log, db).ServeHTTP(w, r)
		case "/helloworld":
			handler.HelloWorld(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}
	})

	return mux
}

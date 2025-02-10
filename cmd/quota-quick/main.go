package main

import (
	"api.quota-quick/api/internal/config"
	"api.quota-quick/api/internal/http-server/handlers/containers/get"
	"api.quota-quick/api/internal/http-server/handlers/containers/save"
	"api.quota-quick/api/internal/lib/logger/handlers/slogpretty"
	"api.quota-quick/api/internal/lib/logger/sl"
	"api.quota-quick/api/internal/storage/postgresql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"log/slog"
	"net/http"
	"os"
	//mwLogger "api.quota-quick/api/internal/http-server/middleware/logger"
)

const (
	envLocal = "local"
	endDev   = "dev"
	envProd  = "prod"
)

func main() {
	conf := config.MustLoad()

	log := setupLogger(conf.Env)

	log.Info("starting quota quick ^_^ ;)", slog.String("env", conf.Env))
	log.Debug("debug messages enabled")

	str, err := postgresql.GetConnStr(conf)
	if err != nil {

		log.Error("failed getting db conn str", sl.Err(err))
	}

	// TODO: Реализовать factory для разных субд
	storage, err := postgresql.New(str)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	err = storage.GetDb().Ping()
	if err != nil {
		log.Error("failed to ping db", sl.Err(err))
	}
	log.Info("Ping good!")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	//Кастомный логгер хендлеров Коли в logger http-server
	//router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/containers", save.New(log, storage))
	//router.Get("/containers", get.GetAll(log, storage))
	router.Get("/containers/{id}", get.GetById(log))
	//router.Get("/containers/user/{id}", get.GetByUserId(log, storage))

	log.Info("starting server", slog.String("address", conf.Addr))

	srv := &http.Server{
		Addr:         conf.Addr,
		Handler:      router,
		ReadTimeout:  conf.Timeout,
		WriteTimeout: conf.Timeout,
		IdleTimeout:  conf.IdleTimeout,
	}

	/* Run server */
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed starting server", sl.Err(err))
	}

	log.Info("Conf: ", conf)
	//fmt.Printf("%#v", conf)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
		//slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case endDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

func initStorage() {

}

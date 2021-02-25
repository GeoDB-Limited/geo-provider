package app

import (
	"net/http"

	"github.com/GeoDB-Limited/geo-provider/app/ctx"
	"github.com/GeoDB-Limited/geo-provider/app/handlers"
	"github.com/GeoDB-Limited/geo-provider/app/logging"
	"github.com/GeoDB-Limited/geo-provider/config"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type App interface {
	Run() error
}

type app struct {
	log    *logrus.Logger
	config config.Config
}

func New(cfg config.Config) App {
	return &app{
		config: cfg,
		log:    cfg.Logger(),
	}
}

func (a *app) Run() error {
	defer func() {
		if rvr := recover(); rvr != nil {
			a.log.Error("app panicked\n", rvr)
		}
	}()

	a.log.WithField("port", a.config.Listener()).Info("Starting server")
	if err := http.ListenAndServe(a.config.Listener(), a.router()); err != nil {
		return errors.Wrap(err, "listener failed")
	}
	return nil
}

func (a *app) router() chi.Router {
	router := chi.NewRouter()

	router.Use(
		logging.Middleware(a.log),
		ctx.Middleware(
			ctx.CtxLog(a.log),
			ctx.CtxConfig(a.config),
		),
	)

	router.Get("/geo/sources", handlers.GetSources)
	router.Get("/geo/data/{source}", handlers.GetData)

	return router
}

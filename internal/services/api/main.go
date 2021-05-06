package api

import (
	"github.com/geo-provider/internal/config"
	"github.com/geo-provider/internal/data"
	"github.com/geo-provider/internal/data/postgres"
	"github.com/geo-provider/internal/services/api/ctx"
	"github.com/geo-provider/internal/services/api/handlers"
	"github.com/geo-provider/internal/services/api/logging"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type API interface {
	Run() error
}

type api struct {
	log    *logrus.Logger
	config config.Config
	db     data.Storage
}

func New(cfg config.Config) API {
	return &api{
		config: cfg,
		log:    cfg.Logger(),
		db:     postgres.New(cfg.DB()),
	}
}

func (a *api) Run() error {
	defer func() {
		if rvr := recover(); rvr != nil {
			a.log.Error("api panicked\n", rvr)
		}
	}()
	a.log.WithField("port", a.config.Listener()).Info("Starting server")
	err := http.ListenAndServe(a.config.Listener(), a.router())
	return errors.Wrap(err, "listener failed")
}

func (a *api) router() chi.Router {
	router := chi.NewRouter()

	router.Use(
		logging.Middleware(a.log),
		ctx.Middleware(
			ctx.CtxLog(a.log),
			ctx.CtxConfig(a.config),
			ctx.CtxLocations(a.db.Locations()),
			ctx.CtxDevices(a.db.Devices()),
		),
	)

	router.Get("/geo/sources", handlers.GetSources)
	router.Get("/geo/data/{owner}/{source}", handlers.GetData)

	return router
}

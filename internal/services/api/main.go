package api

import (
	config2 "github.com/geo-provider/internal/config"
	"github.com/geo-provider/internal/data"
	"github.com/geo-provider/internal/data/postgres"
	ctx2 "github.com/geo-provider/internal/services/api/ctx"
	handlers2 "github.com/geo-provider/internal/services/api/handlers"
	logging2 "github.com/geo-provider/internal/services/api/logging"
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
	config config2.Config
	db     data.Storage
}

func New(cfg config2.Config) API {
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
		logging2.Middleware(a.log),
		ctx2.Middleware(
			ctx2.CtxLog(a.log),
			ctx2.CtxConfig(a.config),
			ctx2.CtxLocations(a.db.Locations()),
			ctx2.CtxDevices(a.db.Devices()),
		),
	)

	router.Get("/geo/sources", handlers2.GetSources)
	router.Get("/geo/data/{owner}/{source}", handlers2.GetData)

	return router
}

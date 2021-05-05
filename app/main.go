package app

import (
	"github.com/geo-provider/app/ctx"
	"github.com/geo-provider/app/data"
	"github.com/geo-provider/app/data/postgres"
	"github.com/geo-provider/app/handlers"
	"github.com/geo-provider/app/logging"
	"github.com/geo-provider/app/storage"
	"github.com/geo-provider/config"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type App interface {
	Run() error
}

type app struct {
	log    *logrus.Logger
	config config.Config
	db     data.Storage
	csv    storage.Storage
}

func New(cfg config.Config) App {
	return &app{
		config: cfg,
		log:    cfg.Logger(),
		db:     postgres.New(cfg.DB()),
		csv:    storage.New(cfg),
	}
}

func (a *app) Run() error {
	defer func() {
		if rvr := recover(); rvr != nil {
			a.log.Error("app panicked\n", rvr)
		}
	}()
	a.log.WithField("port", a.config.Listener()).Info("Starting server")
	a.migrateDataFromCSV()
	err := http.ListenAndServe(a.config.Listener(), a.router())
	return errors.Wrap(err, "listener failed")
}

func (a *app) router() chi.Router {
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

func (a *app) migrateDataFromCSV() {
	go a.migrateDevicesFromCSV()
	go a.migrateLocationsFromCSV()
}

func (a *app) migrateLocationsFromCSV() {
	locations, err := a.csv.SelectLocationsFromCSV()
	if err != nil {
		panic(errors.Wrap(err, "failed to parse locations"))
	}
	for _, location := range locations {
		if err := a.db.Locations().Insert(location); err != nil {
			panic(errors.Wrap(err, "failed to insert location data"))
		}
	}
	a.log.Info("Finished migrating locations data")
}

func (a *app) migrateDevicesFromCSV() {
	devices, err := a.csv.SelectDevicesFromCSV()
	if err != nil {
		panic(errors.Wrap(err, "failed to parse devices"))
	}
	for _, device := range devices {
		if err := a.db.Devices().Insert(device); err != nil {
			panic(errors.Wrap(err, "failed to insert device data"))
		}
	}
	a.log.Info("Finished migrating devices data")
}

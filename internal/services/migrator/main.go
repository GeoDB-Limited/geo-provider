package migrator

import (
	config2 "github.com/geo-provider/internal/config"
	"github.com/geo-provider/internal/data"
	"github.com/geo-provider/internal/data/postgres"
	"github.com/geo-provider/internal/storage"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Run()
}

type service struct {
	log    *logrus.Logger
	config config2.Config
	db     data.Storage
	csv    storage.Storage
}

func New(cfg config2.Config) Service {
	return &service{
		config: cfg,
		log:    cfg.Logger(),
		db:     postgres.New(cfg.DB()),
		csv:    storage.New(cfg),
	}
}

func (s *service) Run() {
	defer func() {
		if rvr := recover(); rvr != nil {
			s.log.Error("service panicked\n", rvr)
		}
	}()
	s.log.Info("Starting migrator service...")
	go s.migrateDevicesFromCSV()
	s.migrateLocationsFromCSV()
}

func (s *service) migrateLocationsFromCSV() {
	locations, err := s.csv.SelectLocationsFromCSV()
	if err != nil {
		panic(errors.Wrap(err, "failed to parse locations"))
	}
	for _, location := range locations {
		if err := s.db.Locations().Insert(location); err != nil {
			panic(errors.Wrap(err, "failed to insert location data"))
		}
	}
	s.log.Info("Finished migrating locations data")
}

func (s *service) migrateDevicesFromCSV() {
	devices, err := s.csv.SelectDevicesFromCSV()
	if err != nil {
		panic(errors.Wrap(err, "failed to parse devices"))
	}
	for _, device := range devices {
		if err := s.db.Devices().Insert(device); err != nil {
			panic(errors.Wrap(err, "failed to insert device data"))
		}
	}
	s.log.Info("Finished migrating devices data")
}

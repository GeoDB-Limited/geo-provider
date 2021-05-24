package migrator

import (
	"github.com/geo-provider/internal/config"
	"github.com/geo-provider/internal/data"
	"github.com/geo-provider/internal/data/postgres"
	"github.com/geo-provider/internal/storage"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	PaginationRowsCount = 9000
)

type Service interface {
	Run()
}

type service struct {
	log    *logrus.Logger
	config config.Config
	db     data.Storage
	csv    storage.Storage
}

func New(cfg config.Config) Service {
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
	s.migrateDataFromCSV()
}

func (s *service) migrateDataFromCSV() {
	var wg sync.WaitGroup
	wg.Add(1)
	go s.migrateLocationsFromCSV(&wg)
	s.migrateDevicesFromCSV()
	wg.Wait()
}

func (s *service) migrateLocationsFromCSV(wg *sync.WaitGroup) {
	defer wg.Done()
	s.log.Info("Started migrating locations data")
	locations, err := s.csv.SelectLocationsFromCSV()
	if err != nil {
		panic(errors.Wrap(err, "failed to parse locations"))
	}
	locationsCount := len(locations)
	for i := 0; i < locationsCount; i += PaginationRowsCount {
		var toInsert []data.Location
		if i+PaginationRowsCount > locationsCount {
			toInsert = locations[i:locationsCount]
		} else {
			toInsert = locations[i : i+PaginationRowsCount]
		}
		if err := s.db.Locations().Insert(toInsert...); err != nil {
			panic(errors.Wrap(err, "failed to insert locations data"))
		}
	}
	s.log.Info("Finished migrating locations data")
}

func (s *service) migrateDevicesFromCSV() {
	s.log.Info("Started migrating devices data")
	devices, err := s.csv.SelectDevicesFromCSV()
	if err != nil {
		panic(errors.Wrap(err, "failed to parse devices"))
	}
	devicesCount := len(devices)
	for i := 0; i < devicesCount; i += PaginationRowsCount {
		var toInsert []data.Device
		if i+PaginationRowsCount > devicesCount {
			toInsert = devices[i:devicesCount]
		} else {
			toInsert = devices[i : i+PaginationRowsCount]
		}
		if err := s.db.Devices().Insert(toInsert...); err != nil {
			panic(errors.Wrap(err, "failed to insert devices data"))
		}
	}
	s.log.Info("Finished migrating devices data")
}

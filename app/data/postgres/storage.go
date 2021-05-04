package postgres

import (
	"database/sql"
	"github.com/geo-provider/app/data"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Devices() data.DevicesStorage {
	return &DevicesStorage{db: s.db}
}

func (s *Storage) Locations() data.LocationsStorage {
	return &LocationsStorage{db: s.db}
}

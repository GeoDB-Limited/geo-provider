package data

import (
	"time"
)

type LocationsStorage interface {
	New() LocationsStorage
	Select() ([]Location, error)
	Insert(Location) error
}

type Location struct {
	Address   string    `db:"address"`
	Latitude  float64   `db:"latitude"`
	Longitude float64   `db:"longitude"`
	Altitude  float64   `db:"altitude"`
	Time      time.Time `db:"time"`
	Timestamp time.Time `db:"timestamp"`
	Date      time.Time `db:"date"`
}

package data

import (
	"time"
)

type LocationsStorage interface {
	New() LocationsStorage
	Select(limit, offset uint64) ([]Location, error)
	Insert(location ...Location) error
}

type Location struct {
	ID        int64     `db:"id" json:"id"`
	Address   string    `db:"address" json:"address"`
	Latitude  float64   `db:"latitude" json:"latitude"`
	Longitude float64   `db:"longitude" json:"longitude"`
	Altitude  float64   `db:"altitude" json:"altitude"`
	Time      time.Time `db:"time" json:"time"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
	Date      time.Time `db:"date" json:"date"`
}

func (l Location) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"address":   l.Address,
		"latitude":  l.Latitude,
		"longitude": l.Longitude,
		"altitude":  l.Altitude,
		"time":      l.Time,
		"timestamp": l.Timestamp,
		"date":      l.Date,
	}
	return result
}

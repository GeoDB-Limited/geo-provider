package data

import (
	"github.com/google/uuid"
	"time"
)

type DevicesStorage interface {
	New() DevicesStorage
	Select(limit, offset uint64) ([]Device, error)
	Insert(device Device) error
}

type Device struct {
	ID             int64     `db:"id" json:"id"`
	Address        string    `db:"address" json:"address"`
	UUID           uuid.UUID `db:"uuid" json:"uuid"`
	OS             string    `db:"os" json:"os"`
	Model          string    `db:"model" json:"model"`
	Locale         string    `db:"locale" json:"locale"`
	Apps           string    `db:"apps" json:"apps"`
	Version        string    `db:"version" json:"version"`
	Time           time.Time `db:"time" json:"time"`
	Timestamp      time.Time `db:"timestamp" json:"timestamp"`
	Date           time.Time `db:"date" json:"date"`
	GeocashVersion string    `db:"geocash_version" json:"geocash_version"`
}

func (d Device) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"address":         d.Address,
		"uuid":            d.UUID,
		"os":              d.OS,
		"model":           d.Model,
		"locale":          d.Locale,
		"apps":            d.Apps,
		"version":         d.Version,
		"time":            d.Time,
		"timestamp":       d.Timestamp,
		"date":            d.Date,
		"geocash_version": d.GeocashVersion,
	}
	return result
}

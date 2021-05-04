package data

import (
	"github.com/google/uuid"
	"time"
)

type DevicesStorage interface {
	New() DevicesStorage
	Select() ([]Device, error)
	Insert(Device) error
}

type Device struct {
	Address        string    `db:"address"`
	UUID           uuid.UUID `db:"uuid"`
	OS             string    `db:"os"`
	Model          string    `db:"model"`
	Locale         string    `db:"locale"`
	Apps           string    `db:"apps"`
	Version        string    `db:"version"`
	Time           time.Time `db:"time"`
	Timestamp      time.Time `db:"timestamp"`
	Date           time.Time `db:"date"`
	GeocashVersion string    `db:"geocash_version"`
}

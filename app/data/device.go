package data

import "time"

type Device struct {
	Address string `db:"address"`

	UUID    string `db:"uuid"`
	OS      string `db:"os"`
	Model   string `db:"model"`
	Locale  string `db:"locale"`
	Apps    string `db:"apps"`
	Version string `db:"version"`

	Time      time.Time `db:"time"`
	Timestamp time.Time `db:"timestamp"`
	Date      time.Time `db:"date"`

	GeocashVersion string `db:"geocash_version"`
}

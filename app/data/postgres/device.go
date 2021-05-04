package postgres

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/geo-provider/app/data"
	"github.com/pkg/errors"
)

const (
	devicesTable = "devices"
	all          = "*"
)

type DevicesStorage struct {
	db *sql.DB
}

func (s *DevicesStorage) New() data.DevicesStorage {
	return NewDevicesStorage(s.db)
}

func NewDevicesStorage(db *sql.DB) data.DevicesStorage {
	return &DevicesStorage{
		db: db,
	}
}

func (s *DevicesStorage) Insert(device data.Device) error {
	query := squirrel.Insert(devicesTable).PlaceholderFormat(squirrel.Dollar).Columns(
		"address",
		"uuid",
		"os",
		"model",
		"locale",
		"apps",
		"version",
		"time",
		"timestamp",
		"date",
		"geocash_version",
	).Values(
		device.Address,
		device.UUID,
		device.OS,
		device.Model,
		device.Locale,
		device.Apps,
		device.Version,
		device.Time,
		device.Timestamp,
		device.Date,
		device.GeocashVersion,
	).RunWith(s.db)

	_, err := query.Exec()
	return errors.Wrap(err, "failed to insert device")
}

func (s *DevicesStorage) Select() ([]data.Device, error) {
	query := squirrel.Select(all).From(devicesTable).PlaceholderFormat(squirrel.Dollar).RunWith(s.db)
	rows, err := query.Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var models []data.Device

	for rows.Next() {
		model := data.Device{}
		err := rows.Scan(
			&model.Address,
			&model.UUID,
			&model.OS,
			&model.Model,
			&model.Locale,
			&model.Apps,
			&model.Version,
			&model.Time,
			&model.Timestamp,
			&model.Date,
			&model.GeocashVersion,
		)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, nil
}

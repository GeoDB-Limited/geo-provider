package postgres

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/geo-provider/app/data"
	"github.com/geo-provider/config"
)

const (
	devicesTable = "devices"
	all          = "*"
)

type devicesStorage struct {
	db  *sql.DB
	sql sq.SelectBuilder
}

type DevicesStorage interface {
	Select() ([]data.Device, error)
}

var devicesSelect = sq.Select(all).From(devicesTable).PlaceholderFormat(sq.Dollar)

func NewDevicesStorage(cfg config.Config) DevicesStorage {
	return &devicesStorage{
		db:  cfg.Databaser(),
		sql: devicesSelect.RunWith(cfg.Databaser()),
	}
}

func (s *devicesStorage) Select() ([]data.Device, error) {
	rows, err := s.sql.Query()
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

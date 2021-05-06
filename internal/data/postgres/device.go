package postgres

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/geo-provider/internal/data"
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

func (s *DevicesStorage) Insert(value data.Device) error {
	query := squirrel.Insert(devicesTable).PlaceholderFormat(squirrel.Dollar).SetMap(value.ToMap())
	_, err := query.RunWith(s.db).Exec()
	return errors.Wrap(err, "failed to insert device")
}

func (s *DevicesStorage) Select(limit, offset uint64) ([]data.Device, error) {
	query := squirrel.Select(all).From(devicesTable).PlaceholderFormat(squirrel.Dollar).Limit(limit).Offset(offset)
	rows, err := query.RunWith(s.db).Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var models []data.Device

	for rows.Next() {
		model := data.Device{}
		err := rows.Scan(
			&model.ID,
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

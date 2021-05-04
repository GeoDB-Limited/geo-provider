package postgres

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/geo-provider/app/data"
	"github.com/pkg/errors"
)

const locationsTable = "Locations"

type LocationsStorage struct {
	db *sql.DB
}

func (s *LocationsStorage) New() data.LocationsStorage {
	return NewLocationsStorage(s.db)
}

func NewLocationsStorage(db *sql.DB) data.LocationsStorage {
	return &LocationsStorage{
		db: db,
	}
}

func (s *LocationsStorage) Insert(location data.Location) error {
	query := squirrel.Insert(locationsTable).PlaceholderFormat(squirrel.Dollar).Columns(
		"address",
		"latitude",
		"longitude",
		"altitude",
		"time",
		"timestamp",
		"date",
	).Values(
		location.Address,
		location.Latitude,
		location.Longitude,
		location.Altitude,
		location.Time,
		location.Timestamp,
		location.Date,
	).RunWith(s.db)

	_, err := query.Exec()
	return errors.Wrap(err, "failed to insert location")
}

func (s *LocationsStorage) Select() ([]data.Location, error) {
	query := squirrel.Select(all).From(locationsTable).PlaceholderFormat(squirrel.Dollar).RunWith(s.db)
	rows, err := query.Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var models []data.Location

	for rows.Next() {
		model := data.Location{}
		err := rows.Scan(
			&model.Address,
			&model.Latitude,
			&model.Longitude,
			&model.Altitude,
			&model.Time,
			&model.Timestamp,
			&model.Date,
		)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, nil
}

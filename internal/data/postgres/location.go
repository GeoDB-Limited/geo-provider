package postgres

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/geo-provider/internal/data"
	"github.com/pkg/errors"
)

const locationsTable = "locations"

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

func (s *LocationsStorage) Insert(locations ...data.Location) error {
	if len(locations) == 0 {
		return nil
	}
	query := squirrel.Insert(locationsTable).PlaceholderFormat(squirrel.Dollar).Columns(
		"address",
		"latitude",
		"longitude",
		"altitude",
		"time",
		"timestamp",
		"date",
	)
	for _, location := range locations {
		query = query.Values(
			location.Address,
			location.Latitude,
			location.Longitude,
			location.Altitude,
			location.Time,
			location.Timestamp,
			location.Date,
		)
	}
	_, err := query.RunWith(s.db).Exec()
	return errors.Wrap(err, "failed to insert location")
}

func (s *LocationsStorage) Select(limit, offset uint64) ([]data.Location, error) {
	query := squirrel.Select(all).From(locationsTable).PlaceholderFormat(squirrel.Dollar).Limit(limit).Offset(offset)
	rows, err := query.RunWith(s.db).Query()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var models []data.Location

	for rows.Next() {
		model := data.Location{}
		err := rows.Scan(
			&model.ID,
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

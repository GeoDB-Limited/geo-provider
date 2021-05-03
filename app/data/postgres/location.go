package postgres

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/geo-provider/app/data"
	"github.com/geo-provider/config"
)

const locationsTable = "locations"

type locationsStorage struct {
	db  *sql.DB
	sql sq.SelectBuilder
}

type LocationsStorage interface {
	Select() ([]data.Location, error)
}

var locationsSelect = sq.Select(all).From(locationsTable).PlaceholderFormat(sq.Dollar)

func NewLocationsStorage(cfg config.Config) LocationsStorage {
	return &locationsStorage{
		db:  cfg.Databaser(),
		sql: locationsSelect.RunWith(cfg.Databaser()),
	}
}

func (s *locationsStorage) Select() ([]data.Location, error) {
	rows, err := s.sql.Query()
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

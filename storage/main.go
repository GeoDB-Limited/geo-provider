package storage

import (
	"encoding/csv"
	"encoding/json"
	"github.com/geo-provider/utils"
	"github.com/pkg/errors"
	"os"
)

const (
	MaxReadRows = 100
)

type Keeper interface {
	// Read reads storage file [0 + offset:0 + offset + count)
	Read(offset, count int) ([]json.RawMessage, error)
}

type keeper struct {
	Path string
}

func Open(path string) (Keeper, error) {
	if !utils.FileExists(path) {
		return nil, errors.New("file for the specified source not found")
	}
	return &keeper{
		Path: path,
	}, nil
}

func (k *keeper) Read(offset, count int) ([]json.RawMessage, error) {
	file, err := os.Open(k.Path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	rows, err := utils.CSVToMap(reader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read csv to map")
	}

	start := utils.Min(len(rows)-1, offset)
	end := utils.Max(0, utils.Min(len(rows)-offset-1, count))
	rows = rows[start : start+end]

	res := make([]json.RawMessage, len(rows))
	for i, row := range rows {
		res[i], err = json.Marshal(row)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal row to json")
		}
	}

	return res, nil
}

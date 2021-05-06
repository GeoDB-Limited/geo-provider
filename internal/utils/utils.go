package utils

import (
	"encoding/csv"
	"github.com/pkg/errors"
	"os"
)

func Keys(src map[string]string) []string {
	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}

	return keys
}

func ReadCsv(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return [][]string{}, errors.Wrapf(err, "failed to open file: %s", path)
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, errors.Wrapf(err, "failed to read file: %s", path)
	}

	return lines, nil
}

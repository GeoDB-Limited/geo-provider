package storage

import (
	"encoding/csv"
	"github.com/geo-provider/app/data"
	"github.com/geo-provider/config"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"time"
)

const (
	timeLayout = "2006-01-02 15:04:05 MST"
	dateLayout = "2006-01-02"

	locationsStorageKey = "locations"
	devicesStorageKey   = "devices"
)

type Storage interface {
	SelectLocationsFromCSV() ([]data.Location, error)
	SelectDevicesFromCSV() ([]data.Device, error)
}

type storage struct {
	config config.Config
}

func New(cfg config.Config) Storage {
	return &storage{
		config: cfg,
	}
}

func (s *storage) SelectLocationsFromCSV() ([]data.Location, error) {
	filePath := s.config.Source(locationsStorageKey)
	lines, err := s.ReadCsv(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read locations from %s", filePath)
	}
	locations, err := s.ParseLocations(lines[1:])
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse locations")
	}
	return locations, nil
}

func (s *storage) ParseLocations(lines [][]string) ([]data.Location, error) {
	locations := make([]data.Location, 0, len(lines))
	for _, line := range lines {
		latitude, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse latitude: %s", line[1])
		}
		longitude, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse longitude: %s", line[2])
		}
		altitude, err := strconv.ParseFloat(line[3], 64)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse altitude: %s", line[3])
		}
		locationTime, err := time.Parse(timeLayout, line[4])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse time: %s", line[4])
		}
		timestamp, err := time.Parse(timeLayout, line[5])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse timestamp: %s", line[5])
		}
		date, err := time.Parse(dateLayout, line[6])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse date: %s", line[6])
		}
		location := data.Location{
			Address:   line[0],
			Latitude:  latitude,
			Longitude: longitude,
			Altitude:  altitude,
			Time:      locationTime,
			Timestamp: timestamp,
			Date:      date,
		}
		locations = append(locations, location)
	}

	return locations, nil
}

func (s *storage) SelectDevicesFromCSV() ([]data.Device, error) {
	filePath := s.config.Source(devicesStorageKey)
	lines, err := s.ReadCsv(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read devices from %s", filePath)
	}
	devices, err := s.ParseDevices(lines[1:])
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse devices")
	}
	return devices, nil
}

func (s *storage) ParseDevices(lines [][]string) ([]data.Device, error) {
	devices := make([]data.Device, 0, len(lines))
	for _, line := range lines {
		deviceUUID, err := uuid.Parse(line[5])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse uuid: %s", line[5])
		}
		deviceTime, err := time.Parse(timeLayout, line[7])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse time: %s", line[7])
		}
		timestamp, err := time.Parse(timeLayout, line[8])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse timestamp: %s", line[8])
		}
		date, err := time.Parse(dateLayout, line[9])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse date: %s", line[9])
		}
		device := data.Device{
			Address:        line[0],
			OS:             line[1],
			Model:          line[2],
			Locale:         line[3],
			Version:        line[4],
			UUID:           deviceUUID,
			Apps:           line[6],
			Time:           deviceTime,
			Timestamp:      timestamp,
			Date:           date,
			GeocashVersion: line[10],
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func (s *storage) ReadCsv(path string) ([][]string, error) {
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

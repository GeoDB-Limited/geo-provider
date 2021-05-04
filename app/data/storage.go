package data

type Storage interface {
	Devices() DevicesStorage
	Locations() LocationsStorage
}

package ctx

import (
	"context"
	"github.com/geo-provider/internal/config"
	"github.com/geo-provider/internal/data"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	ctxLog = iota
	ctxConfig
	locationCtxKey
	deviceCtxKey
)

func CtxConfig(cfg config.Config) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ctxConfig, cfg)
	}
}

func Config(r *http.Request) config.Config {
	return r.Context().Value(ctxConfig).(config.Config)
}

func CtxLog(log *logrus.Logger) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ctxLog, log)
	}
}

func Log(r *http.Request) *logrus.Logger {
	return r.Context().Value(ctxLog).(*logrus.Logger)
}

func CtxLocations(entry data.LocationsStorage) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, locationCtxKey, entry)
	}
}

func Locations(r *http.Request) data.LocationsStorage {
	return r.Context().Value(locationCtxKey).(data.LocationsStorage).New()
}

func CtxDevices(entry data.DevicesStorage) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, deviceCtxKey, entry)
	}
}

func Devices(r *http.Request) data.DevicesStorage {
	return r.Context().Value(deviceCtxKey).(data.DevicesStorage).New()
}

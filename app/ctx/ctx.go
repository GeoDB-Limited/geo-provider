package ctx

import (
	"context"
	"net/http"

	"github.com/GeoDB-Limited/geo-provider/config"
	"github.com/sirupsen/logrus"
)

const (
	ctxLog    = "ctxLog"
	ctxConfig = "ctxConfig"
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

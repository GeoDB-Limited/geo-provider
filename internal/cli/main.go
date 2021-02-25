package cli

import (
	"github.com/GeoDB-Limited/geo-provider/internal/api"
	"github.com/GeoDB-Limited/geo-provider/internal/config"
	"github.com/urfave/cli/v2"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func Run(args []string) bool {
	var cfg config.Config
	log := logan.New()

	defer func() {
		if rvr := recover(); rvr != nil {
			log.WithRecover(rvr).Error("app panicked")
		}
	}()

	app := cli.NewApp()

	before := func(_ *cli.Context) error {
		getter, err := kv.FromEnv()
		if err != nil {
			return errors.Wrap(err, "failed to get config")
		}
		cfg = config.New(getter)
		log = cfg.Log()
		return nil
	}

	app.Commands = cli.Commands{
		{
			Name:   "run",
			Before: before,
			Action: func(_ *cli.Context) error {
				return api.Run(cfg)
			},
		},
		{
			Subcommands: cli.Commands{
				{
					Name: "up",
					Action: func(ctx *cli.Context) error {
						return MigrateUp(cfg)
					},
				},
				{
					Name: "down",
					Action: func(ctx *cli.Context) error {
						return MigrateDown(cfg)
					},
				},
			},
			Name:   "migrate",
			Before: before,
		},
	}

	if err := app.Run(args); err != nil {
		log.WithError(err).Error("app finished")
		return false
	}
	return true
}

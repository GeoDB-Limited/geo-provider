package cli

import (
	"github.com/alecthomas/kingpin"
	"github.com/geo-provider/internal/config"
	"github.com/geo-provider/internal/services/api"
	"github.com/geo-provider/internal/services/migrator"
	"os"
)

func Run(args []string) bool {
	cfg := config.New(os.Getenv("CONFIG"))
	log := cfg.Logger()

	defer func() {
		if rvr := recover(); rvr != nil {
			log.Error("internal panicked", rvr)
		}
	}()

	app := kingpin.New("geo-provider", "")
	runCmd := app.Command("run", "run command")
	apiCmd := runCmd.Command("api", "run api")
	migratorCmd := runCmd.Command("migrator", "run service to migrate data from csv")

	migrateDBCmd := app.Command("migrate", "migrate command")
	migrateDBUpCmd := migrateDBCmd.Command("up", "migrate db up")
	migrateDBDownCmd := migrateDBCmd.Command("down", "migrate db down")

	cmd, err := app.Parse(args[1:])
	if err != nil {
		log.WithError(err).Error("failed to parse arguments")
		return false
	}

	switch cmd {
	case apiCmd.FullCommand():
		if err := api.New(cfg).Run(); err != nil {
			log.Error("failed to start api", err)
			return false
		}
	case migratorCmd.FullCommand():
		migrator.New(cfg).Run()
	case migrateDBUpCmd.FullCommand():
		err = MigrateUp(cfg)
	case migrateDBDownCmd.FullCommand():
		err = MigrateDown(cfg)
	default:
		log.Errorf("unknown command %s", cmd)
		return false
	}

	if err != nil {
		log.WithError(err).Error("failed to exec cmd")
		return false
	}

	return true
}

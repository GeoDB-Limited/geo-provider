package config

import (
	"database/sql"
	"fmt"
	"github.com/geo-provider/app/data/migrate"
	"github.com/geo-provider/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config interface {
	IsOwner(string) bool
	Source(string) string
	ListSources() []string
	Listener() string
	Logger() *logrus.Logger
	Databaser() *sql.DB
}

type config struct {
	Sources  map[string]string `yaml:"sources"`
	Owners   []string          `yaml:"owners"`
	Addr     string            `yaml:"addr"`
	Log      string            `yaml:"log"`
	Database struct {
		URL     string `yaml:"url"`
		Migrate string `yaml:"migrate"`
	} `yaml:"db"`
}

func New(path string) Config {
	cfg := config{}

	yamlConfig, err := ioutil.ReadFile(path)
	if err != nil {
		panic(errors.New(fmt.Sprintf("failed to read config: %s", path)))
	}

	err = yaml.Unmarshal(yamlConfig, &cfg)
	if err != nil {
		panic(errors.New(fmt.Sprintf("failed to unmarshal config: %s", path)))
	}

	return &cfg
}

func (c *config) Source(name string) string {
	if _, ok := c.Sources[name]; !ok {
		return ""
	}
	return c.Sources[name]
}

func (c *config) ListSources() []string {
	return utils.Keys(c.Sources)
}

func (c *config) Listener() string {
	return c.Addr
}

func (c *config) Logger() *logrus.Logger {
	level, err := logrus.ParseLevel(c.Log)
	if err != nil {
		panic(errors.Wrapf(err, "failed to parse logging level %s", c.Log))
	}

	logger := logrus.New()
	logger.SetLevel(level)
	return logger
}

func (c *config) IsOwner(owner string) bool {
	for _, o := range c.Owners {
		if o == owner {
			return true
		}
	}

	return false
}

func (c *config) Databaser() *sql.DB {
	db, err := sql.Open("postgres", c.Database.URL)
	if err != nil {
		panic(err)
	}

	switch c.Database.Migrate {
	case migrate.Up:
		applied, err := migrate.MigrateUp(db)
		if err != nil {
			panic(err)
		}
		c.Logger().WithField("applied", applied).Info("Migrations up applied")
	case migrate.Down:
		applied, err := migrate.MigrateDown(db)
		if err != nil {
			panic(err)
		}
		c.Logger().WithField("applied", applied).Info("Migrations down applied")
	default:
		panic("Unknown migration method")
	}

	if err := db.Ping(); err != nil {
		panic(errors.Wrap(err, "database unavailable"))
	}

	return db
}

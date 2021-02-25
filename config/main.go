package config

import (
	"fmt"
	"io/ioutil"

	"github.com/GeoDB-Limited/geo-provider/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config interface {
	Source(string) string
	ListSources() []string
	Listener() string
	Logger() *logrus.Logger
}

type config struct {
	Sources map[string]string `yaml:"sources"`
	Addr    string            `yaml:"addr"`
	Log     string            `yaml:"log"`
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

package api

import (
	"net"
	"net/http"

	"github.com/GeoDB-Limited/geo-provider/internal/config"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type service struct {
	log      *logan.Entry
	listener net.Listener
}

func (s *service) run() error {
	// TODO implement custom logic here and rm panic below
	r := s.router()

	if err := http.Serve(s.listener, r); err != nil {
		return errors.Wrap(err, "router crashed")
	}
	panic("implement me!") // this panic will not be raised if service is generated with http-handling features but consider removing it
}

func newService(cfg config.Config) *service {
	return &service{
		log:      cfg.Log(),
		listener: cfg.Listener(),
	}
}

func Run(cfg config.Config) error {
	if err := newService(cfg).run(); err != nil {
		return err
	}

	return nil
}

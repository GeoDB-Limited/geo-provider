package api

import (
	"github.com/GeoDB-Limited/geo-provider/internal/api/handlers"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log), // this line may cause compilation error but in general case `dep ensure -v` will fix it
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
		),
	)

	// configure endpoints here

	return r
}

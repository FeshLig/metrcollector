package handler

import (
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handlers struct {
	Root       *RootHandler
	UpdateJSON *UpdateJSONHandler
	ValueJSON  *ValueJSONHandler
	UpdateURL  *UpdateURLHandler
	ValueURL   *ValueURLHandler
	Ping       *PingHandler
}

func NewHandlers(service service.MetricsService, db *pgxpool.Pool) *Handlers {

	checker := repository.NewDBChecker(db)

	return &Handlers{
		Root:       NewRootHandler(service),
		UpdateJSON: NewUpdateJSONHandler(service),
		ValueJSON:  NewValueJSONHandler(service),
		UpdateURL:  NewUpdateURLHandler(service),
		ValueURL:   NewValueURLHandler(service),
		Ping:       NewPingHandler(checker),
	}

}

package handler

import (
	"github.com/FeshLig/metrcollector/internal/service"
)

type Handlers struct {
	Root       *RootHandler
	UpdateJSON *UpdateJSONHandler
	ValueJSON  *ValueJSONHandler
	UpdateURL  *UpdateURLHandler
	ValueURL   *ValueURLHandler
	Ping       *PingHandler
	Updates    *UpdatesHandler
}

func NewHandlers(service service.MetricsService) *Handlers {

	return &Handlers{
		Root:       NewRootHandler(service),
		UpdateJSON: NewUpdateJSONHandler(service),
		ValueJSON:  NewValueJSONHandler(service),
		UpdateURL:  NewUpdateURLHandler(service),
		ValueURL:   NewValueURLHandler(service),
		Ping:       NewPingHandler(service),
		Updates:    NewUpdatesHandler(service),
	}

}

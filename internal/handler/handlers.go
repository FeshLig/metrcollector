package handler

import (
	"github.com/FeshLig/metrcollector/internal/audit"
	"github.com/FeshLig/metrcollector/internal/service"
)

// Handlers contains all HTTP handlers of the application.
type Handlers struct {
	Root       *RootHandler
	UpdateJSON *UpdateJSONHandler
	ValueJSON  *ValueJSONHandler
	UpdateURL  *UpdateURLHandler
	ValueURL   *ValueURLHandler
	Ping       *PingHandler
	Updates    *UpdatesHandler
}

// NewHandlers creates and initializes all application handlers.
func NewHandlers(service service.MetricsService, audit *audit.Publisher) *Handlers {

	return &Handlers{
		Root:       NewRootHandler(service),
		UpdateJSON: NewUpdateJSONHandler(service, audit),
		ValueJSON:  NewValueJSONHandler(service),
		UpdateURL:  NewUpdateURLHandler(service, audit),
		ValueURL:   NewValueURLHandler(service),
		Ping:       NewPingHandler(service),
		Updates:    NewUpdatesHandler(service, audit),
	}

}

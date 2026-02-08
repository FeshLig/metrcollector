package main

import (
	"github.com/FeshLig/metrcollector/internal/config"
	"github.com/FeshLig/metrcollector/internal/persister"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/FeshLig/metrcollector/internal/router"
	"github.com/FeshLig/metrcollector/internal/service"
)

func main() {
	Run()
}

func Run() {

	flags := config.GetOptions()

	memStorage := repository.NewMemStorage()

	p := persister.NewFilePersister(
		flags.FileStoragePath.String(),
		memStorage,
		flags.StoreInterval.Duration,
	)

	if flags.Restore {
		_ = p.Load()
	}

	p.Start()
	defer p.Stop()

	service := service.NewMetricService(
		memStorage,
		p,
		flags.StoreInterval.Duration == 0,
	)

	router := router.NewRouter(service)

	router.Run(flags.Address.String())

}

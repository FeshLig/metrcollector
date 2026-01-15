package main

import (
	"github.com/FeshLig/metrcollector/internal/config"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/FeshLig/metrcollector/internal/router"
)

func main() {
	Run()
}

func Run() {

	flags := config.GetOptions()

	memStorage := repository.NewMemStorage()

	router := router.NewRouter(memStorage)

	router.Run(flags.Address.String())

}

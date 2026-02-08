package main

import (
	"github.com/FeshLig/metrcollector/internal/agent"
	"github.com/FeshLig/metrcollector/internal/repository"
)

func main() {

	options := agent.GetOptions()

	storage := repository.NewMemStorage()

	agent.RunSender(storage, options)

}

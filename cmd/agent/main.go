package main

import (
	"log"

	"github.com/FeshLig/metrcollector/internal/agent"
	"github.com/FeshLig/metrcollector/internal/repository"
)

func main() {

	options := agent.GetOptions()

	storage := repository.NewMemStorage()

	if err := agent.RunSender(storage, options); err != nil {
		log.Printf("agent error: %v\n", err)
	}

}

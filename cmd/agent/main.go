package main

import (
	"context"
	"log"

	"github.com/FeshLig/metrcollector/internal/agent"
	"github.com/FeshLig/metrcollector/internal/repository"
)

func main() {

	options := agent.GetOptions()

	storage := repository.NewMemStorage()

	ctx := context.Background()

	if err := agent.RunSender(ctx, storage, options); err != nil {
		log.Printf("agent error: %v\n", err)
	}

}

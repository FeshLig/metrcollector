package main

import (
	"context"

	"github.com/FeshLig/metrcollector/internal/agent"
	"github.com/FeshLig/metrcollector/internal/repository"
)

func main() {

	options := agent.GetOptions()

	storage := repository.NewMemStorage()

	ctx := context.Background()

	agent.RunSender(ctx, storage, options)

}

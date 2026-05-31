package main

import (
	"context"
	"fmt"

	"github.com/FeshLig/metrcollector/internal/agent"
	"github.com/FeshLig/metrcollector/internal/repository"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {

	printBuildInfo()

	options := agent.GetOptions()

	storage := repository.NewMemStorage()

	ctx := context.Background()

	agent.RunSender(ctx, storage, options)

}

func printBuildInfo() {
	version := buildVersion
	if version == "" {
		version = "N/A"
	}

	date := buildDate
	if date == "" {
		date = "N/A"
	}

	commit := buildCommit
	if commit == "" {
		commit = "N/A"
	}

	fmt.Printf("Build version: %s\n", version)
	fmt.Printf("Build date: %s\n", date)
	fmt.Printf("Build commit: %s\n", commit)
}

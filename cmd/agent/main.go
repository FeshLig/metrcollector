package main

import (
	"net/http"
	"time"

	"github.com/FeshLig/metrcollector/internal/agent"
	"github.com/FeshLig/metrcollector/internal/repository"
)

func main() {
	pollInterval := time.Duration(2)
	reportInterval := time.Duration(10)

	storage := repository.NewMemStorage()
	client := &http.Client{}

	collector := agent.NewMetricCollector(storage)
	sender := agent.NewSender("http://localhost:8080", client)

	go func() {
		for {
			collector.CollectMetrics()
			time.Sleep(pollInterval * time.Second)
		}
	}()

	go func() {
		for {
			sender.SendMetrics(storage)
			time.Sleep(reportInterval * time.Second)
		}
	}()

	select {}

}

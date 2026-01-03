package main

import (
	"net/http"
	"time"

	"github.com/FeshLig/metrcollector/internal/agent"
	"github.com/FeshLig/metrcollector/internal/repository"
)

func main() {

	flags := ParseFlags()

	storage := repository.NewMemStorage()
	client := &http.Client{}

	collector := agent.NewMetricCollector(storage)
	sender := agent.NewSender("http://"+flags.address.String(), client)

	go func() {
		for {
			collector.CollectMetrics()
			time.Sleep(flags.pollInterval * time.Second)
		}
	}()

	go func() {
		for {
			sender.SendMetrics(storage)
			time.Sleep(flags.reportInterval * time.Second)
		}
	}()

	select {}

}

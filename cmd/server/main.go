package main

import (
	"fmt"
	"net/http"

	"github.com/FeshLig/metrcollector/internal/handler"
	"github.com/FeshLig/metrcollector/internal/metric"
)

func main() {

	memStorage := metric.NewMemStorage()
	mux := http.NewServeMux()

	gaugeHandler := handler.NewGaugeHandler(memStorage.GetGauges())
	counterHandler := handler.NewCounterHandler(memStorage.GetCounters())

	updateHandler := handler.NewUpdateHandler(gaugeHandler, counterHandler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain")

		fmt.Fprintln(w, "GAUGES:")
		for name, gauge := range memStorage.GetGauges() {
			fmt.Fprintf(w, "%s = %f\n", name, gauge)
		}

		fmt.Fprintln(w)
		fmt.Fprintln(w, "COUNTERS:")
		for name, counter := range memStorage.GetCounters() {
			fmt.Fprintf(w, "%s = %d\n", name, counter)
		}

	})

	mux.HandleFunc("/update/", http.HandlerFunc(updateHandler.UpdatePage))

	err := http.ListenAndServe(`localhost:8080`, mux)

	if err != nil {
		panic(err)
	}
}

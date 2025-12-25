package main

import (
	"net/http"

	"github.com/FeshLig/metrcollector/internal/handler"
	"github.com/FeshLig/metrcollector/internal/repository"
)

func main() {

	memStorage := repository.NewMemStorage()
	mux := http.NewServeMux()

	gaugeHandler := handler.NewGaugeHandler(memStorage)
	counterHandler := handler.NewCounterHandler(memStorage)

	updateHandler := handler.NewUpdateHandler(gaugeHandler, counterHandler)
	rootHandler := handler.NewRootHandler(memStorage)

	mux.HandleFunc("/update/", http.HandlerFunc(updateHandler.UpdatePage))
	mux.HandleFunc("/", http.HandlerFunc(rootHandler.RootPage))

	err := http.ListenAndServe(`localhost:8080`, mux)

	if err != nil {
		panic(err)
	}
}

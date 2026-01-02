package main

import (
	"net/http"

	"github.com/FeshLig/metrcollector/internal/handler"
	"github.com/FeshLig/metrcollector/internal/repository"
)

func main() {
	Run()
}

// TODO:
// Добавить возврат ошибки
// Убрать панику
func Run() {
	memStorage := repository.NewMemStorage()
	mux := http.NewServeMux()

	updateHandler := handler.NewUpdateHandler(memStorage)
	rootHandler := handler.NewRootHandler(memStorage)

	mux.HandleFunc("/update/", http.HandlerFunc(updateHandler.UpdatePage))
	mux.HandleFunc("/", http.HandlerFunc(rootHandler.RootPage))

	err := http.ListenAndServe(`localhost:8080`, mux)

	if err != nil {
		panic(err)
	}
}

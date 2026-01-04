package main

import (
	"github.com/FeshLig/metrcollector/internal/handler"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	Run()
}

// TODO:
// Добавить возврат ошибки
// Убрать панику
func Run() {

	flags := ParseFlags()

	r := gin.Default()

	memStorage := repository.NewMemStorage()
	updateHandler := handler.NewUpdateHandler(memStorage)
	rootHandler := handler.NewRootHandler(memStorage)
	valueHandler := handler.NewValueHandler(memStorage)

	r.GET("/", rootHandler.RootPage)
	r.POST("/update/:type/:name/:value/", updateHandler.UpdatePage)
	r.GET("/value/:type/:name/", valueHandler.ValuePage)

	r.Run(flags.address.String())

}

package router

import (
	"github.com/FeshLig/metrcollector/internal/handler"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/gin-gonic/gin"
)

func NewRouter(memStorage *repository.MemStorage) *gin.Engine {

	r := gin.Default()
	r.LoadHTMLGlob("./internal/templates/*")

	updateHandler := handler.NewUpdateHandler(memStorage)
	rootHandler := handler.NewRootHandler(memStorage)
	valueHandler := handler.NewValueHandler(memStorage)

	r.GET("/", rootHandler.RootPage)
	r.POST("/update/:type/:name/:value/", updateHandler.UpdatePage)
	r.GET("/value/:type/:name/", valueHandler.ValuePage)

	return r

}

package router

import (
	"github.com/FeshLig/metrcollector/internal/handler"
	"github.com/FeshLig/metrcollector/internal/middleware"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter(memStorage *repository.MemStorage) *gin.Engine {

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}

	defer logger.Sync()

	r := gin.New()
	r.Use(middleware.Logger(logger), gin.Recovery())

	r.LoadHTMLGlob("./internal/templates/*")

	updateHandler := handler.NewUpdateHandler(memStorage)
	rootHandler := handler.NewRootHandler(memStorage)
	valueHandler := handler.NewValueHandler(memStorage)

	r.GET("/", rootHandler.RootPage)
	r.POST("/update/:type/:name/:value/", updateHandler.UpdatePage)
	r.GET("/value/:type/:name/", valueHandler.ValuePage)

	return r

}

package router

import (
	"github.com/FeshLig/metrcollector/internal/handler"
	"github.com/FeshLig/metrcollector/internal/middleware"
	"github.com/FeshLig/metrcollector/internal/repository"
	"github.com/FeshLig/metrcollector/internal/service"
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

	service := service.NewMetricService(memStorage)

	rootHandler := handler.NewRootHandler(memStorage)

	updateJSONHandler := handler.NewUpdateJSONHandler(service)
	valueJSONHandler := handler.NewValueJSONHandler(service)

	updateURLHandler := handler.NewUpdateURLHandler(service)
	valueURLHandler := handler.NewValueURLHandler(service)

	r.GET("/", rootHandler.RootPage)
	r.POST("/update/", updateJSONHandler.UpdateFromJSON)
	r.POST("/value/", valueJSONHandler.UpdateFromJSON)
	r.POST("/update/:type/:name/:value/", updateURLHandler.UpdateFromURL)
	r.GET("/value/:type/:name/", valueURLHandler.ValueFromURL)

	return r

}

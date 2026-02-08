package router

import (
	"github.com/FeshLig/metrcollector/internal/handler"
	"github.com/FeshLig/metrcollector/internal/middleware"
	"github.com/FeshLig/metrcollector/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter(service service.MetricsService) *gin.Engine {

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}

	defer logger.Sync()

	r := gin.New()
	r.Use(middleware.Logger(logger), middleware.Gzip(), gin.Recovery())

	r.LoadHTMLGlob("./internal/templates/*")

	rootHandler := handler.NewRootHandler(service)

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

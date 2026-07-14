// Package router configures HTTP routes
// and middleware for the application.
package router

import (
	"crypto/rsa"

	"github.com/FeshLig/metrcollector/internal/config"
	"github.com/FeshLig/metrcollector/internal/handler"
	"github.com/FeshLig/metrcollector/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewRouter creates and configures HTTP router.
func NewRouter(handlers *handler.Handlers, cfg config.Options, logger *zap.Logger, privateKey *rsa.PrivateKey) *gin.Engine {

	r := gin.New()

	middlewares := []gin.HandlerFunc{
		middleware.Logger(logger),
	}
	if privateKey != nil {
		middlewares = append(middlewares, middleware.Decrypt(privateKey))
	}
	middlewares = append(
		middlewares,
		middleware.Gzip(),
		middleware.Hash(cfg.Key.String()),
		middleware.TrustedSubnet(cfg.TrustedSubnet),
		gin.Recovery(),
	)

	r.Use(middlewares...)

	if gin.Mode() != gin.TestMode {
		r.LoadHTMLGlob("./internal/templates/*")
	} else {
		r.LoadHTMLGlob("../../internal/templates/*")
	}

	r.GET("/", handlers.Root.RootPage)
	r.POST("/update/", handlers.UpdateJSON.UpdateFromJSON)
	r.POST("/value/", handlers.ValueJSON.UpdateFromJSON)
	r.POST("/update/:type/:name/:value/", handlers.UpdateURL.UpdateFromURL)
	r.GET("/value/:type/:name/", handlers.ValueURL.ValueFromURL)
	r.GET("/ping/", handlers.Ping.PingPage)
	r.POST("/updates/", handlers.Updates.Updates)

	return r

}

package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FeshLig/metrcollector/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logs bytes.Buffer

	encoderCfg := zap.NewProductionEncoderConfig()

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(&logs),
		zap.InfoLevel,
	)

	logger := zap.New(core)

	r := gin.New()
	r.Use(middleware.Logger(logger))

	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "hello")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	logOutput := logs.String()

	require.Contains(t, logOutput, "HTTP Request")
	require.Contains(t, logOutput, "HTTP Response")

	require.Contains(t, logOutput, `"method":"GET"`)
	require.Contains(t, logOutput, `"uri":"/test"`)

	require.Contains(t, logOutput, `"status":200`)
}

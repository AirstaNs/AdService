package httpgin

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoggerMiddleware(t *testing.T) {
	logOutput := bytes.Buffer{}
	logger := log.New(&logOutput, "", 0)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/test", nil)

	middleware := LoggerMiddleware(logger)
	middleware(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/test", c.Request.URL.Path)
	assert.NotEqual(t, len(logOutput.String()), 0)
}

func TestRecoveryMiddleware(t *testing.T) {
	logOutput := bytes.Buffer{}
	logger := log.New(&logOutput, "", 0)

	router := gin.New()
	router.Use(RecoveryMiddleware(logger))
	router.GET("/panic", func(c *gin.Context) {
		panic("something went wrong")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "/panic", req.URL.Path)
	assert.NotEqual(t, len(logOutput.String()), 0)
}
